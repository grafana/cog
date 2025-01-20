---
weight: 5
---

# Applying schema transformations

`cog` supports modifying its [intermediate representation](../reference/glossary.md#intermediate-representation)
of types – built as a result of parsing the inputs – to tweak it, add missing
elements or fix incorrect ones.

This feature can be useful when the original schemas can not be fixed
directly and should be used with caution.

!!! note
    Since transformations are applied on the [intermediate representation](../reference/glossary.md#intermediate-representation),
    their effects automatically impact the code generated in every language.

Schema transformations are defined within their own configuration files and are
then referenced by codegen pipeline configurations.

## Globally

When defined within the `transformations` block of a codegen pipeline,
transformations are applied to every input in that pipeline:

```yaml hl_lines="9-11"
# yaml-language-server: $schema=https://raw.githubusercontent.com/grafana/cog/main/schemas/pipeline.json
# cog.yaml

debug: true

inputs:
  - …

transformations:
  schemas:
    - '%__config_dir%/transformations/schemas/common.yaml'

output:
  directory: './generated/%l'

  # …
```

## Per schema

When defined within `input` block of a codegen pipeline, transformations are
applied to that input only:

```yaml hl_lines="10-11"
# yaml-language-server: $schema=https://raw.githubusercontent.com/grafana/cog/main/schemas/pipeline.json
# cog.yaml

debug: true

inputs:
  - jsonschema:
      url: 'https://example.com/schema.json'
      package: example
      transformations:
        - '%__config_dir%/transformations/schemas/example.yaml'

output:
  directory: './generated/%l'

  # …
```

## Reference

See [the schema transformations reference](../reference/schema_transformations.md) for a list of all supported transformations.

## Example

??? tip "Recommended: configuration validation and auto-complete"

    In order to help with configuring schema transformations, a
    [schema][schema.json] is provided. If your editor supports YAML schema
    validation, it is definitely recommended to set it up:

    === "Visual Studio Code"

        1.  Install [`vscode-yaml`][vscode-yaml] for YAML language support.
        2.  Add the following lines at the top of your file:

            ``` yaml
            # yaml-language-server: $schema=https://raw.githubusercontent.com/grafana/cog/main/schemas/compiler_passes.json
            ```

    === "Other"

        1.  Ensure your editor of choice has support for YAML schema validation.
        2.  Add the following lines at the top of your file:

            ``` yaml
            # yaml-language-server: $schema=https://raw.githubusercontent.com/grafana/cog/main/schemas/compiler_passes.json
            ```

[schema.json]: https://raw.githubusercontent.com/grafana/cog/main/schemas/compiler_passes.json
[vscode-yaml]: https://marketplace.visualstudio.com/items?itemName=redhat.vscode-yaml

```yaml
# yaml-language-server: $schema=https://raw.githubusercontent.com/grafana/cog/main/schemas/compiler_passes.json
# transformations/schemas/example.yaml

passes:
  # Rename the AlertRuleGroup object to RuleGroup to avoid stuttering.
  - rename_object:
      from: alerting.AlertRuleGroup
      to: RuleGroup

  # The original type is incorrect, let's change it.
  - retype_object:
      object: alerting.Duration
      as:
        kind: scalar
        scalar: { scalar_kind: int64 }
      comments:
        - 'Duration in seconds.'

  # These fields are incorrectly set as non-required in the original schema.
  - fields_set_required:
      fields:
        - alerting.RuleGroup.interval
        - alerting.RuleGroup.name

  # Add a missing field.
  # This is meant to be removed once the schema is updated.
  - add_fields:
      to: alerting.Query
      fields:
        - name: timeout
          type:
            kind: scalar
            scalar: { scalar_kind: int64 }
```
