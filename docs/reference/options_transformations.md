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

AddAssignmentAction adds an assignment to an existing option.
### Usage

```yaml
add_assignment:
  assignment: Assignment
```

## `add_comments`

AddCommentsAction adds comments to an option.
### Usage

```yaml
add_comments:
  comments: []string
```

## `array_to_append`

ArrayToAppendAction updates the option to perform an "append" assignment.

Example:

	```
	func Tags(tags []string) {
		this.resource.tags = tags
	}
	```

Will become:

	```
	func Tags(tags string) {
		this.resource.tags.append(tags)
	}
	```

This action returns the option unchanged if:
  - it doesn't have exactly one argument
  - the argument is not an array
### Usage

```yaml
array_to_append: {}
```

## `debug`

DebugAction prints debugging information about an option.
### Usage

```yaml
debug: {}
```

## `disjunction_as_options`

DisjunctionAsOptionsAction uses the branches of the first argument's disjunction (assuming it is one) and turns them
into options.

Example:

	```
	func Panel(panel Panel|Row) {
		this.resource.panels.append(panel)
	}
	```

Will become:

	```
	func Panel(panel Panel) {
		this.resource.panels.append(panel)
	}

	func Row(row Row) {
		this.resource.panels.append(row)
	}
	```

This action returns the option unchanged if:
  - it has no arguments
  - the given argument is not a disjunction or a reference to one
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

MapToIndexAction updates the option to perform an "index" assignment.

Example:

	```
	func Elements(elements map[string]Element) {
		this.resource.elements = elements
	}
	```

Will become:

	```
	func Elements(key string, elements Element) {
		this.resource.elements[key] = tags
	}
	```

This action returns the option unchanged if:
  - it doesn't have exactly one argument
  - the argument is not a map
### Usage

```yaml
map_to_index: {}
```

## `omit`

OmitAction removes an option.
### Usage

```yaml
omit: {}
```

## `rename`

RenameAction renames an option.
### Usage

```yaml
rename:
  as: string
```

## `rename_arguments`

RenameArgumentsAction renames the arguments of an options.
### Usage

```yaml
rename_arguments:
  as: []string
```

## `struct_fields_as_arguments`

StructFieldsAsArgumentsAction uses the fields of the first argument's struct (assuming it is one) and turns them
into arguments.

Optionally, an explicit list of fields to turn into arguments can be given.

Example:

	```
	func Time(time {from string, to string) {
		this.resource.time = time
	}
	```

Will become:

	```
	func Time(from string, to string) {
		this.resource.time.from = from
		this.resource.time.to = to
	}
	```

This action returns the option unchanged if:
  - it has no arguments
  - the first argument is not a struct or a reference to one

FIXME: considers the first argument only.
### Usage

```yaml
struct_fields_as_arguments:
  fields: []string
```

## `struct_fields_as_options`

StructFieldsAsOptionsAction uses the fields of the first argument's struct (assuming it is one) and turns them
into options.

Optionally, an explicit list of fields to turn into options can be given.

Example:

	```
	func GridPos(gridPos {x int, y int) {
		this.resource.gridPos = gridPos
	}
	```

Will become:

	```
	func X(x int) {
		this.resource.gridPos.x = x
	}

	func Y(y int) {
		this.resource.gridPos.y = y
	}
	```

This action returns the option unchanged if:
  - it has no arguments
  - the first argument is not a struct or a reference to one

FIXME: considers the first argument only.
### Usage

```yaml
struct_fields_as_options:
  fields: []string
```

## `unfold_boolean`

UnfoldBooleanAction transforms an option accepting a boolean argument into two argument-less options.

Example:

	```
	func Editable(editable bool) {
		this.resource.editable = editable
	}
	```

Will become:

	```
	func Editable() {
		this.resource.editable = true
	}

	func ReadOnly() {
		this.resource.editable = false
	}
	```
### Usage

```yaml
unfold_boolean:
  true_as: string
  false_as: string
```

