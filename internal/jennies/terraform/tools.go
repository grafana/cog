package terraform

import (
	"regexp"
	"strings"
)

func formatPackageName(pkg string) string {
	splitPath := strings.Split(pkg, "/")
	if len(splitPath) > 1 {
		pkg = splitPath[len(splitPath)-1]
	}
	rgx := regexp.MustCompile("[^a-zA-Z0-9_]+")

	return strings.ToLower(rgx.ReplaceAllString(pkg, ""))
}
