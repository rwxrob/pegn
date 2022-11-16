package is

// SP / TAB / LF / CR
var C_ws = func(r rune) bool {
	switch r {
	case ' ', '\t', '\n', '\r':
		return true
	default:
		return false
	}
}
