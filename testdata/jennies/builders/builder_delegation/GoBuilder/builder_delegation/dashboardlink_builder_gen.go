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
	resource := &DashboardLink{}
	builder := &DashboardLinkBuilder{
		internal: resource,
		errors: make(map[string]cog.BuildErrors),
	}

	builder.applyDefaults()

	return builder
}

func (builder *DashboardLinkBuilder) Build() (DashboardLink, error) {
	var errs cog.BuildErrors

	for _, err := range builder.errors {
		errs = append(errs, cog.MakeBuildErrors("DashboardLink", err)...)
	}

	if len(errs) != 0 {
		return DashboardLink{}, errs
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

func (builder *DashboardLinkBuilder) applyDefaults() {
}
