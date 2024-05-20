package structwithdefaults

import (
	cog "github.com/grafana/cog/generated/cog"
)

var _ cog.Builder[Struct] = (*StructBuilder)(nil)

type StructBuilder struct {
    internal *Struct
    errors map[string]cog.BuildErrors
}

func NewStructBuilder() *StructBuilder {
	resource := &Struct{}
	builder := &StructBuilder{
		internal: resource,
		errors: make(map[string]cog.BuildErrors),
	}

	builder.applyDefaults()

	return builder
}

func (builder *StructBuilder) Build() (Struct, error) {
	var errs cog.BuildErrors

	for _, err := range builder.errors {
		errs = append(errs, cog.MakeBuildErrors("Struct", err)...)
	}

	if len(errs) != 0 {
		return Struct{}, errs
	}

	return *builder.internal, nil
}

func (builder *StructBuilder) AllFields(allFields NestedStruct) *StructBuilder {
    builder.internal.AllFields = allFields

    return builder
}

func (builder *StructBuilder) PartialFields(partialFields NestedStruct) *StructBuilder {
    builder.internal.PartialFields = partialFields

    return builder
}

func (builder *StructBuilder) EmptyFields(emptyFields NestedStruct) *StructBuilder {
    builder.internal.EmptyFields = emptyFields

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
	AVal int64 `json:"aVal"`
}) *StructBuilder {
    builder.internal.PartialComplexField = partialComplexField

    return builder
}

func (builder *StructBuilder) applyDefaults() {
    builder.AllFields(unknown)
    builder.PartialFields(unknown)
    builder.ComplexField(struct {
 Uid string `json:"uid"`
Nested struct {
 NestedVal string `json:"nestedVal"`
 } `json:"nested"`
Array []string `json:"array"`
 } {
 Array: []string{"hello"},
Nested: struct {
 NestedVal string `json:"nestedVal"`
 } {
 NestedVal: "nested",
},
Uid: "myUID",
 })
}
