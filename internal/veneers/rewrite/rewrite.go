package rewrite

import (
	"fmt"
	"log/slog"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/logs"
	"github.com/grafana/cog/internal/tools"
	"github.com/grafana/cog/internal/veneers/builder"
	"github.com/grafana/cog/internal/veneers/option"
)

const AllLanguages = "all"

type Config struct {
	Debug bool
}

type RuleSet struct {
	Languages    []string
	BuilderRules []*builder.Rule
	OptionRules  []option.RewriteRule
}

type Rewriter struct {
	logger *slog.Logger
	config Config

	// Rules applied to `Builder` objects, grouped by language
	builderRules map[string][]*builder.Rule
	// Rules applied to `Option` objects, grouped by language
	optionRules map[string][]option.RewriteRule
}

func NewRewrite(logger *slog.Logger, rules []RuleSet, config Config) *Rewriter {
	builderRules := make(map[string][]*builder.Rule)
	optionRules := make(map[string][]option.RewriteRule)

	for _, set := range rules {
		for _, language := range set.Languages {
			builderRules[language] = append(builderRules[language], set.BuilderRules...)
			optionRules[language] = append(optionRules[language], set.OptionRules...)
		}
	}

	return &Rewriter{
		logger:       logger,
		config:       config,
		builderRules: builderRules,
		optionRules:  optionRules,
	}
}

func (engine *Rewriter) ApplyTo(schemas ast.Schemas, builders []ast.Builder, language string) ([]ast.Builder, error) {
	var err error
	// TODO: should we deepCopy the builders instead?
	newBuilders := make([]ast.Builder, 0, len(builders))
	newBuilders = append(newBuilders, builders...)

	// start by applying veneers common to all languages, then
	// apply language-specific ones.
	for _, l := range []string{AllLanguages, language} {
		newBuilders, err = engine.applyBuilderRules(language, schemas, newBuilders, engine.builderRules[l])
		if err != nil {
			return nil, err
		}

		newBuilders = engine.applyOptionRules(schemas, newBuilders, engine.optionRules[l])
	}

	// and optionally, apply "debug" veneers
	if engine.config.Debug {
		newBuilders, err = engine.applyBuilderRules("debug", schemas, newBuilders, engine.debugBuilderRules())
		if err != nil {
			return nil, err
		}

		newBuilders = engine.applyOptionRules(schemas, newBuilders, engine.debugOptionRules())
	}

	return newBuilders, nil
}

func (engine *Rewriter) applyBuilderRules(language string, schemas ast.Schemas, builders []ast.Builder, rules []*builder.Rule) ([]ast.Builder, error) {
	var err error
	var transformedBuilders []ast.Builder

	for _, rule := range rules {
		unselectedBuilders := make([]ast.Builder, 0, len(builders))
		var matches ast.Builders
		for _, builder := range builders {
			if !rule.Selector.Matches(schemas, builder) {
				unselectedBuilders = append(unselectedBuilders, builder)
				continue
			}

			matches = append(matches, builder)
		}

		if len(matches) == 0 {
			engine.logger.Warn("builder rule matched nothing", slog.String("language", language), slog.String("rule", rule.String()))
			continue
		}

		ctx := builder.RuleCtx{
			Schemas:  schemas,
			Builders: builders,
		}
		transformedBuilders, err = rule.Action.Run(ctx, matches)
		if err != nil {
			engine.logger.Error("builder rule failed", slog.String("language", language), slog.String("rule", rule.String()), logs.Err(err))
			return nil, fmt.Errorf("builder rule failed: err=%w", err)
		}

		builders = append(unselectedBuilders, transformedBuilders...)
	}

	return builders, nil
}

func (engine *Rewriter) applyOptionRules(schemas ast.Schemas, builders []ast.Builder, rules []option.RewriteRule) []ast.Builder {
	for _, rule := range rules {
		for i, b := range builders {
			processedOptions := make([]ast.Option, 0, len(b.Options))

			for _, opt := range b.Options {
				if !rule.Selector(b, opt) {
					processedOptions = append(processedOptions, opt)
					continue
				}

				processedOptions = append(processedOptions, rule.Action(schemas, b, opt)...)
			}

			builders[i].Options = processedOptions
		}
	}

	return tools.Filter(builders, func(builder ast.Builder) bool {
		// "no options" means that the builder was dismissed.
		return len(builder.Options) != 0
	})
}

func (engine *Rewriter) debugBuilderRules() []*builder.Rule {
	return []*builder.Rule{
		&builder.Rule{
			Selector: builder.EveryBuilder(),
			Action:   builder.VeneerTrailAsComments(),
		},
	}
}

func (engine *Rewriter) debugOptionRules() []option.RewriteRule {
	return []option.RewriteRule{
		option.VeneerTrailAsComments(option.EveryOption()),
	}
}
