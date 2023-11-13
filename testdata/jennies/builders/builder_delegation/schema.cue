package builder_delegation

#DashboardLink: {
		title: string
		url: string
}

Dashboard: {
	id: int64
	title: string
	links: [...#DashboardLink] // we currently don't expand that to []cog.Builder<DashboardLink>
	singleLink: #DashboardLink // will be expanded to cog.Builder<DashboardLink>
}
