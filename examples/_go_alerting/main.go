package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/grafana/cog/generated/go/alerting"
	"github.com/grafana/cog/generated/go/cog/plugins"
	"github.com/grafana/cog/generated/go/expr"
	"github.com/grafana/cog/generated/go/prometheus"
)

func main() {
	// Required to correctly unmarshal dataqueries
	plugins.RegisterDefaultPlugins()

	interval := 5 * time.Minute

	alert, err := alerting.NewRuleGroupBuilder("cog rule group").
		Interval(alerting.Duration(interval.Seconds())). // type is confusing
		WithRule(alerting.NewRuleBuilder("cog alert rule name").
			For("5m").
			RuleGroup("group").
			Annotations(map[string]string{
				"summary":     "summary",
				"description": "desc",
				"runbook_url": "https://foo.local",
			}).
			Labels(map[string]string{
				"owner": "team-name",
			}).
			Condition("B").
			ExecErrState(alerting.RuleExecErrStateAlerting).
			NoDataState(alerting.RuleNoDataStateOK).
			WithQuery(alerting.NewQueryBuilder("A").
				DatasourceUid("cdjvy8rnmh0qof").
				RelativeTimeRange(600, 0). // TODO
				Model(
					prometheus.NewDataqueryBuilder().
						// RefId("A"). // needed?
						Expr("go_memstats_alloc_bytes_total").
						Instant().
						LegendFormat("__auto"),
				)).
			WithQuery(alerting.NewQueryBuilder("B").
				DatasourceUid("__expr__").
				RelativeTimeRange(600, 0). // TODO
				Model(
					expr.NewTypeReduceBuilder().
						Reducer("last").
						Expression("A"),
				)),
		).
		Build()
	if err != nil {
		panic(err)
	}

	alertAsBytes, err := json.MarshalIndent(alert, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(alertAsBytes))
}
