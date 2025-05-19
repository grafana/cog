package initialization_safeguards

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

	return builder
}



func (builder *SomePanelBuilder) Build() (SomePanel, error) {
	if err := builder.internal.Validate(); err != nil {
		return SomePanel{}, err
	}
	
	if len(builder.errors) > 0 {
	    return SomePanel{}, cog.MakeBuildErrors("initialization_safeguards.somePanel", builder.errors)
	}

	return *builder.internal, nil
}

func (builder *SomePanelBuilder) Title(title string) *SomePanelBuilder {
    builder.internal.Title = title

    return builder
}

func (builder *SomePanelBuilder) ShowLegend(show bool) *SomePanelBuilder {
if builder.internal.Options == nil {
    builder.internal.Options = NewOptions()
}
    builder.internal.Options.Legend.Show = show

    return builder
}

