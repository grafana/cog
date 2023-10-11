package compiler

import (
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/tools"
)

var _ Pass = (*AnonymousEnumToExplicitType)(nil)

type AnonymousEnumToExplicitType struct {
	newObjects     []ast.Object
	currentPackage string
}

func (pass *AnonymousEnumToExplicitType) Process(schemas []*ast.Schema) ([]*ast.Schema, error) {
	newSchemas := make([]*ast.Schema, 0, len(schemas))

	for _, schema := range schemas {
		newSchema, err := pass.processSchema(schema)
		if err != nil {
			return nil, err
		}

		newSchemas = append(newSchemas, newSchema)
	}

	return newSchemas, nil
}

func (pass *AnonymousEnumToExplicitType) processSchema(schema *ast.Schema) (*ast.Schema, error) {
	pass.newObjects = nil
	pass.currentPackage = schema.Package

	processedObjects := make([]ast.Object, 0, len(schema.Objects))
	for _, object := range schema.Objects {
		processedObjects = append(processedObjects, pass.processObject(object))
	}

	newSchema := schema.DeepCopy()
	newSchema.Objects = processedObjects
	newSchema.Objects = append(newSchema.Objects, pass.newObjects...)

	return &newSchema, nil
}

func (pass *AnonymousEnumToExplicitType) processObject(object ast.Object) ast.Object {
	if object.Type.Kind == ast.KindEnum {
		return object
	}

	newObject := object
	newObject.Type = pass.processType(object.Name, object.Type)

	return newObject
}

func (pass *AnonymousEnumToExplicitType) processType(parentName string, def ast.Type) ast.Type {
	if def.Kind == ast.KindArray {
		return pass.processArray(parentName, def.AsArray())
	}

	if def.Kind == ast.KindStruct {
		return pass.processStruct(def.AsStruct())
	}

	if def.Kind == ast.KindEnum {
		return pass.processAnonymousEnum(parentName, def.AsEnum())
	}

	// TODO: do the same for disjunctions?

	return def
}

func (pass *AnonymousEnumToExplicitType) processArray(parentName string, def ast.ArrayType) ast.Type {
	return ast.NewArray(pass.processType(parentName, def.ValueType))
}

func (pass *AnonymousEnumToExplicitType) processStruct(def ast.StructType) ast.Type {
	processedFields := make([]ast.StructField, 0, len(def.Fields))
	for _, field := range def.Fields {
		processedFields = append(processedFields, ast.StructField{
			Name:     field.Name,
			Comments: field.Comments,
			Type:     pass.processType(field.Name, field.Type),
			Required: field.Required,
		})
	}

	return ast.NewStruct(processedFields...)
}

func (pass *AnonymousEnumToExplicitType) processAnonymousEnum(parentName string, def ast.EnumType) ast.Type {
	enumTypeName := tools.UpperCamelCase(parentName) + "Enum"

	values := make([]ast.EnumValue, 0, len(def.Values))
	for _, val := range def.Values {
		values = append(values, ast.EnumValue{
			Type:  val.Type,
			Name:  parentName + tools.UpperCamelCase(val.Name),
			Value: val.Value,
		})
	}

	pass.newObjects = append(pass.newObjects, ast.Object{
		Name: enumTypeName,
		Type: ast.NewEnum(values),
	})

	return ast.NewRef(pass.currentPackage, enumTypeName)
}
