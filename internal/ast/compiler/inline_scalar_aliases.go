package compiler

import (
	"fmt"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/orderedmap"
)

var _ Pass = (*InlineScalarAliases)(nil)

// InlineScalarAliases inlines objects defined as an alias to a scalar type.
// Inlined types are: ast.KindScalar, ast.KindArray and ast.KindMap.
// This compiler pass is meant to be used to generate code in languages that
// don't support type aliases on scalars.
//
// Note: constants are not impacted.
//
// Example:
//
//	```
//	TimeZone string
//	Details map[string, any]
//	Targets []string
//
//	Foo struct {
//	  TimezoneField TimeZone
//	  DetailsField Details
//	  TargetsField Targets
//	}
//	```
//
// Will become:
//
//	```
//	Foo struct {
//	  TimezoneField string
//	  DetailsField map[string, any]
//	  TargetsField []string
//	}
//	```
type InlineScalarAliases struct {
	objectsToInline *orderedmap.Map[string, ast.Type]
}

func (pass *InlineScalarAliases) Process(schemas []*ast.Schema) ([]*ast.Schema, error) {
	pass.objectsToInline = orderedmap.New[string, ast.Type]()

	for _, schema := range schemas {
		schema.Objects.Iterate(func(_ string, object ast.Object) {
			if !object.Type.IsAnyOf(ast.KindScalar, ast.KindArray, ast.KindMap) {
				return
			}

			// do not inline constants
			if object.Type.IsConcreteScalar() {
				return
			}

			pass.objectsToInline.Set(object.SelfRef.String(), object.Type)
		})
	}

	visitor := &Visitor{
		OnRef: pass.processRef,
	}

	newSchemas, err := visitor.VisitSchemas(schemas)
	if err != nil {
		return nil, err
	}

	for i, schema := range newSchemas {
		newSchemas[i].Objects = schema.Objects.Filter(func(_ string, object ast.Object) bool {
			return !pass.objectsToInline.Has(object.SelfRef.String())
		})
	}

	return newSchemas, nil
}

func (pass *InlineScalarAliases) processRef(_ *Visitor, _ *ast.Schema, def ast.Type) (ast.Type, error) {
	if !pass.objectsToInline.Has(def.Ref.String()) {
		return def, nil
	}

	typeDef := pass.objectsToInline.Get(def.Ref.String()).DeepCopy()
	typeDef.AddToPassesTrail(fmt.Sprintf("InlineScalarAliases[original=%s]", def.Ref.String()))

	return typeDef, nil
}