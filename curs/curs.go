/*

Package cursor is a basic struct containing the current rune and
position in a bytes slice buffer (to which it keeps a pointer as well).

PEG scanning algorithms usually snap back to previous positions in the
buffer. Therefore, it is important that the last scanned rune (Rune),
pointer position (P), and last pointer position (LP) be restored when
snapping back. The Mark/Goto Scanner interface methods are required to
facilitate these operations consistently across pegn.Scanner
implementations.

*/
package curs

import (
	"fmt"
)

// R contains a cursor pointer to a bytes slice buffer and information
// pointing to a specific location in the bytes buffer. The order of
// fields is guaranteed to never change.
type R struct {
	Buf *[]byte // pointer to actual bytes buffer
	R   rune    // last rune scanned
	B   int     // beginning of last rune scanned
	E   int     // effective end of last rune scanned (beginning of next)
}

// String implements fmt.Stringer with the last rune scanned (R/Rune),
// and the beginning and ending byte positions joined with
// a dash.
func (c R) String() string {
	return fmt.Sprintf("%q %v-%v", c.R, c.B, c.E)
}
