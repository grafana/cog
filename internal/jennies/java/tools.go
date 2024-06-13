package java

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/grafana/cog/internal/ast"
)

func formatPackageName(pkg string) string {
	rgx := regexp.MustCompile("[^a-zA-Z0-9_]+")

	return strings.ToLower(rgx.ReplaceAllString(pkg, ""))
}

func formatScalar(val any) any {
	newVal := fmt.Sprintf("%#v", val)
	if len(strings.Split(newVal, ".")) > 1 {
		return val
	}
	return newVal
}

func formatType(t ast.ScalarKind, val any) string {
	// When the default is 0, is detected as integer even if it's a float.
	parseFloatVal := func(val any) any {
		if v, ok := val.(int64); ok {
			val = float64(v)
		} else {
			val = val.(float64)
		}
		return val
	}

	switch t {
	case ast.KindInt64:
		return fmt.Sprintf("%dL", val.(int64))
	case ast.KindUint64:
		return fmt.Sprintf("%dL", val.(int64))
	case ast.KindInt32:
		return fmt.Sprintf("%d", val.(int64))
	case ast.KindFloat32:
		return fmt.Sprintf("%.1ff", parseFloatVal(val))
	case ast.KindFloat64:
		return fmt.Sprintf("%.1f", parseFloatVal(val))
	}

	return fmt.Sprintf("%#v", val)
}

// TODO: Need to say to the serializer the correct name.
func escapeVarName(varName string) string {
	if isReservedJavaKeyword(varName) {
		return varName + "Arg"
	}

	return varName
}

func lastPathIdentifier(fieldPath ast.Path) string {
	lastPath := make([]string, 0)
	shouldAddPath := false
	for _, path := range fieldPath {
		if shouldAddPath {
			lastPath = append(lastPath, path.Identifier)
		}
		if path.Type.IsAny() {
			shouldAddPath = true
		}
	}
	return strings.Join(lastPath, ".")
}

// nolint: gocyclo
func isReservedJavaKeyword(input string) bool {
	// see https://docs.oracle.com/javase/tutorial/java/nutsandbolts/_keywords.html
	switch input {
	case "static", "abstract", "enum", "class", "if", "else", "switch", "final", "public", "private", "protected", "package", "continue", "new", "for", "assert",
		"do", "default", "goto", "synchronized", "boolean", "double", "int", "short", "char", "float", "long", "byte", "break", "throw", "throws", "this",
		"implements", "transient", "return", "catch", "extends", "case", "try", "void", "volatile", "super", "native", "finally", "instanceof", "import", "while":
		return true
	}
	return false
}
