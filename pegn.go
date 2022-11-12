/*

Package pegn provides tools for working with grammars written in Parsing
Expression Grammar Notation (Rob Muhlestein's combination of ABNF with
Bryan Ford's original PEG design as found in the examples in his paper).

Tools

	  * scan/     - scan functions:  func(s Scanner) bool
		* parse/    - parse functions: func(s Scanner) *Node
		* is/       - class functions: func(a rune) bool
    * tk/       - token constants
    * types.go  - type constants
		* docs/     - multi-lingual rule documentation
    * scanner/  - pegn.Scanner implementation (DefaultScanner)
    * curs/     - curs.R struct for position within bytes buffer
    * cmd.go    - Bonzai stateful command tree composition
    * cmd/pegn  - helper utility with code generation

Types of Rules

A PEGN rule is what appears on the right of the arrow (<-). There are
two subtypes of a PEGN rule: class rules and a token rules. A class is
a set of runes. A token is a specific rune or string constant.

All rules must have a unique integer constant type.

PEGN reserved integer types (types.go) are guaranteed to never change
and always be negative. This frees grammar developers to use positive
integers. On day, the greater PEGN community may wish to organize range
reservations for different common grammar rules to maintain grammar and
AST interoperability.

No guarantee is made about what numeric range a rule, class, or token
integer will be only that any specific type integer will never be reused
for a different type. Use of such value ranges is strongly discouraged
over creation of proper range maps (as is used in the unicode package).

The integer 0 is reserved as Untyped.

Every PEGN rule must have a ScanFunc (scan) and a ParseFunc (parse).
Class rules must also have a ClassFunc (is). Token rules must also have
a token constant (tk).

Balanced Simplicity and Performance

A balance between performance and simplicity has been the primary goal
with this package. The thought is that the time saved quickly writing
recursive descent parser prototypes using this package can be applied
later to convert these into a more performant forms having done the hard
initial work of figuring out the grammar notation and parsing algorithm
required.

*/
package pegn

import (
	"github.com/rwxrob/pegn/curs"
	"github.com/rwxrob/pegn/gr"
	"github.com/rwxrob/pegn/rule"
	"github.com/rwxrob/pegn/scanner"
)

var DefaultScanner = scanner.New()

// ScanFunc uses the Scanner to scan for rule returning false and
// pushing at least one error to the Scanner.ErrStack. When working with
// PEGN grammars (or derivatives) the error must be of pegn.Error type.
type ScanFunc func(s Scanner) bool

// ParseFunc depends on one or more ScanFuncs to create a node pointer
// to return if successful pushing at least one error to the
// Scanner.ErrStack. When working with PEGN grammars (or derivatives)
// the error must be of pegn.Error type.
type ParseFunc func(s Scanner) bool

type Scanner interface {
	ScannerCore
	ScannerState
	ScannerCursor
	ScannerRangeCopy
	ScannerObservability
	ScannerErrors
}

// A Scanner implements a buffered rune scanner and must employ design
// principles outlined by PEG including the buffering of all bytes that
// will be used during the scan. See pegn/scanner for one usable
// implementation.
//
// Bytes() *[]bytes
//
// Returns the pointer to the bytes buffer being scanned providing
// direct access.  PEG/N assumes infinite memory unlike older
// single-byte lookahead designs. The byte slice can be modified without
// fear of disrupting anything else within the Scanner so long as the
// current position is updated appropriately. Bytes are most efficiently
// set this way. Use Buffer for convenience at a higher-level.
//
// Buffer(input any) error
//
// Must minimally accept a string, []byte, or io.Reader as input
// parameter and load that into the *[]bytes return by Bytes method.
//
// Scan() bool
//
// Scans the next UNICODE code point (rune) beginning at position RuneE
// in the Bytes buffer storing it into Rune and advancing RuneB and
// RuneE appropriately depending on the byte size of the rune (uint32,
// up to 4 bytes maximum). Scan must return false if failed to scan
// a rune for whatever reason. Scanning the end of buffer should never
// push an error to ErrStack. Scan is frequently used in the idiomatic
// loop fashion.
//
//     for s.Scan() {
//         ...
//     }
//
type ScannerCore interface {
	Bytes() *[]byte
	Buffer(input any) error
	Scan() bool
}

// The ScannerState interface provides convenience methods for writing
// grammar scan rules.
//
// Peek(a string) bool
//
// Peek returns true if the passed string matches from current position
// in the buffer (s.RuneB) forward. Returns false if the string
// would go beyond the length of buffer (len(s.Buf)). Peek does not
// advance the Scanner.
//
// Finished() bool
//
// Returns true if Scan would fail because there is nothing left to
// scan.
//
// Beginning() bool
//
// Returns true if no Scan has yet been called (identical to Rune ==
// `\x00` or RuneB == 0 && RuneE == 0).
//
type ScannerState interface {
	Peek(a string) bool
	Finished() bool
	Beginning() bool
}

