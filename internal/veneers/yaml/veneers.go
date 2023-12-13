package yaml

import (
	"fmt"
	"io"
	"os"

	"github.com/grafana/cog/internal/veneers/builder"
	"github.com/grafana/cog/internal/veneers/option"
	"github.com/grafana/cog/internal/veneers/rewrite"
	"gopkg.in/yaml.v3"
)

type Veneers struct {
	Language string        `yaml:"language"`
	Package  string        `yaml:"package"`
	Builders []BuilderRule `yaml:"builders"`
	Options  []OptionRule  `yaml:"options"`
}

type Loader struct {
}

func NewLoader() *Loader {
	return &Loader{}
}

func (loader *Loader) RewriterFrom(filenames []string) (*rewrite.Rewriter, error) {
	readers := make([]io.Reader, 0, len(filenames))
	for _, filename := range filenames {
		reader, err := os.Open(filename)
		if err != nil {
			return nil, err
		}

		readers = append(readers, reader)
	}

	rules, err := loader.LoadAll(readers)
	if err != nil {
		return nil, err
	}

	return rewrite.NewRewrite(rules), nil
}

func (loader *Loader) LoadAll(readers []io.Reader) ([]rewrite.LanguageRules, error) {
	languageRules := make([]rewrite.LanguageRules, 0, len(readers))

	for _, filename := range readers {
		rules, err := loader.Load(filename)
		if err != nil {
			return nil, err
		}

		languageRules = append(languageRules, rules)
	}

	return languageRules, nil
}

func (loader *Loader) Load(reader io.Reader) (rewrite.LanguageRules, error) {
	var builderRules []builder.RewriteRule
	var optionRules []option.RewriteRule

	veneers := &Veneers{}

	decoder := yaml.NewDecoder(reader)
	decoder.KnownFields(true)

	if err := decoder.Decode(&veneers); err != nil {
		return rewrite.LanguageRules{}, err
	}

	if veneers.Package == "" {
		return rewrite.LanguageRules{}, fmt.Errorf("missing 'package' statement in veneers file '%s'", reader)
	}

	builderRules = make([]builder.RewriteRule, 0, len(veneers.Builders))
	optionRules = make([]option.RewriteRule, 0, len(veneers.Options))

	// convert builder rules
	for _, rule := range veneers.Builders {
		builderRule, err := rule.AsRewriteRule(veneers.Package)
		if err != nil {
			return rewrite.LanguageRules{}, err
		}

		builderRules = append(builderRules, builderRule)
	}

	// convert option rules
	for _, rule := range veneers.Options {
		optionRule, err := rule.AsRewriteRule(veneers.Package)
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
