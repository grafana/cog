package yaml

import (
	"os"

	"github.com/grafana/cog/internal/veneers/builder"
	"github.com/grafana/cog/internal/veneers/option"
	"github.com/grafana/cog/internal/veneers/rewrite"
	"gopkg.in/yaml.v3"
)

type Veneers struct {
	Language string        `yaml:"language"`
	Builders []BuilderRule `yaml:"builders"`
	Options  []OptionRule  `yaml:"options"`
}

type Loader struct {
}

func NewLoader() *Loader {
	return &Loader{}
}

func (loader *Loader) RewriterFromFiles(filenames []string) (*rewrite.Rewriter, error) {
	rules, err := loader.LoadFiles(filenames)
	if err != nil {
		return nil, err
	}

	return rewrite.NewRewrite(rules), nil
}

func (loader *Loader) LoadFiles(filenames []string) ([]rewrite.LanguageRules, error) {
	languageRules := make([]rewrite.LanguageRules, 0, len(filenames))

	for _, filename := range filenames {
		rules, err := loader.LoadFile(filename)
		if err != nil {
			return nil, err
		}

		languageRules = append(languageRules, rules)
	}

	return languageRules, nil
}

func (loader *Loader) LoadFile(filename string) (rewrite.LanguageRules, error) {
	var builderRules []builder.RewriteRule
	var optionRules []option.RewriteRule

	veneers := &Veneers{}

	// read and parse the input file
	data, err := os.ReadFile(filename)
	if err != nil {
		return rewrite.LanguageRules{}, err
	}

	if err := yaml.Unmarshal(data, &veneers); err != nil {
		return rewrite.LanguageRules{}, err
	}

	builderRules = make([]builder.RewriteRule, 0, len(veneers.Builders))
	optionRules = make([]option.RewriteRule, 0, len(veneers.Options))

	// convert builder rules
	for _, rule := range veneers.Builders {
		builderRule, err := rule.AsRewriteRule()
		if err != nil {
			return rewrite.LanguageRules{}, err
		}

		builderRules = append(builderRules, builderRule)
	}

	// convert option rules
	for _, rule := range veneers.Options {
		optionRule, err := rule.AsRewriteRule()
		if err != nil {
			return rewrite.LanguageRules{}, err
		}

		optionRules = append(optionRules, optionRule)
	}

	return rewrite.LanguageRules{
		Language:     veneers.Language,
		BuilderRules: builderRules,
		OptionRules:  optionRules,
	}, nil
}
