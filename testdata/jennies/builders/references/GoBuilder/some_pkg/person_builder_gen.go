package some_pkg

import (
	cog "github.com/grafana/cog/generated/cog"
	other_pkg "github.com/grafana/cog/generated/other_pkg"
)

var _ cog.Builder[Person] = (*PersonBuilder)(nil)

type PersonBuilder struct {
    internal *Person
    errors map[string]cog.BuildErrors
}

func NewPersonBuilder() *PersonBuilder {
	resource := NewPerson()
	builder := &PersonBuilder{
		internal: resource,
		errors: make(map[string]cog.BuildErrors),
	}

	return builder
}



func (builder *PersonBuilder) Build() (Person, error) {
	if err := builder.internal.Validate(); err != nil {
		return Person{}, err
	}

	return *builder.internal, nil
}

func (builder *PersonBuilder) Name(name other_pkg.Name) *PersonBuilder {
    builder.internal.Name = name

    return builder
}

