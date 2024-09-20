package compiler

import (
	"fmt"

	"github.com/grafana/cog/internal/ast"
)

type TypedDefaults struct{}

func (t *TypedDefaults) Process(schemas []*ast.Schema) ([]*ast.Schema, error) {
	visitor := Visitor{
		OnScalar: t.processScalar,
		OnEnum:   t.processEnum,
	}

	return visitor.VisitSchemas(schemas)
}

func (t *TypedDefaults) processScalar(_ *Visitor, _ *ast.Schema, def ast.Type) (ast.Type, error) {
	if def.Default == nil {
		return def, nil
	}

	typedDef := ast.NewScalar(def.AsScalar().ScalarKind, ast.Default(def.Default))
	def.TypedDefault = &typedDef
	return def, nil
}

func (t *TypedDefaults) processEnum(_ *Visitor, _ *ast.Schema, def ast.Type) (ast.Type, error) {
	if def.Default != nil {
		return def, nil
	}

	var enumValue ast.EnumValue
	valueFound := false
	for _, v := range def.AsEnum().Values {
		if v.Value == def.Default {
			enumValue = ast.EnumValue{
				Type:  v.Type,
				Name:  v.Name,
				Value: v.Value,
			}
			valueFound = true
		}
	}

	if !valueFound {
		return ast.Type{}, fmt.Errorf("could not find enum value for %s", def.Default)
	}

	typedDef := ast.NewEnum([]ast.EnumValue{enumValue}, ast.Default(def.Default))
	def.TypedDefault = &typedDef
	return def, nil
}
