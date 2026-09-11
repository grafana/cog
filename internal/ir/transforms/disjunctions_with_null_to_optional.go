package transforms

import (
	"fmt"

	"github.com/grafana/cog/internal/ir"
)

var _ Transform = (*DisjunctionWithNullToOptional)(nil)

// DisjunctionWithNullToOptional simplifies disjunctions with two branches, where one is `null`. For those,
// it transforms `type | null` into `*type` (optional, nullable reference to `type`).
//
// Example:
//
//	```
//	MaybeString: string | null
//	```
//
// Will become:
//
//	```
//	MaybeString?: string
//	```
type DisjunctionWithNullToOptional struct {
}

func (pass *DisjunctionWithNullToOptional) Process(schemas ir.Schemas) (ir.Schemas, error) {
	visitor := &Visitor{
		OnDisjunction: pass.processDisjunction,
	}

	return visitor.VisitSchemas(schemas)
}

func (pass *DisjunctionWithNullToOptional) processDisjunction(_ *Visitor, _ *ir.Schema, def ir.Type) (ir.Type, error) {
	disjunction := def.AsDisjunction()

	if len(disjunction.Branches) != 2 || !disjunction.Branches.HasNullType() {
		return def, nil
	}

	// type | null
	finalType := disjunction.Branches.NonNullTypes()[0]
	finalType.Nullable = true
	finalType.AddToPassesTrail(fmt.Sprintf("DisjunctionWithNullToOptional[%[1]s|null → %[1]s?]", ir.TypeName(finalType)))

	return finalType, nil
}
