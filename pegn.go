package pegn

import "github.com/rwxrob/pegn/ast"

// NEVER REMOVE FROM LIST!
// Append to list only (even if deprecated or not supported)
const (
	Untyped int = -iota
	C_ws
)

/*
// Token Definitions
const (
	// PEGN-tokens (pegn.dev/spec/tokens.pegn)
	TAB       = 0x0009 // "\t"
	LF        = 0x000A // "\n" (line feed)
	CR        = 0x000D // "\r" (carriage return)
	CRLF      = "\r\n" // "\r\n"
	SP        = 0x0020 // " "
	VT        = 0x000B // "\v" (vertical tab)
	FF        = 0x000C // "\f" (form feed)
	NOT       = 0x0021 // !
	BANG      = 0x0021 // !
	DQ        = 0x0022 // "
	HASH      = 0x0023 // #
	DOLLAR    = 0x0024 // $
	PERCENT   = 0x0025 // %
	AND       = 0x0026 // &
	SQ        = 0x0027 // '
	LPAREN    = 0x0028 // (
	RPAREN    = 0x0029 // )
	STAR      = 0x002A // *
	PLUS      = 0x002B // +
	COMMA     = 0x002C // ,
	DASH      = 0x002D // -
	MINUS     = 0x002D // -
	DOT       = 0x002E // .
	SLASH     = 0x002F // /
	COLON     = 0x003A // :
	SEMI      = 0x003B // ;
	LT        = 0x003C // <
	EQ        = 0x003D // =
	GT        = 0x003E // >
	QUERY     = 0x003F // ?
	QUESTION  = 0x003F // ?
	AT        = 0x0040 // @
	LBRAKT    = 0x005B // [
	BKSLASH   = 0x005C // \
	RBRAKT    = 0x005D // ]
	CARET     = 0x005E // ^
	UNDER     = 0x005F // _
	BKTICK    = 0x0060 // `
	LCURLY    = 0x007B // {
	LBRACE    = 0x007B // {
	BAR       = 0x007C // |
	PIPE      = 0x007C // |
	RCURLY    = 0x007D // }
	RBRACE    = 0x007D // }
	TILDE     = 0x007E // ~
	UNKNOWN   = 0xFFFD
	REPLACE   = 0xFFFD
	MAXRUNE   = 0x0010FFFF
	ENDOFDATA = 134217727 // largest int32
	MAXASCII  = 0x007F
	MAXLATIN  = 0x00FF
	RARROWF   = "=>"
	LARROWF   = "<="
	LARROW    = "<-"
	RARROW    = "->"
	LLARROW   = "<--"
	RLARROW   = "-->"
	LFAT      = "<="
	RFAT      = "=>"
	WALRUS    = ":="
)
*/

// -------------------------------- ws --------------------------------

var Is_ws = func(r rune) bool {
	switch r {
	case ' ', '\t', '\n', '\r':
		return true
	default:
		return false
	}
}

func Scan_ws(s Scanner, buf *[]rune) bool {
	m := s.Mark()
	if !s.Scan() {
		return false
	}
	r := s.Rune()
	if Is_ws(s.Rune()) {
		if buf != nil {
			*buf = append(*buf, r)
		}
		return true
	}
	return s.Revert(m, C_ws)
}

func Parse_ws(s Scanner) *ast.Node {
	buf := make([]rune, 0, 1)
	if !Scan_ws(s, &buf) {
		return nil
	}
	return &ast.Node{T: C_ws, V: string(buf)}
}

/*
var (
	Is_ucontrol = u.IsControl
	Is_udigit   = u.IsDigit
	Is_ugraphic = u.IsGraphic
	Is_uletter  = u.IsLetter
	Is_ulower   = u.IsLower
	Is_umark    = u.IsMark
	Is_unumber  = u.IsNumber
	Is_uprint   = u.IsPrint
	Is_upunct   = u.IsPunct
	Is_uspace   = u.IsSpace
	Is_usymbol  = u.IsSymbol
	Is_utitle   = u.IsTitle
	Is_uupper   = u.IsUpper
)

var Is_uc_cc = func(r rune) bool { return u.IsOneOf([]*u.RangeTable{u.Cc}, r) }
*/
