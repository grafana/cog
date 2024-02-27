package builder_delegation_in_disjunction

#DashboardLink: {
		title: string
		url: string
}

#ExternalLink: {
	url: string
}

Dashboard: {
	singleLinkOrString: #DashboardLink | string // will be expanded to cog.Builder<DashboardLink> | string
	linksOrStrings: [...(#DashboardLink | string)] // will be expanded to [](cog.Builder<DashboardLink> | string)
	disjunctionOfBuilders: #DashboardLink | #ExternalLink
}
