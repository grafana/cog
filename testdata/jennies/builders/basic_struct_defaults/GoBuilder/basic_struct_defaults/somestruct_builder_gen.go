package basic_struct_defaults

import (
	cog "github.com/grafana/cog/generated/cog"
)

var _ cog.Builder[SomeStruct] = (*SomeStructBuilder)(nil)

type SomeStructBuilder struct {
    internal *SomeStruct
    errors map[string]cog.BuildErrors
}

func NewSomeStructBuilder() *SomeStructBuilder {
	resource := &SomeStruct{}
	builder := &SomeStructBuilder{
		internal: resource,
		errors: make(map[string]cog.BuildErrors),
	}

	builder.applyDefaults()

	return builder
}

func (builder *SomeStructBuilder) Build() (SomeStruct, error) {
	if err := builder.internal.Validate(); err != nil {
		return SomeStruct{}, err
	}

	return *builder.internal, nil
}

func (builder *SomeStructBuilder) Id(id int64) *SomeStructBuilder {
    builder.internal.Id = id

    return builder
}

func (builder *SomeStructBuilder) Uid(uid string) *SomeStructBuilder {
    builder.internal.Uid = uid

    return builder
}

func (builder *SomeStructBuilder) Tags(tags []string) *SomeStructBuilder {
    builder.internal.Tags = tags

    return builder
}

func (builder *SomeStructBuilder) LiveNow(liveNow bool) *SomeStructBuilder {
    builder.internal.LiveNow = liveNow

    return builder
}

func (builder *SomeStructBuilder) applyDefaults() {
    builder.Id(42)
    builder.Uid("default-uid")
    builder.Tags([]string{"generated", "cog"})
    builder.LiveNow(true)
}
