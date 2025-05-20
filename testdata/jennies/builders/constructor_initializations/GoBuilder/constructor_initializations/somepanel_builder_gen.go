package constructor_initializations

import (
	cog "github.com/grafana/cog/generated/cog"
)

var _ cog.Builder[SomePanel] = (*SomePanelBuilder)(nil)

type SomePanelBuilder struct {
    internal *SomePanel
    errors cog.BuildErrors
}

func NewSomePanelBuilder() *SomePanelBuilder {
	resource := NewSomePanel()
	builder := &SomePanelBuilder{
		internal: resource,
		errors: make(cog.BuildErrors, 0),
	}
    builder.internal.Type = "panel_type"
    builder.internal.Cursor = "tooltip"

	return builder
}



func (builder *SomePanelBuilder) Build() (SomePanel, error) {
	if err := builder.internal.Validate(); err != nil {
		return SomePanel{}, err
	}
	
	if len(builder.errors) > 0 {
	    return SomePanel{}, cog.MakeBuildErrors("constructor_initializations.somePanel", builder.errors)
	}

	return *builder.internal, nil
}

func (builder *SomePanelBuilder) Title(title string) *SomePanelBuilder {
    builder.internal.Title = title

    return builder
}

