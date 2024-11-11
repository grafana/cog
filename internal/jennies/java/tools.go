package java

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/languages"
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

func cleanString(s string) string {
	if strings.Contains(s, "\n") {
		return strings.Replace(s, "\n", `\n`, -1)
	}

	return s
}

func formatType(t ast.ScalarKind, val any) string {
	// When the default is 0, is detected as integer even if it's a float.
	parseFloatVal := func(val any) any {
		if v, ok := val.(int64); ok {
			return float64(v)
		}
		return val.(float64)
	}

	// Integers could be floats in JSON
	parseIntVal := func(val any) any {
		if v, ok := val.(float64); ok {
			return int64(v)
		}

		return val.(int64)
	}

	switch t {
	case ast.KindInt64:
		return fmt.Sprintf("%dL", parseIntVal(val))
	case ast.KindUint64:
		return fmt.Sprintf("%dL", parseIntVal(val))
	case ast.KindInt32:
		return fmt.Sprintf("%d", parseIntVal(val))
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

func fillAnnotationPattern(input string, value string) string {
	if strings.Contains(input, "%#v") {
		return fmt.Sprintf(input, value)
	}
	return input
}

func containsValue(value string, list []DataqueryUnmarshalling) bool {
	for _, v := range list {
		if v.FieldName == value {
			return true
		}
	}

	return false
}

func getJavaFieldTypeCheck(t ast.Type) string {
	switch t.Kind {
	case ast.KindArray:
		return "isArray()"
	case ast.KindMap:
		return "isObject()"
	case ast.KindScalar:
		switch t.AsScalar().ScalarKind {
		case ast.KindString:
			return "isTextual()"
		case ast.KindBool:
			return "isBoolean()"
		case ast.KindInt8, ast.KindUint8, ast.KindInt32, ast.KindUint32:
			return "isInt()"
		case ast.KindFloat64:
			return "isDouble()"
		default:
			return "isObject()"
		}
	default:
		return "isObject()"
	}
}

func objectNeedsCustomDeserialiser(context languages.Context, obj ast.Object) bool {
	// an object needs a custom unmarshal if:
	// - it is a struct that was generated from a disjunction by the `DisjunctionToType` compiler pass.
	// - it is a struct and one or more of its fields is a KindComposableSlot, or an array of KindComposableSlot

	if !obj.Type.IsStruct() {
		return false
	}

	// is it a struct generated from a disjunction?
	if obj.Type.IsStructGeneratedFromDisjunction() {
		return true
	}

	// is there a KindComposableSlot field somewhere?
	for _, field := range obj.Type.AsStruct().Fields {
		if _, ok := context.ResolveToComposableSlot(field.Type); ok {
			return true
		}
	}

	return false
}
