package anonymous_struct

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

	return *builder.internal, cog.MakeBuildErrors("anonymous_struct.someStruct", builder.errors)
}

func (builder *SomeStructBuilder) Time(time struct {
    From string `json:"from"`
    To string `json:"to"`
}) *SomeStructBuilder {
    builder.internal.Time = &time

    return builder
}

