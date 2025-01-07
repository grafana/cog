package compiler

import (
	"testing"

	"github.com/grafana/cog/internal/ast"
)

func TestTrimEnumValues(t *testing.T) {
	objects := []ast.Object{
		ast.NewObject("enum_with_leading_and_trailing_spaces", "MyEnum", ast.NewEnum([]ast.EnumValue{
			{Type: ast.String(), Name: "Leading", Value: " Leading"},
			{Type: ast.String(), Name: "Trailing", Value: "Trailing "},
			{Type: ast.String(), Name: "Both", Value: " Both "},
			{Type: ast.String(), Name: "SpacesInMiddle", Value: "Spaces in middle"},
		})),
	}

	expected := []ast.Object{
		ast.NewObject("enum_with_leading_and_trailing_spaces", "MyEnum", ast.NewEnum([]ast.EnumValue{
			{Type: ast.String(), Name: "Leading", Value: "Leading"},
			{Type: ast.String(), Name: "Trailing", Value: "Trailing"},
			{Type: ast.String(), Name: "Both", Value: "Both"},
			{Type: ast.String(), Name: "SpacesInMiddle", Value: "Spaces in middle"},
		})),
	}

	runPassOnObjects(t, &TrimEnumValues{}, objects, expected)
}
