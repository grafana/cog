package composable_slot

import (
	cog "github.com/grafana/cog/generated/cog"
	variants "github.com/grafana/cog/generated/cog/variants"
)

var _ cog.Builder[Dashboard] = (*LokiBuilderBuilder)(nil)

type LokiBuilderBuilder struct {
    internal *Dashboard
    errors map[string]cog.BuildErrors
}

func NewLokiBuilderBuilder() *LokiBuilderBuilder {
	resource := &Dashboard{}
	builder := &LokiBuilderBuilder{
		internal: resource,
		errors: make(map[string]cog.BuildErrors),
	}

	builder.applyDefaults()

	return builder
}

func (builder *LokiBuilderBuilder) Build() (Dashboard, error) {
	if err := builder.internal.Validate(); err != nil {
		return Dashboard{}, err
	}

	return *builder.internal, nil
}

func (builder *LokiBuilderBuilder) Target(target cog.Builder[variants.Dataquery]) *LokiBuilderBuilder {
    targetResource, err := target.Build()
    if err != nil {
        builder.errors["target"] = err.(cog.BuildErrors)
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
                    builder.errors["targets"] = err.(cog.BuildErrors)
                    return builder
                }
                targetsResources = append(targetsResources, targetsDepth1)
        }
    builder.internal.Targets = targetsResources

    return builder
}

func (builder *LokiBuilderBuilder) applyDefaults() {
}
