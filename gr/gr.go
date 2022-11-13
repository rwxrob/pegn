/*

Package gr (grammar) is a collection of common grammars for convenience including the PEGN specification itself. These grammars have the advantage of being pre-generated saving cost to otherwise parse and compile a dynamic grammar from passed PEGN (similar to regular expression compilation requirements).

*/
package gr

import "go/ast"

var PEGN = _pegn{}

type _pegn struct{}

func (g _pegn) Scan(in any, spec string) (bool, []error) {
	// TODO
	return false, nil
}

func (g _pegn) Parse(in any, spec string) (*ast.Node, []error) {
	// TODO
	return nil, nil
}
