package compiler

import (
	"github.com/grafana/cog/internal/ast"
)

var _ Pass = (*ConstantToEnum)(nil)

type ConstantToEnum struct {
	Objects ObjectReferences
}

func (pass *ConstantToEnum) Process(schemas []*ast.Schema) ([]*ast.Schema, error) {
	visitor := &Visitor{
		OnObject: pass.processObject,
	}

	return visitor.VisitSchemas(schemas)
}

func (pass *ConstantToEnum) processObject(_ *Visitor, _ *ast.Schema, object ast.Object) (ast.Object, error) {
	if !pass.Objects.Matches(object) {
		return object, nil
	}

	if !object.Type.IsConcreteScalar() || object.Type.Scalar.ScalarKind != ast.KindString {
		return object, nil
	}

	trailMessage := "ConstantToEnum"

	object.Type = ast.NewEnum([]ast.EnumValue{
		{
			Type:  object.Type,
			Name:  object.Type.Scalar.Value.(string),
			Value: object.Type.Scalar.Value.(string),
		},
	})
	object.AddToPassesTrail(trailMessage)

	return object, nil
}
