# Grafana Foundation SDK – Go

A set of tools, types and *builder libraries* for building and manipulating Grafana objects in Go.

> [!NOTE]
> This branch contains **types and builders generated for Grafana {{ .Extra.GrafanaVersion }}.**
> Other supported versions of Grafana can be found at [this repository's root](https://github.com/grafana/grafana-foundation-sdk/).

## Installing

```shell
go get github.com/grafana/grafana-foundation-sdk/go@{{ .Extra.ReleaseBranch }}
```

## Example usage

### Building a dashboard

```go
package main

import (
	"encoding/json"
	"fmt"

	"github.com/grafana/grafana-foundation-sdk/go/common"
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
	"github.com/grafana/grafana-foundation-sdk/go/prometheus"
	"github.com/grafana/grafana-foundation-sdk/go/timeseries"
)

func main() {
	builder := dashboard.NewDashboardBuilder("Sample dashboard").
		Uid("generated-from-go").
		Tags([]string{"generated", "from", "go"}).
		Refresh("1m").
		Time("now-30m", "now").
		Timezone(common.TimeZoneBrowser).
		WithRow(dashboard.NewRowBuilder("Overview")).
		WithPanel(
			timeseries.NewPanelBuilder().
				Title("Network Received").
				Unit("bps").
				Min(0).
				WithTarget(
					prometheus.NewDataqueryBuilder().
						Expr(`rate(node_network_receive_bytes_total{job="integrations/raspberrypi-node", device!="lo"}[$__rate_interval]) * 8`).
						LegendFormat({{ `"{{ device }}"` }}),
				),
		)

	sampleDashboard, err := builder.Build()
	if err != nil {
		panic(err)
	}
	dashboardJson, err := json.MarshalIndent(sampleDashboard, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(dashboardJson))
}
```

### Unmarshaling a dashboard

```go
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/grafana/grafana-foundation-sdk/go/cog/plugins"
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
)

func main() {
	// Required to correctly unmarshal panels and dataqueries
	plugins.RegisterDefaultPlugins()

	dashboardJSON, err := os.ReadFile("dashboard.json")
	if err != nil {
		panic(err)
	}

	sampleDashboard := &dashboard.Dashboard{}
	if err := json.Unmarshal(dashboardJSON, sampleDashboard); err != nil {
		panic(fmt.Sprintf("%s", err))
	}

	fmt.Printf("%#v\n", sampleDashboard)
}
```

### Defining a custom query type

While the SDK ships with support for all core datasources and their query types,
it can be extended for private/third-party plugins.

To do so, define a type and a builder for the custom query:

```go
package main

import (
	"encoding/json"

	"github.com/grafana/grafana-foundation-sdk/go/cog"
	cogvariants "github.com/grafana/grafana-foundation-sdk/go/cog/variants"
)

type CustomQuery struct {
	// RefId and Hide are expected on all queries
	RefId        *string `json:"refId,omitempty"`
	Hide         *bool   `json:"hide,omitempty"`

	// Query is specific to the CustomQuery type
	Query        string  `json:"query,omitempty"`
}

// Let cog know that CustomQuery is a Dataquery variant
func (resource CustomQuery) ImplementsDataqueryVariant() {}

func CustomQueryVariantConfig() cogvariants.DataqueryConfig {
	return cogvariants.DataqueryConfig{
		Identifier: "custom", // datasource plugin ID
		DataqueryUnmarshaler: func(raw []byte) (cogvariants.Dataquery, error) {
			dataquery := &CustomQuery{}

			if err := json.Unmarshal(raw, dataquery); err != nil {
				return nil, err
			}

			return dataquery, nil
		},
	}
}

// Compile-time check to ensure that CustomQueryBuilder indeed is
// a builder for cogvariants.Dataquery
var _ cog.Builder[cogvariants.Dataquery] = (*CustomQueryBuilder)(nil)

type CustomQueryBuilder struct {
	internal *CustomQuery
}

func NewCustomQueryBuilder(query string) *CustomQueryBuilder {
	return &CustomQueryBuilder{
		internal: &CustomQuery{Query: query},
	}
}

func (builder *CustomQueryBuilder) Build() (cogvariants.Dataquery, error) {
	return *builder.internal, nil
}

func (builder *CustomQueryBuilder) RefId(refId string) *CustomQueryBuilder {
	builder.internal.RefId = &refId
	return builder
}

func (builder *CustomQueryBuilder) Hide(hide bool) *CustomQueryBuilder {
	builder.internal.Hide = &hide
	return builder
}
```

Register the type with cog, and use it as usual to build a dashboard:

```go
package main

import (
	"encoding/json"
	"fmt"

	"github.com/grafana/grafana-foundation-sdk/go/cog"
	"github.com/grafana/grafana-foundation-sdk/go/cog/plugins"
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
	"github.com/grafana/grafana-foundation-sdk/go/timeseries"
)

func main() {
	// Required to correctly unmarshal panels and dataqueries
	plugins.RegisterDefaultPlugins()

	// This lets cog know about the newly created query type and how to unmarshal it.
	cog.NewRuntime().RegisterDataqueryVariant(CustomQueryVariantConfig())

	sampleDashboard, err := dashboard.NewDashboardBuilder("Custom query type").
		Uid("test-custom-query-type").
		Refresh("1m").
		Time("now-30m", "now").
		WithRow(dashboard.NewRowBuilder("Overview")).
		WithPanel(
			timeseries.NewPanelBuilder().
				Title("Sample panel").
				WithTarget(
					NewCustomQueryBuilder("query here").LegendFormat("{{ cpu }}"),
				),
		).
		Build()
	if err != nil {
		panic(err)
	}

	dashboardJson, err := json.MarshalIndent(sampleDashboard, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(dashboardJson))
}
```

## Maturity

> [!WARNING]
> The code in this repository should be considered experimental. Documentation is only
available alongside the code. It comes with no support, but we are keen to receive
feedback and suggestions on how to improve it, though we cannot commit
to resolution of any particular issue.

Grafana Labs defines experimental features as follows:

> Projects and features in the Experimental stage are supported only by the Engineering
teams; on-call support is not available. Documentation is either limited or not provided
outside of code comments. No SLA is provided.
>
> Experimental projects or features are primarily intended for open source engineers who
want to participate in ensuring systems stability, and to gain consensus and approval
for open source governance projects.
>
> Projects and features in the Experimental phase are not meant to be used in production
environments, and the risks are unknown/high.

## License

[Apache 2.0 License](./LICENSE)
