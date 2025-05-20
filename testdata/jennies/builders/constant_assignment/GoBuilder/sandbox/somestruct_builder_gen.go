package sandbox

import (
	cog "github.com/grafana/cog/generated/cog"
)

var _ cog.Builder[SomeStruct] = (*SomeStructBuilder)(nil)

type SomeStructBuilder struct {
    internal *SomeStruct
    errors cog.BuildErrors
}

func NewSomeStructBuilder() *SomeStructBuilder {
	resource := NewSomeStruct()
	builder := &SomeStructBuilder{
		internal: resource,
		errors: make(cog.BuildErrors, 0),
	}

	return builder
}



func (builder *SomeStructBuilder) Build() (SomeStruct, error) {
	if err := builder.internal.Validate(); err != nil {
		return SomeStruct{}, err
	}
	
	if len(builder.errors) > 0 {
	    return SomeStruct{}, cog.MakeBuildErrors("sandbox.someStruct", builder.errors)
	}

	return *builder.internal, nil
}

func (builder *SomeStructBuilder) Editable() *SomeStructBuilder {
    builder.internal.Editable = true

    return builder
}

func (builder *SomeStructBuilder) Readonly() *SomeStructBuilder {
    builder.internal.Editable = false

    return builder
}

func (builder *SomeStructBuilder) AutoRefresh() *SomeStructBuilder {
            valAutoRefresh := true
    builder.internal.AutoRefresh = &valAutoRefresh

    return builder
}

func (builder *SomeStructBuilder) NoAutoRefresh() *SomeStructBuilder {
            valAutoRefresh := false
    builder.internal.AutoRefresh = &valAutoRefresh

    return builder
}

