package rust

import (
	"testing"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/languages"
	"github.com/stretchr/testify/require"
)

// TestModuleInit_Generate asserts the crate-scaffolding jenny ties the
// per-package type and builder modules into a compilable crate. It uses a
// multi-package context (including a dashed package name) so the test exercises
// the formatPackageName sanitization and the deterministic ordering, and proves
// the type-package set and builder-package set are derived independently.
func TestModuleInit_Generate(t *testing.T) {
	req := require.New(t)

	context := languages.Context{
		Schemas: ast.Schemas{
			&ast.Schema{Package: "dashboard"},
			&ast.Schema{Package: "with-dashes"},
			&ast.Schema{Package: "common"},
		},
		Builders: ast.Builders{
			// Builders cover only a subset of, and a different set than, the
			// type packages: there is no "common" builder, and there is a
			// "playlist" builder package with no matching schema.
			{Package: "dashboard"},
			{Package: "playlist"},
			{Package: "with-dashes"},
		},
	}

	jenny := ModuleInit{config: Config{GenerateCargoToml: true}}
	files, err := jenny.Generate(context)
	req.NoError(err)

	byName := map[string]string{}
	for _, f := range files {
		byName[f.RelativePath] = string(f.Data)
	}

	// Cargo.toml is gated on GenerateCargoToml and carries crate identity + deps.
	cargo, ok := byName["Cargo.toml"]
	req.True(ok, "Cargo.toml should be emitted when GenerateCargoToml is set")
	req.Contains(cargo, `name = "grafana_foundation_sdk"`)
	req.Contains(cargo, `version = "0.0.1"`)
	req.Contains(cargo, `edition = "2021"`)
	req.Contains(cargo, "serde")
	req.Contains(cargo, "serde_json")
	req.Contains(cargo, "serde_repr")

	lib, ok := byName["src/lib.rs"]
	req.True(ok)
	req.Contains(lib, "pub mod cog;")
	req.Contains(lib, "pub mod types;")
	req.Contains(lib, "pub mod builders;")
	req.Contains(lib, "#![allow(clippy::enum_variant_names)]")

	// types/mod.rs lists every schema package, sanitized and sorted.
	typesMod, ok := byName["src/types/mod.rs"]
	req.True(ok)
	req.Equal("pub mod common;\npub mod dashboard;\npub mod with_dashes;\n", typesMod)

	// builders/mod.rs lists only packages with builders, sanitized and sorted.
	buildersMod, ok := byName["src/builders/mod.rs"]
	req.True(ok)
	req.Equal("pub mod dashboard;\npub mod playlist;\npub mod with_dashes;\n", buildersMod)
}

// TestModuleInit_Generate_NoCargoToml verifies the manifest is omitted when the
// flag is off, while scaffolding modules are still emitted.
func TestModuleInit_Generate_NoCargoToml(t *testing.T) {
	req := require.New(t)

	jenny := ModuleInit{config: Config{}}
	files, err := jenny.Generate(languages.Context{
		Schemas: ast.Schemas{{Package: "dashboard"}},
	})
	req.NoError(err)

	for _, f := range files {
		req.NotEqual("Cargo.toml", f.RelativePath, "Cargo.toml must be gated off by default")
	}
}

// TestModuleInit_Generate_NoBuilders omits the builders module from lib.rs and
// emits no builders/mod.rs when the context carries no builders, so lib.rs never
// declares a module that has no backing files.
func TestModuleInit_Generate_NoBuilders(t *testing.T) {
	req := require.New(t)

	jenny := ModuleInit{config: Config{}}
	files, err := jenny.Generate(languages.Context{
		Schemas: ast.Schemas{{Package: "dashboard"}},
	})
	req.NoError(err)

	byName := map[string]string{}
	for _, f := range files {
		byName[f.RelativePath] = string(f.Data)
	}

	req.NotContains(byName["src/lib.rs"], "pub mod builders;")
	_, ok := byName["src/builders/mod.rs"]
	req.False(ok, "builders/mod.rs should not be emitted without builders")
}
