package builderpkg

import (
	cog "github.com/grafana/cog/generated/cog"
	somepkg "github.com/grafana/cog/generated/somepkg"
)

var _ cog.Builder[somepkg.SomeStruct] = (*SomeNiceBuilderBuilder)(nil)

type SomeNiceBuilderBuilder struct {
    internal *somepkg.SomeStruct
    errors map[string]cog.BuildErrors
}

func NewSomeNiceBuilderBuilder() *SomeNiceBuilderBuilder {
	resource := &somepkg.SomeStruct{}
	builder := &SomeNiceBuilderBuilder{
		internal: resource,
		errors: make(map[string]cog.BuildErrors),
	}

	builder.applyDefaults()

	return builder
}

func (builder *SomeNiceBuilderBuilder) Build() (somepkg.SomeStruct, error) {
	var errs cog.BuildErrors

	for _, err := range builder.errors {
		errs = append(errs, cog.MakeBuildErrors("SomeNiceBuilder", err)...)
	}

	if len(errs) != 0 {
		return somepkg.SomeStruct{}, errs
	}

	return *builder.internal, nil
}

func (builder *SomeNiceBuilderBuilder) Title(title string) *SomeNiceBuilderBuilder {
    builder.internal.Title = title

    return builder
}

func (builder *SomeNiceBuilderBuilder) applyDefaults() {
}
