package map_of_builders

Panel: {
	title: string
}

Dashboard: {
	panels: {
		[string]: Panel
	}

	facets: [string]: [...#FacetValue]
}

#FacetValue: {
	value: string
	count: int64
}
