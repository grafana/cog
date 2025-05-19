package sandbox

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
	    return Dashboard{}, cog.MakeBuildErrors("sandbox.dashboard", builder.errors)
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

