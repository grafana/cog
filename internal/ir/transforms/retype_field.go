package transforms

import (
	"fmt"

	"github.com/grafana/cog/internal/ir"
)

var _ Transform = (*RetypeField)(nil)

type RetypeField struct {
	Field      FieldReference
	As         ir.Type
	Comments   []string
	fieldFound bool
}

func (pass *RetypeField) Process(schemas ir.Schemas) (ir.Schemas, error) {
	pass.fieldFound = false

	visitor := &Visitor{
		OnObject: pass.processObject,
	}

	return visitor.VisitSchemas(schemas)
}

func (pass *RetypeField) processObject(_ *Visitor, _ *ir.Schema, object ir.Object) (ir.Object, error) {
	if !object.Type.IsStruct() {
		return object, nil
	}

	for i, field := range object.Type.Struct.Fields {
		if !pass.Field.Matches(object, field) {
			continue
		}

		object.Type.Struct.Fields[i].AddToPassesTrail(fmt.Sprintf("RetypeField[%s → %s]", ir.TypeName(field.Type), ir.TypeName(pass.As)))
		object.Type.Struct.Fields[i].Type = pass.As

		if pass.Comments != nil {
			object.Type.Struct.Fields[i].Comments = pass.Comments
		}

		pass.fieldFound = true

		break
	}

	return object, nil
}

func (pass *RetypeField) Diagnostics() []string {
	if pass.fieldFound {
		return nil
	}

	return []string{
		fmt.Sprintf("field '%s' not found", pass.Field),
	}
}
