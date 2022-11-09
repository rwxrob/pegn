package pegn

var Cws Rule = ws{}

type ws struct{}

func (ws) Name() string        { return `ws` }
func (ws) PEGN() string        { return `` }
func (ws) Description() string { return `` }

func (ws) Scan(s Scanner) bool {
	// TODO
	return false
}

func (ws) Parse(s Scanner) Node {
	// TODO
	return nil
}
