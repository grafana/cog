package map_of_disjunctions

import (
	cog "github.com/grafana/cog/generated/cog"
)

var _ cog.Builder[PanelOrLibraryPanel] = (*PanelOrLibraryPanelBuilder)(nil)

type PanelOrLibraryPanelBuilder struct {
    internal *PanelOrLibraryPanel
    errors cog.BuildErrors
}

func NewPanelOrLibraryPanelBuilder() *PanelOrLibraryPanelBuilder {
	resource := NewPanelOrLibraryPanel()
	builder := &PanelOrLibraryPanelBuilder{
		internal: resource,
		errors: make(cog.BuildErrors, 0),
	}

	return builder
}



func (builder *PanelOrLibraryPanelBuilder) Build() (PanelOrLibraryPanel, error) {
	if err := builder.internal.Validate(); err != nil {
		return PanelOrLibraryPanel{}, err
	}
	
	if len(builder.errors) > 0 {
	    return PanelOrLibraryPanel{}, cog.MakeBuildErrors("map_of_disjunctions.panelOrLibraryPanel", builder.errors)
	}

	return *builder.internal, nil
}

func (builder *PanelOrLibraryPanelBuilder) Panel(panel cog.Builder[Panel]) *PanelOrLibraryPanelBuilder {
    panelResource, err := panel.Build()
    if err != nil {
        builder.errors = append(builder.errors, err.(cog.BuildErrors)...)
        return builder
    }
    builder.internal.Panel = &panelResource

    return builder
}

func (builder *PanelOrLibraryPanelBuilder) LibraryPanel(libraryPanel cog.Builder[LibraryPanel]) *PanelOrLibraryPanelBuilder {
    libraryPanelResource, err := libraryPanel.Build()
    if err != nil {
        builder.errors = append(builder.errors, err.(cog.BuildErrors)...)
        return builder
    }
    builder.internal.LibraryPanel = &libraryPanelResource

    return builder
}
