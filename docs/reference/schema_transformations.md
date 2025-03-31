---
weight: 10
---
<!-- Generated with `make docs` -->
# Schema transformations

## `add_fields`

AddFields rewrites the definition of an object to add new fields.
Note: existing fields will not be overwritten.

### Usage

```yaml
add_fields:
  # Expected format: [package].[object]
  to: string
  fields: []ast.StructField
```

## `add_object`

AddObject adds a new object to a schema.

### Usage

```yaml
add_object:
  object: string
  as: Type
  comments: []string
```

## `anonymous_structs_to_named`

AnonymousStructsToNamed turns "anonymous structs" into a named object.

Example:

	```
	Panel struct {
		Options struct {
			Title string
		}
	}
	```

Will become:

	```
	Panel struct {
		Options PanelOptions
	}

	PanelOptions struct {
		Title string
	}
	```

### Usage

```yaml
anonymous_structs_to_named: {}
```

## `constant_to_enum`

ConstantToEnum turns `string` constants into an enum definition with a
single member.
This is useful to "future-proof" a schema where a type can have a single
value for now but is expected to allow more in the future.

### Usage

```yaml
constant_to_enum:
  objects: []string
```

## `dataquery_identification`

N/A

### Usage

```yaml
dataquery_identification: {}
```

## `disjunction_infer_mapping`

DisjunctionInferMapping infers the discriminator field and mapping used to
describe a disjunction of references.
See https://swagger.io/docs/specification/data-models/inheritance-and-polymorphism/

### Usage

```yaml
disjunction_infer_mapping: {}
```

## `disjunction_of_anonymous_structs_to_explicit`

DisjunctionOfAnonymousStructsToExplicit looks for anonymous structs used as
branches of disjunctions and turns them into explicitly named types.

### Usage

```yaml
disjunction_of_anonymous_structs_to_explicit: {}
```

## `disjunction_to_type`

DisjunctionToType transforms disjunction into a struct, mapping disjunction branches to
an optional and nullable field in that struct.

Example:

		```
		SomeType: {
			type: "some-type"
	 	}
		SomeOtherType: {
			type: "other-type"
	 	}
		SomeStruct: {
			foo: string | bool
		}
		OtherStruct: {
			bar: SomeType | SomeOtherType
		}
		```

Will become:

		```
		SomeType: {
			type: "some-type"
	 	}
		SomeOtherType: {
			type: "other-type"
	 	}
		StringOrBool: {
			string: *string
			bool: *string
		}
		SomeStruct: {
			foo: StringOrBool
		}
		SomeTypeOrSomeOtherType: {
			SomeType: *SomeType
			SomeOtherType: *SomeOtherType
		}
		OtherStruct: {
			bar: SomeTypeOrSomeOtherType
		}
		```

### Usage

```yaml
disjunction_to_type: {}
```

## `disjunction_with_constant_to_default`

N/A

### Usage

```yaml
disjunction_with_constant_to_default: {}
```

## `duplicate_object`

DuplicateObject duplicates the source object. The duplicate is created under
a different name, possibly in a different package.

Note: if the source object isn't found, this pass does nothing.

### Usage

```yaml
duplicate_object:
  object: string
  as: string
  omit_fields: []string
```

## `entrypoint_identification`

N/A

### Usage

```yaml
entrypoint_identification: {}
```

## `fields_set_default`

FieldsSetDefault sets the default value for the given fields.

### Usage

```yaml
fields_set_default:
  defaults: map[string]interface {}
```

## `fields_set_not_required`

FieldsSetNotRequired rewrites the definition of given fields to mark them as nullable and not required.

### Usage

```yaml
fields_set_not_required:
  fields: []string
```

## `fields_set_required`

FieldsSetRequired rewrites the definition of given fields to mark them as not nullable and required.

### Usage

```yaml
fields_set_required:
  fields: []string
```

## `hint_object`

N/A

### Usage

```yaml
hint_object:
  object: string
  hints: JenniesHints
```

## `name_anonymous_struct`

NameAnonymousStruct rewrites the definition of a struct field typed as an
anonymous struct to instead refer to a named type.

### Usage

```yaml
name_anonymous_struct:
  field: string
  as: string
```

## `omit`

Omit rewrites schemas to omit the configured objects.

### Usage

```yaml
omit:
  objects: []string
```

## `omit_fields`

OmitFields removes the selected fields from their object definition.

### Usage

```yaml
omit_fields:
  fields: []string
```

## `rename_object`

N/A

### Usage

```yaml
rename_object:
  from: string
  to: string
```

## `replace_reference`

ReplaceReference replaces any usage of the `From` reference by the one given in `To`.

### Usage

```yaml
replace_reference:
  from: string
  to: string
```

## `retype_field`

N/A

### Usage

```yaml
retype_field:
  field: string
  as: Type
  comments: []string
```

## `retype_object`

N/A

### Usage

```yaml
retype_object:
  object: string
  as: Type
  comments: []string
```

## `schema_set_entry_point`

N/A

### Usage

```yaml
schema_set_entry_point:
  package: string
  entry_point: string
```

## `schema_set_identifier`

SchemaSetIdentifier overwrites the Metadata.Identifier field of a schema.

### Usage

```yaml
schema_set_identifier:
  package: string
  identifier: string
```

## `trim_enum_values`

TrimEnumValues removes leading and trailing spaces from string values.
It could happen when they add them by mistake in jsonschema/openapi when they define the enums

### Usage

```yaml
trim_enum_values: {}
```

## `unspec`

Unspec removes the Kubernetes-style envelope added by kindsys.

Objects named "spec" will be renamed, using the package as new name.

### Usage

```yaml
unspec: {}
```

