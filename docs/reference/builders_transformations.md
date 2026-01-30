---
weight: 10
---
<!-- Generated with `make docs` -->
# Builder transformations

Each builder transformation requires the use of one of the following selectors, explicitly and unambiguously stating on which builder(s) the transformation should apply.
```yaml
by_object: string
by_name: string
by_variant: string
generated_from_disjunction: bool
```

Example:
```yaml
- rename:
    by_object: RowPanel
    as: Row
```

## `add_factory`

AddFactory adds a builder factory to the selected builders.
These factories are meant to be used to simplify the instantiation of
builders for common use-cases.
### Usage

```yaml
add_factory:
  factory: BuilderFactory
```

## `add_option`

AddOption adds a completely new option to the selected builders.
### Usage

```yaml
add_option:
  option: Option
```

## `compose`

N/A

### Usage

```yaml
compose:
  source_builder_name: string
  plugin_discriminator_field: string
  exclude_options: []string
  composition_map: map[string]string
  composed_builder_name: string
  preserve_original_builders: bool
  rename_options: map[string]string
```

## `duplicate`

Duplicate duplicates a builder.
The name of the duplicated builder has to be specified and some options can
be excluded.
### Usage

```yaml
duplicate:
  as: string
  exclude_options: []string
```

## `initialize`

N/A

### Usage

```yaml
initialize:
  set: []yaml.Initialization
```

## `merge_into`

N/A

### Usage

```yaml
merge_into:
  destination: string
  source: string
  under_path: string
  exclude_options: []string
  rename_options: map[string]string
```

## `omit`

Omit removes a builder.
### Usage

```yaml
omit: {}
```

## `promote_options_to_constructor`

PromoteOptionsToConstructor promotes the given options as constructor
parameters. Both arguments and assignments described by the options
will be exposed in the builder's constructor.
### Usage

```yaml
promote_options_to_constructor:
  options: []string
```

## `properties`

N/A

### Usage

```yaml
properties:
  set: []ast.StructField
```

## `rename`

Rename renames a builder.
### Usage

```yaml
rename:
  as: string
```

