package discriminator_without_option

import (
	cog "github.com/grafana/cog/generated/cog"
)

var _ cog.Builder[NoShowFieldOption] = (*NoShowFieldOptionBuilder)(nil)

type NoShowFieldOptionBuilder struct {
    internal *NoShowFieldOption
    errors cog.BuildErrors
}

func NewNoShowFieldOptionBuilder() *NoShowFieldOptionBuilder {
	resource := NewNoShowFieldOption()
	builder := &NoShowFieldOptionBuilder{
		internal: resource,
		errors: make(cog.BuildErrors, 0),
	}

	return builder
}



func (builder *NoShowFieldOptionBuilder) Build() (NoShowFieldOption, error) {
	if err := builder.internal.Validate(); err != nil {
		return NoShowFieldOption{}, err
	}
	
	if len(builder.errors) > 0 {
	    return NoShowFieldOption{}, cog.MakeBuildErrors("discriminator_without_option.noShowFieldOption", builder.errors)
	}

	return *builder.internal, nil
}

func (builder *NoShowFieldOptionBuilder) Text(text string) *NoShowFieldOptionBuilder {
    builder.internal.Text = text

    return builder
}
