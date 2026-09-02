package csharp

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
		Name:         "CSharpRawTypes",
	}

	cfg := Config{}
	cfg.InterpolateParameters(func(input string) string { return input })

	jenny := RawTypes{config: cfg, tmpl: initTemplates(cfg, common.NewAPIReferenceCollector())}
	compilerPasses := New(cfg).CompilerPasses()

	test.Run(t, func(tc *testutils.Test[ast.Schema]) {
		req := require.New(tc)

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
