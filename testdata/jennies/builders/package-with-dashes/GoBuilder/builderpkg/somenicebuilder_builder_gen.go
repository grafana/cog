package builderpkg

import (
	cog "github.com/grafana/cog/generated/cog"
	withdashes "github.com/grafana/cog/generated/with-dashes"
)

var _ cog.Builder[withdashes.SomeStruct] = (*SomeNiceBuilderBuilder)(nil)

type SomeNiceBuilderBuilder struct {
    internal *withdashes.SomeStruct
    errors map[string]cog.BuildErrors
}

func NewSomeNiceBuilderBuilder() *SomeNiceBuilderBuilder {
	resource := &withdashes.SomeStruct{}
	builder := &SomeNiceBuilderBuilder{
		internal: resource,
		errors: make(map[string]cog.BuildErrors),
	}

	builder.applyDefaults()

	return builder
}

func (builder *SomeNiceBuilderBuilder) Build() (withdashes.SomeStruct, error) {
	if err := builder.internal.Validate(); err != nil {
		return withdashes.SomeStruct{}, err
	}

	return *builder.internal, nil
}

func (builder *SomeNiceBuilderBuilder) Title(title string) *SomeNiceBuilderBuilder {
    builder.internal.Title = title

    return builder
}

func (builder *SomeNiceBuilderBuilder) applyDefaults() {
}
