package java

import (
	"testing"

	"github.com/grafana/cog/internal/testutils"
	"github.com/grafana/cog/pkg/apiref"
	"github.com/grafana/cog/pkg/ir"
	"github.com/grafana/cog/pkg/languages"
	"github.com/grafana/cog/pkg/logs"
	"github.com/stretchr/testify/require"
)

func TestDeserializers_Generate(t *testing.T) {
	test := testutils.GoldenFilesTestSuite[ir.Schema]{
		TestDataRoot: "../../../testdata/jennies/deserializers",
		Name:         "JavaRawTypes",
	}

	cfg := Config{}

	jenny := Deserializers{
		config: cfg,
		tmpl:   initTemplates(cfg, apiref.NewAPIReferenceCollector()),
	}
	transforms := New(logs.NoopLogger(), cfg).Transform

	test.Run(t, func(tc *testutils.Test[ir.Schema]) {
		req := require.New(tc)

		// We run the compiler passes defined fo Java since without them, we
		// might not be able to translate some of the IR's semantics into Java.
		// Example: disjunctions.
		schema := tc.UnmarshalJSONInput(testutils.RawTypesIRInputFile)
		processedAsts, err := transforms(ir.Schemas{&schema})
		req.NoError(err)

		req.Len(processedAsts, 1, "we somehow got more ast.Schema than we put in")

		files, err := jenny.Generate(languages.Context{
			Schemas: processedAsts,
		})
		req.NoError(err)

		tc.WriteFiles(files)
	})
}
