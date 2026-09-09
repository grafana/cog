package rust

import (
	"fmt"
	"sort"
	"strings"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/languages"
)

// Plugins emits the generated variant-registration module (src/cog/plugins.rs),
// the Rust analog of the Go `cog/plugins` package. It scans every schema for
// objects that implement a composable variant (the dataquery and panelcfg
// variants) and builds a `variants::Registry` that maps each variant's wire
// discriminator to a decoder. The generated `register_default_plugins` installs
// that registry as the process-wide default via `variants::set_default_registry`,
// mirroring how the Go/TypeScript SDKs populate their global runtime from a
// `plugins` package `init()`.
//
// Registration is gated exactly like the Go target: an object must carry the
// `implements_variant` hint and must not carry `skip_variant_plugin_registration`.
type Plugins struct {
	config Config
}

func (jenny Plugins) JennyName() string {
	return "RustVariantsPlugins"
}

// dataqueryVariant pairs a dataquery variant's wire identifier with the fully
// qualified Rust type that decodes it.
type dataqueryVariant struct {
	identifier string
	typePath   string
}

// panelcfgVariant pairs a panel type identifier with the Rust types of its
// options and (custom) fieldConfig payloads, either of which may be absent.
type panelcfgVariant struct {
	identifier      string
	optionsType     string
	fieldConfigType string
	hasOptions      bool
	hasFieldConfig  bool
}

func (jenny Plugins) Generate(context languages.Context) (codejen.Files, error) {
	if jenny.config.SkipRuntime {
		return nil, nil
	}

	dataqueries := make([]dataqueryVariant, 0)
	panelcfgs := map[string]*panelcfgVariant{}

	for _, schema := range context.Schemas {
		identifier := schema.Metadata.Identifier
		schema.Objects.Iterate(func(_ string, object ast.Object) {
			if !object.Type.ImplementsVariant() || object.Type.IsRef() {
				return
			}
			if object.Type.HasHint(ast.HintSkipVariantPluginRegistration) {
				return
			}

			typePath := fmt.Sprintf("crate::types::%s::%s", formatPackageName(schema.Package), formatTypeName(object.Name))

			if object.Type.IsDataqueryVariant() {
				dataqueries = append(dataqueries, dataqueryVariant{
					identifier: identifier,
					typePath:   typePath,
				})
				return
			}

			// Panelcfg variant: the panel type identifier keys a config whose
			// options and fieldConfig.custom payloads are decoded by the named
			// objects. A panelcfg schema names its options object "Options" and its
			// custom fieldConfig object "FieldConfig"; either may be absent.
			variant, ok := panelcfgs[identifier]
			if !ok {
				variant = &panelcfgVariant{identifier: identifier}
				panelcfgs[identifier] = variant
			}
			switch object.Name {
			case "FieldConfig":
				variant.fieldConfigType = typePath
				variant.hasFieldConfig = true
			default:
				variant.optionsType = typePath
				variant.hasOptions = true
			}
		})
	}

	// Sort for deterministic output.
	sort.Slice(dataqueries, func(i, j int) bool {
		if dataqueries[i].identifier == dataqueries[j].identifier {
			return dataqueries[i].typePath < dataqueries[j].typePath
		}
		return dataqueries[i].identifier < dataqueries[j].identifier
	})

	panelcfgKeys := make([]string, 0, len(panelcfgs))
	for key := range panelcfgs {
		panelcfgKeys = append(panelcfgKeys, key)
	}
	sort.Strings(panelcfgKeys)

	output := jenny.render(dataqueries, panelcfgKeys, panelcfgs)
	return codejen.Files{*codejen.NewFile("src/cog/plugins.rs", []byte(output), jenny)}, nil
}

