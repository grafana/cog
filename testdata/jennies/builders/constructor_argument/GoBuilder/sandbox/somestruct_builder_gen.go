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
	resource := &SomeStruct{}
	builder := &SomeStructBuilder{
		internal: resource,
		errors: make(map[string]cog.BuildErrors),
	}

	builder.applyDefaults()
    builder.internal.Title = title

	return builder
}

func (builder *SomeStructBuilder) Build() (SomeStruct, error) {
	var errs cog.BuildErrors

	for _, err := range builder.errors {
		errs = append(errs, cog.MakeBuildErrors("SomeStruct", err)...)
	}

	if len(errs) != 0 {
		return SomeStruct{}, errs
	}

	return *builder.internal, nil
}

func (builder *SomeStructBuilder) Title(title string) *SomeStructBuilder {
    builder.internal.Title = title

    return builder
}

func (builder *SomeStructBuilder) applyDefaults() {
}
