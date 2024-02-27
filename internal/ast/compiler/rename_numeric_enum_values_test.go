package compiler

import (
	"testing"

	"github.com/grafana/cog/internal/ast"
)

func TestRenameNumericEnumValues(t *testing.T) {
	// Prepare test input
	objects := []ast.Object{
		ast.NewObject("pkg", "NotAnEnumStruct", ast.String()),

		ast.NewObject("pkg", "AnEnumWithNumericValues", ast.NewEnum([]ast.EnumValue{
			{Type: ast.NewScalar(ast.KindInt64), Name: "-1", Value: -1},
			{Type: ast.NewScalar(ast.KindInt64), Name: "1", Value: 1},
			{Type: ast.NewScalar(ast.KindInt64), Name: "2", Value: 2},
		})),

		ast.NewObject("pkg", "AnEnumWithNoNumericValues", ast.NewEnum([]ast.EnumValue{
			{Type: ast.String(), Name: "Hide", Value: "hide"},
			{Type: ast.String(), Name: "DontHide", Value: "dont_hide"},
		})),
	}

	// Prepare expected output
	expected := []ast.Object{
		ast.NewObject("pkg", "NotAnEnumStruct", ast.String()),

		ast.NewObject("pkg", "AnEnumWithNumericValues", ast.NewEnum([]ast.EnumValue{
			{Type: ast.NewScalar(ast.KindInt64), Name: "Negative1", Value: -1},
			{Type: ast.NewScalar(ast.KindInt64), Name: "N1", Value: 1},
			{Type: ast.NewScalar(ast.KindInt64), Name: "N2", Value: 2},
		}), "RenameNumericEnumValues"),

		ast.NewObject("pkg", "AnEnumWithNoNumericValues", ast.NewEnum([]ast.EnumValue{
			{Type: ast.String(), Name: "Hide", Value: "hide"},
			{Type: ast.String(), Name: "DontHide", Value: "dont_hide"},
		})),
	}

	// Run the compiler pass
	runPassOnObjects(t, &RenameNumericEnumValues{}, objects, expected)
}
