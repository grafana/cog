package somepkg

import (
	cog "github.com/grafana/cog/generated/cog"
	otherpkg "github.com/grafana/cog/generated/otherpkg"
)

var _ cog.Builder[Person] = (*PersonBuilder)(nil)

type PersonBuilder struct {
    internal *Person
    errors map[string]cog.BuildErrors
}

func NewPersonBuilder() *PersonBuilder {
	resource := &Person{}
	builder := &PersonBuilder{
		internal: resource,
		errors: make(map[string]cog.BuildErrors),
	}

	builder.applyDefaults()

	return builder
}

func (builder *PersonBuilder) Build() (Person, error) {
	var errs cog.BuildErrors

	for _, err := range builder.errors {
		errs = append(errs, cog.MakeBuildErrors("Person", err)...)
	}

	if len(errs) != 0 {
		return Person{}, errs
	}

	return *builder.internal, nil
}

func (builder *PersonBuilder) Name(name otherpkg.Name) *PersonBuilder {
    builder.internal.Name = name

    return builder
}

func (builder *PersonBuilder) applyDefaults() {
}
