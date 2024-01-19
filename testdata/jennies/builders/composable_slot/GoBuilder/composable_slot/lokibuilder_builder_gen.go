package composable_slot

import (
	cog "github.com/grafana/cog/generated/cog"
	cogvariants "github.com/grafana/cog/generated/cog/variants"
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
	var errs cog.BuildErrors

	for _, err := range builder.errors {
		errs = append(errs, cog.MakeBuildErrors("LokiBuilder", err)...)
	}

	if len(errs) != 0 {
		return Dashboard{}, errs
	}

	return *builder.internal, nil
}

func (builder *LokiBuilderBuilder) Target(target cog.Builder[cogvariants.Dataquery]) *LokiBuilderBuilder {
        targetResource, err := target.Build()
        if err != nil {
            builder.errors["target"] = err.(cog.BuildErrors)
            return builder
        }
    builder.internal.Target = targetResource

    return builder
}

func (builder *LokiBuilderBuilder) Targets(targets []cog.Builder[cogvariants.Dataquery]) *LokiBuilderBuilder {
        targetsResources := make([]cogvariants.Dataquery, 0, len(targets))
        for _, r := range targets {
            targetsResource, err := r.Build()
            if err != nil {
                builder.errors["targets"] = err.(cog.BuildErrors)
                return builder
            }
            targetsResources = append(targetsResources, targetsResource)
        }
    builder.internal.Targets = targetsResources

    return builder
}

func (builder *LokiBuilderBuilder) applyDefaults() {
}
