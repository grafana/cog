package struct_with_defaults

import (
	cog "github.com/grafana/cog/generated/cog"
)

var _ cog.Builder[NestedStruct] = (*NestedStructBuilder)(nil)

type NestedStructBuilder struct {
    internal *NestedStruct
    errors cog.BuildErrors
}

func NewNestedStructBuilder() *NestedStructBuilder {
	resource := NewNestedStruct()
	builder := &NestedStructBuilder{
		internal: resource,
		errors: make(cog.BuildErrors, 0),
	}

	return builder
}



func (builder *NestedStructBuilder) Build() (NestedStruct, error) {
	if err := builder.internal.Validate(); err != nil {
		return NestedStruct{}, err
	}
	
	if len(builder.errors) > 0 {
	    return NestedStruct{}, cog.MakeBuildErrors("struct_with_defaults.nestedStruct", builder.errors)
	}

	return *builder.internal, nil
}

func (builder *NestedStructBuilder) StringVal(stringVal string) *NestedStructBuilder {
    builder.internal.StringVal = stringVal

    return builder
}

func (builder *NestedStructBuilder) IntVal(intVal int64) *NestedStructBuilder {
    builder.internal.IntVal = intVal

    return builder
}

