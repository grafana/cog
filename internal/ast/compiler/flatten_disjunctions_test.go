package compiler

import (
	"testing"

	"github.com/grafana/cog/internal/ast"
)

func TestFlattenDisjunctions_WithNestedDisjunctionOfRefs_AsAnObject(t *testing.T) {
	// Prepare test input
	objects := []ast.Object{
		ast.NewObject("test", "ADisjunctionOfRefs", ast.NewDisjunction([]ast.Type{
			ast.NewRef("test", "SomeOrOther"),
			ast.NewRef("test", "LastStruct"),
		})),

		ast.NewObject("test", "SomeOrOther", ast.NewDisjunction([]ast.Type{
			ast.NewRef("test", "SomeStruct"),
			ast.NewRef("test", "OtherStruct"),
		})),

		ast.NewObject("test", "SomeStruct", ast.NewStruct(
			ast.NewStructField("Type", ast.String(ast.Value("some-struct"))),
			ast.NewStructField("FieldFoo", ast.String()),
		)),
		ast.NewObject("test", "OtherStruct", ast.NewStruct(
			ast.NewStructField("FieldBar", ast.NewMap(ast.String(), ast.String())),
			ast.NewStructField("Type", ast.String(ast.Value("other-struct"))),
		)),
		ast.NewObject("test", "LastStruct", ast.NewStruct(
			ast.NewStructField("FieldLast", ast.NewMap(ast.String(), ast.String())),
			ast.NewStructField("Type", ast.String(ast.Value("last-struct"))),
		)),
	}

	// Prepare expected output
	expectedObjects := []ast.Object{
		ast.NewObject("test", "ADisjunctionOfRefs", ast.NewDisjunction([]ast.Type{
			ast.NewRef("test", "SomeStruct"),
			ast.NewRef("test", "OtherStruct"),
			ast.NewRef("test", "LastStruct"),
		})),

		objects[1],
		objects[2],
		objects[3],
		objects[4],
	}

	// Call the compiler pass
	runPassOnObjects(t, &FlattenDisjunctions{}, objects, expectedObjects)
}
