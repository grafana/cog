package builderdelegation

import (
	cog "github.com/grafana/cog/generated/cog"
)

var _ cog.Builder[Dashboard] = (*DashboardBuilder)(nil)

type DashboardBuilder struct {
    internal *Dashboard
    errors map[string]cog.BuildErrors
}

func NewDashboardBuilder() *DashboardBuilder {
	resource := &Dashboard{}
	builder := &DashboardBuilder{
		internal: resource,
		errors: make(map[string]cog.BuildErrors),
	}

	builder.applyDefaults()

	return builder
}

func (builder *DashboardBuilder) Build() (Dashboard, error) {
	var errs cog.BuildErrors

	for _, err := range builder.errors {
		errs = append(errs, cog.MakeBuildErrors("Dashboard", err)...)
	}

	if len(errs) != 0 {
		return Dashboard{}, errs
	}

	return *builder.internal, nil
}

func (builder *DashboardBuilder) Id(id int64) *DashboardBuilder {
    builder.internal.Id = id

    return builder
}

func (builder *DashboardBuilder) Title(title string) *DashboardBuilder {
    builder.internal.Title = title

    return builder
}

// will be expanded to []cog.Builder<DashboardLink>
func (builder *DashboardBuilder) Links(links []DashboardLink) *DashboardBuilder {
    builder.internal.Links = links

    return builder
}

// will be expanded to [][]cog.Builder<DashboardLink>
func (builder *DashboardBuilder) LinksOfLinks(linksOfLinks [][]DashboardLink) *DashboardBuilder {
    builder.internal.LinksOfLinks = linksOfLinks

    return builder
}

// will be expanded to cog.Builder<DashboardLink>
func (builder *DashboardBuilder) SingleLink(singleLink DashboardLink) *DashboardBuilder {
    builder.internal.SingleLink = singleLink

    return builder
}

func (builder *DashboardBuilder) applyDefaults() {
}
