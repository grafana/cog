package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/grafana/cog/generated/dashboard"
	"github.com/grafana/cog/generated/role"
)

func toPtr[T any](input T) *T {
	return &input
}

func main() {
	someRole := &role.Role{
		Name:        "Role-name",
		DisplayName: toPtr("display name"),
		GroupName:   toPtr("group name"),
		Description: toPtr("description"),
		Hidden:      true,
	}

	converted := role.RoleConverter(someRole)

	fmt.Println(converted)

	dashboardJSON, err := os.ReadFile("/home/kevin/sandbox/work/cog/examples/converter/dashboard.json")
	if err != nil {
		panic(err)
	}

	dash := &dashboard.Dashboard{}
	if err := json.Unmarshal(dashboardJSON, dash); err != nil {
		panic(err)
	}

	convertedDash := dashboard.DashboardConverter(dash)
	fmt.Println(convertedDash)
}

func foo() {
	dashboard.NewDashboardBuilder("[TEST] blocky").
		Uid("test-dashboard-blocky").
		Title("[TEST] blocky").
		Timezone("browser").
		Tooltip(1).
		Time("now-30m", "now").
		Timepicker(dashboard.NewTimePickerBuilder().
			Hidden(false).
			RefreshIntervals([]string{"5s", "10s", "30s", "1m", "5m", "15m", "30m", "1h", "2h", "1d"}).
			TimeOptions([]string{"5m", "15m", "1h", "6h", "12h", "24h", "2d", "7d", "30d"})).
		FiscalYearStartMonth(0x0).
		Refresh("30s").
		Annotations(dashboard.NewAnnotationContainerBuilder())

}
