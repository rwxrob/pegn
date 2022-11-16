package scan

import (
	"github.com/rwxrob/pegn"
	"github.com/rwxrob/pegn/is"
	"github.com/rwxrob/pegn/rule/id"
)

func MajorVer(s pegn.Scanner) bool {
	m := s.Mark()
	var c int
	for !s.Peek(" ") && s.Scan() && is.C_digit(s.Rune()) {
		c++
	}
	if c > 0 {
		return true
	}
	return s.Revert(m, id.MajorVer)
}
