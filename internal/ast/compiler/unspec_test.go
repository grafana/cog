package compiler

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/grafana/cog/internal/ast"
	"github.com/stretchr/testify/require"
)

func TestUnspec(t *testing.T) {
	// Prepare test input
	schemas := ast.Schemas{
		&ast.Schema{
			Package: "without_spec",
			Objects: []ast.Object{
				ast.NewObject("without_spec", "NotAStruct", ast.String()),

				ast.NewObject("without_spec", "AStruct", ast.NewStruct(
					ast.NewStructField("AString", ast.String()),
				)),
			},
		},

		&ast.Schema{
			Package: "with_spec_no_meta_id",
			Objects: []ast.Object{
				ast.NewObject("with_spec_no_meta_id", "Metadata", ast.NewStruct(
					ast.NewStructField("SomeMeta", ast.String()),
				)),
				ast.NewObject("with_spec_no_meta_id", "Spec", ast.NewStruct(
					ast.NewStructField("title", ast.String()),
				)),
			},
		},

		&ast.Schema{
			Package: "with_spec_and_meta_id",
			Metadata: ast.SchemaMeta{
				Identifier: "Dashboard",
			},
			Objects: []ast.Object{
				ast.NewObject("with_spec_and_meta_id", "Metadata", ast.NewStruct(
					ast.NewStructField("SomeMeta", ast.String()),
				)),
				ast.NewObject("with_spec_and_meta_id", "Spec", ast.NewStruct(
					ast.NewStructField("title", ast.String()),
				)),
			},
		},
	}

	// Prepare expected output
	expected := ast.Schemas{
		// Unchanged
		schemas[0],

		// No identifier defined in schema metadata: the package is used as name instead of "Spec"
		&ast.Schema{
			Package: "with_spec_no_meta_id",
			Objects: []ast.Object{
				ast.NewObject("with_spec_no_meta_id", "with_spec_no_meta_id", ast.NewStruct(
					ast.NewStructField("title", ast.String()),
				)),
			},
		},

		// Identifier defined in the schema metadata: it's used as object name instead of "Spec"
		&ast.Schema{
			Package: "with_spec_and_meta_id",
			Metadata: ast.SchemaMeta{
				Identifier: "Dashboard",
			},
			Objects: []ast.Object{
				ast.NewObject("with_spec_and_meta_id", "Dashboard", ast.NewStruct(
					ast.NewStructField("title", ast.String()),
				)),
			},
		},
	}

	// Run the compiler pass
	runUnspecPass(t, schemas, expected)
}

func runUnspecPass(t *testing.T, input ast.Schemas, expectedOutput ast.Schemas) {
	t.Helper()

	req := require.New(t)

	compilerPass := &Unspec{}
	processedFiles, err := compilerPass.Process(input)
	req.NoError(err)
	req.Len(processedFiles, len(input))
	for i := range input {
		req.Empty(cmp.Diff(expectedOutput[i], processedFiles[i]))
	}
}
