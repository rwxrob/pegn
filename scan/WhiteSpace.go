package scan

import "github.com/rwxrob/pegn/T"

func WhiteSpace(s Scanner) bool {
	m := s.Mark()
	if !s.Scan() {
		return false
	}
	switch s.Rune() {
	case ' ', '\t', '\n', '\r':
		return true
	}
	return s.Revert(m, T.WhiteSpace)
}
