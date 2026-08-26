package ir

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

// Get returns a schema by its package name and true to indicate that the schema was found.
// Or it returns nil, and false to indicate that nothing was found.
// Note: the search is case-sensitive.
func (schemas Schemas) Get(pkg string) (*Schema, bool) {
	for _, schema := range schemas {
		if schema.Package == pkg {
			return schema, true
		}
	}

	return nil, false
}

// Resolve ensures that `def` refers to an actual type — as opposed to a Ref type —
// by resolving references until a non-ref type is found.
func (schemas Schemas) Resolve(def Type) Type {
	if !def.IsRef() {
		return def
	}

	resolved, found := schemas.GetObjectByRef(def.AsRef())
	if !found {
		return def
	}

	return schemas.Resolve(resolved.Type)
}

// GetObject returns the object named `name` in the package `pkg`. The second
// return value indicates whether the object was found or not.
// Note: the search is case-sensitive.
func (schemas Schemas) GetObject(pkg string, name string) (Object, bool) {
	for _, schema := range schemas {
		if schema.Package != pkg {
			continue
		}

		return schema.GetObject(name)
	}

	return Object{}, false
}

// GetObjectByRef returns the object described by the given reference. The second
// return value indicates whether the object was found or not.
// Note: the search is case-sensitive.
func (schemas Schemas) GetObjectByRef(ref RefType) (Object, bool) {
	return schemas.GetObject(ref.ReferredPkg, ref.ReferredType)
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

type Schema struct { // nolint: musttag
	Package        string
	Metadata       SchemaMeta
	EntryPoint     string `json:",omitempty"`
	EntryPointType Type
	Objects        *orderedmap.Map[string, Object]
}

func NewSchema(pkg string, metadata SchemaMeta) *Schema {
	return &Schema{
		Package:  pkg,
		Metadata: metadata,
		Objects:  orderedmap.New[string, Object](),
	}
}

func (schema *Schema) AddObject(object Object) {
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

	if schema.EntryPoint != other.EntryPoint && (schema.EntryPoint == "" || other.EntryPoint == "") {
		if schema.EntryPoint == "" {
			schema.EntryPoint = other.EntryPoint
			schema.EntryPointType = other.EntryPointType
		}
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
		Package:        schema.Package,
		Metadata:       schema.Metadata,
		EntryPoint:     schema.EntryPoint,
		EntryPointType: schema.EntryPointType,
		Objects: schema.Objects.Map(func(_ string, object Object) Object {
			return object.DeepCopy()
		}),
	}
}

// GetObject returns the object named `name` in the current package. The second
// return value indicates whether the object was found or not.
// Note: the search is case-sensitive.
func (schema *Schema) GetObject(name string) (Object, bool) {
	if !schema.HasObject(name) {
		return Object{}, false
	}

	return schema.Objects.Get(name), true
}

// HasObject indicates whether the object named `name` exists in the current package.
// Note: the search is case-sensitive.
func (schema *Schema) HasObject(name string) bool {
	return schema.Objects.Has(name)
}

// Resolve ensures that `def` refers to an actual type — as opposed to a Ref type —
// by resolving references until a non-ref type is found. The second
// return value indicates whether the type could be resolved or not.
func (schema *Schema) Resolve(typeDef Type) (Type, bool) {
	if !typeDef.IsRef() {
		return typeDef, true
	}

	if typeDef.AsRef().ReferredPkg != schema.Package {
		return Type{}, false
	}

	referredObj, found := schema.GetObject(typeDef.AsRef().ReferredType)
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
