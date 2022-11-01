package scan

import "github.com/rwxrob/pegn"

func SomeWS(s pegn.Scanner) bool {
	m := s.Mark()
	var found bool
	for s.Scan() {
		switch s.Rune() {
		case ' ', '\t', '\r', '\n':
			found = true
			continue
		}
		break
	}
	if !found {
		s.Goto(m)
	}
	return found
}

// EndLine <- LF / CRLF
func EndLine(s pegn.Scanner) bool {
	m := s.Mark()
	var found bool
	s.Scan()
	switch s.Rune() {
	case '\n':
		found = true
	case '\r':
		if s.Scan() && s.Rune() == '\n' {
			found = true
		}
	}
	if !found {
		s.Goto(m)
	}
	return found
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
