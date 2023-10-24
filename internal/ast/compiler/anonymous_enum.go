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

	newSchema := schema.DeepCopy()
	for i, object := range schema.Objects {
		newSchema.Objects[i] = pass.processObject(object)
	}

	newSchema.Objects = append(newSchema.Objects, pass.newObjects...)

	return &newSchema, nil
}

func (pass *AnonymousEnumToExplicitType) processObject(object ast.Object) ast.Object {
	if object.Type.Kind == ast.KindEnum {
		return object
	}

	object.Type = pass.processType(object.Name, tools.UpperCamelCase(object.Name)+"Enum", object.Type)

	return object
}

func (pass *AnonymousEnumToExplicitType) processType(currentObjectName string, suggestedEnumName string, def ast.Type) ast.Type {
	if def.Kind == ast.KindArray {
		return pass.processArray(currentObjectName, suggestedEnumName, def)
	}

	if def.Kind == ast.KindStruct {
		return pass.processStruct(currentObjectName, def)
	}

	if def.Kind == ast.KindEnum {
		return pass.processAnonymousEnum(suggestedEnumName, def.AsEnum())
	}

	// TODO: do the same for disjunctions?

	return def
}

func (pass *AnonymousEnumToExplicitType) processArray(currentObjectName string, suggestedEnumName string, def ast.Type) ast.Type {
	def.Array.ValueType = pass.processType(currentObjectName, suggestedEnumName, def.Array.ValueType)

	return def
}

func (pass *AnonymousEnumToExplicitType) processStruct(parentName string, def ast.Type) ast.Type {
	for i, field := range def.Struct.Fields {
		newField := field
		newField.Type = pass.processType(parentName, tools.UpperCamelCase(parentName)+tools.UpperCamelCase(field.Name), field.Type)

		def.Struct.Fields[i] = newField
	}

	return def
}

func (pass *AnonymousEnumToExplicitType) processAnonymousEnum(parentName string, def ast.EnumType) ast.Type {
	enumTypeName := tools.UpperCamelCase(parentName)

	values := make([]ast.EnumValue, 0, len(def.Values))
	for _, val := range def.Values {
		values = append(values, ast.EnumValue{
			Type:  val.Type,
			Name:  tools.UpperCamelCase(val.Name),
			Value: val.Value,
		})
	}

	pass.newObjects = append(pass.newObjects, ast.Object{
		Name: enumTypeName,
		Type: ast.NewEnum(values),
	})

	return ast.NewRef(pass.currentPackage, enumTypeName)
}