func (jenny Plugins) render(dataqueries []dataqueryVariant, panelcfgKeys []string, panelcfgs map[string]*panelcfgVariant) string {
	var buffer strings.Builder

	buffer.WriteString("//! Generated registration of every known composable variant.\n")
	buffer.WriteString("//!\n")
	buffer.WriteString("//! This module is the Rust analog of the Go `cog/plugins` package. It builds\n")
	buffer.WriteString("//! a [`variants::Registry`] mapping each variant's wire discriminator (a\n")
	buffer.WriteString("//! datasource type for dataqueries, a panel type for panelcfgs) to a decoder,\n")
	buffer.WriteString("//! and installs it as the process-wide default. Call\n")
	buffer.WriteString("//! [`register_default_plugins`] once at startup so the composable-slot\n")
	buffer.WriteString("//! deserialize helpers can resolve concrete variants.\n\n")
	buffer.WriteString("use crate::cog::variants;\n\n")

	buffer.WriteString("/// Builds the registry of every known composable variant.\n")
	buffer.WriteString("pub fn registry() -> variants::Registry {\n")

	// With no variants there is nothing to mutate, so return a fresh registry
	// directly: a `let mut` binding that is never mutated trips clippy's
	// let_and_return / unused_mut lints under -D warnings.
	if len(dataqueries) == 0 && len(panelcfgKeys) == 0 {
		buffer.WriteString("    variants::Registry::new()\n")
		buffer.WriteString("}\n\n")
		buffer.WriteString(jenny.renderRegisterFn())
		return buffer.String()
	}

	buffer.WriteString("    let mut registry = variants::Registry::new();\n")

	if len(dataqueries) > 0 {
		buffer.WriteString("\n    // Dataquery variants.\n")
		for _, dq := range dataqueries {
			buffer.WriteString("    registry.register_dataquery_variant(variants::DataqueryConfig {\n")
			fmt.Fprintf(&buffer, "        identifier: %q.to_string(),\n", dq.identifier)
			fmt.Fprintf(&buffer, "        decoder: |raw| Ok(Box::new(serde_json::from_value::<%s>(raw.clone())?)),\n", dq.typePath)
			// Converter registration is a veneer-level concern; cog itself never
			// populates it, mirroring the Go target's unset GoConverter hook.
			buffer.WriteString("        converter: None,\n")
			buffer.WriteString("    });\n")
		}
	}

	if len(panelcfgKeys) > 0 {
		buffer.WriteString("\n    // Panelcfg variants.\n")
		for _, key := range panelcfgKeys {
			variant := panelcfgs[key]
			buffer.WriteString("    registry.register_panelcfg_variant(variants::PanelcfgConfig {\n")
			fmt.Fprintf(&buffer, "        identifier: %q.to_string(),\n", variant.identifier)
			if variant.hasOptions {
				fmt.Fprintf(&buffer, "        options_decoder: Some(|raw| serde_json::to_value(serde_json::from_value::<%s>(raw.clone())?)),\n", variant.optionsType)
			} else {
				buffer.WriteString("        options_decoder: None,\n")
			}
			if variant.hasFieldConfig {
				fmt.Fprintf(&buffer, "        field_config_decoder: Some(|raw| serde_json::to_value(serde_json::from_value::<%s>(raw.clone())?)),\n", variant.fieldConfigType)
			} else {
				buffer.WriteString("        field_config_decoder: None,\n")
			}
			// See the dataquery configs above: converters are registered by
			// veneers, never by cog itself.
			buffer.WriteString("        converter: None,\n")
			buffer.WriteString("    });\n")
		}
	}

	buffer.WriteString("\n    registry\n")
	buffer.WriteString("}\n\n")
	buffer.WriteString(jenny.renderRegisterFn())

	return buffer.String()
}

// renderRegisterFn emits the `register_default_plugins` installer that is common
// to both the populated and the empty registry output.
func (jenny Plugins) renderRegisterFn() string {
	var buffer strings.Builder
	buffer.WriteString("/// Installs the registry of every known variant as the process-wide default.\n")
	buffer.WriteString("///\n")
	buffer.WriteString("/// Idempotent: the first caller wins, matching the once-only `init()`\n")
	buffer.WriteString("/// semantics of the other SDKs. Returns whether this call performed the\n")
	buffer.WriteString("/// installation.\n")
	buffer.WriteString("pub fn register_default_plugins() -> bool {\n")
	buffer.WriteString("    variants::set_default_registry(registry()).is_ok()\n")
	buffer.WriteString("}\n")
	return buffer.String()
}
