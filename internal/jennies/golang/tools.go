package golang

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
	return input == "string" ||
		input == "uint8" ||
		input == "uint16" ||
		input == "uint32" ||
		input == "uint64" ||
		input == "int8" ||
		input == "int16" ||
		input == "int32" ||
		input == "int64" ||
		input == "float32" ||
		input == "float64" ||
		input == "complex64" ||
		input == "complex128" ||
		input == "byte" ||
		input == "rune" ||
		input == "uint" ||
		input == "int" ||
		input == "uintptr" ||
		input == "bool" ||
		// see: https://go.dev/ref/spec#Keywords
		input == "break" ||
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
