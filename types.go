package pegn

import (
	"fmt"
	"io"
	"os"
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

	// bytes buffer being scanned (PEG/N assumes infinite memory unlike
	// older single-byte lookahead designs).

	Bytes() []byte
	SetBytes(b []byte) // MUST replace bytes buffer and reset Rune*

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

	Finished() bool

	// MUST return true if nothing has yet been scanned.

	Beginning() bool

	// MUST provide a way to print the current position and state using
	// the fmt package (Print) and the package log (Log) to facilitate
	// ScanFunc development and example-based testing

	Print()
	Log()
}

// Errors allow Scanner to keep track of errors and decide how many to
// allow before stopping (whatever that means for the struct
// implementing ErrStack).
type ErrStack interface {
	MaxErr() int                // return max count of errs before quit
	SetMaxErr(i int)            // sets max at which scanner will
	Errors() *[]error           // returns pointer to internal errors slice
	Error(a string)             // creates error and pushes
	Errorf(fm string, v ...any) // creates formatted error and pushes
	ErrPush(e error)            // push error onto stack
	ErrPop() error              // pop last error from stack
	ErrShift() error            // shift oldest error from bottom of stack
	ErrUnshift(e error)         // add new oldest error to bottom of stack
}

// TypeMap allows AST types to be represented by one or two digits
// making mandatory JSON String() output shorter.
type TypeMap map[string]int

// Node abstracts the nodes parsed with a Parser for use in any AST
// implementation that supports PEGN TypeMaps and attribute-less Node
// design. All PEGN Nodes MUST NOT have attributes. A branch (node with
// nodes (N) under it) SHOULD NOT have a value (V) as well. The String
// (see fmt.Stringer) MUST follow https://pegn.dev/schema-node which
// requires "T" and omits empty "V" and "N" if values are empty. Unlike
// other designs that might combine "V" and "N" into a single "C", this
// design does not mix types requiring a type check when using the
// values. See tests in this package for examples.
type Node struct {
	T string // node type
	V string // node value (leafs only)
	N []Node // nodes under (branches only)
}

// ScanFunc validates an element of grammar optionally capturing runes
// as it scans into the *[]rune if defined. A rune slice is used instead
// of a simple string to allow optimized sizing with make before passing
// which is not possible with a string type. PEGN uses angle brackets to
// indicate the beginning and ending of parsed regions. Sometimes
// everything being scanned will also be parsed, other times just one
// range, and somethings -- albeit rarely -- multiple ranges.
type ScanFunc func(s Scanner, parsed *[]rune) bool

// ParseFunc returns a parsed Node or nil if unable to parse.
type ParseFunc func(s Scanner) Node

// Parser contains both specific and user-friendly instruction as to
// what it will scan and optionally parse and why it might fail
// providing the most detail possible to the end user. Each Parser
// handles a PEGN grammar definition. Parsers are combined and
// aggregated to create larger grammars.
//
// Implementations MUST define Scan and Parse. (Usually, Parse is
// a simple Scan wrapper with details about the size of parsed buffer to
// be passed for capture.)
//
// Even if Scan doesn't advance the scanner or parse any runes (i.e.
// "rhetorical" assertions, look aheads, etc.) a Parser MUST always be
// defined.
//
// Implementations SHOULD define a Node function for returning nodes
// suitable for inclusion in an abstract syntax tree. An AST is usually
// far more valuable than simple returning the parsed string, but is not
// always needed or justified. Therefore, defining a Node function is
// not required (but strongly recommended).
//
// The Ident  would normally be the identifier for a given definition
// (on the left side of the PEGN <- or <-- arrows). Identifiers MUST NOT
// conflict with other identifiers within a given grammar scope, but
// specifying the method of grouping into a grammar is left up to the
// specification creation and developer to decide. It is, however,
// strongly recommended to define all Parsers for a given grammar into the
// same package so that naming conflicts can easily be identified.
//
// For international language support the Info simply need be set for
// the target language.
type Parser struct {
	Type  int    //
	Ident string // optional PEGN identifier (use naming conventions)
	PEGN  string // PEGN definition (ex: '#' SP < rune{1,70} >)
	Info  string // user-friendly explanation of PEGN definition
	Scanner
	Parser
}

type Parser interface {
	Parse(s Scanner) Node
}

type Scanner interface {
	Scan(s Scanner) bool
}

// Error returns a ScanError for self.
func (p Parser) Error(s Scanner) error { return ScanError{s, p} }

// ScanFile reads the specified file and SetBytes it into memory buffer,
// then scans and returns true if scan passed (advancing Scanner). Note
// that an error is only returned with problems loading the file.
// A false scan result also produces a ScanError.
func (p *Parser) ScanFile(path string, s Scanner) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	b, err := io.ReadAll(f)
	if err != nil {
		return false, err
	}
	s.SetBytes(b)
	if !p.Scan(s, nil) {
		return false, p.Error(s)
	}
	return true, nil
}

// ParseFile reads the specified file into Scanner memory buffer,
// then parses and returns the string.
func (p *Parser) ParseFile(path string, s Scanner) (string, error) {
	// TODO
	return "", nil
}

// NodeFile(path string, s Scanner) (Node, error)

// ScanReader(in io.Reader, s Scanner) bool

// ParseReader(in io.Reader, s Scanner) (string, error)

// NodeReader(in io.Reader, s Scanner) (Node, error)

type ScanError struct {
	S Scanner // scanner in state when error occurred
	P Parser  // which parser was being used during scan
}

var ScanErrorFormat = "syntax error at %v: expected %v \n\n%v"

func (e ScanError) Error() string {
	return fmt.Sprintf(ScanErrorFormat, e.S.Mark(), e.P.PEGN, e.P.Info)
}
