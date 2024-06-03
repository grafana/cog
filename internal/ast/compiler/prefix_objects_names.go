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
		OnRef:    pass.processRef,
	}

	return visitor.VisitSchemas(schemas)
}

func (pass *PrefixObjectNames) processObject(visitor *Visitor, schema *ast.Schema, object ast.Object) (ast.Object, error) {
	var err error

	originalName := object.Name
	object.Name = pass.Prefix + originalName
	object.SelfRef.ReferredType = object.Name
	object.AddToPassesTrail(fmt.Sprintf("PrefixObjectNames[%s → %s]", originalName, object.Name))

	object.Type, err = visitor.VisitType(schema, object.Type)
	if err != nil {
		return ast.Object{}, err
	}

	return object, nil
}

func (pass *PrefixObjectNames) processRef(_ *Visitor, _ *ast.Schema, ref ast.Type) (ast.Type, error) {
	originalName := ref.Ref.ReferredType
	ref.Ref.ReferredType = pass.Prefix + originalName
	ref.AddToPassesTrail(fmt.Sprintf("PrefixObjectNames[%s → %s]", originalName, ref.Ref.ReferredType))

	return ref, nil
}
