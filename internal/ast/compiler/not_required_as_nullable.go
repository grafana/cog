package compiler

import (
	"github.com/grafana/cog/internal/ast"
)

var _ Pass = (*NotRequiredFieldAsNullableType)(nil)

type NotRequiredFieldAsNullableType struct {
}

func (pass *NotRequiredFieldAsNullableType) Process(files []*ast.File) ([]*ast.File, error) {
	newFiles := make([]*ast.File, 0, len(files))

	for _, file := range files {
		newFiles = append(newFiles, pass.processFile(file))
	}

	return newFiles, nil
}

func (pass *NotRequiredFieldAsNullableType) processFile(file *ast.File) *ast.File {
	processedObjects := make([]ast.Object, 0, len(file.Definitions))
	for _, object := range file.Definitions {
		processedObjects = append(processedObjects, pass.processObject(object))
	}

	return &ast.File{
		Package:     file.Package,
		Definitions: processedObjects,
	}
}

func (pass *NotRequiredFieldAsNullableType) processObject(object ast.Object) ast.Object {
	if object.Type.Kind != ast.KindStruct {
		return object
	}

	newObject := object
	newObject.Type = pass.processType(object.Type)

	return newObject
}

func (pass *NotRequiredFieldAsNullableType) processType(def ast.Type) ast.Type {
	if def.Kind == ast.KindArray {
		return pass.processArray(def.AsArray())
	}

	if def.Kind == ast.KindMap {
		return pass.processMap(def.AsMap())
	}

	if def.Kind == ast.KindStruct {
		return pass.processStruct(def.AsStruct())
	}

	return def
}

func (pass *NotRequiredFieldAsNullableType) processArray(def ast.ArrayType) ast.Type {
	return ast.NewArray(pass.processType(def.ValueType))
}

func (pass *NotRequiredFieldAsNullableType) processMap(def ast.MapType) ast.Type {
	return ast.NewMap(
		pass.processType(def.IndexType),
		pass.processType(def.ValueType),
	)
}

func (pass *NotRequiredFieldAsNullableType) processStruct(def ast.StructType) ast.Type {
	processedFields := make([]ast.StructField, 0, len(def.Fields))
	for _, field := range def.Fields {
		fieldType := pass.processType(field.Type)
		if !field.Required {
			fieldType.Nullable = true
		}

		newField := field
		newField.Type = fieldType

		processedFields = append(processedFields, newField)
	}

	return ast.NewStruct(processedFields...)
}
