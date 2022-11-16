package is

import u "unicode"

var (
	C_ucontrol = u.IsControl
	C_udigit   = u.IsDigit
	C_ugraphic = u.IsGraphic
	C_uletter  = u.IsLetter
	C_ulower   = u.IsLower
	C_umark    = u.IsMark
	C_unumber  = u.IsNumber
	C_uprint   = u.IsPrint
	C_upunct   = u.IsPunct
	C_uspace   = u.IsSpace
	C_usymbol  = u.IsSymbol
	C_utitle   = u.IsTitle
	C_uupper   = u.IsUpper
)

var C_uc_cc = func(r rune) bool { return u.IsOneOf([]*u.RangeTable{u.Cc}, r) }
