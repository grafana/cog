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
		OnArray:  t.processArray,
		OnRef:    t.processRef,
	}

	return visitor.VisitSchemas(schemas)
}

func (t *TypedDefaults) processScalar(_ *Visitor, _ *ast.Schema, def ast.Type) (ast.Type, error) {
	if def.Default == nil {
		return def, nil
	}

	def.TypedDefault = &def
	return def, nil
}

func (t *TypedDefaults) processEnum(_ *Visitor, _ *ast.Schema, def ast.Type) (ast.Type, error) {
	if def.Default == nil {
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

	typedDef := ast.NewEnum([]ast.EnumValue{enumValue})
	def.TypedDefault = &typedDef
	return def, nil
}

func (t *TypedDefaults) processArray(visitor *Visitor, schema *ast.Schema, def ast.Type) (ast.Type, error) {
	if def.Default == nil {
		return def, nil
	}

	if def.AsArray().IsArrayOfScalars() {
		typedDef := ast.NewArray(def.AsArray().ValueType, ast.TypedDefault(def))
		def.TypedDefault = &typedDef
		return def, nil
	}

	valueType := def.AsArray().ValueType
	valueType.Default = def.Default

	var visitFunc func(schema *ast.Schema, def ast.Type) (ast.Type, error)

	switch valueType.Kind {
	case ast.KindEnum:
		visitFunc = visitor.VisitEnum
	case ast.KindRef:
		visitFunc = visitor.VisitRef
	case ast.KindMap:
		visitFunc = visitor.VisitMap
	case ast.KindArray:
		visitFunc = visitor.VisitArray
	case ast.KindScalar:
		visitFunc = visitor.VisitScalar
	case ast.KindStruct:
		visitFunc = visitor.VisitType
	case ast.KindIntersection:
		visitFunc = visitor.VisitIntersection
	case ast.KindDisjunction:
		visitFunc = visitor.VisitDisjunction
	}

	d, err := visitFunc(schema, def)
	if err != nil {
		return ast.Type{}, err
	}

	typedDef := ast.NewArray(def.AsArray().ValueType, ast.TypedDefault(d))
	def.TypedDefault = &typedDef

	return def, nil
}

func (t *TypedDefaults) processRef(_ *Visitor, _ *ast.Schema, def ast.Type) (ast.Type, error) {
	if def.Default == nil {
		return def, nil
	}

	ref := def.AsRef()
	typedDefault := ast.NewRef(ref.ReferredPkg, ref.ReferredType, ast.TypedDefault(def))
	def.TypedDefault = &typedDefault

	return def, nil
}

func (t *TypedDefaults) processMap(visitor *Visitor, schema *ast.Schema, def ast.Type) (ast.Type, error) {
	if def.Default == nil {
		return def, nil
	}

	fmt.Println(def)
	return def, nil
}
