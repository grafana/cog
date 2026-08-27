package transforms

import (
	"fmt"

	"github.com/grafana/cog/internal/tools"
	"github.com/grafana/cog/pkg/ir"
)

var _ Transform = (*FieldsSetRequired)(nil)

// FieldsSetRequired rewrites the definition of given fields to mark them as not nullable and required.
type FieldsSetRequired struct {
	Fields      []FieldReference
	fieldsFound []string
}

func (pass *FieldsSetRequired) Process(schemas ir.Schemas) (ir.Schemas, error) {
	pass.fieldsFound = nil

	visitor := &Visitor{
		OnObject: pass.processObject,
	}

	return visitor.VisitSchemas(schemas)
}

func (pass *FieldsSetRequired) processObject(_ *Visitor, _ *ir.Schema, object ir.Object) (ir.Object, error) {
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
