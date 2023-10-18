package yaml

import (
	"fmt"
	"strings"

	"github.com/grafana/cog/internal/veneers/option"
)

/******************************************************************************
 * Rules
 *****************************************************************************/

type OptionRule struct {
	Omit                 *OptionSelector `yaml:"omit"`
	PromoteToConstructor *OptionSelector `yaml:"promote_to_constructor"`
	Rename               *Rename         `yaml:"rename"`
	UnfoldBoolean        *UnfoldBoolean  `yaml:"unfold_boolean"`
}

func (rule OptionRule) AsRewriteRule() (option.RewriteRule, error) {
	if rule.Omit != nil {
		selector, err := rule.Omit.AsSelector()
		if err != nil {
			return option.RewriteRule{}, err
		}

		return option.Omit(selector), nil
	}

	if rule.PromoteToConstructor != nil {
		selector, err := rule.PromoteToConstructor.AsSelector()
		if err != nil {
			return option.RewriteRule{}, err
		}

		return option.PromoteToConstructor(selector), nil
	}

	if rule.Rename != nil {
		return rule.Rename.AsRewriteRule()
	}

	if rule.UnfoldBoolean != nil {
		return rule.UnfoldBoolean.AsRewriteRule()
	}

	return option.RewriteRule{}, fmt.Errorf("empty rule")
}

type Rename struct {
	OptionSelector `yaml:",inline"`

	As string `yaml:"as"`
}

func (rule Rename) AsRewriteRule() (option.RewriteRule, error) {
	selector, err := rule.AsSelector()
	if err != nil {
		return option.RewriteRule{}, err
	}

	return option.Rename(selector, rule.As), nil
}

type UnfoldBoolean struct {
	OptionSelector `yaml:",inline"`

	TrueAs  string `yaml:"true_as"`
	FalseAs string `yaml:"false_as"`
}

func (rule UnfoldBoolean) AsRewriteRule() (option.RewriteRule, error) {
	selector, err := rule.AsSelector()
	if err != nil {
		return option.RewriteRule{}, err
	}

	return option.UnfoldBoolean(selector, option.BooleanUnfold{
		OptionTrue:  rule.TrueAs,
		OptionFalse: rule.FalseAs,
	}), nil
}

/******************************************************************************
 * Selectors
 *****************************************************************************/

type OptionSelector struct {
	ByName *string `yaml:"by_name"`
}

func (selector OptionSelector) AsSelector() (option.Selector, error) {
	if selector.ByName != nil {
		objectName, optionName, found := strings.Cut(*selector.ByName, ".")
		if !found {
			return nil, fmt.Errorf("option name '%s' is incorrect: no object name found", *selector.ByName)
		}

		return option.ByName(objectName, optionName), nil
	}

	return nil, fmt.Errorf("empty selector")
}
