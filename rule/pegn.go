package rule

import (
	"go/ast"

	"github.com/rwxrob/pegn/curs"
)

/*

This file provides 1-1 interface and type compatibility with the main pegn package preventing cyclical import loops and ensuring things work as expected when and if pegn interfaces undergo updates.

*/

// ScanFunc is from pegn.ScanFunc
type ScanFunc func(s Scanner) bool

// ParseFunc is from pegn.ParseFunc
type ParseFunc func(s Scanner) *ast.Node

// Scanner is from pegn.Scanner
type Scanner interface {
	Bytes() *[]byte
	Buffer(input any) error
	Scan() bool
	Peek(a string) bool
	Finished() bool
	Beginning() bool
	Mark() curs.R
	Goto(a curs.R)
	Rune() rune
	RuneB() int
	RuneE() int
	SetViewLen(a int)
	ViewLen() int
	Print()
	Log()
	TraceOn()
	TraceOff()
	CopyEE(to curs.R) string
	CopyBE(to curs.R) string
	CopyBB(to curs.R) string
	CopyEB(to curs.R) string
	SetMaxErr(i int)
	Errors() *[]error
	ErrPush(e error)
	ErrPop() error
	Expected(t int) bool
}
