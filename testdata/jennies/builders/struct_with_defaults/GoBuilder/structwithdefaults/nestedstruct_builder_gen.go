package structwithdefaults

import (
	cog "github.com/grafana/cog/generated/cog"
)

var _ cog.Builder[NestedStruct] = (*NestedStructBuilder)(nil)

type NestedStructBuilder struct {
    internal *NestedStruct
    errors map[string]cog.BuildErrors
}

func NewNestedStructBuilder() *NestedStructBuilder {
	resource := &NestedStruct{}
	builder := &NestedStructBuilder{
		internal: resource,
		errors: make(map[string]cog.BuildErrors),
	}

	builder.applyDefaults()

	return builder
}

func (builder *NestedStructBuilder) Build() (NestedStruct, error) {
	var errs cog.BuildErrors

	for _, err := range builder.errors {
		errs = append(errs, cog.MakeBuildErrors("NestedStruct", err)...)
	}

	if len(errs) != 0 {
		return NestedStruct{}, errs
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

func (builder *NestedStructBuilder) applyDefaults() {
}
