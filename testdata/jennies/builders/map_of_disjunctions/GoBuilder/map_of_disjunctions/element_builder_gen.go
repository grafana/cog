package map_of_disjunctions

import (
	cog "github.com/grafana/cog/generated/cog"
)

var _ cog.Builder[Element] = (*ElementBuilder)(nil)

type ElementBuilder struct {
    internal *Element
    errors cog.BuildErrors
}

func NewElementBuilder() *ElementBuilder {
	resource := NewElement()
	builder := &ElementBuilder{
		internal: resource,
		errors: make(cog.BuildErrors, 0),
	}

	return builder
}



func (builder *ElementBuilder) Build() (Element, error) {
	if err := builder.internal.Validate(); err != nil {
		return Element{}, err
	}
	
	if len(builder.errors) > 0 {
	    return Element{}, cog.MakeBuildErrors("map_of_disjunctions.element", builder.errors)
	}

	return *builder.internal, nil
}

func (builder *ElementBuilder) Panel(panel cog.Builder[Panel]) *ElementBuilder {
    panelResource, err := panel.Build()
    if err != nil {
        builder.errors = append(builder.errors, err.(cog.BuildErrors)...)
        return builder
    }
    builder.internal.Panel = &panelResource

    return builder
}

func (builder *ElementBuilder) LibraryPanel(libraryPanel cog.Builder[LibraryPanel]) *ElementBuilder {
    libraryPanelResource, err := libraryPanel.Build()
    if err != nil {
        builder.errors = append(builder.errors, err.(cog.BuildErrors)...)
        return builder
    }
    builder.internal.LibraryPanel = &libraryPanelResource

    return builder
}
