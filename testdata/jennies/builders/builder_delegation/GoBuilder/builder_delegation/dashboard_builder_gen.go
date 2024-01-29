package builder_delegation

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

func (builder *DashboardBuilder) Links(links []cog.Builder[DashboardLink]) *DashboardBuilder {
        linksResources := make([]DashboardLink, 0, len(links))
        for _, r := range links {
            linksResource, err := r.Build()
            if err != nil {
                builder.errors["links"] = err.(cog.BuildErrors)
                return builder
            }
            linksResources = append(linksResources, linksResource)
        }
    builder.internal.Links = linksResources

    return builder
}

func (builder *DashboardBuilder) SingleLink(singleLink cog.Builder[DashboardLink]) *DashboardBuilder {
        singleLinkResource, err := singleLink.Build()
        if err != nil {
            builder.errors["singleLink"] = err.(cog.BuildErrors)
            return builder
        }
    builder.internal.SingleLink = singleLinkResource

    return builder
}

func (builder *DashboardBuilder) applyDefaults() {
}
