package terraform

import (
	"testing"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/languages"
	"github.com/grafana/cog/internal/testutils"
	"github.com/stretchr/testify/require"
)

func TestRawTypes_Generate(t *testing.T) {
	test := testutils.GoldenFilesTestSuite[ast.Schema]{
		TestDataRoot: "../../../testdata/jennies/rawtypes",
		Name:         "TerraformRawTypes",
		Skip: map[string]string{
			"enums": "Enums aren't supported for types.",
		},
	}

	cfg := Config{}

	jenny := RawTypes{config: cfg}
	compilerPasses := New(cfg).CompilerPasses()

	test.Run(t, func(tc *testutils.Test[ast.Schema]) {
		req := require.New(tc)

		// We run the compiler passes defined for Terraform since without them, we
		// might not be able to translate some of the IR's semantics into Java.
		// Example: disjunctions.
		schema := tc.UnmarshalJSONInput(testutils.RawTypesIRInputFile)
		processedAsts, err := compilerPasses.Process(ast.Schemas{&schema})
		req.NoError(err)

		req.Len(processedAsts, 1, "we somehow got more ast.Schema than we put in")

		files, err := jenny.Generate(languages.Context{
			Schemas: processedAsts,
		})
		req.NoError(err)

		tc.WriteFiles(files)
	})
}
