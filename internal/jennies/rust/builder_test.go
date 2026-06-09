package rust

import (
	"testing"

	"github.com/grafana/cog/internal/languages"
	"github.com/grafana/cog/internal/testutils"
	"github.com/stretchr/testify/require"
)

// TestBuilder_Generate wires the golden-file harness for the Rust Builder
// jenny. Phase 4a implements direct constant/argument assignment to single-level
// field paths; fixtures requiring later-chunk features (builder delegation,
// factories, envelopes, variants, array-append, constraints) are skipped.
func TestBuilder_Generate(t *testing.T) {
	test := testutils.GoldenFilesTestSuite[languages.Context]{
		TestDataRoot: "../../../testdata/jennies/builders",
		Name:         "RustBuilder",
		// Phase 4a only generates goldens for basic_struct, basic_struct_defaults,
		// properties, constructor_argument and constructor_initializations. Every
		// other fixture is skipped because it exercises a builder feature
		// implemented in a later chunk (delegation, factories, envelopes, variants,
		// array-append, constraints, anonymous structs, etc).
		Skip: map[string]string{
			"anonymous_struct":                  "anonymous structs out of scope for Phase 4a",
			"builder_delegation":                "builder delegation out of scope for Phase 4a",
			"builder_delegation_in_disjunction": "builder delegation out of scope for Phase 4a",
			"composable_slot":                   "composable slots out of scope for Phase 4a",
			"dashboard_panel":                   "out of scope for Phase 4a",
			"dataquery_variant_builder":         "variants out of scope for Phase 4a",
			"envelope_assignment":               "envelopes out of scope for Phase 4a",
			"factories":                         "factories out of scope for Phase 4a",
			"foreign_builder":                   "out of scope for Phase 4a",
			"initialization_safeguards":         "nil-check guards out of scope for Phase 4a",
			"map_of_builders":                   "builder delegation out of scope for Phase 4a",
			"map_of_disjunctions":               "out of scope for Phase 4a",
			"package-with-dashes":               "out of scope for Phase 4a",
			"panel_builders":                    "out of scope for Phase 4a",
			"references":                        "out of scope for Phase 4a",
			"struct_fields_as_args_assignment":  "out of scope for Phase 4a",
			"struct_with_defaults":              "builder delegation + anonymous struct args out of scope for Phase 4a",
		},
	}

	language := New(Config{})
	jenny := Builder{config: language.config, apiRefCollector: language.apiRefCollector}

	test.Run(t, func(tc *testutils.Test[languages.Context]) {
		var err error
		req := require.New(tc)

		context := tc.UnmarshalJSONInput(testutils.BuildersContextInputFile)
		context, err = languages.GenerateBuilderNilChecks(language, context)
		req.NoError(err)

		files, err := jenny.Generate(context)
		req.NoError(err)

		tc.WriteFiles(files)
	})
}
