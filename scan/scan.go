package scan

import "github.com/rwxrob/pegn/curs"

// Recreation of pegn.Scanner interface.
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
