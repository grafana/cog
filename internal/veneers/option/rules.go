package option

type RewriteRule struct {
	Selector Selector
	Action   RewriteAction
}

func Rename(selector Selector, newName string) RewriteRule {
	return RewriteRule{
		Selector: selector,
		Action:   RenameAction(newName),
	}
}

func ArrayToAppend(selector Selector) RewriteRule {
	return RewriteRule{
		Selector: selector,
		Action:   ArrayToAppendAction(),
	}
}

func Omit(selector Selector) RewriteRule {
	return RewriteRule{
		Selector: selector,
		Action:   OmitAction(),
	}
}

func UnfoldBoolean(selector Selector, unfoldOpts BooleanUnfold) RewriteRule {
	return RewriteRule{
		Selector: selector,
		Action:   UnfoldBooleanAction(unfoldOpts),
	}
}

func PromoteToConstructor(selector Selector) RewriteRule {
	return RewriteRule{
		Selector: selector,
		Action:   PromoteToConstructorAction(),
	}
}

func StructFieldsAsArguments(selector Selector, explicitFields ...string) RewriteRule {
	return RewriteRule{
		Selector: selector,
		Action:   StructFieldsAsArgumentsAction(explicitFields...),
	}
}

func StructFieldsAsOptions(selector Selector, explicitFields ...string) RewriteRule {
	return RewriteRule{
		Selector: selector,
		Action:   StructFieldsAsOptionsAction(explicitFields...),
	}
}

func DisjunctionAsOptions(selector Selector) RewriteRule {
	return RewriteRule{
		Selector: selector,
		Action:   DisjunctionAsOptionsAction(),
	}
}
