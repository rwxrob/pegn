package scanner

// Cursor points to the beginning and ending of the last rune scanned
// and includes a copy of that rune. Fulfills the pegn.Cursor interface.
type Cursor struct {
	P
}
