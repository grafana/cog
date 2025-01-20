---
weight: 10
---

# Applying builder transformations

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

language: all

package: package_name

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
