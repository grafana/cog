package transforms

import (
	"fmt"

	"github.com/grafana/cog/internal/tools"
	"github.com/grafana/cog/pkg/ir"
)

var _ Transform = (*ConstantToEnum)(nil)

// ConstantToEnum turns `string` constants into an enum definition with a
// single member.
// This is useful to "future-proof" a schema where a type can have a single
// value for now but is expected to allow more in the future.
//
// Example:
//
//	```
//	ElementType: "thing"
//	```
//
// Will become:
//
//	```
//	ElementType enum {
//		thing = "thing"
//	}
//	```
type ConstantToEnum struct {
	Objects      ObjectReferences
	objectsFound []string
	warnings     []string
}

func (pass *ConstantToEnum) Process(schemas ir.Schemas) (ir.Schemas, error) {
	pass.objectsFound = nil
	pass.warnings = nil

	visitor := &Visitor{
		OnObject: pass.processObject,
	}

	return visitor.VisitSchemas(schemas)
}

func (pass *ConstantToEnum) processObject(_ *Visitor, _ *ir.Schema, object ir.Object) (ir.Object, error) {
	if !pass.Objects.Matches(object) {
		return object, nil
	}

	if !object.Type.IsConcreteScalar() || object.Type.Scalar.ScalarKind != ir.KindString {
		pass.warnings = append(pass.warnings, fmt.Sprintf("object '%s' is not a concrete string", object.SelfRef))
		return object, nil
	}

	object.Type = ir.NewEnum([]ir.EnumValue{
		{
			Type:  ir.String(),
			Name:  object.Type.Scalar.Value.(string),
			Value: object.Type.Scalar.Value.(string),
		},
	})
	object.AddToPassesTrail("ConstantToEnum")

	pass.objectsFound = append(pass.objectsFound, object.SelfRef.String())

	return object, nil
}

func (pass *ConstantToEnum) Diagnostics() []string {
	var diags []string
	diags = append(diags, pass.warnings...)

	if len(pass.objectsFound) == len(pass.Objects) {
		return diags
	}

	expected := tools.Map(pass.Objects, func(ref ObjectReference) string {
		return ref.String()
	})
	missing := tools.Map(tools.SliceFindMissing(pass.objectsFound, expected), func(ref string) string {
		return fmt.Sprintf("object not found '%s'", ref)
	})

	return append(diags, missing...)
}
