package transforms

import (
	"errors"
	"fmt"

	"github.com/grafana/cog/internal/orderedmap"
	"github.com/grafana/cog/pkg/ir"
)

type VisitSchemaFunc func(visitor *Visitor, schema *ir.Schema) (*ir.Schema, error)
type VisitObjectFunc func(visitor *Visitor, schema *ir.Schema, object ir.Object) (ir.Object, error)
type VisitTypeFunc func(visitor *Visitor, schema *ir.Schema, def ir.Type) (ir.Type, error)
type VisitStructFieldFunc func(visitor *Visitor, schema *ir.Schema, field ir.StructField) (ir.StructField, error)

type Visitor struct {
	OnSchema       VisitSchemaFunc
	OnObject       VisitObjectFunc
	OnStructField  VisitStructFieldFunc
	OnArray        VisitTypeFunc
	OnMap          VisitTypeFunc
	OnStruct       VisitTypeFunc
	OnDisjunction  VisitTypeFunc
	OnIntersection VisitTypeFunc
	OnEnum         VisitTypeFunc
	OnScalar       VisitTypeFunc
	OnRef          VisitTypeFunc
	OnConstantRef  VisitTypeFunc

	newObjects *orderedmap.Map[string, ir.Object]
}

func (visitor *Visitor) RegisterNewObject(object ir.Object) {
	visitor.newObjects.Set(object.SelfRef.String(), object)
}

func (visitor *Visitor) HasNewObject(ref ir.RefType) bool {
	return visitor.newObjects.Has(ref.String())
}

func (visitor *Visitor) VisitSchemas(schemas ir.Schemas) (ir.Schemas, error) {
	var err error

	for i, schema := range schemas {
		schemas[i], err = visitor.VisitSchema(schema)
		if err != nil {
			return nil, fmt.Errorf("[%s] %w", schema.Package, err)
		}
	}

	return schemas, nil
}

func (visitor *Visitor) VisitSchema(schema *ir.Schema) (*ir.Schema, error) {
	newSchema := schema.DeepCopy()
	newSchema.Objects = orderedmap.New[string, ir.Object]()

	visitor.newObjects = orderedmap.New[string, ir.Object]()
	defer func() {
		visitor.newObjects.Iterate(func(_ string, object ir.Object) {
			newSchema.AddObject(object)
		})
	}()

	if visitor.OnSchema != nil {
		visitedSchema, err := visitor.OnSchema(visitor, schema)
		if err != nil {
			return nil, err
		}

		// to ensure that new objects will be added to the visited schema by
		// the defer statement above.
		newSchema = *visitedSchema

		return &newSchema, nil
	}

	var err error
	var obj ir.Object

	newSchema.EntryPointType, err = visitor.VisitType(schema, schema.EntryPointType)
	if err != nil {
		return nil, fmt.Errorf("could not process entrypoint type")
	}

	schema.Objects.Iterate(func(_ string, object ir.Object) {
		if err != nil {
			return
		}

		obj, err = visitor.VisitObject(schema, object)
		if err != nil {
			err = errors.Join(
				fmt.Errorf("could not process object '%s'", object.Name),
				err,
			)
		}

		newSchema.AddObject(obj)
	})

	return &newSchema, err
}

func (visitor *Visitor) VisitObject(schema *ir.Schema, object ir.Object) (ir.Object, error) {
	if visitor.OnObject != nil {
		return visitor.OnObject(visitor, schema, object)
	}

	var err error

	object.Type, err = visitor.VisitType(schema, object.Type)
	if err != nil {
		return ir.Object{}, err
	}

	return object, nil
}

func (visitor *Visitor) VisitType(schema *ir.Schema, def ir.Type) (ir.Type, error) {
	if def.IsArray() {
		return visitor.VisitArray(schema, def)
	}

	if def.IsMap() {
		return visitor.VisitMap(schema, def)
	}

	if def.IsStruct() {
		return visitor.VisitStruct(schema, def)
	}

	if def.IsDisjunction() {
		return visitor.VisitDisjunction(schema, def)
	}

	if def.IsIntersection() {
		return visitor.VisitIntersection(schema, def)
	}

	if def.IsEnum() {
		return visitor.VisitEnum(schema, def)
	}

	if def.IsScalar() {
		return visitor.VisitScalar(schema, def)
	}

	if def.IsRef() {
		return visitor.VisitRef(schema, def)
	}

	if def.IsConstantRef() {
		return visitor.VisitConsantRef(schema, def)
	}

	return def, nil
}

