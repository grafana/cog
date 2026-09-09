package rust

import (
	"testing"

	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/languages"
	"github.com/grafana/cog/internal/testutils"
	"github.com/stretchr/testify/require"
)

// TestConverter_Generate wires the golden-file harness for the Rust Converter
// jenny. The skip set is the Rust Builder skip set: a converter emits calls to
// a builder, so every fixture without a Rust builder has no converter either.
func TestConverter_Generate(t *testing.T) {
	test := testutils.GoldenFilesTestSuite[languages.Context]{
		TestDataRoot: "../../../testdata/jennies/builders",
		Name:         "RustConverter",
		Skip: map[string]string{
			// See builder_test.go for the full rationale on each of these: the Rust
			// Builder jenny skips them, so there is no builder for a converter to
			// target. The first three are also skipped by the Go converter suite.
			"anonymous_struct":                  "no Rust builder: inline anonymous struct has no idiomatic Rust argument type (Go converter also skips)",
			"builder_delegation_in_disjunction": "no Rust builder: inline-disjunction delegation has no idiomatic Rust mapping (Go converter also skips)",
			"dashboard_panel":                   "Java-generics-specific builder shape (Go converter also skips)",
			"struct_fields_as_args_assignment":  "no Rust builder: deep assignment into an un-hoisted inline anonymous struct is not type-safe Rust",
		},
	}

	language := New(Config{})
	jenny := Converter{config: language.config, apiRefCollector: common.NewAPIReferenceCollector()}

	test.Run(t, func(tc *testutils.Test[languages.Context]) {
		var err error
		req := require.New(tc)

		context := tc.UnmarshalJSONInput(testutils.BuildersContextInputFile)
		context, err = languages.GenerateBuilderNilChecks(language, context)
		req.NoError(err)

		files, err := jenny.Generate(context)
		req.NoError(err)

		// See rawtypes_test.go: the harness does not run postprocessors, so format
		// the generated files with rustfmt here to match the rustfmt-formatted
		// goldens the FormatRustFiles postprocessor produces in real generation.
		files = formatRustGoldenFiles(tc, files)

		tc.WriteFiles(files)
	})
}
