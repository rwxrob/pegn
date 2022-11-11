/*

Package scan contains a large library of scan functions for use by those
developing PEG-centric recursive descent parsers.

Simple and Adequately Performant

A balance between performance and simplicity has been the primary goal
with this package. The thought is that the time saved quickly writing
recursive descent parser prototypes using this package can be applied
later to convert these into a more performant forms having done the hard
initial work of figuring out the grammar notation and parsing algorithm
required.

The pegn.Scanner contains the last rune scanned and both it's beginning
and end index in the bytes buffer. Using a cursor in this way adds
several steps of indirection beyond the function calls themselves, but
this makes for scanners that are very easy to quickly write, trace, test,
and maintain.

*/
package pegn

import (
	"encoding/json"
	"fmt"
)

// Cursor points a specific location in the bytes buffer. The order of
// fields is guaranteed to never change.
//
// PEG scanning algorithms usually snap back to previous positions in
// the buffer. Therefore, it is important that the last scanned rune
// (Rune), pointer position (P), and last pointer position (LP) be
// restored when snapping back. The Mark/Goto methods are required to
// facilitate these operations consistently across implementations
// (although they may other changes depending on other state such
// implementations choose to add in addition to this).
//
type Cursor struct {
	Buf *[]byte // pointer to actual bytes buffer
	R   rune    // last rune scanned
	B   int     // beginning of last rune scanned
	E   int     // effective end of last rune scanned (beginning of next)
}

