package sandbox

import (
	cog "github.com/grafana/cog/generated/cog"
)

var _ cog.Builder[SomeStructWithDefaultEnum] = (*SomeStructWithDefaultEnumBuilder)(nil)

type SomeStructWithDefaultEnumBuilder struct {
    internal *SomeStructWithDefaultEnum
    errors cog.BuildErrors
}

func NewSomeStructWithDefaultEnumBuilder() *SomeStructWithDefaultEnumBuilder {
	resource := NewSomeStructWithDefaultEnum()
	builder := &SomeStructWithDefaultEnumBuilder{
		internal: resource,
		errors: make(cog.BuildErrors, 0),
	}

	return builder
}


func (builder *SomeStructWithDefaultEnumBuilder) Build() (SomeStructWithDefaultEnum, error) {
	if err := builder.internal.Validate(); err != nil {
		return SomeStructWithDefaultEnum{}, err
	}
	
	if len(builder.errors) > 0 {
	    return SomeStructWithDefaultEnum{}, cog.MakeBuildErrors("sandbox.someStructWithDefaultEnum", builder.errors)
	}

	return *builder.internal, nil
}

func (builder *SomeStructWithDefaultEnumBuilder) RecordError(path string, err error) *SomeStructWithDefaultEnumBuilder {
	builder.errors = append(builder.errors, cog.MakeBuildErrors(path, err)...)
	return builder
}

func (builder *SomeStructWithDefaultEnumBuilder) Data(key StringEnumWithDefault,value string) *SomeStructWithDefaultEnumBuilder {
if builder.internal.Data == nil {
    builder.internal.Data = map[StringEnumWithDefault]string{}
}
    builder.internal.Data[key] = value

    return builder
}

