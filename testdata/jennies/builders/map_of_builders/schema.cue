package map_of_builders

Panel: {
	title: string
}

Dashboard: {
	panels: {
		[string]: Panel
	}
}
