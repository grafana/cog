package rust

import (
	"testing"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/languages"
	"github.com/grafana/cog/internal/testutils"
	"github.com/stretchr/testify/require"
)

// TestRawTypes_Generate runs the golden-file harness for the Rust RawTypes
// jenny. Phases 3a-3f cover scalars, structs, enums, arrays/maps, cross-package
// refs, untagged disjunctions, constant references, constraints, time hints,
// nullable fields and cross-object struct-literal defaults. Only intersections
// (matching the Python target, which also omits them) and the variant/dashboard
// fixtures (composable slots, Phase 6) remain skipped.
func TestRawTypes_Generate(t *testing.T) {
	test := testutils.GoldenFilesTestSuite[ast.Schema]{
		TestDataRoot: "../../../testdata/jennies/rawtypes",
		Name:         "RustRawTypes",
		Skip: map[string]string{
			"dashboard":     "Phase 6: composable slots not implemented",
			"intersections": "Intersections not implemented, matching the python target",
		},
	}

	config := Config{}
	jenny := RawTypes{
		config:          config,
		apiRefCollector: common.NewAPIReferenceCollector(),
	}
	compilerPasses := New(config).CompilerPasses()

	test.Run(t, func(tc *testutils.Test[ast.Schema]) {
		req := require.New(tc)

		schema := tc.UnmarshalJSONInput(testutils.RawTypesIRInputFile)
		processedAsts, err := compilerPasses.Process(ast.Schemas{&schema})
		req.NoError(err)
		req.Len(processedAsts, 1, "we somehow got more ast.Schema than we put in")

		files, err := jenny.Generate(languages.Context{Schemas: processedAsts})
		req.NoError(err)

		// The Builder/RawTypes jennies emit the simple single-line form and let the
		// FormatRustFiles postprocessor (rustfmt) own all wrapping. The golden-file
		// harness does not run postprocessors, so run rustfmt here too, otherwise
		// the goldens (which are rustfmt-formatted) would never match.
		files = formatRustGoldenFiles(tc, files)

		tc.WriteFiles(files)
	})
}
