package compiler

import (
	"errors"

	"github.com/grafana/cog/internal/ast"
)

type AddEnumValue struct {
	ObjectRef ObjectReference
	FieldRef  FieldReference
	Name      string
	Value     any
}

func (pass *AddEnumValue) Process(schemas []*ast.Schema) ([]*ast.Schema, error) {
	if pass.Name == "" || pass.Value == nil {
		return nil, errors.New("name and value are required")
	}

	visitor := &Visitor{
		OnObject: pass.onObject,
		OnEnum:   pass.onEnum,
	}
	return visitor.VisitSchemas(schemas)
}

func (pass *AddEnumValue) onObject(visitor *Visitor, schema *ast.Schema, object ast.Object) (ast.Object, error) {
	if object.Type.IsEnum() && pass.ObjectRef.Matches(object) {
		enum, err := visitor.OnEnum(visitor, schema, object.Type)
		if err != nil {
			return ast.Object{}, err
		}

		object.Type = enum
		object.AddToPassesTrail("AddEnumValue")
		return object, nil
	}

	if !object.Type.IsStruct() {
		return object, nil
	}

	for i, field := range object.Type.AsStruct().Fields {
		if pass.FieldRef.Matches(object, field) {
			if field.Type.IsEnum() {
				updatedType, err := visitor.OnEnum(visitor, schema, field.Type)
				if err != nil {
					return ast.Object{}, err
				}
				object.Type.AsStruct().Fields[i].Type = updatedType
				object.Type.AsStruct().Fields[i].AddToPassesTrail("AddEnumValue")
				return object, nil
			}

			if field.Type.IsRef() {
				if enum, ok := pass.updateEnumObject(visitor, schema, field.Type); ok {
					schema.Objects.Set(object.Name, enum)
				}
				return object, nil
			}
		}
	}

	return object, nil
}

func (pass *AddEnumValue) onEnum(_ *Visitor, _ *ast.Schema, def ast.Type) (ast.Type, error) {
	enumString := def.AsEnum().Values[0].Type.AsScalar().ScalarKind == ast.KindString
	_, isString := pass.Value.(string)

	if enumString && !isString {
		return ast.Type{}, errors.New("enum value must be of type string")
	}
	if !enumString && isString {
		return ast.Type{}, errors.New("enum value must be of type integer")
	}

	def.Enum.Values = append(def.Enum.Values, ast.EnumValue{
		Name:  pass.Name,
		Type:  def.AsEnum().Values[0].Type,
		Value: pass.Value,
	})

	return def, nil
}

func (pass *AddEnumValue) updateEnumObject(visitor *Visitor, schema *ast.Schema, def ast.Type) (ast.Object, bool) {
	obj, ok := schema.LocateObject(def.AsRef().ReferredType)
	if !ok {
		return ast.Object{}, false
	}

	if !obj.Type.IsEnum() {
		return ast.Object{}, false
	}

	enum, err := visitor.OnEnum(visitor, schema, obj.Type)
	if err != nil {
		return ast.Object{}, false
	}

	obj.Type = enum
	return obj, true
}
