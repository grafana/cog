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

func (builder *SomeStructBuilder) Annotations(key string,value string) *SomeStructBuilder {
if builder.internal.Annotations == nil {
    builder.internal.Annotations = map[string]string{}
}
    builder.internal.Annotations[key] = value

    return builder
}

