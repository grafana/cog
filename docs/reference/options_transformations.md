---
weight: 10
---
<!-- Generated with `make docs` -->
# Option transformations

Each option transformation requires the use of one of the following selectors, explicitly and unambiguously stating on which option(s) the transformation should apply.
```yaml
by_name: string
by_builder: string
by_names: yaml.ByNamesSelector
```

Example:
```yaml
# H() → Height()
- rename:
    by_name: Panel.h
    as: height
```

## `add_assignment`

N/A

### Usage

```yaml
add_assignment:
  assignment: Assignment
```

## `add_comments`

N/A

### Usage

```yaml
add_comments:
  comments: []string
```

## `array_to_append`

N/A

### Usage

```yaml
array_to_append: {}
```

## `disjunction_as_options`

N/A

### Usage

```yaml
disjunction_as_options:
  argument_index: int
```

## `duplicate`

N/A

### Usage

```yaml
duplicate:
  as: string
```

## `map_to_index`

N/A

### Usage

```yaml
map_to_index: {}
```

## `omit`

N/A

### Usage

```yaml
omit: {}
```

## `rename`

Rename does things.
### Usage

```yaml
rename:
  as: string
```

## `rename_arguments`

N/A

### Usage

```yaml
rename_arguments:
  as: []string
```

## `struct_fields_as_arguments`

N/A

### Usage

```yaml
struct_fields_as_arguments:
  fields: []string
```

## `struct_fields_as_options`

N/A

### Usage

```yaml
struct_fields_as_options:
  fields: []string
```

## `unfold_boolean`

N/A

### Usage

```yaml
unfold_boolean:
  true_as: string
  false_as: string
```

