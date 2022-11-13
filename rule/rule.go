/*

Package rule provides the glue that brings the scan and parse functions together with the documentation and PEGN that describes them.

*/
package rule

import (
	"fmt"

	"github.com/rwxrob/pegn/curs"
)

// LangMap is a map that contains strings associated with a specific
// language identifier.
type LangMap map[string]string

// Rule structs bring together all the required parts in a way that can
// also be marshalled as JSON (see docs). All PEGN rules have a unique
// integer ID and type (0 rule, 1 token, and 2 class) that
// corresponds to the PEGN case conventions (Mixed, CAPS, lower). No two
// rules must ever have the same case-insensitive name. Rules
// often have their descriptions (Desc) omitted until needed when they
// can be dynamically loaded based on the languages needed.  The rest of
// the properties are language agnostic.
type Rule struct {
	ID    int       `json:"id,omitempty"`   // uniq type identifier
	Name  string    `json:"name,omitempty"` // RuleName, TokenName, ClassName
	Type  int       `json:"type"`           // 0 rule, 1 token, 2, class
	PEGN  string    `json:"pegn,omitempty"` // specific PEGN notation
	Desc  LangMap   `json:"desc,omitempty"` // human-friendly descriptions
	Scan  ScanFunc  `json:"-"`              // func(s Scanner) bool
	Parse ParseFunc `json:"-"`              // func(s Scanner) *Node
}

// Error wraps the type (T) and current scanner position (C)
// such that it can be located and displayed with help information by
// looking up those things from other sources when displayed to the end
// user. The position of fields is guaranteed never to change allowing
// for short-form instantiation (ex: pegn.Error{1,s.Mark()}). See
// ScannerErrors interface for more.
//
type Error struct {
	T int
	C curs.R
}

var DefaultErrFmt = `expecting %v at %v`

func (e Error) Error() string {
	return fmt.Sprintf(DefaultErrFmt, e.T, e.C)
}
