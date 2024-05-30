package compiler

import (
	"fmt"

	"github.com/grafana/cog/internal/ast"
)

var _ Pass = (*PrefixObjectNames)(nil)

// PrefixObjectNames appends the given comment to every object definition.
type PrefixObjectNames struct {
	Prefix string
}

func (pass *PrefixObjectNames) Process(schemas []*ast.Schema) ([]*ast.Schema, error) {
	if pass.Prefix == "" {
		return schemas, nil
	}

	visitor := &Visitor{
		OnObject: pass.processObject,
	}

	return visitor.VisitSchemas(schemas)
}

func (pass *PrefixObjectNames) processObject(_ *Visitor, _ *ast.Schema, object ast.Object) (ast.Object, error) {
	originalName := object.Name
	object.Name = pass.Prefix + originalName
	object.AddToPassesTrail(fmt.Sprintf("PrefixObjectNames[%s → %s]", originalName, object.Name))

	return object, nil
}
