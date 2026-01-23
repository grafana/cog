package struct_with_defaults

import (
	cog "github.com/grafana/cog/generated/cog"
)

var _ cog.Builder[Struct] = (*StructBuilder)(nil)

type StructBuilder struct {
    internal *Struct
    errors cog.BuildErrors
}

func NewStructBuilder() *StructBuilder {
	resource := NewStruct()
	builder := &StructBuilder{
		internal: resource,
		errors: make(cog.BuildErrors, 0),
	}

	return builder
}



func (builder *StructBuilder) Build() (Struct, error) {
	if err := builder.internal.Validate(); err != nil {
		return Struct{}, err
	}
	
	if len(builder.errors) > 0 {
	    return Struct{}, cog.MakeBuildErrors("struct_with_defaults.struct", builder.errors)
	}

	return *builder.internal, nil
}

func (builder *StructBuilder) AllFields(allFields cog.Builder[NestedStruct]) *StructBuilder {
    allFieldsResource, err := allFields.Build()
    if err != nil {
        builder.errors = append(builder.errors, err.(cog.BuildErrors)...)
        return builder
    }
    builder.internal.AllFields = allFieldsResource

    return builder
}

func (builder *StructBuilder) PartialFields(partialFields cog.Builder[NestedStruct]) *StructBuilder {
    partialFieldsResource, err := partialFields.Build()
    if err != nil {
        builder.errors = append(builder.errors, err.(cog.BuildErrors)...)
        return builder
    }
    builder.internal.PartialFields = partialFieldsResource

    return builder
}

func (builder *StructBuilder) EmptyFields(emptyFields cog.Builder[NestedStruct]) *StructBuilder {
    emptyFieldsResource, err := emptyFields.Build()
    if err != nil {
        builder.errors = append(builder.errors, err.(cog.BuildErrors)...)
        return builder
    }
    builder.internal.EmptyFields = emptyFieldsResource

    return builder
}

func (builder *StructBuilder) ComplexField(complexField struct {
    Uid string `json:"uid"`
    Nested struct {
        NestedVal string `json:"nestedVal"`
    } `json:"nested"`
    Array []string `json:"array"`
}) *StructBuilder {
    builder.internal.ComplexField = complexField

    return builder
}

func (builder *StructBuilder) PartialComplexField(partialComplexField struct {
    Uid string `json:"uid"`
    IntVal int64 `json:"intVal"`
}) *StructBuilder {
    builder.internal.PartialComplexField = partialComplexField

    return builder
}

