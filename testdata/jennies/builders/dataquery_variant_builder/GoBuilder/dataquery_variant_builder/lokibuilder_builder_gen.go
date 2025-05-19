package dataquery_variant_builder

import (
	cog "github.com/grafana/cog/generated/cog"
	variants "github.com/grafana/cog/generated/cog/variants"
)

var _ cog.Builder[variants.Dataquery] = (*LokiBuilderBuilder)(nil)

type LokiBuilderBuilder struct {
    internal *Loki
    errors cog.BuildErrors
}

func NewLokiBuilderBuilder() *LokiBuilderBuilder {
	resource := NewLoki()
	builder := &LokiBuilderBuilder{
		internal: resource,
		errors: make(cog.BuildErrors, 0),
	}

	return builder
}



func (builder *LokiBuilderBuilder) Build() (variants.Dataquery, error) {
	if err := builder.internal.Validate(); err != nil {
		return Loki{}, err
	}

	return *builder.internal, cog.MakeBuildErrors("dataquery_variant_builder.lokiBuilder", builder.errors)
}

func (builder *LokiBuilderBuilder) Expr(expr string) *LokiBuilderBuilder {
    builder.internal.Expr = expr

    return builder
}

