package compiler

import (
	"fmt"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/tools"
)

var _ Pass = (*PrefixObjectNames)(nil)

// PrefixObjectNames adds the given prefix to every object's name.
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
		OnEnum:   pass.processEnum,
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

func (pass *PrefixObjectNames) processEnum(_ *Visitor, _ *ast.Schema, enum ast.Type) (ast.Type, error) {
	values := make([]ast.EnumValue, 0, len(enum.AsEnum().Values))
	for _, val := range enum.AsEnum().Values {
		values = append(values, ast.EnumValue{
			Type:  val.Type,
			Name:  tools.UpperCamelCase(pass.Prefix) + tools.UpperCamelCase(val.Name),
			Value: val.Value,
		})
	}

	enum.Enum.Values = values

	return enum, nil
}
