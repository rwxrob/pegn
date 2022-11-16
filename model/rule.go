/*

Package model provides the data model and schema files associated with the PEGN domain-specific language itself. This provides a way to associate a PEGN rule by ID to its full description in multiple languages.

*/
package model

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
	ID   int     `json:"id,omitempty"`   // uniq type identifier
	Name string  `json:"name,omitempty"` // RuleName, TokenName, ClassName
	Type int     `json:"type"`           // 0 rule, 1 token, 2, class
	PEGN string  `json:"pegn,omitempty"` // specific PEGN notation
	Desc LangMap `json:"desc,omitempty"` // human-friendly descriptions
}
