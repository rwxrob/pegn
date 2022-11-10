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

	Scan() bool

	// MUST return false if ErrStack.MaxErr has been reached.

	//ErrStack

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

	//Buffers
}

type Buffers interface {
	BufBeg(name string)        // opens and begins caching buffer
	BufEnd(name string)        // closes and frees resources for buffer
	BufMap() map[string][]byte // return internal buffer map
}

// Errors allow Scanner to keep track of errors and decide how many to
// allow before stopping (whatever that means for the struct
// implementing ErrStack). SetMaxErr is called by the highest level
// caller in order to trigger a panic once that many errors have been
// pushed onto the stack. Generally, implementations of ErrStack will
// never panic unless MaxErr is reached and will expect and defer such
// panics by design.
type ErrStack interface {
	MaxErr() int        // return max count of errs before quit
	SetMaxErr(i int)    // sets max at which scanner will fail
	Errors() *[]error   // returns pointer to internal errors stack
	ErrPush(e error)    // push new error onto stack
	ErrPop() error      // pop most recent error from stack
	ErrShift() error    // shift first error from beginning of stack
	ErrUnshift(e error) // add new first error to beginning of stack
	Error() string      // formatted error text (detects language)
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
// Alias, or NodeType properties. Also, no two rules must ever have the
// same Ident, Alias, or NodeType. The PEGN, Scan, and Parse
// implementations may change, however, provided the results are
// consistent.
//
type Rule interface {
	Ident() string         // identifier from PEGN following conventions
	Alias() string         // mostly for classes (all-lower) names
	NodeType() int         // associate a node to its parsing rule
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
	T int    // node type (linked to Rule.NodeType)
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
