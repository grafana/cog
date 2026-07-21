package csharp

import (
	"testing"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/languages"
	"github.com/grafana/cog/internal/testutils"
	"github.com/stretchr/testify/require"
)

// TestConverters_Generate covers the JsonConverter<T> output emitted
// by the Converters jenny for the three supported disjunction shapes
// (scalars, discriminated refs, scalars+refs).
//
// Reuses the IR fixtures under testdata/jennies/serializers/, which
// are raw disjunctions that the compiler passes turn into structs.
// Goldens land under testdata/jennies/serializers/<case>/CSharpJsonConverters/.
func TestConverters_Generate(t *testing.T) {
	test := testutils.GoldenFilesTestSuite[ast.Schema]{
		TestDataRoot: "../../../testdata/jennies/serializers",
		Name:         "CSharpJsonConverters",
	}

	cfg := Config{GenerateJSONConverters: true}
	cfg.InterpolateParameters(func(input string) string { return input })

	tmpl := initTemplates(cfg, common.NewAPIReferenceCollector())
	jenny := &Converters{config: cfg, tmpl: tmpl}
	compilerPasses := New(cfg).CompilerPasses()

	test.Run(t, func(tc *testutils.Test[ast.Schema]) {
		req := require.New(tc)

		schema := tc.UnmarshalJSONInput(testutils.RawTypesIRInputFile)
		processedAsts, err := compilerPasses.Process(ast.Schemas{&schema})
		req.NoError(err)
		req.Len(processedAsts, 1)

		files, err := jenny.Generate(languages.Context{Schemas: processedAsts})
		req.NoError(err)

		tc.WriteFiles(files)
	})
}
