package compiler

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/grafana/cog/internal/ast"
	"github.com/stretchr/testify/require"
)

func TestNotRequiredFieldAsNullableType(t *testing.T) {
	// Prepare test input
	objects := []ast.Object{
		ast.NewObject("NotAStruct", ast.String()),

		ast.NewObject("AStruct", ast.NewStruct(
			ast.NewStructField("RequiredString", ast.String(), ast.Required()),
			ast.NewStructField("RequiredNullableString", ast.String(ast.Nullable()), ast.Required()),
			ast.NewStructField("NotRequiredString", ast.String()),

			ast.NewStructField("RequiredRef", ast.NewRef("SomeStruct"), ast.Required()),
			ast.NewStructField("RequiredNullableRef", ast.NewRef("SomeStruct", ast.Nullable()), ast.Required()),
			ast.NewStructField("NotRequiredRef", ast.NewRef("SomeStruct")),

			ast.NewStructField("NotRequiredArray", ast.NewArray(ast.String())),
			ast.NewStructField("RequiredArray", ast.NewArray(ast.String()), ast.Required()),

			ast.NewStructField("NotRequiredMap", ast.NewMap(
				ast.String(),
				ast.Bool(),
			)),
			ast.NewStructField("RequiredMap", ast.NewMap(
				ast.String(),
				ast.Bool(),
			), ast.Required()),
		)),
	}

	// Prepare expected output
	expected := []ast.Object{
		ast.NewObject("NotAStruct", ast.String()),

		ast.NewObject("AStruct", ast.NewStruct(
			ast.NewStructField("RequiredString", ast.String(), ast.Required()),
			ast.NewStructField("RequiredNullableString", ast.String(ast.Nullable()), ast.Required()),
			ast.NewStructField("NotRequiredString", ast.String(ast.Nullable())), // should become nullable

			ast.NewStructField("RequiredRef", ast.NewRef("SomeStruct"), ast.Required()),
			ast.NewStructField("RequiredNullableRef", ast.NewRef("SomeStruct", ast.Nullable()), ast.Required()),
			ast.NewStructField("NotRequiredRef", ast.NewRef("SomeStruct", ast.Nullable())), // should become nullable

			ast.NewStructField("NotRequiredArray", ast.NewArray(ast.String(), ast.Nullable())), // should become nullable
			ast.NewStructField("RequiredArray", ast.NewArray(ast.String()), ast.Required()),

			ast.NewStructField("NotRequiredMap", ast.NewMap(
				ast.String(),
				ast.Bool(),
				ast.Nullable(), // should become nullable
			)),
			ast.NewStructField("RequiredMap", ast.NewMap(
				ast.String(),
				ast.Bool(),
			), ast.Required()),
		)),
	}

	// Run the compiler pass
	runNotRequiredAsNullablePass(t, objects, expected)
}

func runNotRequiredAsNullablePass(t *testing.T, input []ast.Object, expectedOutput []ast.Object) {
	t.Helper()

	req := require.New(t)

	compilerPass := &NotRequiredFieldAsNullableType{}
	processedFiles, err := compilerPass.Process([]*ast.File{
		{
			Package:     "test",
			Definitions: input,
		},
	})
	req.NoError(err)
	req.Len(processedFiles, 1)
	req.Empty(cmp.Diff(expectedOutput, processedFiles[0].Definitions))
}
