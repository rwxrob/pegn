package scan

import (
	"github.com/rwxrob/pegn/T"
	"github.com/rwxrob/pegn/is"
)

func Field(s Scanner) bool {
	m := s.Mark()
	var c int
	for !s.Peek(" ") && is.Uprint(s.Rune()) {
		c++
	}
	if c > 0 {
		return true
	}
	return s.Revert(m, T.C_uprint)
}
