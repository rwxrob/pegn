package pegn

import "fmt"

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

// A Scanner must employ design principles outlined by PEG including the
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
	// older single-byte lookahead designs), use Buffer or New to load
	// a new bytes buffer

	Bytes() []byte

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

	// MUST return true if there is nothing left to scan.

	Finished() bool

	// MUST provide a way to print the current position and state using
	// the fmt package (Print) and the package log (Log) to facilitate
	// ScanFunc development and example-based testing

	Print()
	Log()
}

type ScanFunc func(a Scanner) bool
