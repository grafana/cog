package compiler

import (
	"fmt"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/tools"
)

var _ Pass = (*ConstantToEnum)(nil)

// ConstantToEnum turns `string` constants into an enum definition with a
// single member.
// This is useful to "future-proof" a schema where a type can have a single
// value for now but is expected to allow more in the future.
type ConstantToEnum struct {
	Objects      ObjectReferences
	objectsFound []string
	warnings     []string
}

func (pass *ConstantToEnum) Process(schemas []*ast.Schema) ([]*ast.Schema, error) {
	pass.objectsFound = nil
	pass.warnings = nil

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
		pass.warnings = append(pass.warnings, fmt.Sprintf("object '%s' is not a concrete string", object.SelfRef))
		return object, nil
	}

	object.Type = ast.NewEnum([]ast.EnumValue{
		{
			Type:  ast.String(),
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
