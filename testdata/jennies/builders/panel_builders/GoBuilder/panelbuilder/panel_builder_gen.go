package panelbuilder

import (
	cog "github.com/grafana/cog/generated/cog"
)

var _ cog.Builder[Options] = (*PanelBuilder)(nil)

type PanelBuilder struct {
    internal *Options
    errors map[string]cog.BuildErrors
}

func NewPanelBuilder() *PanelBuilder {
	resource := &Options{}
	builder := &PanelBuilder{
		internal: resource,
		errors: make(map[string]cog.BuildErrors),
	}

	builder.applyDefaults()

	return builder
}

func (builder *PanelBuilder) Build() (Options, error) {
	var errs cog.BuildErrors

	for _, err := range builder.errors {
		errs = append(errs, cog.MakeBuildErrors("Panel", err)...)
	}

	if len(errs) != 0 {
		return Options{}, errs
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

func (builder *PanelBuilder) applyDefaults() {
    builder.OnlyFromThisDashboard(false)
    builder.OnlyInTimeRange(false)
    builder.Limit(10)
    builder.ShowUser(true)
    builder.ShowTime(true)
    builder.ShowTags(true)
    builder.NavigateToPanel(true)
    builder.NavigateBefore("10m")
    builder.NavigateAfter("10m")
}
