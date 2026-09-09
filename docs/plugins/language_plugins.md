# Language plugins

In addition to *core* languages, `cog` can support additional languages through plugins.

Plugins are executables written in Go that communicate with `cog` over an RPC interface.

## Using language plugins

Languages exposed by plugins are configured in the [codegen pipeline](../pipelines/creating_pipeline.md):

```yaml hl_lines="13-15"
# yaml-language-server: $schema=https://raw.githubusercontent.com/grafana/cog/main/schemas/pipeline.json

inputs: [] # …

output:
  directory: './generated/%l'

  types: true
  builders: true
  converters: true
  api_reference: true

  language_plugins:
    rust:
      create_name: 'crate'
```

`language_plugins` is a map of plugins to use, associating their name with their configuration.

For each of these plugins, `cog` will look for an executable named `cog-{language}` in the `PATH`.

!!! tip "Manually locating plugins"

    The `--plugin-directory` flag can also be used to explicitly tell `cog` where to look for plugins.

## Writing language plugins

At their core, languages plugins are binaries that rely on [hashicorp/go-plugin](https://github.com/hashicorp/go-plugin) to expose the following interface to `cog`:

```go
// Language describes the interface that must be implemented by language plugins.
// Note: the methods are listed in the order in which they will be called.
type Language interface {
	// ValidateConfig receives the configuration for the current language
	// plugin and validates it.
	// Invalid configurations will result in an error while this method will
	// return nil for valid ones.
	ValidateConfig(config map[string]any) error

	// NullableConfig describes some properties of nullable types for a given language.
	// See also: [languages.NullableConfig]
	NullableConfig(config map[string]any) (languages.NullableConfig, error)

	// TransformSchemas modifies the input schemas to make them suitable to
	// the current language specifically.
	//
	// These transformations should only alter the schemas to make them
	// "compatible" with the target language, NOT to add missing elements or
	// fix incorrect ones.
	//
	// Examples of possible transformations:
	//   * [transforms.AnonymousStructsToNamed] can be used for languages that
	//     don't support anonymous structures
	//   * [transforms.DisjunctionOfConstantsToEnum] for languages that don't
	//     support disjunctions (or: union/sum types). See also other
	//     Disjunction* transformations.
	//   * …
	//
	// Returns the transformed schemas, or an error.
	TransformSchemas(config map[string]any, schemas ir.Schemas) (ir.Schemas, error)

	// Generate performs the code generation.
	//
	// In addition to the codegen context, it receives both the "global"
	// configuration (what should be generated) and the language-specific
	// configuration.
	//
	// The context contains intermediate representation (IR) of the schemas and
	// builders.
	//
	// Note: the builders IR is in the context only if builders generation is
	// enabled (ie: codegenConfig.Builders == true)
	Generate(codegenConfig languages.Config, config map[string]any, context languages.Context) (codejen.Files, error)
}
```

`cog` provides a command that generates a skeleton that can be used to implement a plugin:

```shell
cog create-plugin --go-module-path github.com/org/plugin rust
```

This command will bootstrap a plugin in the `./plugin` directory. It will also
include a `README.md` file that will explain how to build the project and what
to do next.
