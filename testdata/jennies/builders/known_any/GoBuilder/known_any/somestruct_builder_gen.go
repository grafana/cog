package known_any

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
	    return SomeStruct{}, cog.MakeBuildErrors("known_any.someStruct", builder.errors)
	}

	return *builder.internal, nil
}

func (builder *SomeStructBuilder) Title(title string) *SomeStructBuilder {
if builder.internal.Config == nil {
    builder.internal.Config = NewConfig()
}
    builder.internal.Config.(*Config).Title = title

    return builder
}

