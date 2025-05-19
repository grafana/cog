package panelbuilder

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
	
	if len(builder.errors) > 0 {
	    return Panel{}, cog.MakeBuildErrors("panelbuilder.panel", builder.errors)
	}

	return *builder.internal, nil
}

func (builder *PanelBuilder) OnlyFromThisDashboard(onlyFromThisDashboard bool) *PanelBuilder {
    builder.internal.OnlyFromThisDashboard = onlyFromThisDashboard

    return builder
}

func (builder *PanelBuilder) OnlyInTimeRange(onlyInTimeRange bool) *PanelBuilder {
    builder.internal.OnlyInTimeRange = onlyInTimeRange

    return builder
}

func (builder *PanelBuilder) Tags(tags []string) *PanelBuilder {
    builder.internal.Tags = tags

    return builder
}

func (builder *PanelBuilder) Limit(limit uint32) *PanelBuilder {
    builder.internal.Limit = limit

    return builder
}

func (builder *PanelBuilder) ShowUser(showUser bool) *PanelBuilder {
    builder.internal.ShowUser = showUser

    return builder
}

func (builder *PanelBuilder) ShowTime(showTime bool) *PanelBuilder {
    builder.internal.ShowTime = showTime

    return builder
}

func (builder *PanelBuilder) ShowTags(showTags bool) *PanelBuilder {
    builder.internal.ShowTags = showTags

    return builder
}

func (builder *PanelBuilder) NavigateToPanel(navigateToPanel bool) *PanelBuilder {
    builder.internal.NavigateToPanel = navigateToPanel

    return builder
}

func (builder *PanelBuilder) NavigateBefore(navigateBefore string) *PanelBuilder {
    builder.internal.NavigateBefore = navigateBefore

    return builder
}

func (builder *PanelBuilder) NavigateAfter(navigateAfter string) *PanelBuilder {
    builder.internal.NavigateAfter = navigateAfter

    return builder
}

