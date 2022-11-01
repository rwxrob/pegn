package pegn

// PEG scanning algorithms usually snap back to previous positions in
// the buffer. Therefore, it is important that the last scanned rune
// (Rune), pointer position (P), and last pointer position (LP) be
// restored when snapping back. The Mark/Snap methods are required to
// facilitate these operations consistently across implementations
// (although they may other changes depending on other state such
// implementations choose to add in addition to this).
//
type Cursor interface {

	// current rune (last scanned)

	Rune() rune
	SetRune(a rune)

	// pointer to current byte index position in bytes buffer (which will
	// be the beginning of *next* Rune scanned)

	End() int
	SetEnd(i int)

	// pointer to beginning of previous scanned rune in bytes buffer
	// (which is the position of the current value assigned to Rune)

	Beg() int
	SetBeg(i int)

	// bookmark a position by creating a copy of a Cursor to jump to later

	Mark() Cursor
	Goto(a Cursor)
}

// A Scanner must employ design principles outlined by PEG including the
// buffering of all bytes that will be used during the scan. This
// requires a Cursor pointing to a specific location within the bytes.
//
// Note that scanner implementations are allowed --- even encouraged ---
// to expose their internal struct fields in addition to fulfilling this
// interface in order to provide higher performance at the cost of
// interface abstraction when such makes sense.
//
// See pegn/scanner for one available implementation.
//
type Scanner interface {
	Cursor

	// bytes buffer being scanned (PEG/N assumes infinite memory unlike
	// older single-byte lookahead designs)

	Bytes() []byte
	SetBytes(buf []byte)

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
}

type ScanFunc func(a Scanner) bool
