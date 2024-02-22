package compiler

import (
	"testing"

	"github.com/grafana/cog/internal/ast"
)

func TestDisjunctionWithNullToOptional_WithDisjunctionOfTypeAndNull_AsAnObject(t *testing.T) {
	// Prepare test input
	objects := []ast.Object{
		ast.NewObject("test", "ScalarWithNull", ast.NewDisjunction([]ast.Type{
			ast.String(),
			ast.Null(),
		})),
		ast.NewObject("test", "RefWithNull", ast.NewDisjunction([]ast.Type{
			ast.NewRef("test", "SomeType"),
			ast.Null(),
		})),
	}

	expectedObjects := []ast.Object{
		ast.NewObject("test", "ScalarWithNull", ast.String(ast.Nullable(), ast.Trail("DisjunctionWithNullToOptional[String|null → String?]"))),
		ast.NewObject("test", "RefWithNull", ast.NewRef("test", "SomeType", ast.Nullable(), ast.Trail("DisjunctionWithNullToOptional[SomeType|null → SomeType?]"))),
	}

	// Call the compiler pass
	runPassOnObjects(t, &DisjunctionWithNullToOptional{}, objects, expectedObjects)
}

func TestDisjunctionWithNullToOptional_WithDisjunctionOfTypeAndNull_AsAStructField(t *testing.T) {
	// Prepare test input
	objects := []ast.Object{
		ast.NewObject("test", "StructWithScalarWithNull", ast.NewStruct(
			ast.NewStructField("Field", ast.NewDisjunction([]ast.Type{
				ast.String(),
				ast.Null(),
			})),
		)),
		ast.NewObject("test", "StructWithRefWithNull", ast.NewStruct(
			ast.NewStructField("Field", ast.NewDisjunction([]ast.Type{
				ast.NewRef("test", "SomeType"),
				ast.Null(),
			})),
		)),
	}

	expectedObjects := []ast.Object{
		ast.NewObject("test", "StructWithScalarWithNull", ast.NewStruct(
			ast.NewStructField("Field", ast.String(ast.Nullable(), ast.Trail("DisjunctionWithNullToOptional[String|null → String?]"))),
		)),
		ast.NewObject("test", "StructWithRefWithNull", ast.NewStruct(
			ast.NewStructField("Field", ast.NewRef("test", "SomeType", ast.Nullable(), ast.Trail("DisjunctionWithNullToOptional[SomeType|null → SomeType?]"))),
		)),
	}

	// Call the compiler pass
	runPassOnObjects(t, &DisjunctionWithNullToOptional{}, objects, expectedObjects)
}
