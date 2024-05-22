<!-- Generated with `make docs` -->
# Compiler passes

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

## `cloudwatch`

Cloudwatch rewrites a part of the cloudwatch schema.

In that schema, the `QueryEditorExpression` type is defined as a disjunction
for which the discriminator and mapping can not be inferred.
This compiler pass is here to define that mapping.

The `QueryEditorArrayExpression` struct type is also modified to simplify the
definition of its `expression` field from `[...#QueryEditorExpression] | [...#QueryEditorArrayExpression]` to
`[...#QueryEditorExpression]`.
This should be semantically equivalent since `#QueryEditorExpression` is a
union type that includes `#QueryEditorArrayExpression`.

The Cloudwatch pass also alerts the definition of the `#CloudWatchMetricsQuery`, `#CloudWatchLogsQuery` and
`#CloudWatchAnnotationQuery` types.
It removes the "dataquery variant" hint they carry, and defines a `CloudWatchQuery` type instead as a disjunction.
That disjunction serves as "dataquery entrypoint" for cloudwatch.

### Usage

```yaml
cloudwatch: {}
```

## `dashboard_panels`

DashboardPanelsRewrite rewrites the definition of "panels" fields in the "dashboard" package.

In the original schema, panels are defined as follows:

	```
	# In the Dashboard object
	panels?: [...#Panel | #RowPanel | #GraphPanel | #HeatmapPanel]

	# In the RowPanel object
	panels: [...#Panel | #GraphPanel | #HeatmapPanel]
	```

These definitions become:

	```
	# In the Dashboard object
	panels?: [...#Panel | #RowPanel]

	# In the RowPanel object
	panels: [...#Panel]
	```

### Usage

```yaml
dashboard_panels: {}
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

## `google_cloud_monitoring`

GoogleCloudMonitoring rewrites a part of the googlecloudmonitoring schema.

Older schemas (pre 10.2.x) define `CloudMonitoringQuery.timeSeriesList`
as a disjunction that cog can't handle: `timeSeriesList?: #TimeSeriesList | #AnnotationQuery`,
where `AnnotationQuery` is a type that extends `TimeSeriesList` to add two
fields.

This compiler pass checks for the presence of that disjunction, and rewrites
it as a reference to `TimeSeriesList`. It also adds the two missing fields
to this type if they aren't already defined.

### Usage

```yaml
google_cloud_monitoring: {}
```

## `hint_object`

N/A

### Usage

```yaml
hint_object:
  object: string
  hints: JenniesHints
```

## `library_panels`

LibraryPanels rewrites the definition of the "LibraryPanel" object in the "librarypanel" package.

In the original schema, the "model" field is left mainly undefined but a comment indicates
that it should be the same panel schema defined in dashboard with a few fields omitted.

This compiler pass implements the modifications described in that comment to define the
"model" field as:

	```
	# In the LibraryPanel object
	model: Omit<dashboard.Panel, 'gridPos' | 'id' | 'libraryPanel'>
	```

Note: this pass needs the "dashboard.Panel" schema to be parsed. Barring that, it leaves
the schemas untouched.

### Usage

```yaml
library_panels: {}
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

## `rename_object`

N/A

### Usage

```yaml
rename_object:
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

## `schema_set_identifier`

SchemaSetIdentifier overwrites the Metadata.Identifier field of a schema.

### Usage

```yaml
schema_set_identifier:
  package: string
  identifier: string
```

## `unspec`

Unspec removes the Kubernetes-style envelope added by kindsys.

Objects named "spec" will be renamed, using the package as new name.

### Usage

```yaml
unspec: {}
```

