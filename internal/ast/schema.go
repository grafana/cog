package ast

import (
	"github.com/grafana/cog/internal/orderedmap"
)

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

		return schema.LocateObject(name)
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
	Metadata SchemaMeta `json:",omitempty"`
	Objects  *orderedmap.Map[string, Object]
}

func NewSchema(pkg string, metadata SchemaMeta) *Schema {
	return &Schema{
		Package:  pkg,
		Metadata: metadata,
		Objects:  orderedmap.New[string, Object](),
	}
}

func (schema *Schema) AddObject(object Object) {
	if _, exists := schema.LocateObject(object.Name); exists {
		return
	}

	schema.Objects.Set(object.Name, object)
}

func (schema *Schema) AddObjects(objects ...Object) {
	for _, object := range objects {
		schema.AddObject(object)
	}
}

func (schema *Schema) DeepCopy() Schema {
	newSchema := Schema{
		Package:  schema.Package,
		Metadata: schema.Metadata,
		Objects:  orderedmap.New[string, Object](),
	}

	schema.Objects.Iterate(func(_ string, object Object) {
		newSchema.AddObject(object.DeepCopy())
	})

	return newSchema
}

func (schema *Schema) LocateObject(name string) (Object, bool) {
	if !schema.Objects.Has(name) {
		return Object{}, false
	}

	return schema.Objects.Get(name), true
}

func (schema *Schema) Resolve(typeDef Type) (Type, bool) {
	if typeDef.Kind != KindRef {
		return typeDef, true
	}

	referredObj, found := schema.LocateObject(typeDef.AsRef().ReferredType)
	if !found {
		return Type{}, false
	}

	return schema.Resolve(referredObj.Type)
}

type SchemaMeta struct {
	Kind       SchemaKind    `json:",omitempty"`
	Variant    SchemaVariant `json:",omitempty"`
	Identifier string        `json:",omitempty"`
}
