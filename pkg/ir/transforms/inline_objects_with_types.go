package transforms

import (
	"fmt"

	"github.com/grafana/cog/internal/orderedmap"
	"github.com/grafana/cog/pkg/ir"
)

var _ Transform = (*InlineObjectsWithTypes)(nil)

// InlineObjectsWithTypes inlines objects of the given types.
// This compiler pass is meant to be used to generate code in languages that
// don't support type aliases on scalars, top-level disjunctions, ...
//
// Note: constants are not inlined.
//
// Example:
//
//	```
//	TimeZone: string
//	Details: map[string, any]
//	Targets: []string
//
//	Foo: {
//	  TimezoneField: TimeZone
//	  DetailsField: Details
//	  TargetsField: Targets
//	}
//	```
//
// Will become:
//
//	```
//	Foo: {
//	  TimezoneField: string
//	  DetailsField: map[string, any]
//	  TargetsField: []string
//	}
//	```
type InlineObjectsWithTypes struct {
	InlineTypes     []ir.Kind
	objectsToInline *orderedmap.Map[string, ir.Type]
}

func (pass *InlineObjectsWithTypes) Process(schemas ir.Schemas) (ir.Schemas, error) {
	pass.objectsToInline = orderedmap.New[string, ir.Type]()

	for _, schema := range schemas {
		schema.Objects.Iterate(func(_ string, object ir.Object) {
			// follow potential references
			resolvedType := schemas.Resolve(object.Type)

			if !resolvedType.IsAnyOf(pass.InlineTypes...) {
				return
			}

			// do not inline constants
			if object.Type.IsConcreteScalar() {
				return
			}

			pass.objectsToInline.Set(object.SelfRef.String(), resolvedType)
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
		newSchemas[i].Objects = schema.Objects.Filter(func(_ string, object ir.Object) bool {
			return !pass.objectsToInline.Has(object.SelfRef.String())
		})
	}

	return newSchemas, nil
}

func (pass *InlineObjectsWithTypes) processRef(_ *Visitor, _ *ir.Schema, def ir.Type) (ir.Type, error) {
	if !pass.objectsToInline.Has(def.Ref.String()) {
		return def, nil
	}

	typeDef := pass.objectsToInline.Get(def.Ref.String()).DeepCopy()
	typeDef.AddToPassesTrail(fmt.Sprintf("InlineObjectsWithTypes[original=%s]", def.Ref.String()))

	return typeDef, nil
}
