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

}
