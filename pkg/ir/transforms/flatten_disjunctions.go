package transforms

import (
	"fmt"

	"github.com/grafana/cog/pkg/ir"
)

var _ Transform = (*FlattenDisjunctions)(nil)

// FlattenDisjunctions will traverse all the branches every given disjunctions
// and, for each disjunction it finds, flatten it into the top-level type.
//
// Example:
//
//	```
//	SomeStruct: {
//		foo: string
//	}
//	OtherStruct: {
//		bar: string
//	}
//	LastStruct: {
//		hello: string
//	}
//	SomeOrOther: SomeStruct | OtherStruct
//	AnyStruct: SomeOrOther | LastStruct
//	```
//
// Will become:
//
//	```
//	SomeStruct: {
//		foo: string
//	}
//	OtherStruct: {
//		bar: string
//	}
//	LastStruct: {
//		hello: string
//	}
//	SomeOrOther: SomeStruct | OtherStruct
//	AnyStruct: SomeStruct | OtherStruct | LastStruct # this disjunction has been flattened
//	```
type FlattenDisjunctions struct {
}

func (pass *FlattenDisjunctions) Process(schemas ir.Schemas) (ir.Schemas, error) {
	visitor := &Visitor{
		OnDisjunction: pass.processDisjunction,
	}

	return visitor.VisitSchemas(schemas)
}

func (pass *FlattenDisjunctions) processDisjunction(_ *Visitor, schema *ir.Schema, def ir.Type) (ir.Type, error) {
	def.Disjunction = pass.flattenDisjunction(schema, def.AsDisjunction())

	return def, nil
}

func (pass *FlattenDisjunctions) flattenDisjunction(schema *ir.Schema, disjunction ir.DisjunctionType) *ir.DisjunctionType {
	newDisjunction := disjunction.DeepCopy()
	newDisjunction.Branches = nil

	branchMap := make(map[string]struct{})
	addBranch := func(typeName string, typeDef ir.Type) {
		if _, exists := branchMap[typeName]; exists {
			return
		}

		branchMap[typeName] = struct{}{}
		newDisjunction.Branches = append(newDisjunction.Branches, typeDef)
	}

	for i, branch := range disjunction.Branches {
		typeName := ir.TypeName(branch)
		if branch.IsStruct() {
			typeName = fmt.Sprintf("branch_%d", i)
		}

		if branch.IsConcreteScalar() {
			typeName = fmt.Sprintf("concrete_%s_%v", typeName, branch.Scalar.Value)
		}

		if !branch.IsRef() {
			addBranch(typeName, branch)
			continue
		}

		resolved, found := schema.Resolve(branch)
		if !found {
			// FIXME: error here?
			continue
		}

		if !resolved.IsDisjunction() {
			addBranch(typeName, branch)
			continue
		}

		for innerI, resolvedBranch := range resolved.AsDisjunction().Branches {
			innerTypeName := ir.TypeName(resolvedBranch)
			if branch.IsStruct() {
				innerTypeName = fmt.Sprintf("inner_branch_%d", innerI)
			}
			addBranch(innerTypeName, resolvedBranch)
		}
	}

	return &newDisjunction
}
