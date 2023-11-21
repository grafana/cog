package compiler

import (
	"testing"

	"github.com/grafana/cog/internal/ast"
)

func TestPrefixEnumValues(t *testing.T) {
	// Prepare test input
	objects := []ast.Object{
		ast.NewObject("pkg", "VariableRefresh", ast.NewEnum([]ast.EnumValue{
			{Name: "Never", Value: "never", Type: ast.String()},
			{Name: "Always", Value: "always", Type: ast.String()},
		})),

		ast.NewObject("pkg", "SomeType", ast.String()),
	}

	// Prepare expected output
	expected := []ast.Object{
		ast.NewObject("pkg", "VariableRefresh", ast.NewEnum([]ast.EnumValue{
			{Name: "VariableRefreshNever", Value: "never", Type: ast.String()},
			{Name: "VariableRefreshAlways", Value: "always", Type: ast.String()},
		})),

		ast.NewObject("pkg", "SomeType", ast.String()),
	}

	// Run the compiler pass
	runPassOnObjects(t, &PrefixEnumValues{}, objects, expected)
}
