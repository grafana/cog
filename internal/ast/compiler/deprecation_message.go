package compiler

import "github.com/grafana/cog/internal/ast"

var _ Pass = (*DeprecationMessage)(nil)

type DeprecationMessage struct {
	Object  ObjectReference
	Message string
}

func (pass *DeprecationMessage) Process(schemas []*ast.Schema) ([]*ast.Schema, error) {
	visitor := &Visitor{
		OnObject: pass.processObject,
	}

	return visitor.VisitSchemas(schemas)
}

func (pass *DeprecationMessage) processObject(_ *Visitor, _ *ast.Schema, object ast.Object) (ast.Object, error) {
	if !pass.Object.Matches(object) {
		return object, nil
	}

	object.DeprecationMessage = pass.Message
	return object, nil
}
