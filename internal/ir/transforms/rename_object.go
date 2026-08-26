package transforms

import (
	"fmt"

	"github.com/grafana/cog/internal/ir"
)

var _ Transform = (*RenameObject)(nil)

type RenameObject struct {
	From        ObjectReference
	To          string
	objectFound bool
}

func (pass *RenameObject) Process(schemas ir.Schemas) (ir.Schemas, error) {
	pass.objectFound = false

	visitor := &Visitor{
		OnObject: pass.processObject,
		OnRef:    pass.processRef,
	}

	for _, schema := range schemas {
		if schema.Package == pass.From.Package && schema.EntryPoint == pass.From.Object {
			schema.EntryPoint = pass.To
		}
	}

	return visitor.VisitSchemas(schemas)
}

func (pass *RenameObject) processObject(visitor *Visitor, schema *ir.Schema, object ir.Object) (ir.Object, error) {
	var err error

	if pass.From.Matches(object) {
		pass.objectFound = true

		originalName := object.Name
		object.Name = pass.To
		object.SelfRef.ReferredType = pass.To
		object.AddToPassesTrail(fmt.Sprintf("RenameObject[%s → %s]", originalName, object.Name))
	}

	object.Type, err = visitor.VisitType(schema, object.Type)
	if err != nil {
		return ir.Object{}, err
	}

	return object, nil
}

func (pass *RenameObject) processRef(_ *Visitor, _ *ir.Schema, def ir.Type) (ir.Type, error) {
	if def.Ref.ReferredPkg == pass.From.Package && def.Ref.ReferredType == pass.From.Object {
		def.Ref.ReferredType = pass.To
	}

	return def, nil
}

func (pass *RenameObject) Diagnostics() []string {
	if pass.objectFound {
		return nil
	}

	return []string{
		fmt.Sprintf("object '%s' not found", pass.From),
	}
}
