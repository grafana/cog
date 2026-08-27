package transforms

import (
	"fmt"

	"github.com/grafana/cog/internal/ir"
	"github.com/grafana/cog/internal/tools"
)

var _ Transform = (*FieldsSetDefault)(nil)

// FieldsSetDefault sets the default value for the given fields.
// Invalid field references will be ignored and existing default values will be
// overridden.
type FieldsSetDefault struct {
	DefaultValues map[FieldReference]any
	fieldsFound   []string
}

func (pass *FieldsSetDefault) Process(schemas ir.Schemas) (ir.Schemas, error) {
	pass.fieldsFound = nil

	visitor := &Visitor{
		OnObject: pass.processObject,
	}

	return visitor.VisitSchemas(schemas)
}

func (pass *FieldsSetDefault) processObject(_ *Visitor, _ *ir.Schema, object ir.Object) (ir.Object, error) {
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
