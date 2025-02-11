package builder_delegation

import (
	cog "github.com/grafana/cog/generated/cog"
)

var _ cog.Builder[DashboardLink] = (*DashboardLinkBuilder)(nil)

type DashboardLinkBuilder struct {
    internal *DashboardLink
    errors map[string]cog.BuildErrors
}

func NewDashboardLinkBuilder() *DashboardLinkBuilder {
	resource := NewDashboardLink()
	builder := &DashboardLinkBuilder{
		internal: resource,
		errors: make(map[string]cog.BuildErrors),
	}

	return builder
}



func (builder *DashboardLinkBuilder) Build() (DashboardLink, error) {
	if err := builder.internal.Validate(); err != nil {
		return DashboardLink{}, err
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

