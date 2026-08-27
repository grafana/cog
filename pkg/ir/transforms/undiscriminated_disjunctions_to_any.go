package transforms

import (
	"github.com/grafana/cog/pkg/ir"
)

var _ Transform = (*UndiscriminatedDisjunctionToAny)(nil)

// UndiscriminatedDisjunctionToAny turns any undiscriminated disjunction into
// the `any` type.
// Disjunctions of scalars are not impacted, disjunctions having a configured
// discriminator field and mapping are not impacted (see DisjunctionInferMapping).
// Note: this pass _should_ run after DisjunctionInferMapping.
type UndiscriminatedDisjunctionToAny struct {
	GenerateUndiscriminatedDisjunctions bool
}

func (pass *UndiscriminatedDisjunctionToAny) Process(schemas ir.Schemas) (ir.Schemas, error) {
	if pass.GenerateUndiscriminatedDisjunctions {
		return schemas, nil
	}

	visitor := &Visitor{
		OnDisjunction: pass.processDisjunction,
	}

	return visitor.VisitSchemas(schemas)
}

func (pass *UndiscriminatedDisjunctionToAny) processDisjunction(_ *Visitor, schema *ir.Schema, def ir.Type) (ir.Type, error) {
	disjunction := def.AsDisjunction()

	// Ex: "some concrete value" | "some other value" | string
	if pass.hasOnlySingleTypeScalars(schema, disjunction) {
		return def, nil
	}

	if disjunction.Branches.HasOnlyScalarOrArrayOrMap() {
		return def, nil
	}

	if disjunction.Branches.HasOnlyRefs() {
		if len(disjunction.Discriminator) == 0 || len(disjunction.DiscriminatorMapping) == 0 {
			return ir.Any(ir.Trail("UndiscriminatedDisjunctionToAny")), nil
		}
	}

	return def, nil
}

func (pass *UndiscriminatedDisjunctionToAny) hasOnlySingleTypeScalars(schema *ir.Schema, disjunction ir.DisjunctionType) bool {
	branches := disjunction.Branches

	if len(branches) == 0 {
		return false
	}

	firstBranchType, found := schema.Resolve(branches[0])
	if !found {
		return false
	}

	if !firstBranchType.IsScalar() {
		return false
	}

	scalarKind := firstBranchType.AsScalar().ScalarKind
	for _, t := range branches {
		resolvedType, found := schema.Resolve(t)
		if !found {
			return false
		}

		if !resolvedType.IsScalar() {
			return false
		}

		if resolvedType.AsScalar().ScalarKind != scalarKind {
			return false
		}
	}

	return true
}
