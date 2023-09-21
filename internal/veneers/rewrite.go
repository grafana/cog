package veneers

import (
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/veneers/builder"
	"github.com/grafana/cog/internal/veneers/option"
)

type Rewriter struct {
	builderRules []builder.RewriteRule
	optionRules  []option.RewriteRule
}

func NewRewrite(builderRules []builder.RewriteRule, optionRules []option.RewriteRule) *Rewriter {
	return &Rewriter{
		builderRules: builderRules,
		optionRules:  optionRules,
	}
}

func (engine *Rewriter) ApplyTo(builders []ast.Builder) []ast.Builder {
	// TODO: should we deepCopy the builders instead?
	newBuilders := make([]ast.Builder, 0, len(builders))
	newBuilders = append(newBuilders, builders...)

	newBuilders = engine.applyBuilderRules(newBuilders)
	newBuilders = engine.applyOptionRules(newBuilders)

	return newBuilders
}

func (engine *Rewriter) applyBuilderRules(builders []ast.Builder) []ast.Builder {
	for _, rule := range engine.builderRules {
		for i, b := range builders {
			// this builder is being discarded
			if len(b.Options) == 0 {
				continue
			}

			if !rule.Selector(b) {
				continue
			}

			builders[i] = rule.Action(builders, b)
		}
	}

	return engine.filterDiscardedBuilders(builders)
}

func (engine *Rewriter) applyOptionRules(builders []ast.Builder) []ast.Builder {
	for _, rule := range engine.optionRules {
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

	return engine.filterDiscardedBuilders(builders)
}

func (engine *Rewriter) filterDiscardedBuilders(builders []ast.Builder) []ast.Builder {
	filteredBuilders := make([]ast.Builder, 0, len(builders))
	for _, b := range builders {
		// the builder was dismissed
		if len(b.Options) == 0 {
			continue
		}

		filteredBuilders = append(filteredBuilders, b)
	}

	return filteredBuilders
}
