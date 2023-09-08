package compiler

import (
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/tools"
)

var _ Pass = (*AnonymousEnumToExplicitType)(nil)

type AnonymousEnumToExplicitType struct {
	newObjects []ast.Object
}

func (pass *AnonymousEnumToExplicitType) Process(files []*ast.File) ([]*ast.File, error) {
	newFiles := make([]*ast.File, 0, len(files))

	for _, file := range files {
		newFile, err := pass.processFile(file)
		if err != nil {
			return nil, err
		}

		newFiles = append(newFiles, newFile)
	}

	return newFiles, nil
}

func (pass *AnonymousEnumToExplicitType) processFile(file *ast.File) (*ast.File, error) {
	pass.newObjects = nil

	processedObjects := make([]ast.Object, 0, len(file.Definitions))
	for _, object := range file.Definitions {
		processedObjects = append(processedObjects, pass.processObject(object))
	}

	return &ast.File{
		Package:     file.Package,
		Definitions: append(processedObjects, pass.newObjects...),
	}, nil
}

func (pass *AnonymousEnumToExplicitType) processObject(object ast.Object) ast.Object {
	if object.Type.Kind() == ast.KindEnum {
		return object
	}

	newObject := object
	newObject.Type = pass.processType(object.Name, object.Type)

	return newObject
}

func (pass *AnonymousEnumToExplicitType) processType(parentName string, def ast.Type) ast.Type {
	if def.Kind() == ast.KindArray {
		return pass.processArray(parentName, def.(ast.ArrayType))
	}

	if def.Kind() == ast.KindStruct {
		return pass.processStruct(def.(ast.StructType))
	}

	if def.Kind() == ast.KindEnum {
		return pass.processAnonymousEnum(parentName, def.(ast.EnumType))
	}

	// TODO: do the same for disjunctions?

	return def
}

func (pass *AnonymousEnumToExplicitType) processArray(parentName string, def ast.ArrayType) ast.ArrayType {
	return ast.ArrayType{
		ValueType: pass.processType(parentName, def.ValueType),
	}
}

func (pass *AnonymousEnumToExplicitType) processStruct(def ast.StructType) ast.StructType {
	newDef := def

	processedFields := make([]ast.StructField, 0, len(def.Fields))
	for _, field := range def.Fields {
		processedFields = append(processedFields, ast.StructField{
			Name:     field.Name,
			Comments: field.Comments,
			Type:     pass.processType(field.Name, field.Type),
			Required: field.Required,
			Default:  field.Default,
		})
	}

	newDef.Fields = processedFields

	return newDef
}

func (pass *AnonymousEnumToExplicitType) processAnonymousEnum(parentName string, def ast.EnumType) ast.RefType {
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
		Type: ast.EnumType{
			Values: values,
		},
	})

	return ast.RefType{
		ReferredType: enumTypeName,
	}
}
