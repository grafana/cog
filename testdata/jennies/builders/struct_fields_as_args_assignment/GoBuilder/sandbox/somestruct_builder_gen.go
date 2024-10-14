package sandbox

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

func (builder *SomeStructBuilder) Time(from string,to string) *SomeStructBuilder {
if builder.internal.Time == nil {
    builder.internal.Time = &struct {
	From string `json:"from"`
	To string `json:"to"`
}{}
}
    builder.internal.Time.From = from
    builder.internal.Time.To = to

    return builder
}

func (builder *SomeStructBuilder) applyDefaults() {
}
