package map_of_builders

import (
	cog "github.com/grafana/cog/generated/cog"
)

var _ cog.Builder[Dashboard] = (*DashboardBuilder)(nil)

type DashboardBuilder struct {
    internal *Dashboard
    errors map[string]cog.BuildErrors
}

func NewDashboardBuilder() *DashboardBuilder {
	resource := NewDashboard()
	builder := &DashboardBuilder{
		internal: resource,
		errors: make(map[string]cog.BuildErrors),
	}

	return builder
}



func (builder *DashboardBuilder) Build() (Dashboard, error) {
	if err := builder.internal.Validate(); err != nil {
		return Dashboard{}, err
	}

	return *builder.internal, nil
}

func (builder *DashboardBuilder) Panels(panels map[string]cog.Builder[Panel]) *DashboardBuilder {
        panelsResource := make(map[string]Panel)
        for key1, val1 := range panels {
                panelsDepth1, err := val1.Build()
                if err != nil {
                    builder.errors["panels"] = err.(cog.BuildErrors)
                    return builder
                }
                panelsResource[key1] = panelsDepth1
        }
    builder.internal.Panels = panelsResource

    return builder
}

