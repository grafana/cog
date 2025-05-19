package builder_delegation

import (
	cog "github.com/grafana/cog/generated/cog"
)

var _ cog.Builder[DashboardLink] = (*DashboardLinkBuilder)(nil)

type DashboardLinkBuilder struct {
    internal *DashboardLink
    errors cog.BuildErrors
}

func NewDashboardLinkBuilder() *DashboardLinkBuilder {
	resource := NewDashboardLink()
	builder := &DashboardLinkBuilder{
		internal: resource,
		errors: make(cog.BuildErrors, 0),
	}

	return builder
}



func (builder *DashboardLinkBuilder) Build() (DashboardLink, error) {
	if err := builder.internal.Validate(); err != nil {
		return DashboardLink{}, err
	}
	
	if len(builder.errors) > 0 {
	    return DashboardLink{}, cog.MakeBuildErrors("builder_delegation.dashboardLink", builder.errors)
	}

	return *builder.internal, nil
}

func (builder *DashboardLinkBuilder) Title(title string) *DashboardLinkBuilder {
    builder.internal.Title = title

    return builder
}

func (builder *DashboardLinkBuilder) Url(url string) *DashboardLinkBuilder {
    builder.internal.Url = url

    return builder
}

