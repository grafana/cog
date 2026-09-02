package transforms

import (
	"fmt"

	"github.com/grafana/cog/internal/ir"
)

var _ Transform = (*AddFields)(nil)

// AddFields rewrites the definition of an object to add new fields.
// Note: existing fields will not be overwritten.
type AddFields struct {
	Object         ObjectReference
	Fields         []ir.StructField
	objectFound    bool
	existingFields []string
}

func (pass *AddFields) Process(schemas ir.Schemas) (ir.Schemas, error) {
	pass.objectFound = false
	pass.existingFields = nil

	visitor := &Visitor{
		OnObject: pass.processObject,
	}

	return visitor.VisitSchemas(schemas)
}

func (pass *AddFields) processObject(_ *Visitor, _ *ir.Schema, object ir.Object) (ir.Object, error) {
	if !pass.Object.Matches(object) {
		return object, nil
	}

	if !object.Type.IsStruct() {
		return object, fmt.Errorf("cannot add fields to a non-struct object: %s", pass.Object.String())
	}

	for _, field := range pass.Fields {
		// let's be safe: if a field with the same name already exists, we do not overwrite it.
		if _, exists := object.Type.AsStruct().FieldByName(field.Name); exists {
			pass.existingFields = append(pass.existingFields, field.Name)
			continue
		}

		field.AddToPassesTrail("AddFields[created]")

		object.Type.Struct.Fields = append(object.Type.Struct.Fields, field)
	}

	pass.objectFound = true

	return object, nil
}

func (pass *AddFields) Diagnostics() []string {
	if pass.objectFound && len(pass.existingFields) == 0 {
		return nil
	}

	var diags []string

	if !pass.objectFound {
		diags = append(diags, fmt.Sprintf("object '%s' not found", pass.Object))
	}

	for _, field := range pass.existingFields {
		diags = append(diags, fmt.Sprintf("object '%s' already has a '%s' field", pass.Object, field))
	}

	return diags
}
