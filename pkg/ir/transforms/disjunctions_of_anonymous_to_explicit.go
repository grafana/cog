package transforms

import (
	"fmt"

	"github.com/grafana/cog/internal/tools"
	"github.com/grafana/cog/pkg/ir"
)

var _ Transform = (*DisjunctionOfAnonymousStructsToExplicit)(nil)

// DisjunctionOfAnonymousStructsToExplicit looks for anonymous structs used as
// branches of disjunctions and turns them into explicitly named types.
type DisjunctionOfAnonymousStructsToExplicit struct {
}

func (pass *DisjunctionOfAnonymousStructsToExplicit) Process(schemas ir.Schemas) (ir.Schemas, error) {
	visitor := &Visitor{
		OnDisjunction: pass.processDisjunction,
	}

	return visitor.VisitSchemas(schemas)
}

func (pass *DisjunctionOfAnonymousStructsToExplicit) processDisjunction(visitor *Visitor, schema *ir.Schema, def ir.Type) (ir.Type, error) {
	scalarCount := 0
	anonymousCount := 0
	for _, branch := range def.Disjunction.Branches {
		if branch.IsScalar() {
			scalarCount++
		} else if branch.IsStruct() {
			anonymousCount++
		}
	}

	if scalarCount == 1 && anonymousCount == 1 {
		return def, nil
	}

	for i, branch := range def.Disjunction.Branches {
		if !branch.IsStruct() {
			continue
		}

		branchName := pass.generateBranchName(branch, i)

		newType, err := visitor.VisitType(schema, branch)
		if err != nil {
			return ir.Type{}, err
		}

		newObject := ir.NewObject(schema.Package, branchName, newType)
		visitor.RegisterNewObject(newObject)

		def.Disjunction.Branches[i] = ir.NewRef(schema.Package, newObject.Name)
	}

	return def, nil
}

func (pass *DisjunctionOfAnonymousStructsToExplicit) generateBranchName(branch ir.Type, index int) string {
	for _, field := range branch.Struct.Fields {
		if field.Type.IsConcreteScalar() {
			val := fmt.Sprintf("%v", field.Type.Scalar.Value)
			return fmt.Sprintf("%s%s", tools.UpperCamelCase(field.Name), tools.UpperCamelCase(val))
		}
	}

	return fmt.Sprintf("Branch%d", index)
}
