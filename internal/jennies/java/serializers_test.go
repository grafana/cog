package java

import (
	"testing"

	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/logs"
	"github.com/grafana/cog/internal/testutils"
	"github.com/grafana/cog/pkg/ir"
	"github.com/grafana/cog/pkg/languages"
	"github.com/stretchr/testify/require"
)

func TestSerializers_Generate(t *testing.T) {
	test := testutils.GoldenFilesTestSuite[ir.Schema]{
		TestDataRoot: "../../../testdata/jennies/serializers",
		Name:         "JavaRawTypes",
	}

	cfg := Config{}

	jenny := Serializers{
		config: cfg,
		tmpl:   initTemplates(cfg, common.NewAPIReferenceCollector()),
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
