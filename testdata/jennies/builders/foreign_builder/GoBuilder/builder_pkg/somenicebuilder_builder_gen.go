package builder_pkg

import (
	cog "github.com/grafana/cog/generated/cog"
	some_pkg "github.com/grafana/cog/generated/some_pkg"
)

var _ cog.Builder[some_pkg.SomeStruct] = (*SomeNiceBuilderBuilder)(nil)

type SomeNiceBuilderBuilder struct {
    internal *some_pkg.SomeStruct
    errors cog.BuildErrors
}

func NewSomeNiceBuilderBuilder() *SomeNiceBuilderBuilder {
	resource := some_pkg.NewSomeStruct()
	builder := &SomeNiceBuilderBuilder{
		internal: resource,
		errors: make(cog.BuildErrors, 0),
	}

	return builder
}



func (builder *SomeNiceBuilderBuilder) Build() (some_pkg.SomeStruct, error) {
	if err := builder.internal.Validate(); err != nil {
		return some_pkg.SomeStruct{}, err
	}

	return *builder.internal, cog.MakeBuildErrors("builder_pkg.someNiceBuilder", builder.errors)
}

func (builder *SomeNiceBuilderBuilder) Title(title string) *SomeNiceBuilderBuilder {
    builder.internal.Title = title

    return builder
}

