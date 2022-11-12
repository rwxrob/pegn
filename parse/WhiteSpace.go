package parse

import (
	"github.com/rwxrob/pegn"
	"github.com/rwxrob/pegn/scan"
)

func WhiteSpace(s Scanner) *pegn.Node {
	if !scan.WhiteSpace(s) {
		return nil
	}
	return &Node{T: pegn.WhiteSpace, V: string(s.Rune())}
}
