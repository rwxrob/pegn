package parse

import (
	"github.com/rwxrob/pegn"
	"github.com/rwxrob/pegn/ast"
	"github.com/rwxrob/pegn/rule/id"
	"github.com/rwxrob/pegn/scan"
)

func C_ws(s pegn.Scanner) *ast.Node {
	if !scan.C_ws(s) {
		return nil
	}
	return &ast.Node{T: id.C_ws, V: string(s.Rune())}
}
