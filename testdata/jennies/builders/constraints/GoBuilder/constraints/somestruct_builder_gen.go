package constraints

import (
	cog "github.com/grafana/cog/generated/cog"
	"errors"
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

func (builder *SomeStructBuilder) Id(id uint64) *SomeStructBuilder {
    if !(id >= 5) {
        builder.errors["id"] = cog.MakeBuildErrors("id", errors.New("id must be >= 5"))
        return builder
    }
    if !(id < 10) {
        builder.errors["id"] = cog.MakeBuildErrors("id", errors.New("id must be < 10"))
        return builder
    }
    builder.internal.Id = id

    return builder
}

func (builder *SomeStructBuilder) Title(title string) *SomeStructBuilder {
    if !(len([]rune(title)) >= 1) {
        builder.errors["title"] = cog.MakeBuildErrors("title", errors.New("len([]rune(title)) must be >= 1"))
        return builder
    }
    builder.internal.Title = title

    return builder
}

func (builder *SomeStructBuilder) applyDefaults() {
}