func (visitor *Visitor) VisitArray(schema *ir.Schema, def ir.Type) (ir.Type, error) {
	if visitor.OnArray != nil {
		return visitor.OnArray(visitor, schema, def)
	}

	var err error

	def.Array.ValueType, err = visitor.VisitType(schema, def.AsArray().ValueType)
	if err != nil {
		return ir.Type{}, err
	}

	return def, nil
}

func (visitor *Visitor) VisitMap(schema *ir.Schema, def ir.Type) (ir.Type, error) {
	if visitor.OnMap != nil {
		return visitor.OnMap(visitor, schema, def)
	}

	var err error

	def.Map.ValueType, err = visitor.VisitType(schema, def.AsMap().ValueType)
	if err != nil {
		return ir.Type{}, err
	}

	return def, nil
}

func (visitor *Visitor) VisitStruct(schema *ir.Schema, def ir.Type) (ir.Type, error) {
	if visitor.OnStruct != nil {
		return visitor.OnStruct(visitor, schema, def)
	}

	var err error

	for i, field := range def.Struct.Fields {
		def.Struct.Fields[i], err = visitor.VisitStructField(schema, field)
		if err != nil {
			return ir.Type{}, errors.Join(
				fmt.Errorf("could not process struct field '%s'", field.Name),
				err,
			)
		}
	}

	return def, nil
}

func (visitor *Visitor) VisitDisjunction(schema *ir.Schema, def ir.Type) (ir.Type, error) {
	if visitor.OnDisjunction != nil {
		return visitor.OnDisjunction(visitor, schema, def)
	}

	var err error

	for i, branch := range def.Disjunction.Branches {
		def.Disjunction.Branches[i], err = visitor.VisitType(schema, branch)
		if err != nil {
			return ir.Type{}, err
		}
	}

	return def, nil
}

func (visitor *Visitor) VisitIntersection(schema *ir.Schema, def ir.Type) (ir.Type, error) {
	if visitor.OnIntersection != nil {
		return visitor.OnIntersection(visitor, schema, def)
	}

	var err error

	for i, branch := range def.Intersection.Branches {
		def.Intersection.Branches[i], err = visitor.VisitType(schema, branch)
		if err != nil {
			return ir.Type{}, err
		}
	}

	return def, nil
}

func (visitor *Visitor) VisitStructField(schema *ir.Schema, field ir.StructField) (ir.StructField, error) {
	if visitor.OnStructField != nil {
		return visitor.OnStructField(visitor, schema, field)
	}

	var err error

	field.Type, err = visitor.VisitType(schema, field.Type)
	if err != nil {
		return ir.StructField{}, err
	}

	return field, nil
}

func (visitor *Visitor) VisitEnum(schema *ir.Schema, def ir.Type) (ir.Type, error) {
	if visitor.OnEnum != nil {
		return visitor.OnEnum(visitor, schema, def)
	}

	return def, nil
}

func (visitor *Visitor) VisitScalar(schema *ir.Schema, def ir.Type) (ir.Type, error) {
	if visitor.OnScalar != nil {
		return visitor.OnScalar(visitor, schema, def)
	}

	return def, nil
}

func (visitor *Visitor) VisitRef(schema *ir.Schema, def ir.Type) (ir.Type, error) {
	if visitor.OnRef != nil {
		return visitor.OnRef(visitor, schema, def)
	}

	return def, nil
}

func (visitor *Visitor) VisitConsantRef(schema *ir.Schema, def ir.Type) (ir.Type, error) {
	if visitor.OnConstantRef != nil {
		return visitor.OnConstantRef(visitor, schema, def)
	}

	return def, nil
}
