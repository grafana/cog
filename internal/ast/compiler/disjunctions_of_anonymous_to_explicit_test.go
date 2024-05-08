package compiler

import (
	"testing"

	"github.com/grafana/cog/internal/ast"
)

func TestDisjunctionOfAnonymousStructsToExplicit(t *testing.T) {
	// Prepare test input
	objects := []ast.Object{
		ast.NewObject("test", "disjunctionOfThings", ast.NewDisjunction([]ast.Type{
			ast.NewRef("test", "someStruct"),
			ast.NewStruct(
				ast.NewStructField("Type", ast.String(ast.Value("anonymous-struct"))),
				ast.NewStructField("FieldFoo", ast.String()),
			),
		})),

		ast.NewObject("test", "someStruct", ast.NewStruct(
			ast.NewStructField("Type", ast.String(ast.Value("some-struct"))),
			ast.NewStructField("FieldFoo", ast.String()),
		)),
	}

	// Prepare expected output
	expectedObjects := []ast.Object{
		ast.NewObject("test", "disjunctionOfThings", ast.NewDisjunction([]ast.Type{
			ast.NewRef("test", "someStruct"),
			ast.NewRef("test", "TypeAnonymousStruct"),
		})),

		objects[1],

		ast.NewObject("test", "TypeAnonymousStruct", ast.NewStruct(
			ast.NewStructField("Type", ast.String(ast.Value("anonymous-struct"))),
			ast.NewStructField("FieldFoo", ast.String()),
		)),
	}

	// Call the compiler pass
	runPassOnObjects(t, &DisjunctionOfAnonymousStructsToExplicit{}, objects, expectedObjects)
}
