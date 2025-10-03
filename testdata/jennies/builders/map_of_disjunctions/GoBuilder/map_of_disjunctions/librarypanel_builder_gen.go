package map_of_disjunctions

import (
	cog "github.com/grafana/cog/generated/cog"
)

var _ cog.Builder[LibraryPanel] = (*LibraryPanelBuilder)(nil)

type LibraryPanelBuilder struct {
    internal *LibraryPanel
    errors cog.BuildErrors
}

func NewLibraryPanelBuilder() *LibraryPanelBuilder {
	resource := NewLibraryPanel()
	builder := &LibraryPanelBuilder{
		internal: resource,
		errors: make(cog.BuildErrors, 0),
	}
    builder.internal.Kind = "Library"

	return builder
}



func (builder *LibraryPanelBuilder) Build() (LibraryPanel, error) {
	if err := builder.internal.Validate(); err != nil {
		return LibraryPanel{}, err
	}
	
	if len(builder.errors) > 0 {
	    return LibraryPanel{}, cog.MakeBuildErrors("map_of_disjunctions.libraryPanel", builder.errors)
	}

	return *builder.internal, nil
}

func (builder *LibraryPanelBuilder) Text(text string) *LibraryPanelBuilder {
    builder.internal.Text = text

    return builder
}
