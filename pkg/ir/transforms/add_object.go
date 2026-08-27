package transforms

import (
	"fmt"

	"github.com/grafana/cog/pkg/ir"
)

var _ Transform = (*AddObject)(nil)

// AddObject adds a new object to a schema.
type AddObject struct {
	Object       ObjectReference
	As           ir.Type
	Comments     []string
	packageFound bool
}

func (pass *AddObject) Process(schemas ir.Schemas) (ir.Schemas, error) {
	pass.packageFound = false

	visitor := &Visitor{
		OnSchema: pass.processSchema,
	}

	return visitor.VisitSchemas(schemas)
}

func (pass *AddObject) processSchema(visitor *Visitor, schema *ir.Schema) (*ir.Schema, error) {
	if schema.Package != pass.Object.Package {
		return schema, nil
	}

	newObject := ir.NewObject(pass.Object.Package, pass.Object.Object, pass.As)
	newObject.Comments = pass.Comments
	newObject.AddToPassesTrail("AddObject[created]")

	visitor.RegisterNewObject(newObject)

	pass.packageFound = true

	return schema, nil
}

func (pass *AddObject) Diagnostics() []string {
	if pass.packageFound {
		return nil
	}

	return []string{
		fmt.Sprintf("package '%s' not found", pass.Object.Package),
	}
}
