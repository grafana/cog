package compiler

import (
	"testing"

	"github.com/grafana/cog/internal/ast"
)

func TestUndiscriminatedDisjunctionToAny(t *testing.T) {
	// Prepare test input
	disjunctionTypeNoMapping := ast.NewDisjunction([]ast.Type{
		ast.NewRef("test", "SomeStruct"),
		ast.NewRef("test", "OtherStruct"),
	})
	disjunctionTypeNoMapping.Disjunction.Discriminator = "Type"
	disjunctionTypeMapping := ast.NewDisjunction([]ast.Type{
		ast.NewRef("test", "SomeStruct"),
		ast.NewRef("test", "OtherStruct"),
	})
	disjunctionTypeMapping.Disjunction.Discriminator = "Type"
	disjunctionTypeMapping.Disjunction.DiscriminatorMapping = map[string]string{
		"some-struct":  "SomeStruct",
		"other-struct": "OtherStruct",
	}

	objects := []ast.Object{
		ast.NewObject("test", "ADisjunctionOfRefsNoMapping", disjunctionTypeNoMapping),
		ast.NewObject("test", "ADisjunctionOfRefsMapping", disjunctionTypeMapping),
		ast.NewObject("test", "ADisjunctionOfScalars", ast.NewDisjunction([]ast.Type{
			ast.String(),
			ast.Bool(),
		})),

		ast.NewObject("test", "SomeStruct", ast.NewStruct(
			ast.NewStructField("Type", ast.String()), // Not a concrete scalar
			ast.NewStructField("FieldFoo", ast.String()),
		)),
		ast.NewObject("test", "OtherStruct", ast.NewStruct(
			ast.NewStructField("Type", ast.String(ast.Value("other-struct"))),
			ast.NewStructField("FieldBar", ast.Bool()),
		)),
	}

	disjunctionOfRefNoMapping := objects[0].DeepCopy()
	disjunctionOfRefNoMapping.Type = ast.Any(ast.Trail("UndiscriminatedDisjunctionToAny"))
	expected := []ast.Object{
		disjunctionOfRefNoMapping,
		objects[1],
		objects[2],
		objects[3],
		objects[4],
	}

	runPassOnObjects(t, &UndiscriminatedDisjunctionToAny{}, objects, expected)
}
