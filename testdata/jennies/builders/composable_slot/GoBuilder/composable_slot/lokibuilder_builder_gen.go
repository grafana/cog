package composable_slot

import (
	cog "github.com/grafana/cog/generated/cog"
	variants "github.com/grafana/cog/generated/cog/variants"
)

var _ cog.Builder[Dashboard] = (*LokiBuilderBuilder)(nil)

type LokiBuilderBuilder struct {
    internal *Dashboard
    errors cog.BuildErrors
}

func NewLokiBuilderBuilder() *LokiBuilderBuilder {
	resource := NewDashboard()
	builder := &LokiBuilderBuilder{
		internal: resource,
		errors: make(cog.BuildErrors, 0),
	}

	return builder
}



func (builder *LokiBuilderBuilder) Build() (Dashboard, error) {
	if err := builder.internal.Validate(); err != nil {
		return Dashboard{}, err
	}
	
	if len(builder.errors) > 0 {
	    return Dashboard{}, cog.MakeBuildErrors("composable_slot.lokiBuilder", builder.errors)
	}

	return *builder.internal, nil
}

func (builder *LokiBuilderBuilder) Target(target cog.Builder[variants.Dataquery]) *LokiBuilderBuilder {
    targetResource, err := target.Build()
    if err != nil {
        builder.errors = append(builder.errors, err.(cog.BuildErrors)...)
        return builder
    }
    builder.internal.Target = targetResource

    return builder
}

func (builder *LokiBuilderBuilder) Targets(targets []cog.Builder[variants.Dataquery]) *LokiBuilderBuilder {
        targetsResources := make([]variants.Dataquery, 0, len(targets))
        for _, r1 := range targets {
                targetsDepth1, err := r1.Build()
                if err != nil {
                    builder.errors = append(builder.errors, err.(cog.BuildErrors)...)
                    return builder
                }
                targetsResources = append(targetsResources, targetsDepth1)
        }
    builder.internal.Targets = targetsResources

    return builder
}

