package transforms

import (
	"github.com/grafana/cog/internal/ir"
)

var _ Transform = (*NotRequiredFieldAsNullableType)(nil)

// NotRequiredFieldAsNullableType identifies all the struct fields marked as not `Required`
// and rewrites their `Type` to be `Nullable`.
type NotRequiredFieldAsNullableType struct {
}

func (pass *NotRequiredFieldAsNullableType) Process(schemas ir.Schemas) (ir.Schemas, error) {
	visitor := &Visitor{
		OnStructField: pass.processStructField,
	}

	return visitor.VisitSchemas(schemas)
}

func (pass *NotRequiredFieldAsNullableType) processStructField(visitor *Visitor, schema *ir.Schema, field ir.StructField) (ir.StructField, error) {
	var err error
	field.Type, err = visitor.VisitType(schema, field.Type)
	if err != nil {
		return field, err
	}

	if !field.Required && !field.Type.Nullable {
		field.Type.Nullable = true
		field.AddToPassesTrail("NotRequiredFieldAsNullableType[nullable=true]")
	}

	return field, nil
}
