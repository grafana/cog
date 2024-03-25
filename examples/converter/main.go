package main

import (
	"fmt"

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
}
