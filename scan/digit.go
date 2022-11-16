package scan

import (
	"github.com/rwxrob/pegn"
	"github.com/rwxrob/pegn/is"
	"github.com/rwxrob/pegn/rule/id"
)

func C_digit(s pegn.Scanner) bool {
	m := s.Mark()
	if !s.Scan() {
		return false
	}
	if is.C_digit(s.Rune()) {
		return true
	}
	return s.Revert(m, id.C_digit)
}
