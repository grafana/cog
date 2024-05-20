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

func (builder *SomeStructBuilder) Time(from string,to string) *SomeStructBuilder {
<<<<<<< HEAD
if builder.internal.Time == nil {
    builder.internal.Time = &struct {
	From string `json:"from"`
	To string `json:"to"`
}{}
}
    builder.internal.Time.From = from
if builder.internal.Time == nil {
    builder.internal.Time = &struct {
	From string `json:"from"`
	To string `json:"to"`
=======
    if builder.internal.Time == nil {
	builder.internal.Time = &struct {
	From string `json:""`
	To string `json:""`
}{}
}
    builder.internal.Time.From = from
    if builder.internal.Time == nil {
	builder.internal.Time = &struct {
	From string `json:""`
	To string `json:""`
>>>>>>> a2493e55 (Start fixing tests)
}{}
}
    builder.internal.Time.To = to

    return builder
}

func (builder *SomeStructBuilder) applyDefaults() {
}
