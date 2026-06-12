package compiler

import (
	"fmt"

	"github.com/grafana/cog/internal/ast"
)

var _ Pass = (*DeprecateObject)(nil)

// DeprecateObject marks an object as deprecated.
// Note: builders generated from this object will be marked as well.
type DeprecateObject struct {
	Object      ObjectReference
	Message     string
	objectFound bool
}

func (pass *DeprecateObject) Process(schemas []*ast.Schema) ([]*ast.Schema, error) {
	pass.objectFound = false

	visitor := &Visitor{
		OnObject: pass.processObject,
	}

	return visitor.VisitSchemas(schemas)
}

func (pass *DeprecateObject) processObject(_ *Visitor, _ *ast.Schema, object ast.Object) (ast.Object, error) {
	if !pass.Object.Matches(object) {
		return object, nil
	}

	object.DeprecationMessage = pass.Message
	object.AddToPassesTrail("DeprecateObject[added]")
	pass.objectFound = true

	return object, nil
}

func (pass *DeprecateObject) Diagnostics() []string {
	if pass.objectFound {
		return nil
	}

	return []string{
		fmt.Sprintf("object '%s' not found", pass.Object),
	}
}
