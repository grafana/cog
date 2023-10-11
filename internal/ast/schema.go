package ast

type SchemaKind string

const (
	SchemaKindCore       SchemaKind = "core"
	SchemaKindComposable SchemaKind = "composable"
)

type SchemaVariant string

const (
	SchemaVariantPanel     SchemaVariant = "panelcfg"
	SchemaVariantDataQuery SchemaVariant = "dataquery"
)

type Schemas []*Schema

func (schemas Schemas) LocateObject(pkg string, name string) (Object, bool) {
	for _, schema := range schemas {
		if schema.Package != pkg {
			continue
		}

		obj := schema.LocateObject(name)

		// TODO: schema.LocateObject() should return a "found" boolean
		return obj, obj.Name != ""
	}

	return Object{}, false
}

func (schemas Schemas) DeepCopy() []*Schema {
	newSchemas := make([]*Schema, 0, len(schemas))

	for _, schema := range schemas {
		newSchema := schema.DeepCopy()
		newSchemas = append(newSchemas, &newSchema)
	}

	return newSchemas
}

type Schema struct { //nolint: musttag
	Package  string
	Metadata SchemaMeta
	Objects  []Object
}

func (schema *Schema) DeepCopy() Schema {
	newSchema := Schema{
		Package:  schema.Package,
		Metadata: schema.Metadata,
		Objects:  make([]Object, 0, len(schema.Objects)),
	}

	for _, def := range schema.Objects {
		newSchema.Objects = append(newSchema.Objects, def.DeepCopy())
	}

	return newSchema
}

func (schema *Schema) LocateObject(name string) Object {
	for _, def := range schema.Objects {
		if def.Name == name {
			return def
		}
	}

	return Object{}
}

type SchemaMeta struct {
	Kind       SchemaKind
	Variant    SchemaVariant
	Identifier string
}
