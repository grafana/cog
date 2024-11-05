package sandbox

import (
	cog "github.com/grafana/cog/generated/cog"
)

var _ cog.Builder[Dashboard] = (*DashboardBuilder)(nil)

type DashboardBuilder struct {
    internal *Dashboard
    errors map[string]cog.BuildErrors
}

func NewDashboardBuilder() *DashboardBuilder {
	builder := &DashboardBuilder{
		internal: NewDashboard(),
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

func (builder *DashboardBuilder) WithVariable(name string,value string) *DashboardBuilder {
    builder.internal.Variables = append(builder.internal.Variables, Variable{
        Name: name,
        Value: value,
    })

    return builder
}

