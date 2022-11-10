package pegn

var WhiteSpace = rule{
	ident: `ws`, alias: `WhiteSpace`, node: 1,
	pegn: `SP / TAB / LF / CR`,
	desc: `space, tab, line feed ('\n') or carriage return ('\r')`,
	scan: ScanWhiteSpace,
	parse: func(s Scanner) *Node {
		if ScanWhiteSpace(s) {
			return &Node{T: 1, V: string(s.Rune())}
		}
		return nil
	},
}

func ScanWhiteSpace(s Scanner) bool {
	m := s.Mark()
	if !s.Scan() {
		return false
	}
	switch s.Rune() {
	case ' ', '\t', '\n', '\r':
		return true
	}
	s.Goto(m)
	return false
}
