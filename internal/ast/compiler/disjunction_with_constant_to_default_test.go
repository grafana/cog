package compiler

import (
	"testing"

	"github.com/grafana/cog/internal/ast"
)

func TestDisjunctionWithConstantToDefault(t *testing.T) {
	// Prepare test input
	objects := []ast.Object{
		ast.NewObject("test", "DisjunctionWithScalarAndConst", ast.NewDisjunction([]ast.Type{
			ast.String(),
			ast.String(ast.Value("foo")),
		})),

		ast.NewObject("test", "DisjunctionWithDifferentKinds", ast.NewDisjunction([]ast.Type{
			ast.String(),
			ast.Bool(ast.Value(false)),
		})),

		ast.NewObject("test", "DisjunctionWithTwoScalarsAndConst", ast.NewDisjunction([]ast.Type{
			ast.String(),
			ast.Bool(),
			ast.String(ast.Value("foo")),
		})),

		ast.NewObject("test", "DisjunctionWithTwoScalars", ast.NewDisjunction([]ast.Type{
			ast.String(),
			ast.Bool(),
		})),

		ast.NewObject("test", "DisjunctionWithTwoScalarConsts", ast.NewDisjunction([]ast.Type{
			ast.String(ast.Value("bar")),
			ast.String(ast.Value("foo")),
		})),

		ast.NewObject("test", "ScalarObject", ast.String()),
	}

	// Prepare expected output
	expectedObjects := []ast.Object{
		ast.NewObject("test", "DisjunctionWithScalarAndConst", ast.String(ast.Default("foo"), ast.Trail("DisjunctionWithConstantToDefault"))),
		objects[1],
		objects[2],
		objects[3],
		objects[4],
		objects[5],
	}

	// Call the compiler pass
	runPassOnObjects(t, &DisjunctionWithConstantToDefault{}, objects, expectedObjects)
}
