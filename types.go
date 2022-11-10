package pegn

// RulesMap is for looking up a rule by its Ident or Alias.
type RulesMap map[string]Rule

var Rules = RulesMap{
	`WhiteSpace`: WhiteSpace,
	`ws`:         WhiteSpace,
}
