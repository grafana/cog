package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/grafana/cog/generated/cog/plugins"
	"github.com/grafana/cog/generated/dashboard"
)

func main() {
	plugins.RegisterDefaultPlugins()

	dashboardJSON, err := os.ReadFile("/home/kevin/sandbox/work/cog/examples/_go/converter/dashboard.json")
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

}
