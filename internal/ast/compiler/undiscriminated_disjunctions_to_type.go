package compiler

import (
	"strings"

	"github.com/grafana/cog/internal/ast"
)

var _ Pass = (*UndiscriminatedDisjunctionToType)(nil)

// UndiscriminatedDisjunctionToType turns undiscriminated disjunctions of refs
// into a struct type, with a nullable field for each branch of the disjunction.
// The struct is annotated with the HintUndiscriminatedDisjunctionOfRefs hint.
//
// Disjunctions of scalars are not impacted, disjunctions having a configured
// discriminator field and mapping are not impacted (see DisjunctionInferMapping).
// Note: this pass _should_ run after DisjunctionInferMapping.
type UndiscriminatedDisjunctionToType struct {
}

func (pass *UndiscriminatedDisjunctionToType) Process(schemas []*ast.Schema) ([]*ast.Schema, error) {
	visitor := &Visitor{
		OnDisjunction: pass.processDisjunction,
	}

	return visitor.VisitSchemas(schemas)
}

func (pass *UndiscriminatedDisjunctionToType) processDisjunction(visitor *Visitor, schema *ast.Schema, def ast.Type) (ast.Type, error) {
	disjunction := def.AsDisjunction()

	// Only process undiscriminated disjunctions of refs
	if !disjunction.Branches.HasOnlyRefs() {
		return def, nil
	}

	// If a discriminator is set, let DisjunctionToType handle it
	if len(disjunction.Discriminator) > 0 {
		return def, nil
	}

	newTypeName := pass.disjunctionTypeName(disjunction)

	// if we already generated a new object for this disjunction, return a reference to it.
	if visitor.HasNewObject(ast.RefType{ReferredPkg: schema.Package, ReferredType: newTypeName}) {
		ref := ast.NewRef(schema.Package, newTypeName, ast.Hints(def.Hints))
		ref.AddToPassesTrail("UndiscriminatedDisjunctionToType[disjunction → ref]")
		if def.Nullable || disjunction.Branches.HasNullType() {
			ref.Nullable = true
		}

		return ref, nil
	}

	fields := make([]ast.StructField, 0, len(disjunction.Branches))
	for _, branch := range disjunction.Branches {
		// Handled below, by allowing the reference to the disjunction struct to be null.
		if branch.IsNull() {
			continue
		}

		processedBranch := branch
		processedBranch.Nullable = true

		fields = append(fields, ast.NewStructField(ast.TypeName(processedBranch), processedBranch))
	}

	structType := ast.NewStruct(fields...)
	structType.Hints[ast.HintUndiscriminatedDisjunctionOfRefs] = disjunction

	newObject := ast.NewObject(schema.Package, newTypeName, structType)
	newObject.AddToPassesTrail("UndiscriminatedDisjunctionToType[created]")

	visitor.RegisterNewObject(newObject)

	ref := ast.NewRef(schema.Package, newTypeName, ast.Hints(def.Hints))
	ref.AddToPassesTrail("UndiscriminatedDisjunctionToType[disjunction → ref]")
	if def.Nullable || disjunction.Branches.HasNullType() {
		ref.Nullable = true
	}

	return ref, nil
}

func (pass *UndiscriminatedDisjunctionToType) disjunctionTypeName(def ast.DisjunctionType) string {
	parts := make([]string, 0, len(def.Branches))

	for _, subType := range def.Branches {
		parts = append(parts, ast.TypeName(subType))
	}

	return strings.Join(parts, "Or")
}
