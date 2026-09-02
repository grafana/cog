package transforms

import (
	"fmt"

	"github.com/grafana/cog/internal/ir"
)

var _ Transform = (*DisjunctionOfConstantsToEnum)(nil)

// DisjunctionOfConstantsToEnum transforms disjunction of constants into an enum.
//
// Example:
//
//	```
//	SomeDisjunction: "first" | "second"
//	```
//
// Will become:
//
//	```
//	SomeDisjunction enum {
//		first = "first"
//		second = "second"
//	}
//	```
type DisjunctionOfConstantsToEnum struct {
	schemas ir.Schemas
}

func (pass *DisjunctionOfConstantsToEnum) Process(schemas ir.Schemas) (ir.Schemas, error) {
	pass.schemas = schemas

	visitor := &Visitor{
		OnDisjunction: pass.processDisjunction,
	}

	return visitor.VisitSchemas(schemas)
}

func (pass *DisjunctionOfConstantsToEnum) processDisjunction(_ *Visitor, _ *ir.Schema, def ir.Type) (ir.Type, error) {
	if len(def.Disjunction.Branches) < 2 {
		return def, nil
	}

	var scalarKindCandidate *ir.ScalarKind
	isScalarValidEnumMember := func(scalar ir.ScalarType) bool {
		accepted := scalar.ScalarKind == ir.KindString || scalar.IsNumeric()
		if !accepted {
			return false
		}

		if scalarKindCandidate == nil {
			scalarKindCandidate = &scalar.ScalarKind
		}

		return *scalarKindCandidate == scalar.ScalarKind
	}

	valueToString := func(val any) string {
		if str, ok := val.(string); ok {
			return str
		}

		return fmt.Sprintf("%v", val)
	}

	var identifiedMembers []ir.EnumValue
	var resolvesToConcreteScalarsOnly func(typeDef ir.Type) bool
	resolvesToConcreteScalarsOnly = func(typeDef ir.Type) bool {
		resolved := pass.schemas.Resolve(typeDef)

		if resolved.IsConcreteScalar() {
			if isScalarValidEnumMember(*resolved.Scalar) {
				identifiedMembers = append(identifiedMembers, ir.EnumValue{
					Type:  ir.NewScalar(*scalarKindCandidate),
					Name:  valueToString(resolved.Scalar.Value),
					Value: resolved.Scalar.Value,
				})
				return true
			}

			return false
		}

		if resolved.IsDisjunction() {
			for _, branch := range resolved.Disjunction.Branches {
				if !resolvesToConcreteScalarsOnly(branch) {
					return false
				}
			}

			return true
		}

		if resolved.IsEnum() {
			for _, member := range resolved.Enum.Values {
				if !isScalarValidEnumMember(*member.Type.Scalar) {
					return false
				}

				identifiedMembers = append(identifiedMembers, member)
			}

			return true
		}

		return false
	}

	if !resolvesToConcreteScalarsOnly(def) {
		return def, nil
	}

	typeOpts := []ir.TypeOption{
		ir.Default(def.Default),
		ir.Trail("DisjunctionOfConstantsToEnum"),
	}
	if def.Nullable {
		typeOpts = append(typeOpts, ir.Nullable())
	}

	return ir.NewEnum(identifiedMembers, typeOpts...), nil
}
