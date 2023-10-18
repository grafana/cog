package yaml

import (
	"os"

	"github.com/grafana/cog/internal/veneers/builder"
	"github.com/grafana/cog/internal/veneers/option"
	"github.com/grafana/cog/internal/veneers/rewrite"
	"gopkg.in/yaml.v3"
)

type Veneers struct {
	Builders []BuilderRule `yaml:"builders"`
	Options  []OptionRule  `yaml:"options"`
}

type Loader struct {
}

func NewLoader() *Loader {
	return &Loader{}
}

func (loader *Loader) RewriterFromFiles(filenames []string) (*rewrite.Rewriter, error) {
	builderRules, optionRules, err := loader.LoadFiles(filenames)
	if err != nil {
		return nil, err
	}

	return rewrite.NewRewrite(builderRules, optionRules), nil
}

func (loader *Loader) LoadFiles(filenames []string) ([]builder.RewriteRule, []option.RewriteRule, error) {
	var allBuilderRules []builder.RewriteRule
	var allOptionRules []option.RewriteRule

	for _, filename := range filenames {
		builderRules, optionRules, err := loader.LoadFile(filename)
		if err != nil {
			return nil, nil, err
		}

		allBuilderRules = append(allBuilderRules, builderRules...)
		allOptionRules = append(allOptionRules, optionRules...)
	}

	return allBuilderRules, allOptionRules, nil
}

func (loader *Loader) LoadFile(filename string) ([]builder.RewriteRule, []option.RewriteRule, error) {
	var builderRules []builder.RewriteRule
	var optionRules []option.RewriteRule

	veneers := &Veneers{}

	// read and parse the input file
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, nil, err
	}

	if err := yaml.Unmarshal(data, &veneers); err != nil {
		return nil, nil, err
	}

	builderRules = make([]builder.RewriteRule, 0, len(veneers.Builders))
	optionRules = make([]option.RewriteRule, 0, len(veneers.Options))

	// convert builder rules
	for _, rule := range veneers.Builders {
		builderRule, err := rule.AsRewriteRule()
		if err != nil {
			return nil, nil, err
		}

		builderRules = append(builderRules, builderRule)
	}

	// convert option rules
	for _, rule := range veneers.Options {
		optionRule, err := rule.AsRewriteRule()
		if err != nil {
			return nil, nil, err
		}

		optionRules = append(optionRules, optionRule)
	}

	return builderRules, optionRules, nil
}
