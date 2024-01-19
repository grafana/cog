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
	var errs cog.BuildErrors

	for _, err := range builder.errors {
		errs = append(errs, cog.MakeBuildErrors("SomeStruct", err)...)
	}

	if len(errs) != 0 {
		return SomeStruct{}, errs
	}

	return *builder.internal, nil
}

func (builder *SomeStructBuilder) Editable() *SomeStructBuilder {
    builder.internal.Editable = true

    return builder
}

func (builder *SomeStructBuilder) Readonly() *SomeStructBuilder {
    builder.internal.Editable = false

    return builder
}

func (builder *SomeStructBuilder) AutoRefresh() *SomeStructBuilder {
            valAutoRefresh := true
    builder.internal.AutoRefresh = &valAutoRefresh

    return builder
}

func (builder *SomeStructBuilder) NoAutoRefresh() *SomeStructBuilder {
            valAutoRefresh := false
    builder.internal.AutoRefresh = &valAutoRefresh

    return builder
}

func (builder *SomeStructBuilder) applyDefaults() {
}
