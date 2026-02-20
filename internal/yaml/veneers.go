package yaml

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
	"github.com/grafana/cog/internal/veneers/builder"
	"github.com/grafana/cog/internal/veneers/option"
	"github.com/grafana/cog/internal/veneers/rewrite"
)

type Veneers struct {
	Languages []string      `yaml:"languages"`
	Package   string        `yaml:"package"`
	Builders  []BuilderRule `yaml:"builders"`
	Options   []OptionRule  `yaml:"options"`
}

type VeneersLoader struct {
}

func NewVeneersLoader() *VeneersLoader {
	return &VeneersLoader{}
}

func (loader *VeneersLoader) RewriterFrom(filenames []string, config rewrite.Config) (*rewrite.Rewriter, error) {
	ruleSet := make([]rewrite.RuleSet, 0, len(filenames))

	for _, filename := range filenames {
		rules, err := loader.load(filename)
		if err != nil {
			return nil, err
		}

		ruleSet = append(ruleSet, rules)
	}

	return rewrite.NewRewrite(ruleSet, config), nil
}

func (loader *VeneersLoader) load(file string) (rewrite.RuleSet, error) {
	var builderRules []builder.RewriteRule
	var optionRules []option.RewriteRule

	contents, err := os.ReadFile(file)
	if err != nil {
		return rewrite.RuleSet{}, err
	}

	veneers := &Veneers{}
	if err := yaml.UnmarshalWithOptions(contents, veneers, yaml.DisallowUnknownField()); err != nil {
		return rewrite.RuleSet{}, fmt.Errorf("can not load veneers: %s\n%s", file, yaml.FormatError(err, true, true))
	}

	if veneers.Package == "" {
		return rewrite.RuleSet{}, fmt.Errorf("missing 'package' statement in veneers file '%s'\n", file)
	}

	builderRules = make([]builder.RewriteRule, 0, len(veneers.Builders))
	optionRules = make([]option.RewriteRule, 0, len(veneers.Options))

	// convert builder rules
	for i, rule := range veneers.Builders {
		builderRule, err := rule.AsRewriteRule(veneers.Package)
		if err != nil {
			path, innerErr := yaml.PathString(fmt.Sprintf("$.builders[%d]", i))
			if innerErr != nil {
				return rewrite.RuleSet{}, err
			}
			source, innerErr := path.AnnotateSource(contents, true)
			if innerErr != nil {
				return rewrite.RuleSet{}, err
			}

			return rewrite.RuleSet{}, fmt.Errorf("%w in %s\n%s", err, file, string(source))
		}

		builderRules = append(builderRules, builderRule)
	}

	// convert option rules
	for i, rule := range veneers.Options {
		optionRule, err := rule.AsRewriteRule(veneers.Package)
		if err != nil {
			path, innerErr := yaml.PathString(fmt.Sprintf("$.options[%d]", i))
			if innerErr != nil {
				return rewrite.RuleSet{}, err
			}
			source, innerErr := path.AnnotateSource(contents, true)
			if innerErr != nil {
				return rewrite.RuleSet{}, err
			}

			return rewrite.RuleSet{}, fmt.Errorf("%w in %s\n%s", err, file, string(source))
		}

		optionRules = append(optionRules, optionRule)
	}

	return rewrite.RuleSet{
		Languages:    veneers.Languages,
		BuilderRules: builderRules,
		OptionRules:  optionRules,
	}, nil
}
