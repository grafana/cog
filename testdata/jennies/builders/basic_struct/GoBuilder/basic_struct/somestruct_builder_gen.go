package basic_struct

import (
	cog "github.com/grafana/cog/generated/cog"
)

var _ cog.Builder[SomeStruct] = (*SomeStructBuilder)(nil)

// SomeStruct, to hold data.
type SomeStructBuilder struct {
    internal *SomeStruct
    errors map[string]cog.BuildErrors
}

func NewSomeStructBuilder() *SomeStructBuilder {
	resource := NewSomeStruct()
	builder := &SomeStructBuilder{
		internal: resource,
		errors: make(map[string]cog.BuildErrors),
	}

	return builder
}



func (builder *SomeStructBuilder) Build() (SomeStruct, error) {
	if err := builder.internal.Validate(); err != nil {
		return SomeStruct{}, err
	}

	return *builder.internal, nil
}

// id identifies something. Weird, right?
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

// This thing could be live.
// Or maybe not.
func (builder *SomeStructBuilder) LiveNow(liveNow bool) *SomeStructBuilder {
    builder.internal.LiveNow = liveNow

    return builder
}

