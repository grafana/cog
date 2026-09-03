package rust

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/grafana/cog/internal/languages"
	"github.com/stretchr/testify/require"
)

// TestPhase7aCrateScaffolding_Integration is a manual end-to-end proof that the
// ModuleInit scaffolding turns the loose per-package .rs files into a compilable
// crate. It merges three single-package builder fixtures (basic_struct, sandbox,
// properties) into one multi-package context so the generated crate exercises
// several type modules and builder modules tied together by the scaffolding. It
// is gated behind COG_RUST_P7A=1 because it writes outside the repo and depends
// on a local cargo toolchain.
//
//	COG_RUST_P7A=1 go test ./internal/jennies/rust/ -run TestPhase7aCrateScaffolding_Integration
func TestPhase7aCrateScaffolding_Integration(t *testing.T) {
	if os.Getenv("COG_RUST_P7A") != "1" {
		t.Skip("set COG_RUST_P7A=1 to run the cargo build integration proof")
	}
	req := require.New(t)

	fixtures := []string{"basic_struct", "constructor_argument", "properties"}
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

	jennies := language.Jennies(languages.Config{Builders: true})
	files, err := jennies.GenerateFS(merged)
	req.NoError(err)

	outDir := "/tmp/p7a"
	req.NoError(os.RemoveAll(outDir))
	req.NoError(files.Write(context.Background(), outDir))

	for _, f := range []string{"Cargo.toml", "src/lib.rs", "src/types/mod.rs", "src/builders/mod.rs"} {
		_, statErr := os.Stat(filepath.Join(outDir, f))
		req.NoError(statErr, "expected scaffolding file %s", f)
	}
}
