package grizzly

#ApiVersion: "grizzly.grafana.com/v1alpha1" @cog(kind="enum",memberNames="v1alpha1")

#Kind: "Dashboard" | "DashboardFolder" | "Datasource" | "Team" @cog(kind="enum")

#Metadata: [string]: _

Resource: {
	apiVersion: #ApiVersion & "grizzly.grafana.com/v1alpha1"
	kind: #Kind
	metadata: #Metadata
	spec: _
}
