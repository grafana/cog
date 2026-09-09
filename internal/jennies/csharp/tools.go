package csharp

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/tools"
)

// pascalCase converts a snake_case, kebab-case, or lower-case identifier
// into PascalCase suitable for C# namespace and type names.
func pascalCase(input string) string {
	if input == "" {
		return ""
	}

	parts := strings.FieldsFunc(input, func(r rune) bool {
		return r == '_' || r == '-' || r == '.' || r == ' '
	})

	var b strings.Builder
	for _, part := range parts {
		if part == "" {
			continue
		}
		b.WriteString(strings.ToUpper(part[:1]))
		b.WriteString(part[1:])
	}
	return b.String()
}

// formatObjectName produces a PascalCased C# type name.
func formatObjectName(name string) string {
	return tools.UpperCamelCase(name)
}

// formatFieldName produces a PascalCased C# field/property name.
// (C# convention is PascalCase for public fields and properties.)
func formatFieldName(name string) string {
	return escapeVarName(tools.UpperCamelCase(name))
}

// formatArgName produces a camelCased C# parameter name.
func formatArgName(name string) string {
	return escapeVarName(tools.LowerCamelCase(name))
}

// formatPackageName converts a schema package identifier into the
// PascalCase segment used in a C# namespace.
//
// Example: "library_panel" -> "LibraryPanel".
func formatPackageName(pkg string) string {
	rgx := regexp.MustCompile(`[^a-zA-Z0-9_-]+`)
	cleaned := rgx.ReplaceAllString(pkg, "")
	return pascalCase(cleaned)
}

// escapeVarName appends a suffix to identifiers that collide with C#
// reserved keywords.
func escapeVarName(varName string) string {
	if isReservedCSharpKeyword(varName) {
		return varName + "Arg"
	}
	return varName
}

// isReservedCSharpKeyword reports whether the given identifier collides
// with a C# reserved keyword. Source:
// https://learn.microsoft.com/en-us/dotnet/csharp/language-reference/keywords/
//
//nolint:gocyclo
func isReservedCSharpKeyword(input string) bool {
	switch input {
	case "abstract", "as", "base", "bool", "break", "byte", "case", "catch",
		"char", "checked", "class", "const", "continue", "decimal", "default",
		"delegate", "do", "double", "else", "enum", "event", "explicit",
		"extern", "false", "finally", "fixed", "float", "for", "foreach",
		"goto", "if", "implicit", "in", "int", "interface", "internal", "is",
		"lock", "long", "namespace", "new", "null", "object", "operator",
		"out", "override", "params", "private", "protected", "public",
		"readonly", "ref", "return", "sbyte", "sealed", "short", "sizeof",
		"stackalloc", "static", "string", "struct", "switch", "this", "throw",
		"true", "try", "typeof", "uint", "ulong", "unchecked", "unsafe",
		"ushort", "using", "virtual", "void", "volatile", "while":
		return true
	}
	return false
}

// formatScalarType returns the C# type expression for a scalar AST type.
//
// Strings, booleans and numerics map to the matching BCL primitive type.
// Bytes are rendered as `byte`. `any` becomes `object`.
func formatScalarType(def ast.ScalarType) string {
	switch def.ScalarKind {
	case ast.KindString:
		return "string"
	case ast.KindBytes:
		return "byte"
	case ast.KindInt8:
		return "sbyte"
	case ast.KindUint8:
		return "byte"
	case ast.KindInt16:
		return "short"
	case ast.KindUint16:
		return "ushort"
	case ast.KindInt32:
		return "int"
	case ast.KindUint32:
		return "uint"
	case ast.KindInt64:
		return "long"
	case ast.KindUint64:
		return "ulong"
	case ast.KindFloat32:
		return "float"
	case ast.KindFloat64:
		return "double"
	case ast.KindBool:
		return "bool"
	case ast.KindAny:
		return "object"
	}
	return "object"
}

// formatScalarValue renders a Go value as a C# literal of the given
// scalar kind. It is used for field defaults and constants.
func formatScalarValue(t ast.ScalarKind, val any) string {
	if list, ok := val.([]any); ok {
		items := make([]string, 0, len(list))
		for _, item := range list {
			items = append(items, formatScalarValue(t, item))
		}
		return strings.Join(items, ", ")
	}

	parseFloatVal := func(val any) float64 {
		if v, ok := val.(int64); ok {
			return float64(v)
		}
		if v, ok := val.(float64); ok {
			return v
		}
		return 0
	}

	switch t {
	case ast.KindBool:
		if b, ok := val.(bool); ok && b {
			return "true"
		}
		return "false"
	case ast.KindString:
		s, _ := val.(string)
		return fmt.Sprintf("%q", s)
	case ast.KindInt64:
		return fmt.Sprintf("%dL", tools.AnyToInt64(val))
	case ast.KindUint64:
		return fmt.Sprintf("%dUL", tools.AnyToInt64(val))
	case ast.KindInt8, ast.KindUint8, ast.KindInt16, ast.KindUint16, ast.KindInt32, ast.KindUint32:
		return fmt.Sprintf("%d", tools.AnyToInt64(val))
	case ast.KindFloat32:
		return fmt.Sprintf("%gf", parseFloatVal(val))
	case ast.KindFloat64:
		return fmt.Sprintf("%gd", parseFloatVal(val))
	}

	return fmt.Sprintf("%#v", val)
}
