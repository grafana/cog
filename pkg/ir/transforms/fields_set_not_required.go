package transforms

import (
	"fmt"

	"github.com/grafana/cog/internal/tools"
	"github.com/grafana/cog/pkg/ir"
)

var _ Transform = (*FieldsSetNotRequired)(nil)

// FieldsSetNotRequired rewrites the definition of given fields to mark them as nullable and not required.
type FieldsSetNotRequired struct {
	Fields      []FieldReference
	fieldsFound []string
}

func (pass *FieldsSetNotRequired) Process(schemas ir.Schemas) (ir.Schemas, error) {
	pass.fieldsFound = nil

	visitor := &Visitor{
		OnObject: pass.processObject,
	}

	return visitor.VisitSchemas(schemas)
}

func (pass *FieldsSetNotRequired) processObject(_ *Visitor, _ *ir.Schema, object ir.Object) (ir.Object, error) {
	if !object.Type.IsStruct() {
		return object, nil
	}

	for i, field := range object.Type.AsStruct().Fields {
		for _, fieldRef := range pass.Fields {
			if !fieldRef.Matches(object, field) {
				continue
			}

			field.Type.Nullable = true
			field.Required = false
			field.AddToPassesTrail("FieldsSetNotRequired[nullable=true, required=false]")

			object.Type.Struct.Fields[i] = field
			pass.fieldsFound = append(pass.fieldsFound, fieldRef.String())
		}
	}

	return object, nil
}

func (pass *FieldsSetNotRequired) Diagnostics() []string {
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
