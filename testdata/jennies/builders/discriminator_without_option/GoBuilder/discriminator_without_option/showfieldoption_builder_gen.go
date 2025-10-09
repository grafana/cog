package discriminator_without_option

import (
	cog "github.com/grafana/cog/generated/cog"
)

var _ cog.Builder[ShowFieldOption] = (*ShowFieldOptionBuilder)(nil)

type ShowFieldOptionBuilder struct {
    internal *ShowFieldOption
    errors cog.BuildErrors
}

func NewShowFieldOptionBuilder() *ShowFieldOptionBuilder {
	resource := NewShowFieldOption()
	builder := &ShowFieldOptionBuilder{
		internal: resource,
		errors: make(cog.BuildErrors, 0),
	}

	return builder
}



func (builder *ShowFieldOptionBuilder) Build() (ShowFieldOption, error) {
	if err := builder.internal.Validate(); err != nil {
		return ShowFieldOption{}, err
	}
	
	if len(builder.errors) > 0 {
	    return ShowFieldOption{}, cog.MakeBuildErrors("discriminator_without_option.showFieldOption", builder.errors)
	}

	return *builder.internal, nil
}

func (builder *ShowFieldOptionBuilder) Field(field AnEnum) *ShowFieldOptionBuilder {
    builder.internal.Field = field

    return builder
}

func (builder *ShowFieldOptionBuilder) Text(text string) *ShowFieldOptionBuilder {
    builder.internal.Text = text

    return builder
}
