package builder_delegation

#DashboardLink: {
		title: string
		url: string
}

Dashboard: {
	id: int64
	title: string
	links: [...#DashboardLink] // will be expanded to []cog.Builder<DashboardLink>
	linksOfLinks: [...[...#DashboardLink]] // will be expanded to [][]cog.Builder<DashboardLink>
	singleLink: #DashboardLink // will be expanded to cog.Builder<DashboardLink>
}
