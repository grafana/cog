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
// jenny. Phase 3a only covers scalar/struct type emission for the fixtures
// listed below; every other rawtypes fixture is skipped until later phases
// implement enums, arrays-as-types, maps-as-types, refs across packages,
// disjunctions, intersections and variants.
func TestRawTypes_Generate(t *testing.T) {
	test := testutils.GoldenFilesTestSuite[ast.Schema]{
		TestDataRoot: "../../../testdata/jennies/rawtypes",
		Name:         "RustRawTypes",
		Skip: map[string]string{
			"arrays":                           "Phase 3b+: array type aliases not implemented",
			"constant_reference_as_default":    "Phase 3b+: constant references not implemented",
			"constant_reference_discriminator": "Phase 3b+: constant references not implemented",
			"constant_references":              "Phase 3b+: constant references not implemented",
			"constraints":                      "Phase 3b+: constraints not implemented",
			"dashboard":                        "Phase 6: composable slots not implemented",
			"disjunction_of_numbers":           "Phase 3b+: disjunctions not implemented",
			"disjunctions":                     "Phase 3b+: disjunctions not implemented",
			"disjunctions_of_refs_without_discriminator": "Phase 3b+: disjunctions not implemented",
			"disjunctions_of_scalars_and_refs":           "Phase 3b+: disjunctions not implemented",
			"enums":                                      "Phase 3b+: enum emission not implemented",
			"field_with_struct_with_defaults":            "Phase 3b+: cross-object defaults not implemented",
			"intersections":                              "Phase 3b+: intersections not implemented",
			"maps":                                       "Phase 3b+: map type aliases not implemented",
			"nullable_fields":                            "Phase 3b: constant references not implemented",
			"package-with-dashes":                        "Phase 3b+: refs not implemented",
			"reference_of_reference":                     "Phase 3b+: refs not implemented",
			"refs":                                       "Phase 3b+: refs not implemented",
			"struct_with_complex_fields":                 "Phase 3b+: complex fields not implemented",
			"struct_with_optional_fields":                "Phase 3b: inline enum field requires enum emission",
			"struct_with_map_and_slice_default":          "Phase 3b+: map/slice defaults not implemented",
			"time_hint":                                  "Phase 3b+: time hints not implemented",
			"variant_dataquery":                          "Phase 6: variants not implemented",
			"variant_panelcfg_full":                      "Phase 6: variants not implemented",
			"variant_panelcfg_only_options":              "Phase 6: variants not implemented",
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

		tc.WriteFiles(files)
	})
}
