package constructor_initializations

import (
	cog "github.com/grafana/cog/generated/cog"
)

var _ cog.Builder[SomePanel] = (*SomePanelBuilder)(nil)

type SomePanelBuilder struct {
    internal *SomePanel
    errors map[string]cog.BuildErrors
}

func NewSomePanelBuilder() *SomePanelBuilder {
	resource := &SomePanel{}
	builder := &SomePanelBuilder{
		internal: resource,
		errors: make(map[string]cog.BuildErrors),
	}

	builder.applyDefaults()
    builder.internal.Type = "panel_type"
    builder.internal.Cursor = "tooltip"

	return builder
}

func (builder *SomePanelBuilder) Build() (SomePanel, error) {
	if err := builder.internal.Validate(); err != nil {
		return SomePanel{}, err
	}

	return *builder.internal, nil
}

func (builder *SomePanelBuilder) Title(title string) *SomePanelBuilder {
    builder.internal.Title = title

    return builder
}

func (builder *SomePanelBuilder) applyDefaults() {
}
