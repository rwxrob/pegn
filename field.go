package pegn

const FieldT = 2

var Field = _Field{}

type _Field struct{}

func (_Field) Type() int     { return FieldT }
func (_Field) Ident() string { return `Field` }
func (_Field) Alias() string { return `Field` }
func (_Field) PEGN() string  { return `(!SP uprint)+` }
func (_Field) Description() string {
	return `one or more printable UNICODE code points except space`
}
func (r _Field) Error() string {
	return "some error"
}

func (r _Field) Scan(s Scanner) bool {
	m := s.Mark()
	var c int
	for !s.Peek(" ") && Uprint.Scan(s) {
		c++
	}
	if c > 0 {
		return true
	}
	// TODO push error s.ErrPush(r)
	s.Goto(m)
	return false
}

func (r _Field) Parse(s Scanner) *Node {
	m := s.Mark()
	if !r.Scan(s) {
		return nil
	}
	return &Node{T: FieldT, V: s.CopyEE(m)}
}
