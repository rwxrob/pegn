package scan

import (
	"github.com/rwxrob/pegn"
)

// ws <- SP / TAB / CR / LF
func Is_ws(s pegn.Scanner) bool {
	m := s.Mark()
	if !s.Scan() {
		return false
	}
	switch s.Rune() {
	case ' ', '\t', '\r', '\n':
		return true
	}
	s.Goto(m)
	return false
}

// ws*
func SomeWS(s pegn.Scanner) bool {
	m := s.Mark()
	if !WS(s) {
		s.Goto(m)
		return false
	}
	for WS(s) {
	}
	return true
}

// EndLine <- LF / CRLF
func EndLine(s pegn.Scanner) bool {
	m := s.Mark()
	s.Scan()
	switch s.Rune() {
	case '\n':
		return true
	case '\r':
		if s.Scan() && s.Rune() == '\n' {
			return true
		}
	}
	s.Goto(m)
	return false
}

// EndPara <- ws* (!. / EndLine !. / EndLine{2})
func EndPara(s pegn.Scanner) bool {
	m := s.Mark()
	var found bool
TOP:
	{
		switch {
		case s.Finished():
			found = true
		case EndLine(s) && s.Finished():
			found = true
		case EndLine(s) && EndLine(s):
			found = true
			goto TOP
		case SomeWS(s):
			goto TOP
		}
	}
	if !found {
		s.Goto(m)
	}
	return found
}
