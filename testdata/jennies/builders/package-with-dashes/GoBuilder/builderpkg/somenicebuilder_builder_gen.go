package builderpkg

import (
	cog "github.com/grafana/cog/generated/cog"
	withdashes "github.com/grafana/cog/generated/with-dashes"
)

var _ cog.Builder[withdashes.SomeStruct] = (*SomeNiceBuilderBuilder)(nil)

type SomeNiceBuilderBuilder struct {
    internal *withdashes.SomeStruct
    errors cog.BuildErrors
}

func NewSomeNiceBuilderBuilder() *SomeNiceBuilderBuilder {
	resource := withdashes.NewSomeStruct()
	builder := &SomeNiceBuilderBuilder{
		internal: resource,
		errors: make(cog.BuildErrors, 0),
	}

	return builder
}



func (builder *SomeNiceBuilderBuilder) Build() (withdashes.SomeStruct, error) {
	if err := builder.internal.Validate(); err != nil {
		return withdashes.SomeStruct{}, err
	}
	
	if len(builder.errors) > 0 {
	    return withdashes.SomeStruct{}, cog.MakeBuildErrors("builderpkg.someNiceBuilder", builder.errors)
	}

	return *builder.internal, nil
}

func (builder *SomeNiceBuilderBuilder) Title(title string) *SomeNiceBuilderBuilder {
    builder.internal.Title = title

    return builder
}

