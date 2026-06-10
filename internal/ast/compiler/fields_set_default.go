package compiler

import (
	"fmt"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/tools"
)

var _ Pass = (*FieldsSetDefault)(nil)

// FieldsSetDefault sets the default value for the given fields.
type FieldsSetDefault struct {
	DefaultValues map[FieldReference]any
	fieldsFound   []string
}

func (pass *FieldsSetDefault) Process(schemas []*ast.Schema) ([]*ast.Schema, error) {
	pass.fieldsFound = nil

	visitor := &Visitor{
		OnObject: pass.processObject,
	}

	return visitor.VisitSchemas(schemas)
}

func (pass *FieldsSetDefault) processObject(_ *Visitor, _ *ast.Schema, object ast.Object) (ast.Object, error) {
	if !object.Type.IsStruct() {
		return object, nil
	}

	for i, field := range object.Type.AsStruct().Fields {
		for fieldRef, value := range pass.DefaultValues {
			if !fieldRef.Matches(object, field) {
				continue
			}

			field.Type.Default = value
			field.AddToPassesTrail(fmt.Sprintf("FieldsSetDefault[default=%v]", value))

			object.Type.Struct.Fields[i] = field
			pass.fieldsFound = append(pass.fieldsFound, fieldRef.String())
		}
	}

	return object, nil
}

func (pass *FieldsSetDefault) Diagnostics() []string {
	if len(pass.fieldsFound) == len(pass.DefaultValues) {
		return nil
	}

	expected := make([]string, 0, len(pass.DefaultValues))
	for fieldRef := range pass.DefaultValues {
		expected = append(expected, fieldRef.String())
	}

	return tools.Map(tools.SliceFindMissing(pass.fieldsFound, expected), func(ref string) string {
		return fmt.Sprintf("field not found '%s'", ref)
	})
}