// The ScannerCursor interface provides a one-rune cursor (1-4 bytes)
// that includes the position of the beginning and ending of the rune
// to allow quick bookmarking and repositioning within the bytes buffer.
//
// Rune() rune
//
// Returns a copy of the last rune scanned (or null `\x00` if nothing yet
// scanned).
//
// RuneB() int
//
// Returns the index in the bytes buffer pointing to the beginning of
// the last rune scanned (Rune)
//
// RuneE() int
//
// Returns the index in the bytes buffer pointing to end of the last
// rune scanned (Rune) and the beginning of the next rune to scan on
// next call to Scan.
//
// Mark() curs.R
//
// Mark returns a cursor pointing to the last Rune, and it's
// location. Pass this to Goto to jump to another position in the bytes
// buffer easily.
//
// Goto(a curs.R)
//
// Jumps to a specific position in the bytes buffer and sets the last
// rune scanned as well.
//
type ScannerCursor interface {
	Mark() curs.R
	Goto(a curs.R)
	Rune() rune
	RuneB() int
	RuneE() int
}

// The ScannerObservability interface ensures that all implementations
// have a consistent way to observe and test the Scanner state.
//
// SetViewLen(a int)
//
// Set the number of bytes from upcoming bytes buffer to display from
// String, Log, and Print.
//
// ViewLen() int
//
// Returns previous SetViewLen
//
// String() string
//
// Fulfills the fmt.Stringer interface. Must return the Cursor as
// a string, followed by a single space, followed by the quoted (%q)
// number of bytes set by ViewLen as a preview of what is next in the
// bytes buffer.
//
//    '\x00' 0-0 "some"
//    's' 0-1 "ome"
//    'e' 2-3 ""
//
// This output must be consistent to provide consistency across test
// code for all PEGN rule Scanner implementations.
//
// Print()
//
// Prints the String() with the fmt package.
//
// Log()
//
// Logs the String() with the log package.
//
// TraceOn()
// TraceOff()
//
// Activate (deactivate) a Log call for ever call to Scan.
//
type ScannerObservability interface {
	SetViewLen(a int)
	ViewLen() int
	Print()
	Log()
	TraceOn()
	TraceOff()
}

// The ScannerRangeCopy interface allow the scanner to return a copy
// of the runes as a string from current scanner position to the cursor
// passed. If the cursor passed is before the current position the first
// cursor is the cursor passed. If the cursor passed is after the
// current position then first cursor is the current position.
//
// For more direct access to the bytes buffer use Bytes() *[]byte
// instead. The cursor passed must point to a valid position within the
// bytes buffer or call will panic.
//
// The qualifiers B (beginning) and E (ending) indicate whether the range
// is to the beginning or ending position (also B and E) of the cursor
// indicating if that the cursor's Rune is included or not:
//
//    (n,m] - EE
//    [n,m] - BE
//    [n,m) - BB
//    (n,m) - EB
//
type ScannerRangeCopy interface {
	ScannerCursor
	CopyEE(to curs.R) string
	CopyBE(to curs.R) string
	CopyBB(to curs.R) string
	CopyEB(to curs.R) string
}

// ScannerErrors allow Scanner to keep track of errors and decide how
// many to allow before stopping. SetMaxErr is called by the highest
// level caller in order to trigger a panic once that many errors have
// been pushed onto the stack. Generally, implementations should not
// panic unless max err is reached
//
// Even though any error type is used for these methods, the errors
// passed and produced should be instances of Error with both
// a type and cursor set when implementing scanners for use with this
// PEGN package or others. This is also why Expected takes a simple
// integer instead of a pegn.Type.
//
type ScannerErrors interface {
	SetMaxErr(i int)             // sets max at which scanner will panic
	Errors() *[]error            // returns pointer to internal errors stack
	ErrPush(e error)             // push new error onto stack
	ErrPop() error               // pop most recent error from stack
	Expected(t int) bool         // ErrPush + return false
	Revert(m curs.C, t int) bool // Goto(m) + Expected(t)
}

// Grammar represents all the rules of a PEGN-compatible specification.
type Grammar interface {
	Scan(in any, s Scanner) (bool, []error)
	Parse(in any, s Scanner) (*Node, []error)
	Rules() []rule.Rule
}

// Scan uses the DefaultScanner to load whatever input the Scanner
// implementation accepts to its Buffer method and then parses the PEGN
// spec string into a Grammar and delegates Scan to it.
func Scan(in any, spec string) (bool, []error) { return gr.PEGN.Scan(in, spec) }

// Parse does the same thing as Scan but produces a tree of parsed nodes
// instead of just validating.
func Parse(in any, spec string) (*Node, []error) { return gr.PEGN.Parse(in, spec) }

func ScanClass(s Scanner, f ClassFunc) bool {
	// TODO
	return false
}

func ParseClass(s Scanner, f ClassFunc) *Node {
	// TODO
	return nil
}
