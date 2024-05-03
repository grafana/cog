package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/grafana/cog/generated/alerting"
	"github.com/grafana/cog/generated/cog"
	"github.com/grafana/cog/generated/cog/plugins"
	"github.com/grafana/cog/generated/prometheus"
)

func main() {
	// Required to correctly unmarshal dataqueries
	plugins.RegisterDefaultPlugins()

	interval := 5 * time.Minute

	alert, err := alerting.NewRuleGroupBuilder().
		Title("rule group").
		Interval(int64(interval.Seconds())). // type is confusing
		Rules([]cog.Builder[alerting.Rule]{
			alerting.NewRuleBuilder().
				Title("alert rule name").
				For("5m").
				Annotations(map[string]string{
					"summary":     "summary",
					"description": "desc",
					"runbook_url": "https://foo.local",
				}).
				Labels(map[string]string{
					"owner": "team-name",
				}).
				Condition("B").
				Data([]cog.Builder[alerting.AlertQuery]{
					alerting.NewAlertQueryBuilder().
						RefId("A").
						DatasourceUid("cdjvy8rnmh0qof").
						RelativeTimeRange("600", "0"). // TODO
						Model(
							prometheus.NewDataqueryBuilder().
								RefId("A").
								Expr("go_memstats_alloc_bytes_total").
								Instant(true).
								LegendFormat("__auto"),
						),
					alerting.NewAlertQueryBuilder().
						RefId("B").
						DatasourceUid("__expr__").
						RelativeTimeRange("600", "0"), // TODO
					/*
						Model(
							expr.NewTypeReduceBuilder().
								Reducer("last").
								Expression("A"),
						)*/
				}),
		}).
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
