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

`cog` provides a command that generates a skeleton that can be used to implement a plugin:

```shell
cog create-plugin -o ./cog-rust-plugin --go-module-path github.com/org/plugin rust
```

This command will write its results in the `./cog-rust-plugin` directory. It
will  also include a `README.md` file that will explain how to build the plugin
and what to do next.
