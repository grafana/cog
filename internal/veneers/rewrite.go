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
	newBuilders := engine.applyBuilderRules(builders)
	newBuilders = engine.applyOptionRules(builders)

	return newBuilders
}

func (engine *Rewriter) applyBuilderRules(builders []ast.Builder) []ast.Builder {
	processedBuilders := make([]ast.Builder, 0, len(builders))
	processedBuilders = append(processedBuilders, builders...)

	for _, rule := range engine.builderRules {
		for i, b := range processedBuilders {
			// this builder is being discarded
			if len(b.Options) == 0 {
				continue
			}

			processedBuilders[i] = rule.Action(processedBuilders, b)
		}
	}

	return engine.filterDiscardedBuilders(processedBuilders)
}

func (engine *Rewriter) applyOptionRules(builders []ast.Builder) []ast.Builder {
	processedBuilders := make([]ast.Builder, 0, len(builders))
	processedBuilders = append(processedBuilders, builders...)

	for _, rule := range engine.optionRules {
		for i, b := range processedBuilders {
			processedOptions := make([]ast.Option, 0, len(b.Options))

			for _, opt := range b.Options {
				if !rule.Selector(b, opt) {
					processedOptions = append(processedOptions, opt)
					continue
				}

				processedOptions = append(processedOptions, rule.Action(b, opt)...)
			}

			processedBuilders[i].Options = processedOptions
		}
	}

	return engine.filterDiscardedBuilders(processedBuilders)
}

func (engine *Rewriter) filterDiscardedBuilders(builders []ast.Builder) []ast.Builder {
	finalBuilders := make([]ast.Builder, 0, len(builders))
	for _, b := range builders {
		// the builder was dismissed
		if len(b.Options) == 0 {
			continue
		}

		finalBuilders = append(finalBuilders, b)
	}

	return finalBuilders
}
