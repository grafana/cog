package compiler

import (
	"github.com/grafana/cog/internal/ast"
	"testing"
)

func TestExtractK8ResourceNames(t *testing.T) {
	objects := []ast.Object{
		ast.NewObject("test", "my.super.large.name.Resource",
			ast.NewStruct(
				ast.NewStructField("aRef", ast.NewRef("test", "my.super.large.name.ARef")),
				ast.NewStructField("aConstantRef", ast.NewConstantReferenceType("test", "other.name.with.ugly.name.AConstantRef", "a")),
				ast.NewStructField("aDisjunction", ast.NewDisjunction([]ast.Type{
					ast.NewRef("test", "my.super.large.name.ARef"),
					ast.NewRef("test", "other.name.with.ugly.name.AConstantRef"),
				})),
			),
		),
		ast.NewObject("test", "my.super.large.name.ARef", ast.NewStruct(
			ast.NewStructField("aString", ast.String()),
		)),
		ast.NewObject("test", "other.name.with.ugly.name.AConstantRef", ast.NewEnum([]ast.EnumValue{
			{Type: ast.String(), Name: "A", Value: "a"},
			{Type: ast.String(), Name: "B", Value: "b"},
		},
		)),
	}

	expected := []ast.Object{
		ast.NewObject("test", "Resource",
			ast.NewStruct(
				ast.NewStructField("aRef", ast.NewRef("test", "ARef")),
				ast.NewStructField("aConstantRef", ast.NewConstantReferenceType("test", "AConstantRef", "a")),
				ast.NewStructField("aDisjunction", ast.NewDisjunction([]ast.Type{
					ast.NewRef("test", "ARef"),
					ast.NewRef("test", "AConstantRef"),
				})),
			),
		),
		ast.NewObject("test", "ARef", ast.NewStruct(
			ast.NewStructField("aString", ast.String()),
		)),
		ast.NewObject("test", "AConstantRef", ast.NewEnum([]ast.EnumValue{
			{Type: ast.String(), Name: "A", Value: "a"},
			{Type: ast.String(), Name: "B", Value: "b"},
		},
		)),
	}

	runPassOnObjects(t, &ExtractK8ResourceNames{}, objects, expected)
}
