package compiler

import (
	"fmt"
	"strings"

	"github.com/grafana/cog/internal/ast"
)

var _ Pass = (*FieldsSetRequired)(nil)

type FieldReference struct {
	Package string
	Object  string
	Field   string
}

func FieldReferenceFromString(ref string) (FieldReference, error) {
	parts := strings.Split(ref, ".")
	if len(parts) != 3 {
		return FieldReference{}, fmt.Errorf("invalid field reference '%s'", ref)
	}

	return FieldReference{
		Package: parts[0],
		Object:  parts[1],
		Field:   parts[2],
	}, nil
}

// FieldsSetRequired rewrites the definition of given fields to mark them as not nullable and required.
type FieldsSetRequired struct {
	Fields []FieldReference
}

func (pass *FieldsSetRequired) Process(schemas []*ast.Schema) ([]*ast.Schema, error) {
	for i, schema := range schemas {
		schemas[i] = pass.processSchema(schema)
	}

	return schemas, nil
}

func (pass *FieldsSetRequired) processSchema(schema *ast.Schema) *ast.Schema {
	schema.Objects = schema.Objects.Map(func(_ string, object ast.Object) ast.Object {
		return pass.processObject(schema, object)
	})

	return schema
}

func (pass *FieldsSetRequired) processObject(schema *ast.Schema, object ast.Object) ast.Object {
	if !object.Type.IsStruct() {
		return object
	}

	for _, fieldRef := range pass.Fields {
		if fieldRef.Package != schema.Package {
			continue
		}

		if fieldRef.Object != object.Name {
			continue
		}

		for i, field := range object.Type.AsStruct().Fields {
			if field.Name == fieldRef.Field {
				field.Type.Nullable = false
				field.Required = true

				object.Type.Struct.Fields[i] = field
			}
		}
	}

	return object
}
