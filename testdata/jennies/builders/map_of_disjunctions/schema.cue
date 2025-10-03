package map_of_disjunctions

Element: Panel | LibraryPanel

Panel: {
	kind: "Panel"
	title: string
}

LibraryPanel: {
	kind: "Library"
	text: string
}

Dashboard: {
	panels: {
		[string]: Element
	}
}
