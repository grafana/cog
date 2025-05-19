package map_of_builders

import (
	cog "github.com/grafana/cog/generated/cog"
)

var _ cog.Builder[Panel] = (*PanelBuilder)(nil)

type PanelBuilder struct {
    internal *Panel
    errors cog.BuildErrors
}

func NewPanelBuilder() *PanelBuilder {
	resource := NewPanel()
	builder := &PanelBuilder{
		internal: resource,
		errors: make(cog.BuildErrors, 0),
	}

	return builder
}



func (builder *PanelBuilder) Build() (Panel, error) {
	if err := builder.internal.Validate(); err != nil {
		return Panel{}, err
	}

	return *builder.internal, cog.MakeBuildErrors("map_of_builders.panel", builder.errors)
}

func (builder *PanelBuilder) Title(title string) *PanelBuilder {
    builder.internal.Title = title

    return builder
}

