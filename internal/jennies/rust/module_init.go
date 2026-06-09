package rust

import (
	"fmt"
	"sort"
	"strings"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/languages"
)

var _ codejen.OneToMany[languages.Context] = ModuleInit{}

// ModuleInit emits the crate scaffolding that ties the per-package type modules
// (src/types/<pkg>.rs), builder modules (src/builders/<pkg>.rs) and the runtime
// (src/cog/) into a single compilable crate:
//
//   - Cargo.toml         (gated on Config.GenerateCargoToml)
//   - src/lib.rs         (declares the top-level cog/types/builders modules)
//   - src/types/mod.rs   (one `pub mod <pkg>;` per schema package)
//   - src/builders/mod.rs (one `pub mod <pkg>;` per builder package)
//
// The type-package set and builder-package set are derived independently from
// the context (schemas drive types, builders drive builders) because they need
// not coincide: a package may have types but no builders, and vice versa. The
// module names are sanitized with formatPackageName so they match the file
// names the RawTypes and Builder jennies emit.
type ModuleInit struct {
	config Config
}

func (jenny ModuleInit) JennyName() string {
	return "RustModuleInit"
}

func (jenny ModuleInit) Generate(context languages.Context) (codejen.Files, error) {
	files := make(codejen.Files, 0, 4)

	typePackages := jenny.sanitizedPackages(schemaPackages(context))
	builderPackages := jenny.sanitizedPackages(builderPackages(context))

	if jenny.config.GenerateCargoToml {
		files = append(files, *codejen.NewFile("Cargo.toml", jenny.cargoToml(), jenny))
	}

	files = append(files, *codejen.NewFile("src/lib.rs", jenny.libRs(len(builderPackages) > 0), jenny))
	files = append(files, *codejen.NewFile("src/types/mod.rs", modFile(typePackages), jenny))

	if len(builderPackages) > 0 {
		files = append(files, *codejen.NewFile("src/builders/mod.rs", modFile(builderPackages), jenny))
	}

	return files, nil
}

// cargoToml renders the crate manifest. serde (with the derive feature) and
// serde_json back every emitted type; serde_repr backs the numeric enums the
// RawTypes jenny derives with #[derive(Serialize_repr, Deserialize_repr)].
func (jenny ModuleInit) cargoToml() []byte {
	return fmt.Appendf(nil, `[package]
name = %q
version = %q
edition = "2021"

[dependencies]
serde = { version = "1.0", features = ["derive"] }
serde_json = "1.0"
serde_repr = "0.1"
`, jenny.config.crateName(), jenny.config.crateVersion())
}

// libRs renders the crate root. The clippy allows are crate-wide because
// generated code legitimately and pervasively triggers them: enum variants
// commonly repeat the enum name (enum_variant_names), discriminated unions hold
// branches of very different sizes (large_enum_variant), and builder
// constructors mirror wide upstream option sets (too_many_arguments). Silencing
// them at the crate root keeps `cargo clippy -D warnings` green without
// peppering #[allow] across every generated item.
func (jenny ModuleInit) libRs(withBuilders bool) []byte {
	var b strings.Builder
	b.WriteString(`// Generated crate root.
//
// The following lints are allowed crate-wide because generated code triggers
// them by construction, not by mistake:
//   - enum_variant_names: union variants often repeat the enum name.
//   - large_enum_variant: union branches can differ widely in size.
//   - too_many_arguments: builder constructors mirror upstream option sets.
#![allow(clippy::enum_variant_names)]
#![allow(clippy::large_enum_variant)]
#![allow(clippy::too_many_arguments)]

`)
	if !jenny.config.SkipRuntime {
		b.WriteString("pub mod cog;\n")
	}
	b.WriteString("pub mod types;\n")
	if withBuilders {
		b.WriteString("pub mod builders;\n")
	}
	return []byte(b.String())
}

// modFile renders a `pub mod <name>;` line per package. The input is assumed
// already sanitized; it is sorted and de-duplicated for deterministic output.
func modFile(packages []string) []byte {
	var b strings.Builder
	for _, pkg := range packages {
		fmt.Fprintf(&b, "pub mod %s;\n", pkg)
	}
	return []byte(b.String())
}

// sanitizedPackages maps raw IR package names through formatPackageName, sorts
// them and removes duplicates so the emitted module list is deterministic and
// each module is declared exactly once. Distinct raw names can sanitize to the
// same identifier (for example "with-dashes" and "with_dashes"), which must not
// produce a duplicate `pub mod`.
func (jenny ModuleInit) sanitizedPackages(raw []string) []string {
	seen := make(map[string]struct{}, len(raw))
	out := make([]string, 0, len(raw))
	for _, pkg := range raw {
		name := formatPackageName(pkg)
		if _, ok := seen[name]; ok {
			continue
		}
		seen[name] = struct{}{}
		out = append(out, name)
	}
	sort.Strings(out)
	return out
}

func schemaPackages(context languages.Context) []string {
	packages := make([]string, 0, len(context.Schemas))
	for _, schema := range context.Schemas {
		packages = append(packages, schema.Package)
	}
	return packages
}

func builderPackages(context languages.Context) []string {
	packages := make([]string, 0, len(context.Builders))
	for _, builder := range context.Builders {
		packages = append(packages, builder.Package)
	}
	return packages
}
