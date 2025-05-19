package builder_delegation

import (
	cog "github.com/grafana/cog/generated/cog"
)

var _ cog.Builder[Dashboard] = (*DashboardBuilder)(nil)

type DashboardBuilder struct {
    internal *Dashboard
    errors cog.BuildErrors
}

func NewDashboardBuilder() *DashboardBuilder {
	resource := NewDashboard()
	builder := &DashboardBuilder{
		internal: resource,
		errors: make(cog.BuildErrors, 0),
	}

	return builder
}



func (builder *DashboardBuilder) Build() (Dashboard, error) {
	if err := builder.internal.Validate(); err != nil {
		return Dashboard{}, err
	}
	
	if len(builder.errors) > 0 {
	    return Dashboard{}, cog.MakeBuildErrors("builder_delegation.dashboard", builder.errors)
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
func (builder *DashboardBuilder) Links(links []cog.Builder[DashboardLink]) *DashboardBuilder {
        linksResources := make([]DashboardLink, 0, len(links))
        for _, r1 := range links {
                linksDepth1, err := r1.Build()
                if err != nil {
                    builder.errors = append(builder.errors, err.(cog.BuildErrors)...)
                    return builder
                }
                linksResources = append(linksResources, linksDepth1)
        }
    builder.internal.Links = linksResources

    return builder
}

// will be expanded to [][]cog.Builder<DashboardLink>
func (builder *DashboardBuilder) LinksOfLinks(linksOfLinks [][]cog.Builder[DashboardLink]) *DashboardBuilder {
        linksOfLinksResources := make([][]DashboardLink, 0, len(linksOfLinks))
        for _, r1 := range linksOfLinks {
                linksOfLinksDepth1 := make([]DashboardLink, 0)
        for _, r2 := range r1 {
                linksOfLinksDepth2, err := r2.Build()
                if err != nil {
                    builder.errors = append(builder.errors, err.(cog.BuildErrors)...)
                    return builder
                }
                linksOfLinksDepth1 = append(linksOfLinksDepth1, linksOfLinksDepth2)
        }

                linksOfLinksResources = append(linksOfLinksResources, linksOfLinksDepth1)
        }
    builder.internal.LinksOfLinks = linksOfLinksResources

    return builder
}

// will be expanded to cog.Builder<DashboardLink>
func (builder *DashboardBuilder) SingleLink(singleLink cog.Builder[DashboardLink]) *DashboardBuilder {
    singleLinkResource, err := singleLink.Build()
    if err != nil {
        builder.errors = append(builder.errors, err.(cog.BuildErrors)...)
        return builder
    }
    builder.internal.SingleLink = singleLinkResource

    return builder
}

