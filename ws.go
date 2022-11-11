package pegn

var WhiteSpace = _WhiteSpace{}

type _WhiteSpace struct{}

func (_WhiteSpace) Type() int     { return 1 }
func (_WhiteSpace) Ident() string { return `ws` }
func (_WhiteSpace) Alias() string { return `WhiteSpace` }
func (_WhiteSpace) PEGN() string  { return `SP / TAB / LF / CR` }

func (_WhiteSpace) Description() string {
	return `space, tab, line feed ('\n') or carriage return ('\r')`
}

func (ws _WhiteSpace) Scan(s Scanner) bool {
	m := s.Mark()
	if !s.Scan() {
		// TODO ErrPush
		return false
	}
	switch s.Rune() {
	case ' ', '\t', '\n', '\r':
		return true
	}
	// TODO ErrPush
	s.Goto(m)
	return false
}

func (ws _WhiteSpace) Parse(s Scanner) *Node {
	if !ws.Scan(s) {
		return nil
	}
	return &Node{T: ws.Type(), V: string(s.Rune())}
}
