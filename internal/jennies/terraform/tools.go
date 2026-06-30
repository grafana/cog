package terraform

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/tools"
)

func formatPackageName(pkg string) string {
	splitPath := strings.Split(pkg, "/")
	if len(splitPath) > 1 {
		pkg = splitPath[len(splitPath)-1]
	}
	rgx := regexp.MustCompile("[^a-zA-Z0-9_]+")

	return strings.ToLower(rgx.ReplaceAllString(pkg, ""))
}

func formatObjectName(name string) string {
	return tools.UpperCamelCase(name)
}

func formatModelName(ref ast.RefType) string {
	return formatObjectName(ref.ReferredType) + "Model"
}

func formatModelFieldName(name string) string {
	return tools.UpperCamelCase(name)
}

func formatTypeName(ref ast.RefType) string {
	return formatObjectName(ref.ReferredType) + "Type"
}

func formatTfSDKAttrName(input string) string {
	return tools.SnakeCase(input)
}

func formatScalar(val any) string {
	if val == nil {
		return "nil"
	}

	if list, ok := val.([]any); ok {
		items := make([]string, 0, len(list))

		for _, item := range list {
			items = append(items, formatScalar(item))
		}

		// FIXME: this is wrong, we can't just assume a list of strings.
		return fmt.Sprintf("[]string{%s}", strings.Join(items, ", "))
	}

	return fmt.Sprintf("%#v", val)
}

func formatEnumValuesAsConstraints(enumValues []ast.EnumValue) []ast.TypeConstraint {
	values := make([]any, len(enumValues))
	for i, v := range enumValues {
		values[i] = v.Value
	}

	return []ast.TypeConstraint{
		{
			Op:   ast.EqualOp,
			Args: values,
		},
	}
}

func formatScalarAsModel(scalar ast.ScalarType) string {
	switch scalar.ScalarKind {
	case ast.KindString, ast.KindBytes, ast.KindNull:
		return "types.String"
	case ast.KindBool:
		return "types.Bool"
	case ast.KindInt32, ast.KindUint32:
		return "types.Int32"
	case ast.KindInt64, ast.KindUint64:
		return "types.Int64"
	case ast.KindFloat32:
		return "types.Float32"
	case ast.KindFloat64:
		return "types.Float64"
	case ast.KindAny:
		return "types.String" // `any` should be represented as a string holding a JSON payload
	case ast.KindInt8, ast.KindUint8, ast.KindInt16, ast.KindUint16:
		return "types.Number" // types.Number can be converted into any numeric type https://developer.hashicorp.com/terraform/plugin/framework/handling-data/types/number#setting-values
	default:
		return fmt.Sprintf("unsupported scalar kind '%s'", scalar.ScalarKind)
	}
}
