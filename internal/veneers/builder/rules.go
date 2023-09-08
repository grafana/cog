package builder

type RewriteRule struct {
	Selector Selector
	Action   RewriteAction
}

func Omit(selector Selector) RewriteRule {
	return RewriteRule{
		Selector: selector,
		Action:   OmitAction(),
	}
}

func MergeInto(selector Selector, sourceBuilderName string, underPath string, excludeOptions []string) RewriteRule {
	return RewriteRule{
		Selector: selector,
		Action:   MergeIntoAction(sourceBuilderName, underPath, excludeOptions),
	}
}
