// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

/*

Package pegn implements the PEGN 2023-01 specification (pegn.dev) and
contains some helper packages and tooling to create grammars using the
PEGN language.

*/
package pegn

import (
	"fmt"
	"go/ast"

	"github.com/rwxrob/pegn/curs"
)

// ScanFunc uses the Scanner to scan for a rule and optionally populate
// a buffer that can then optionally be used by a ParseFunc to create
// and return a node. When a scan fails, it must push at least one
// pegn.Error to the Scanner.ErrStack and return false. If the buffer
// argument is nil, no buffering must occur. When buffer passed is not
// nil buffering must always occur even if the result of the scan is false.
// These functions should contain most of the work of scanning and
// parsing. See ParseFunc for wrapping this function with a buffer
// optimized for a particular node type and size.
type ScanFunc func(s Scanner, buf *[]rune) bool

// ParseFunc prepares a buffer argument optimized for the node type and
// size and calls one or more ScanFuncs passing the buffer to capture
// scan output and return a new implementation of pegn.Node created
// from that buffered output. While separate grammar package creators
// are welcome and encouraged to return *ast.Node from their own
// ParseFunc implementations many may opt to define their own
// variations on pegn.ParseFunc returning instances of their own
// implementations of pegn.Node. All variations, however, must return
// something that fulfills pegn.Node so that those working with nodes
// and node trees can reliably trust them to have the right methods when
// passed as pegn.Node arguments.
type ParseFunc func(s Scanner) *ast.Node

// ClassFunc return true if the rune is a member of the class. Most
// unicode.Is* functions are this type as are everything in the "is"
// subpackage.
type ClassFunc func(r rune) bool

// Node represents a single node in a rooted node tree graph structure.
// implementations will widely vary. Use of exported, single letter
// struct fields is recommended for those wishing to use the
// implementation directly without incurring the indirection hit in
// performance for these interface implemented method calls. See the ast
// package for an example implementation.
//
// Type() int
//
// Returns a unique type as an integer. For PEGN implementations this
// integer must be a valid pegn.Rule.ID (see pegn/rule/ids.go).
//
// Value() string
//
// Returns the value if it has one. Note that implementations should not
// return a value if any nodes have been added to this node.
//
// Node() Node
//
// Returns the node to which this node belongs, which could be nil.
//
// Add(a Node)
//
// Add the specified node under this node. The target node will have
// its Node value set to the caller.
//
// Destroy()
//
// Destroys self updating any internal relations to this node
// appropriately. For example, Nodes called on the parent of the
// destroyed node must no longer include in returned slice.
//
// Nodes() []Node
//
// Returns all nodes that have been added under the current node.
//
// String() string
//
// The fmt.Stringer interface must be implemented and must produce
// predictable, compact JSON output (calling MarshalJSON on self and converting
// to string). This output is critical and mandatory to ensure all node
// tree implementations match the same JSON schema. Tools and
// implementations should let conversion to more human-friendly and
// alternative formats (such as YAML) up to external packages and tools.
//
// Any error returned from MarshalJSON should force the string output to
// be a single JSON string containing the error with the mandatory
// "error: " prefix. JSON schema definitions must allow for this and
// assume a single string (as opposed to an object or array as is
// usually wanted) is potentially an error.
//
// MarshalJSON() ([]byte, error)
//
// All implementations must produce compact JSON that matches the
// following sample implementation default JSON marshaling tags:
//
//     type node struct {
//       T int     `json:"t"`           // type (rule id)
//       V string  `json:"v,omitempty"` // value (if leaf)
//       N []*node `json:"n,omitempty"` // nodes under (if over/parent)
//     }
//
// All implementations must fail and return an error if there is both
// a value (V) and one or more nodes under it (N).
//
// Producing identical, predictable JSON is critical to interoperability
// between applications using this node tree format. Expensive JSON
// schema validation is not needed and discouraged. Consider an
// intermediary struct to hold the values before outputting them as
// a string.
//
// UnmarshalJSON(b []byte) error
//
// Must unmarshal both compact and non-compact (human-friendly) forms.
// Must throw an error incoming data contains both a value (V) and nodes
// under it (N). See MarshalJSON for equivalent holding struct to
// validate before assigning the actual values.
//
type Node interface {
	Rule() int
	Value() string
	Node() Node
	Add(a Node)
	Destroy()
	Nodes() []Node
	String() string
	MarshalJSON() ([]byte, error)
	UnMarshalJSON(b []byte) error
}

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
// Usage
//
// It is recommended that developers copy this interface to their own
// implementations and refer back to it in order to prevent potential
// cyclical imports of other things within this package. This also
// reduces risk of incompatibility if and when changes to the base
// interface are made.
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
	Revert(m curs.R, t int) bool // Goto(m) + Expected(t)
	Error() string               // combine Errors() into single string
}

// Error wraps the type (T) and current scanner position (C)
// such that it can be located and displayed with help information by
// looking up those things from other sources when displayed to the end
// user. The position of fields is guaranteed never to change allowing
// for short-form instantiation (ex: pegn.Error{1,s.Mark()}). See
// ScannerErrors interface for more.
type Error struct {
	T int
	C curs.R
}

var DefaultErrFmt = `expecting %v at %v`

func (e Error) Error() string {
	return fmt.Sprintf(DefaultErrFmt, e.T, e.C)
}
