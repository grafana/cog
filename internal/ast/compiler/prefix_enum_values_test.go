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
		}), "PrefixEnumValues"),

		ast.NewObject("pkg", "SomeType", ast.String()),
	}

	// Run the compiler pass
	runPassOnObjects(t, &PrefixEnumValues{}, objects, expected)
}

func TestPrefixEnumValuesWithNegativeIntegerName(t *testing.T) {
	// Prepare test input
	objects := []ast.Object{
		ast.NewObject("pkg", "BarAlignment", ast.NewEnum([]ast.EnumValue{
			{Name: "1", Value: 1, Type: ast.NewScalar(ast.KindInt64)},
			{Name: "-1", Value: -1, Type: ast.NewScalar(ast.KindInt64)},
		})),
	}

	// Prepare expected output
	expected := []ast.Object{
		ast.NewObject("pkg", "BarAlignment", ast.NewEnum([]ast.EnumValue{
			{Name: "BarAlignment1", Value: 1, Type: ast.NewScalar(ast.KindInt64)},
			{Name: "BarAlignmentNegative1", Value: -1, Type: ast.NewScalar(ast.KindInt64)},
		}), "PrefixEnumValues"),
	}

	// Run the compiler pass
	runPassOnObjects(t, &PrefixEnumValues{}, objects, expected)
}

func TestPrefixEnumValuesWithEmptyStringMember(t *testing.T) {
	// Prepare test input
	objects := []ast.Object{
		ast.NewObject("pkg", "BarAlignment", ast.NewEnum([]ast.EnumValue{
			{Name: "", Value: "", Type: ast.String()},
			{Name: "foo", Value: "foo", Type: ast.String()},
		})),
	}

	// Prepare expected output
	expected := []ast.Object{
		ast.NewObject("pkg", "BarAlignment", ast.NewEnum([]ast.EnumValue{
			{Name: "BarAlignmentNone", Value: "", Type: ast.String()},
			{Name: "BarAlignmentFoo", Value: "foo", Type: ast.String()},
		}), "PrefixEnumValues"),
	}

	// Run the compiler pass
	runPassOnObjects(t, &PrefixEnumValues{}, objects, expected)
}
