package compiler

import (
	"github.com/grafana/cog/internal/ast"
)

var _ Pass = (*AddFormatUint8)(nil)

// AddFormatUint8 adds a `// +format=int32` comment to fields of type uint8 or []uint8.
// This influences Kubernetes OpenAPI generation.
type AddFormatUint8 struct {
}

func (pass *AddFormatUint8) Process(schemas []*ast.Schema) ([]*ast.Schema, error) {
	visitor := &Visitor{
		OnObject: pass.processObject,
	}

	return visitor.VisitSchemas(schemas)
}

func (pass *AddFormatUint8) processObject(visitor *Visitor, schema *ast.Schema, object ast.Object) (ast.Object, error) {
	// We only care about structs
	if object.Type.Kind != ast.KindStruct {
		return object, nil
	}

	processedFields := make([]ast.StructField, 0, len(object.Type.Struct.Fields))
	needsProcessing := false
	for _, field := range object.Type.Struct.Fields {
		newField := field
		if pass.shouldAddFieldComment(field.Type) {
			newField.Comments = append(newField.Comments, "+format=int32")
			needsProcessing = true
		}
		processedFields = append(processedFields, newField)
	}

	// Only update the object if we actually modified something
	if needsProcessing {
		object.Type.Struct.Fields = processedFields
		object.AddToPassesTrail("AddFormatInt32ForUint8")
	}

	return object, nil
}

func (pass *AddFormatUint8) shouldAddFieldComment(fieldType ast.Type) bool {
	if fieldType.IsScalar() && fieldType.AsScalar().ScalarKind == ast.KindUint8 {
		return true
	}

	if fieldType.IsArray() && fieldType.AsArray().ValueType.IsScalar() && fieldType.AsArray().ValueType.AsScalar().ScalarKind == ast.KindUint8 {
		return true
	}

	return false
}
