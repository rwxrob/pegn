package pegn

import "unicode"

const UprintT = 3

var Uprint = _Uprint{}

type _Uprint struct{}

func (_Uprint) Type() int     { return UprintT }
func (_Uprint) Ident() string { return `uprint` }
func (_Uprint) Alias() string { return `Uprint` }
func (_Uprint) PEGN() string {
	return `uletter / umark / unumber / upunct / usymbol`
}
func (_Uprint) Description() string {
	return `printable UNICODE code point (letter, mark, number, punctuation, or symbol)`
}

func (r _Uprint) Scan(s Scanner) bool {
	m := s.Mark()
	if s.Finished() {
		return false
	}
	s.Scan()
	if unicode.IsPrint(s.Rune()) {
		return true
	}
	// TODO push error
	s.Goto(m)
	return false
}

func (r _Uprint) Parse(s Scanner) *Node {
	if !r.Scan(s) {
		return nil
	}
	return &Node{T: UprintT, V: string(s.Rune())}
}
