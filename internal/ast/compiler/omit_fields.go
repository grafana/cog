package compiler

import (
	"fmt"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/tools"
)

var _ Pass = (*OmitFields)(nil)

// OmitFields removes the selected fields from their object definition.
type OmitFields struct {
	Fields      []FieldReference
	fieldsFound []string
}

func (pass *OmitFields) Process(schemas []*ast.Schema) ([]*ast.Schema, error) {
	pass.fieldsFound = nil

	visitor := &Visitor{
		OnObject: pass.processObject,
	}

	return visitor.VisitSchemas(schemas)
}

func (pass *OmitFields) processObject(_ *Visitor, _ *ast.Schema, object ast.Object) (ast.Object, error) {
	if !object.Type.IsStruct() {
		return object, nil
	}

	object.Type.Struct.Fields = tools.Filter(object.Type.Struct.Fields, func(field ast.StructField) bool {
		for _, fieldRef := range pass.Fields {
			if fieldRef.Matches(object, field) {
				pass.fieldsFound = append(pass.fieldsFound, fieldRef.String())
				return false
			}
		}

		return true
	})

	return object, nil
}

func (pass *OmitFields) Diagnostics() []string {
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
