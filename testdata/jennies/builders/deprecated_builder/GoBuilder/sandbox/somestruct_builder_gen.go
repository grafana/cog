package sandbox

import (
	cog "github.com/grafana/cog/generated/cog"
)

var _ cog.Builder[SomeStruct] = (*SomeStructBuilder)(nil)

// Deprecated: This builder is deprecated. Don't use. Please.
type SomeStructBuilder struct {
    internal *SomeStruct
    errors cog.BuildErrors
}
// Deprecated: "This builder is deprecated. Don't use. Please."
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

func (builder *SomeStructBuilder) RecordError(path string, err error) *SomeStructBuilder {
	builder.errors = append(builder.errors, cog.MakeBuildErrors(path, err)...)
	return builder
}

func (builder *SomeStructBuilder) Title(title string) *SomeStructBuilder {
    builder.internal.Title = title

    return builder
}

