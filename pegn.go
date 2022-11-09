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
	"fmt"
)

// Cursor points a specific location in the bytes buffer. The order of
// fields is guaranteed to never change.
//
// PEG scanning algorithms usually snap back to previous positions in
// the buffer. Therefore, it is important that the last scanned rune
// (Rune), pointer position (P), and last pointer position (LP) be
// restored when snapping back. The Mark/Snap methods are required to
// facilitate these operations consistently across implementations
// (although they may other changes depending on other state such
// implementations choose to add in addition to this).
//
type Cursor struct {
	R rune // last rune scanned
	B int  // beginning of last rune scanned
	E int  // effective end of last rune scanned (beginning of next)
}

// String implements fmt.Stringer with the last rune scanned (R/Rune),
// and the beginning and ending byte positions joined with
// a dash.
func (c Cursor) String() string {
	return fmt.Sprintf("%q %v-%v", c.R, c.B, c.E)
}

// A RuneScanner must employ design principles outlined by PEG including the
// buffering of all bytes that will be used during the scan.
//
// Note that while movement around within the scanner bytes buffer is
// meant to be done with Mark and Goto, scanner implementations are
// allowed --- even encouraged --- to expose their internal struct
// fields in addition to fulfilling this interface in order to provide
// higher performance at the cost of interface abstraction when such
// is required.
//
// See pegn/scanner for a usable implementation.
//
type Scanner interface {

	// Buffer being scanned (PEG/N assumes infinite memory unlike
	// older single-byte lookahead designs). Bytes are most efficiently
	// set this way. Use Buffer for convenience at a higher-level.

	Bytes() []byte
	SetBytes(b []byte)      // MUST replace bytes buffer and reset Rune*
	Buffer(input any) error // string=path, []byte, io.Reader

	// current internal cursor values

	Rune() rune // copy of last rune scanned
	RuneB() int // beginning of last rune scanned
	RuneE() int // end of last rune scanned (beginning of next to scan)

	// forces developer to consider implications of movement around within
	// the bytes buffer without necessarily being able to update all the
	// cursor propreties

	Mark() Cursor
	Goto(a Cursor)

	// MUST scan the next UNICODE code point beginning with the current
	// index position and if (and only if) successful capture it (Rune)
	// and preserve the previous index position (P->LP). If unable to scan
	// a rune for any reason, including being out of data to scan, MUST
	// return false. This allows typical scanner loop idioms:
	//
	//     for s.Scan() {
	//         ...
	//     }
	//
	// MUST return false if ErrStack.MaxErr has been reached.

	Scan() bool
	ScanBytes(a []byte) bool
	ScanString(a string) bool
	ScanRunes(a []rune) bool
	ScanRune(a rune) bool
	ErrStack

	// MUST return true if there is nothing left to scan.
	// MUST return true if nothing has yet been scanned.

	Finished() bool
	Beginning() bool

	// MUST provide a way to print the current position and state using
	// the fmt package (Print) and the package log (Log) to facilitate
	// ScanFunc development and example-based testing

	Print()
	Log()

	// Buffering allows implementations of angle-bracket captures and
	// rule-scoped variables.

	Buffers
}

type Buffers interface {
	BufBeg(name string)        // opens and begins caching buffer
	BufEnd(name string)        // closes and frees resources for buffer
	BufMap() map[string][]byte // return internal buffer map
}

// Errors allow Scanner to keep track of errors and decide how many to
// allow before stopping (whatever that means for the struct
// implementing ErrStack).
type ErrStack interface {
	MaxErr() int                // return max count of errs before quit
	SetMaxErr(i int)            // sets max at which scanner will
	Errors() *[]error           // returns pointer to internal errors slice
	Errorf(fm string, v ...any) // creates formatted error and pushes
	ErrPush(e error)            // push error onto stack
	ErrPop() error              // pop last error from stack
	ErrShift() error            // shift oldest error from bottom of stack
	ErrUnshift(e error)         // add new oldest error to bottom of stack
	Error() string              // formatted error text (detects language)
}

// Rule contains both specific and user-friendly instruction as to
// what it will scan and optionally parse and why it might fail
// providing the most detail possible to the end user. Each Rule
// handles a PEGN grammar rule definition. Rules are combined and
// aggregated to create larger grammars.
//
// Implementations MUST define Scan and Parse. Usually, Parse is
// just a Scan with capture enabled on the scanner.
//
// Even if Scan doesn't advance the scanner or parse any runes (i.e.
// "rhetorical" assertions, look aheads, etc.) a Scan and Parse method
// MUST always be defined.
//
// The Name would normally be the identifier for a given definition
// (on the left side of the PEGN <- arrow). Names MUST NOT
// conflict with other identifiers within a given grammar scope, but
// specifying the method of grouping into a grammar is left up to the
// specification creation and developer to decide. It is, however,
// strongly recommended to define all Rules for a given grammar into the
// same package so that naming conflicts can easily be identified.
//
// For international language support Description SHOULD detect the host
// language of the system and/or user. Other fields are not language
// specific.
//
type Rule interface {
	Name() string         // unique name within grammar namespace
	Description() string  // human PEGN with language detection
	PEGN() string         // formal PEGN specification with captures
	Scan(s Scanner) bool  // advance scanner if true
	Parse(s Scanner) Node // advance and return parsed content
}

// Node abstracts a typical node in a rooted node tree as required for
// any abstract syntax tree. Note that PEGN does not allow node
// attributes of any kind and that Nodes with Nodes under them MUST NOT
// also have a Value.
//
// The Type integer often corresponds to the name (identifier) of the
// PEGN rule used to parse the Node, but not necessarily. For example,
// some rules are simply assertions that do not capture nodes. Others
// are not significant at all. Some don't even look at the content but
// examine the state of the scanner itself (Finished, etc.).
//
// Some rules will inject additional or change runes in the Node.Value
// as the rule is being evaluated.
//
type Node interface {
	Type() int        // often corresponds to Rule.Name (but not always)
	Value() string    // value for leaf nodes
	Over() Node       // node immediately above this node
	Under() []Node    // list of nodes immediately under this node
	LastUnder() Node  // most recently added under
	FirstUnder() Node // first added under
	Left() Node       // immediately connected to the left
	Right() Node      // immediate connected to the right
	RightMost() Node  // farthest peer to the right
	LeftMost() Node   // farthest peer to the left
}
