package sandbox

import (
	cog "github.com/grafana/cog/generated/cog"
)

var _ cog.Builder[SomeStruct] = (*SomeStructBuilder)(nil)

type SomeStructBuilder struct {
    internal *SomeStruct
    errors map[string]cog.BuildErrors
}

func NewSomeStructBuilder(title string) *SomeStructBuilder {
	resource := NewSomeStruct()
	builder := &SomeStructBuilder{
		internal: resource,
		errors: make(map[string]cog.BuildErrors),
	}
    builder.internal.Title = title

	return builder
}



func (builder *SomeStructBuilder) Build() (SomeStruct, error) {
	if err := builder.internal.Validate(); err != nil {
		return SomeStruct{}, err
	}

	return *builder.internal, nil
}

func (builder *SomeStructBuilder) Title(title string) *SomeStructBuilder {
    builder.internal.Title = title

    return builder
}

