/*

Package rule provides the glue that brings the scan and parse functions together with the documentation and PEGN that describes them.

*/
package rule

import (
	"github.com/rwxrob/pegn/scan"
)

// LangMap is a map that contains strings associated with a specific
// language identifier.
type LangMap map[string]string

// Rule structs bring together all the required parts in a way that can
// also be marshalled as JSON (see docs). All PEGN rules have a unique
// type and possibly rule subtype (lower = class, UPPER = token) . Rules
// often have their descriptions (Desc) omitted until needed when they
// can be dynamically loaded based on the languages needed.  The rest of
// the properties are language agnostic.
//
type Rule struct {
	ID    string    `json:"id,omitempty"`   // RuleID, TokenID, ClassID
	Type  int       `json:"type,omitempty"` // uniq type identifier
	PEGN  string    `json:"pegn,omitempty"` // specific PEGN notation
	Desc  LangMap   `json:"desc,omitempty"` // human-friendly descriptions
	Scan  ScanFunc  `json:"-"`              // func(s Scanner) bool
	Parse ParseFunc `json:"-"`              // func(s Scanner) *Node
}

// RuleType returns 0 if just a rule (MixedCase), 1 if a token (CAPS),
// and 2 if a class (LOWER). It returns -1 if the Ident is an invalid name.
func (r Rule) RuleType() int {
	switch {
	case scan.ClassId(): // class
		return 2
	case scan.TokenId(): // TOKEN
		return 1
	case scan.RuleId(): // SomeRule
		return 0
	default:
		return -1
	}
}
