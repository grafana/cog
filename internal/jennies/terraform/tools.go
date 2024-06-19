package terraform

import (
	"regexp"
	"strings"

	"github.com/grafana/cog/internal/tools"
)

func formatPackageName(pkg string) string {
	rgx := regexp.MustCompile("[^a-zA-Z0-9_]+")

	return strings.ToLower(rgx.ReplaceAllString(pkg, ""))
}

func formatArgName(name string) string {
	return escapeVarName(tools.LowerCamelCase(name))
}

func escapeVarName(varName string) string {
	if isReservedGoKeyword(varName) {
		return varName + "Arg"
	}

	return varName
}

func isReservedGoKeyword(input string) bool {
	// see: https://go.dev/ref/spec#Keywords
	return input == "break" ||
		input == "case" ||
		input == "chan" ||
		input == "continue" ||
		input == "const" ||
		input == "default" ||
		input == "defer" ||
		input == "else" ||
		input == "error" ||
		input == "fallthrough" ||
		input == "for" ||
		input == "func" ||
		input == "go" ||
		input == "goto" ||
		input == "if" ||
		input == "import" ||
		input == "interface" ||
		input == "map" ||
		input == "package" ||
		input == "range" ||
		input == "return" ||
		input == "select" ||
		input == "struct" ||
		input == "switch" ||
		input == "type" ||
		input == "var"
}
