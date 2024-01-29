package builder_pkg

import (
	cog "github.com/grafana/cog/generated/cog"
	some_pkg "github.com/grafana/cog/generated/some_pkg"
)

var _ cog.Builder[some_pkg.SomeStruct] = (*SomeNiceBuilderBuilder)(nil)

type SomeNiceBuilderBuilder struct {
    internal *some_pkg.SomeStruct
    errors map[string]cog.BuildErrors
}

func NewSomeNiceBuilderBuilder() *SomeNiceBuilderBuilder {
	resource := &some_pkg.SomeStruct{}
	builder := &SomeNiceBuilderBuilder{
		internal: resource,
		errors: make(map[string]cog.BuildErrors),
	}

	builder.applyDefaults()

	return builder
}

func (builder *SomeNiceBuilderBuilder) Build() (some_pkg.SomeStruct, error) {
	var errs cog.BuildErrors

	for _, err := range builder.errors {
		errs = append(errs, cog.MakeBuildErrors("SomeNiceBuilder", err)...)
	}

	if len(errs) != 0 {
		return some_pkg.SomeStruct{}, errs
	}

	return *builder.internal, nil
}

func (builder *SomeNiceBuilderBuilder) Title(title string) *SomeNiceBuilderBuilder {
    builder.internal.Title = title

    return builder
}

func (builder *SomeNiceBuilderBuilder) applyDefaults() {
}