// String implements fmt.Stringer with the last rune scanned (R/Rune),
// and the beginning and ending byte positions joined with
// a dash.
func (c Cursor) String() string {
	return fmt.Sprintf("%q %v-%v", c.R, c.B, c.E)
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
type Scanner interface {
	Bytes() *[]byte
	Buffer(input any) error
	Scan() bool

	ScannerState
	ScannerCursor
	ScannerRangeCopy
	ScannerObservability
	//	ErrStack
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
// Mark() Cursor
//
// Mark returns a Cursor struct pointing to the last Rune, and it's
// location. Pass this to Goto to jump to another position in the bytes
// buffer easily.
//
// Goto(a Cursor)
//
// Jumps to a specific position in the bytes buffer and sets the last
// rune scanned as well.
//
type ScannerCursor interface {
	Mark() Cursor
	Goto(a Cursor)
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
	CopyEE(to Cursor) string
	CopyBE(to Cursor) string
	CopyBB(to Cursor) string
	CopyEB(to Cursor) string
}

// Errors allow Scanner to keep track of errors and decide how many to
// allow before stopping (whatever that means for the struct
// implementing ErrStack). SetMaxErr is called by the highest level
// caller in order to trigger a panic once that many errors have been
// pushed onto the stack. Generally, implementations of ErrStack will
// never panic unless MaxErr is reached and will expect and defer such
// panics by design.
type ErrStack interface {
	SetMaxErr(i int)  // sets max at which scanner will panic
	Errors() *[]error // returns pointer to internal errors stack
	ErrPush(e error)  // push new error onto stack
	ErrPop() error    // pop most recent error from stack
	Error() string    // output all errors as text (detects language)
}

// Rule contains both specific and user-friendly instruction as to
// what it will scan and optionally parse and why it might fail
// providing the most detail possible to the end user. Each Rule
// handles a PEGN grammar rule definition. Rules are combined and
// aggregated to create larger grammars.
//
// Implementations must define Scan and Parse. Scan is for validation
// while advancing the scanner. Parse productively produces a single
// node either with a value or the first node under it. Parse must
// return a node with either a value (V) or the first node under it (U),
// never anything else. When parsing higher-level content the
// first node under (U) will obviously have more of its fields assigned.
// Note that an empty node value (V) is still a value and fulfills the
// requirement for rules that may be simple validations of state at the
// point in the scan. In such cases, the node type (T) is the only
// useful field set. As such, all nodes returned by a Parse function must
// assign a unique type integer to the node corresponding with that
// specific rule even if the rule is similar to another.
//
// Even if Scan doesn't advance the scanner or parse any runes (i.e.
// "rhetorical" assertions, look aheads, etc.) a Scan and Parse method
// must always be defined.
//
// The Ident would normally be the identifier for a given definition
// (on the left side of the PEGN arrow). Idents must not
// conflict with other identifiers within a given grammar scope, but
// specifying the method of grouping into a grammar is left up to the
// specification creation and developer to decide. It is, however,
// strongly recommended to define all rules for a given grammar into the
// same package so that naming conflicts can easily be identified.
//
// When picking a Ident use the same PEGN rule naming conventions:
//
//     * CAPS for tokens
//     * lower for classes (with MixedCase Alias for convenience)
//     * MixedCase for everything else
//
// Sometimes two rules might exist for parsing the same content in
// different ways. For example, a Title rule could return the entire
// text of a title as its value, while a TitleWords rule could return
// a node with that exact same text split into individual Word nodes
// added under it.
//
// When selecting a name for the actual struct instance implementation
// a Rule consider the following conventions:
//
//     * Use Alias for classes (ws -> pegn.WhiteSpace)
//     * CAPS for tokens (SP -> pegn.SP)
//     * Mixed case for everything else (TitleWords -> pegn.TitleWords)
//
// For international language support Description should detect the host
// language of the system and/or user. Other fields are not language
// specific.
//
// Deprecation
//
// In order to maintain maximum compatibility for all dependencies, rule
// implementations must never be removed or change their Ident,
// Alias, or Type properties. Also, no two rules must ever have the
// same Ident, Alias, or Type. The PEGN, Scan, and Parse
// implementations may change, however, provided the results are
// consistent.
//
type Rule interface {
	Type() int             // unique, associate a node to its parsing rule
	Ident() string         // identifier from PEGN following conventions
	Alias() string         // mostly for classes (all-lower) names
	Description() string   // human PEGN with language detection
	PEGN() string          // formal PEGN specification with captures
	Scan(s Scanner) bool   // advance scanner if true, push errors if false
	Parse(s Scanner) *Node // advance and return parsed content
}

// Node is a typical node in a rooted node tree as required for any
// abstract syntax tree. Note that PEGN does not allow node attributes
// of any kind and that Nodes with Nodes under them must not also have
// a value (V).
//
// Minimal Design for Embedding
//
// The struct design is deliberately minimal with only Nodes, and String
// marshalling methods to ensure the least amount of conflict with
// potential embedded dependencies, which are encouraged to provide
// more involved handling of AST and other node trees when needed. Only
// unmarshaling methods are being considered at the moment (JSON, etc.).
//
// Unique Rule Type Association
//
// The Type integer often corresponds to the name (identifier) of the
// PEGN rule used to parse the Node, but not necessarily. For example,
// some rules are simply assertions that do not capture nodes. Others
// are not significant at all. Some don't even look at the content but
// examine the state of the scanner itself (Finished, etc.).
//
type Node struct {
	T int    // node type (linked to Rule.Type)
	V string // value (leaf nodes only)
	O *Node  // node over this node
	U *Node  // first node under this node (child)
	R *Node  // node to immediate right
}

// Nodes walks from the first node under (U) to the last returning a slice.
func (n Node) Nodes() []*Node {
	if n.U == nil {
		return nil
	}
	var nodes []*Node
	for cur := n.U; cur.R != nil; cur = cur.R {
		nodes = append(nodes, cur)
	}
	return nodes
}

// String fulfills the fmt.Stringer interface as JSON omitting any empty
// value (V) or slice of nodes under (N).
func (n Node) String() string {
	s := struct {
		T int
		V string  `json:",omitempty"`
		N []*Node `json:",omitempty"`
	}{n.T, n.V, n.Nodes()}
	byt, _ := json.Marshal(s)
	return string(byt)
}
