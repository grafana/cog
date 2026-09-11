package transforms

import (
	"testing"

	"github.com/grafana/cog/pkg/ir"
)

func TestDisjunctionOfAnonymousStructsToExplicit(t *testing.T) {
	// Prepare test input
	objects := []ir.Object{
		ir.NewObject("test", "disjunctionOfThings", ir.NewDisjunction([]ir.Type{
			ir.NewRef("test", "someStruct"),
			ir.NewStruct(
				ir.NewStructField("Type", ir.String(ir.Value("anonymous-struct"))),
				ir.NewStructField("FieldFoo", ir.String()),
			),
		})),

		ir.NewObject("test", "someStruct", ir.NewStruct(
			ir.NewStructField("Type", ir.String(ir.Value("some-struct"))),
			ir.NewStructField("FieldFoo", ir.String()),
		)),
	}

	// Prepare expected output
	expectedObjects := []ir.Object{
		ir.NewObject("test", "disjunctionOfThings", ir.NewDisjunction([]ir.Type{
			ir.NewRef("test", "someStruct"),
			ir.NewRef("test", "TypeAnonymousStruct"),
		})),

		objects[1],

		ir.NewObject("test", "TypeAnonymousStruct", ir.NewStruct(
			ir.NewStructField("Type", ir.String(ir.Value("anonymous-struct"))),
			ir.NewStructField("FieldFoo", ir.String()),
		)),
	}

	// Call the compiler pass
	runPassOnObjects(t, &DisjunctionOfAnonymousStructsToExplicit{}, objects, expectedObjects)
}
