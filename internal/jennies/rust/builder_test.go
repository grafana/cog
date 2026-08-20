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
			// anonymous_struct stays skipped: the option argument is a bare inline
			// anonymous struct (Type.Kind == struct), and this builder-test context does
			// not run the AnonymousStructsToNamed pass, so no named type exists to name
			// in the signature. Rust models an un-hoisted inline struct as
			// serde_json::Value, which cannot serve as a typed argument; this matches
			// Python, which also skips this fixture ("anonymous structs not supported").
			"anonymous_struct": "inline anonymous struct (not hoisted to a named type) has no idiomatic Rust argument type; matches Python skip",
			// builder_delegation_in_disjunction stays skipped: its delegated arguments
			// are inline disjunctions (e.g. Builder<T> | string), which have no clean
			// idiomatic-Rust delegation bound. The Rust target does not run the
			// DisjunctionToType compiler pass, so these never become named union types.
			// This matches the Go target, which also skips the fixture.
			"builder_delegation_in_disjunction": "inline-disjunction delegation has no idiomatic Rust mapping (Rust skips DisjunctionToType)",
			"dashboard_panel":                   "Java-generics-specific builder shape; Python also skips this fixture",
			// struct_fields_as_args_assignment stays skipped for the same reason as
			// anonymous_struct: the target `time` field is a nullable inline anonymous
			// struct (not hoisted to a named type in this context), which Rust models as
			// Option<serde_json::Value>. A deep assignment into its `from`/`to` fields
			// is not type-safe Rust, so this is unsupportable without the
			// AnonymousStructsToNamed pass. Python tolerates it only via a degenerate
			// `unknown` placeholder.
			"struct_fields_as_args_assignment": "deep assignment into an un-hoisted inline anonymous struct (serde_json::Value) is not type-safe Rust",
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

		// See rawtypes_test.go: the harness does not run postprocessors, so format
		// the generated files with rustfmt here to match the rustfmt-formatted
		// goldens the FormatRustFiles postprocessor produces in real generation.
		files = formatRustGoldenFiles(tc, files)

		tc.WriteFiles(files)
	})
}
