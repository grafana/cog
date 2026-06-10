package compiler

import (
	"fmt"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/tools"
)

var _ Pass = (*FieldsSetRequired)(nil)

// FieldsSetRequired rewrites the definition of given fields to mark them as not nullable and required.
type FieldsSetRequired struct {
	Fields      []FieldReference
	fieldsFound []string
}

func (pass *FieldsSetRequired) Process(schemas []*ast.Schema) ([]*ast.Schema, error) {
	pass.fieldsFound = nil

	visitor := &Visitor{
		OnObject: pass.processObject,
	}

	return visitor.VisitSchemas(schemas)
}

func (pass *FieldsSetRequired) processObject(_ *Visitor, _ *ast.Schema, object ast.Object) (ast.Object, error) {
	if !object.Type.IsStruct() {
		return object, nil
	}

	for i, field := range object.Type.AsStruct().Fields {
		for _, fieldRef := range pass.Fields {
			if !fieldRef.Matches(object, field) {
				continue
			}

			field.Type.Nullable = false
			field.Required = true
			field.AddToPassesTrail("FieldsSetRequired[nullable=false, required=true]")

			object.Type.Struct.Fields[i] = field
			pass.fieldsFound = append(pass.fieldsFound, fieldRef.String())
		}
	}

	return object, nil
}

func (pass *FieldsSetRequired) Diagnostics() []string {
	if len(pass.fieldsFound) == len(pass.Fields) {
		return nil
	}

	expected := tools.Map(pass.Fields, func(ref FieldReference) string {
		return ref.String()
	})
	missing := tools.SliceFindMissing(pass.fieldsFound, expected)

	return tools.Map(missing, func(ref string) string {
		return fmt.Sprintf("field not found '%s'", ref)
	})
}
