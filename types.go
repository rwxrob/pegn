package pegn

// RulesMap is for looking up a rule by its Ident or Alias.
type RulesMap map[string]Rule

type TypesMap map[int]Rule

var Rules = RulesMap{
	`WhiteSpace`: WhiteSpace,
	`ws`:         WhiteSpace,
	`Field`:      Field,
	`Uprint`:     Uprint,
	`uprint`:     Uprint,
}

var Types = TypesMap{
	1: WhiteSpace,
	2: Field,
	3: Uprint,
}
