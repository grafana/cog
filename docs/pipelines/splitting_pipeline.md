---
weight: 2
---
# Splitting a codegen pipeline

When a codegen pipeline depends on many inputs, its configuration file can
become hard to follow.

In this case, it can be a good idea to split this configuration in several
_units_.

A *codegen unit* is a fragment of a codegen pipeline that contains a set of
inputs and their associated schema/builder transformations.

The codegen pipeline can refer to those *units* with the `units_from` directive.
It contains a list of paths (or glob patterns) to codegen unit files:

```yaml hl_lines="6-9"
# yaml-language-server: $schema=https://raw.githubusercontent.com/grafana/cog/main/schemas/pipeline.json
# cog.yaml

debug: true

# `units_from` describs a list of paths where codegen units can be found.
units_from:
  - '%__config_dir%/resources/*/unit.yaml'
  - '%__config_dir%/extra/config.yaml'

# `inputs` can be declared as usual.
inputs: []

# The rest of the configuration doesn't change.
transformations:
  schemas: []
  builders: []

output:
  directory: './generated/%l'

  types: true
  builders: true
  converters: true
  api_reference: true

  languages: []
```

Each codegen unit can have two sections:

* an `inputs` directive, describing a list of input schemas belonging in the unit. These inputs are [configured as they would be in a pipeline](./creating_pipeline.md#inputs).
* a set of builder transformations, given as a list of files or directories containing the transformations. See "[Applying builder transformations](./builder_transformations.md)".

```yaml
# yaml-language-server: $schema=https://raw.githubusercontent.com/grafana/cog/main/schemas/unit.json
# resources/example/unit.yaml

inputs:
  - jsonschema:
      url: 'https://example.com/schemas/jsonschema.json'
      transformations:
        - '%__config_dir%/schema.transforms.yaml'

builder_transformations:
  - '%__config_dir%/common.builder.transforms.yaml'
  - '%__config_dir%/go-builder-transforms/'
```

!!! info

    In codegen units, the `%__config_dir%` variable refers to the directory
    containing the unit file. 


??? tip "Recommended: units validation and auto-complete"

    In order to help with configuring codegen units, a [schema][unit.json] for
    *codegen unit* files is provided. If your editor supports YAML schema
    validation, it is definitely recommended to set it up:

    === "Visual Studio Code"

        1.  Install [`vscode-yaml`][vscode-yaml] for YAML language support.
        2.  Add the schema under the `yaml.schemas` key in your user or
            workspace [`settings.json`][settings.json]:

            ``` json
            {
              "yaml.schemas": {
                "https://raw.githubusercontent.com/grafana/cog/main/schemas/units.json": "unit.yaml"
              }
            }
            ```

    === "Other"

        1.  Ensure your editor of choice has support for YAML schema validation.
        2.  Add the following lines at the top of `unit.yaml`:

            ``` yaml
            # yaml-language-server: $schema=https://raw.githubusercontent.com/grafana/cog/main/schemas/unit.json
            ```

[units.json]: https://raw.githubusercontent.com/grafana/cog/main/schemas/units.json
[vscode-yaml]: https://marketplace.visualstudio.com/items?itemName=redhat.vscode-yaml
[settings.json]: https://code.visualstudio.com/docs/getstarted/settings

## Transformations

See ["Applying schema transformations"](./schema_transformations.md) and
["Applying builder transformations"](./builder_transformations.md).
