package golang

import (
	"testing"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/languages"
	"github.com/grafana/cog/internal/testutils"
	"github.com/stretchr/testify/require"
)

func TestRawTypes_Generate(t *testing.T) {
	test := testutils.GoldenFilesTestSuite[ast.Schema]{
		TestDataRoot: "../../../testdata/jennies/rawtypes",
		Name:         "GoRawTypes",
		Skip: map[string]string{
			"dashboard": "the dashboard test schema includes a composable slot, which rely on external input to be properly supported",
		},
	}

	config := Config{
		PackageRoot:                "github.com/grafana/cog/generated",
		GenerateEqual:              true,
		GenerateJSONMarshaller:     true,
		GenerateStrictUnmarshaller: true,
		GenerateValidate:           true,
	}
	jenny := RawTypes{
		config:          config,
		tmpl:            initTemplates(config, common.NewAPIReferenceCollector()),
		apiRefCollector: common.NewAPIReferenceCollector(),
	}
	compilerPasses := New(config).CompilerPasses()

	test.Run(t, func(tc *testutils.Test[ast.Schema]) {
		req := require.New(tc)

		// We run the compiler passes defined fo Go since without them, we
		// might not be able to translate some of the IR's semantics into Go.
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
