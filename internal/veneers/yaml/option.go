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
	Omit                    *OptionSelector          `yaml:"omit"`
	PromoteToConstructor    *OptionSelector          `yaml:"promote_to_constructor"`
	Rename                  *RenameOption            `yaml:"rename"`
	UnfoldBoolean           *UnfoldBoolean           `yaml:"unfold_boolean"`
	StructFieldsAsArguments *StructFieldsAsArguments `yaml:"struct_fields_as_arguments"`
	ArrayToAppend           *ArrayToAppend           `yaml:"array_to_append"`
	DisjunctionAsOptions    *DisjunctionAsOptions    `yaml:"disjunction_as_options"`
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

	if rule.StructFieldsAsArguments != nil {
		return rule.StructFieldsAsArguments.AsRewriteRule()
	}

	if rule.ArrayToAppend != nil {
		return rule.ArrayToAppend.AsRewriteRule()
	}

	if rule.DisjunctionAsOptions != nil {
		return rule.DisjunctionAsOptions.AsRewriteRule()
	}

	return option.RewriteRule{}, fmt.Errorf("empty rule")
}

type RenameOption struct {
	OptionSelector `yaml:",inline"`

	As string `yaml:"as"`
}

func (rule RenameOption) AsRewriteRule() (option.RewriteRule, error) {
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

type StructFieldsAsArguments struct {
	OptionSelector `yaml:",inline"`
	Fields         []string `yaml:"fields"`
}

func (rule StructFieldsAsArguments) AsRewriteRule() (option.RewriteRule, error) {
	selector, err := rule.AsSelector()
	if err != nil {
		return option.RewriteRule{}, err
	}

	return option.StructFieldsAsArguments(selector, rule.Fields...), nil
}

type ArrayToAppend struct {
	OptionSelector `yaml:",inline"`
}

func (rule ArrayToAppend) AsRewriteRule() (option.RewriteRule, error) {
	selector, err := rule.AsSelector()
	if err != nil {
		return option.RewriteRule{}, err
	}

	return option.ArrayToAppend(selector), nil
}

type DisjunctionAsOptions struct {
	OptionSelector `yaml:",inline"`
}

func (rule DisjunctionAsOptions) AsRewriteRule() (option.RewriteRule, error) {
	selector, err := rule.AsSelector()
	if err != nil {
		return option.RewriteRule{}, err
	}

	return option.DisjunctionAsOptions(selector), nil
}

/******************************************************************************
 * Selectors
 *****************************************************************************/

type OptionSelector struct {
	// objectName.optionName
	ByName *string `yaml:"by_name"`

	// objectName.optionName
	ByNameCaseInsensitive *string `yaml:"by_name_case_insensitive"`

	ByNames *ByNamesSelector `yaml:"by_names"`
}

func (selector OptionSelector) AsSelector() (option.Selector, error) {
	if selector.ByName != nil {
		objectName, optionName, found := strings.Cut(*selector.ByName, ".")
		if !found {
			return nil, fmt.Errorf("option name '%s' is incorrect: no object name found", *selector.ByName)
		}

		return option.ByName(objectName, optionName), nil
	}

	if selector.ByNameCaseInsensitive != nil {
		objectName, optionName, found := strings.Cut(*selector.ByNameCaseInsensitive, ".")
		if !found {
			return nil, fmt.Errorf("option name '%s' is incorrect: no object name found", *selector.ByNameCaseInsensitive)
		}

		return option.ByNameCaseInsensitive(objectName, optionName), nil
	}

	if selector.ByNames != nil {
		return selector.ByNames.AsSelector()
	}

	return nil, fmt.Errorf("empty selector")
}

type ByNamesSelector struct {
	Object  string   `yaml:"object"`
	Options []string `yaml:"options"`
}

func (selector ByNamesSelector) AsSelector() (option.Selector, error) {
	if selector.Object == "" {
		return nil, fmt.Errorf("`object` is required")
	}

	return option.ByName(selector.Object, selector.Options...), nil
}
