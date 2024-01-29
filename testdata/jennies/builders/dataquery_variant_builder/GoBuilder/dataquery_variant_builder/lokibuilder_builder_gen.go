package dataquery_variant_builder

import (
	cog "github.com/grafana/cog/generated/cog"
	cogvariants "github.com/grafana/cog/generated/cog/variants"
)

var _ cog.Builder[cogvariants.Dataquery] = (*LokiBuilderBuilder)(nil)

type LokiBuilderBuilder struct {
    internal *Loki
    errors map[string]cog.BuildErrors
}

func NewLokiBuilderBuilder() *LokiBuilderBuilder {
	resource := &Loki{}
	builder := &LokiBuilderBuilder{
		internal: resource,
		errors: make(map[string]cog.BuildErrors),
	}

	builder.applyDefaults()

	return builder
}

func (builder *LokiBuilderBuilder) Build() (cogvariants.Dataquery, error) {
	var errs cog.BuildErrors

	for _, err := range builder.errors {
		errs = append(errs, cog.MakeBuildErrors("LokiBuilder", err)...)
	}

	if len(errs) != 0 {
		return Loki{}, errs
	}

	return *builder.internal, nil
}

func (builder *LokiBuilderBuilder) Expr(expr string) *LokiBuilderBuilder {
    builder.internal.Expr = expr

    return builder
}

func (builder *LokiBuilderBuilder) applyDefaults() {
}
