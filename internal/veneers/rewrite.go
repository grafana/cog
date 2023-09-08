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
	newBuilders := make([]ast.Builder, 0, len(builders))

	for _, b := range builders {
		processed := engine.processBuilder(builders, b)
		// the builder was dismissed
		if len(processed.Options) == 0 {
			continue
		}

		newBuilders = append(newBuilders, processed)
	}

	return newBuilders
}

func (engine *Rewriter) processBuilder(builders ast.Builders, builder ast.Builder) ast.Builder {
	processedBuilder := builder

	for _, rule := range engine.builderRules {
		if rule.Selector(processedBuilder) {
			// FIXME: passing `builders` here means that rules only get access to a "pre modification"
			// set of builders. We should probably pass the most up-to-date list of builders
			processedBuilder = rule.Action(builders, processedBuilder)
		}

		// this builder is dismissed, let's return early
		if len(processedBuilder.Options) == 0 {
			return processedBuilder
		}
	}

	processedOptions := make([]ast.Option, 0, len(processedBuilder.Options))
	for _, opt := range processedBuilder.Options {
		processedOptions = append(processedOptions, engine.processOption(processedBuilder, opt)...)
	}

	processedBuilder.Options = processedOptions

	return processedBuilder
}

func (engine *Rewriter) processOption(parentBuilder ast.Builder, opt ast.Option) []ast.Option {
	toProcess := []ast.Option{opt}
	for _, rule := range engine.optionRules {
		if !rule.Selector(parentBuilder, opt) {
			continue
		}

		var wip []ast.Option
		for _, modifiedField := range toProcess {
			wip = append(wip, rule.Action(modifiedField)...)
		}
		toProcess = wip
	}

	return toProcess
}
