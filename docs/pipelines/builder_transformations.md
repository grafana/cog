---
weight: 10
---

# Applying builder transformations

`cog` supports modifying its [intermediate representation](../reference/glossary.md#intermediate-representation)
of [builders](../reference/glossary.md#builder).

These transformations can be useful to ensure that [builders](../reference/glossary.md#builder)
and [options](../reference/glossary.md#builder-option) expose meaningful names,
make their APIs easier to understand and use, …

Builder transformations are defined within their own configuration files and are
then referenced by codegen pipeline configurations.

```yaml hl_lines="9-11"
# yaml-language-server: $schema=https://raw.githubusercontent.com/grafana/cog/main/schemas/pipeline.json
# cog.yaml

debug: true

inputs:
  - …

transformations:
  builders:
    - '%__config_dir%/transformations/builders/'

output:
  directory: './generated/%l'

  # …
```

!!! note
    Contrary to schema transformations, builder transformations always
    target a specific schema, and can target a specific language.

## Structure of builder transformations

A typical configuration file for builder transformations looks like this:


```yaml
# yaml-language-server: $schema=https://raw.githubusercontent.com/grafana/cog/main/schemas/veneers.json
# transformations/builders/dashboard.yaml

# Which language is being targetted.
# `all` targets every language configured by the codegen pipeline.
# To apply the transformations to a specific language, use its name as value.
# Required.
language: all

# Package for which the transformation apply.
# Required.
package: dashboard

# List of transformation applied to builder objects.
builders: [ … ]

# List of transformation applied to builder options.
options: [ … ]
```

## Reference

See [the builder transformations reference](../reference/builders_transformations.md) for a list of all supported transformations.

## Example

??? tip "Recommended: configuration validation and auto-complete"

    In order to help with configuring builder transformations, a
    [schema][schema.json] is provided. If your editor supports YAML schema
    validation, it is definitely recommended to set it up:

    === "Visual Studio Code"

        1.  Install [`vscode-yaml`][vscode-yaml] for YAML language support.
        2.  Add the following lines at the top of your file:

            ``` yaml
            # yaml-language-server: $schema=https://raw.githubusercontent.com/grafana/cog/main/schemas/veneers.json
            ```

    === "Other"

        1.  Ensure your editor of choice has support for YAML schema validation.
        2.  Add the following lines at the top of your file:

            ``` yaml
            # yaml-language-server: $schema=https://raw.githubusercontent.com/grafana/cog/main/schemas/veneers.json
            ```

[schema.json]: https://raw.githubusercontent.com/grafana/cog/main/schemas/veneers.json
[vscode-yaml]: https://marketplace.visualstudio.com/items?itemName=redhat.vscode-yaml

```yaml
# yaml-language-server: $schema=https://raw.githubusercontent.com/grafana/cog/main/schemas/veneers.json
# transformations/builders/dashboard.yaml

language: all

package: dashboard

builders:
  # We don't want this builder at all
  - omit: { by_object: GridPos }

  # Rename the builder for the TimePickerConfig object to TimePicker
  - rename:
      by_object: TimePickerConfig
      as: TimePicker

  # Dashboards must have a title, so let's set it via the constructor
  - promote_options_to_constructor:
      by_object: Dashboard
      options: [title]

options:
  # `Tooltip` looks better than `GraphTooltip`
  - rename:
      by_name: Dashboard.graphTooltip
      as: tooltip

  # Editable() + Readonly() instead of Editable(val bool)
  - unfold_boolean:
      by_name: Dashboard.editable
      true_as: editable
      false_as: readonly
```
