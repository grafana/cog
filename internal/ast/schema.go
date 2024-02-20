package ast

import (
	"errors"
	"fmt"

	"github.com/grafana/cog/internal/orderedmap"
)

var ErrCannotMergeSchemas = errors.New("can not merge schemas")

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

func (schemas Schemas) Consolidate() (Schemas, error) {
	byPackage := make(map[string]Schemas, len(schemas))

	for _, schema := range schemas {
		byPackage[schema.Package] = append(byPackage[schema.Package], schema)
	}

	newSchemas := make([]*Schema, 0, len(schemas))
	for pkg, groupedSchemas := range byPackage {
		newSchema := NewSchema(pkg, groupedSchemas[0].Metadata)
		for _, schema := range groupedSchemas {
			if err := newSchema.Merge(schema); err != nil {
				return nil, err
			}
		}

		newSchemas = append(newSchemas, newSchema)
	}

	return newSchemas, nil
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

func (schema *Schema) Merge(other *Schema) error {
	if schema.Package != other.Package {
		return fmt.Errorf("schemas originate from different packages ('%s', '%s'): %w", schema.Package, other.Package, ErrCannotMergeSchemas)
	}

	if !schema.Metadata.Equal(other.Metadata) {
		return fmt.Errorf("conflicting metadata: %w", ErrCannotMergeSchemas)
	}

	var err error
	other.Objects.Iterate(func(objectName string, remoteObject Object) {
		if !schema.Objects.Has(objectName) {
			schema.AddObject(remoteObject)
			return
		}

		object := schema.Objects.Get(objectName)

		if !object.Equal(remoteObject) {
			err = fmt.Errorf("conflicting definition for object '%s': %w", object.SelfRef.String(), ErrCannotMergeSchemas)
		}
	})
	if err != nil {
		return err
	}

	return nil
}

func (schema *Schema) DeepCopy() Schema {
	return Schema{
		Package:  schema.Package,
		Metadata: schema.Metadata,
		Objects: schema.Objects.Map(func(_ string, object Object) Object {
			return object.DeepCopy()
		}),
	}
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

func (meta SchemaMeta) Equal(other SchemaMeta) bool {
	return meta.Identifier == other.Identifier &&
		meta.Kind == other.Kind &&
		meta.Variant == other.Variant
}
