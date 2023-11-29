package rewrite

import (
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/tools"
	"github.com/grafana/cog/internal/veneers/builder"
	"github.com/grafana/cog/internal/veneers/option"
)

const AllLanguages = "all"

type LanguageRules struct {
	Language     string
	BuilderRules []builder.RewriteRule
	OptionRules  []option.RewriteRule
}

type Rewriter struct {
	// Rules applied to `Builder` objects, grouped by language
	builderRules map[string][]builder.RewriteRule
	// Rules applied to `Option` objects, grouped by language
	optionRules map[string][]option.RewriteRule
}

func NewRewrite(languageRules []LanguageRules) *Rewriter {
	builderRules := make(map[string][]builder.RewriteRule)
	optionRules := make(map[string][]option.RewriteRule)

	for _, languageConfig := range languageRules {
		builderRules[languageConfig.Language] = append(builderRules[languageConfig.Language], languageConfig.BuilderRules...)
		optionRules[languageConfig.Language] = append(optionRules[languageConfig.Language], languageConfig.OptionRules...)
	}

	return &Rewriter{
		builderRules: builderRules,
		optionRules:  optionRules,
	}
}

func (engine *Rewriter) ApplyTo(builders []ast.Builder, language string) ([]ast.Builder, error) {
	var err error
	// TODO: should we deepCopy the builders instead?
	newBuilders := make([]ast.Builder, 0, len(builders))
	newBuilders = append(newBuilders, builders...)

	// start by applying veneers common to all languages, then
	// apply language-specific ones.
	languages := []string{
		AllLanguages,
		language,
	}

	for _, l := range languages {
		newBuilders, err = engine.applyBuilderRules(newBuilders, l)
		if err != nil {
			return nil, err
		}

		newBuilders = engine.applyOptionRules(newBuilders, l)
	}

	return newBuilders, nil
}

func (engine *Rewriter) applyBuilderRules(builders []ast.Builder, language string) ([]ast.Builder, error) {
	var err error

	for _, rule := range engine.builderRules[language] {
		builders, err = rule(builders)
		if err != nil {
			return nil, err
		}
	}

	return builders, nil
}

func (engine *Rewriter) applyOptionRules(builders []ast.Builder, language string) []ast.Builder {
	for _, rule := range engine.optionRules[language] {
		for i, b := range builders {
			processedOptions := make([]ast.Option, 0, len(b.Options))

			for _, opt := range b.Options {
				if !rule.Selector(b, opt) {
					processedOptions = append(processedOptions, opt)
					continue
				}

				processedOptions = append(processedOptions, rule.Action(b, opt)...)
			}

			builders[i].Options = processedOptions
		}
	}

	return tools.Filter(builders, func(builder ast.Builder) bool {
		// "no options" means that the builder was dismissed.
		return len(builder.Options) != 0
	})
}
