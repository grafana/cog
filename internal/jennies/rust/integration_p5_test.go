package rust

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/languages"
	"github.com/stretchr/testify/require"
)

// TestPhase5Converters_Integration is a manual end-to-end proof that the
// converter output is compilable Rust and emits the expected builder calls. It
// merges several converter-exercising fixtures into one context, generates the
// full crate with converters enabled, and drops a behavior test into tests/ so
// `cargo test` can assert on the emitted source strings. It is gated behind
// COG_RUST_P5=1 because it writes outside the repo and depends on a local
// cargo toolchain:
//
//	COG_RUST_P5=1 go test ./internal/jennies/rust/ -run TestPhase5Converters_Integration
//	cd /tmp/p5-converters && cargo build && cargo clippy --all-targets -- -D warnings && cargo fmt --check && cargo test
func TestPhase5Converters_Integration(t *testing.T) {
	if os.Getenv("COG_RUST_P5") != "1" {
		t.Skip("set COG_RUST_P5=1 to run the cargo build integration proof")
	}
	req := require.New(t)

	// Fixture notes:
	//   - constant_assignment is excluded: its IR carries the nonstandard
	//     "boolean" scalar kind (instead of ast.KindBool's "bool"), which the
	//     Rust rawtypes jenny defensively maps to serde_json::Value; the
	//     fixture's *builders* do not compile as a crate either, so this is a
	//     pre-existing fixture quirk rather than a converter issue.
	//   - composable_slot is excluded: its Dashboard has a bare (non-Vec)
	//     `Box<dyn Dataquery>` field, and derive(PartialEq) on such a field does
	//     not compile (operator resolution tries to move the unsized trait
	//     object). This is a pre-existing rawtypes gap; the real SDK only uses
	//     Vec<Box<dyn Dataquery>> slots, which work. The converter runtime path
	//     (cog::convert_dataquery_to_code in src/cog/tools.rs) still compiles as
	//     part of every crate.
	//   - struct_with_defaults is excluded: its schema carries map-valued field
	//     defaults the rawtypes Default impl cannot render (pre-existing
	//     rawtypes limitation, unrelated to converters).
	fixtures := []string{"basic_struct", "builder_delegation", "map_of_builders", "known_any", "envelope_assignment", "references", "discriminator_without_option"}
	var merged languages.Context
	for _, name := range fixtures {
		raw, err := os.ReadFile(filepath.Join("../../../testdata/jennies/builders", name, "builders_context.json"))
		req.NoError(err)
		var ctx languages.Context
		req.NoError(json.Unmarshal(raw, &ctx))
		merged.Schemas = append(merged.Schemas, ctx.Schemas...)
		merged.Builders = append(merged.Builders, ctx.Builders...)
	}

	language := New(Config{GenerateCargoToml: true})
	merged, err := languages.GenerateBuilderNilChecks(language, merged)
	req.NoError(err)

	jennies := language.Jennies(languages.Config{Builders: true, Converters: true})
	files, err := jennies.GenerateFS(merged)
	req.NoError(err)

	outDir := "/tmp/p5-converters"
	req.NoError(os.RemoveAll(outDir))
	req.NoError(files.Write(context.Background(), outDir))

	for _, f := range []string{"src/converters/mod.rs", "src/converters/basic_struct.rs", "src/cog/tools.rs"} {
		_, statErr := os.Stat(filepath.Join(outDir, f))
		req.NoError(statErr, "expected converter file %s", f)
	}

	behaviorTest := `use grafana_foundation_sdk::converters::basic_struct::some_struct_converter;
use grafana_foundation_sdk::converters::builder_delegation::dashboard_converter;
use grafana_foundation_sdk::types::basic_struct::SomeStruct;
use grafana_foundation_sdk::types::builder_delegation::{Dashboard, DashboardLink};

#[test]
fn some_struct_converter_emits_builder_calls() {
    let input = SomeStruct {
        id: 42,
        uid: "abc".to_string(),
        tags: vec!["a".to_string(), "b".to_string()],
        live_now: true,
    };

    let code = some_struct_converter(&input);

    assert!(code.contains("basic_struct::SomeStructBuilder::new()"), "{code}");
    assert!(code.contains("id(42)"), "{code}");
    assert!(code.contains(r#"uid("abc".to_string())"#), "{code}");
    assert!(
        code.contains(r#"tags(vec!["a".to_string(), "b".to_string()])"#),
        "{code}"
    );
    assert!(code.contains("live_now(true)"), "{code}");
}

#[test]
fn dashboard_converter_delegates_to_link_converter() {
    let input = Dashboard {
        id: 1,
        title: "my dashboard".to_string(),
        links: vec![DashboardLink {
            title: "link".to_string(),
            url: "https://example.com".to_string(),
        }],
        links_of_links: Vec::new(),
        single_link: DashboardLink::default(),
    };

    let code = dashboard_converter(&input);

    assert!(
        code.contains("builder_delegation::DashboardBuilder::new()"),
        "{code}"
    );
    assert!(
        code.contains(r#"links(vec![builder_delegation::DashboardLinkBuilder::new()"#),
        "{code}"
    );
    assert!(code.contains(r#"title("link".to_string())"#), "{code}");
    assert!(
        code.contains(r#"url("https://example.com".to_string())"#),
        "{code}"
    );
}
`
	// rustfmt the behavior test so `cargo fmt --check` stays green for the
	// whole crate.
	formatted, err := FormatRustFiles(*codejen.NewFile("tests/converters.rs", []byte(behaviorTest), nil))
	req.NoError(err)
	req.NoError(os.MkdirAll(filepath.Join(outDir, "tests"), 0o755))
	req.NoError(os.WriteFile(filepath.Join(outDir, "tests", "converters.rs"), formatted.Data, 0o600))
}
