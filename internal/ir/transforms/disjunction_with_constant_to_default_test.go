package transforms

import (
	"testing"

	"github.com/grafana/cog/internal/ir"
)

func TestDisjunctionWithConstantToDefault(t *testing.T) {
	// Prepare test input
	objects := []ir.Object{
		ir.NewObject("test", "DisjunctionWithScalarAndConst", ir.NewDisjunction([]ir.Type{
			ir.String(),
			ir.String(ir.Value("foo")),
		})),

		ir.NewObject("test", "DisjunctionWithDifferentKinds", ir.NewDisjunction([]ir.Type{
			ir.String(),
			ir.Bool(ir.Value(false)),
		})),

		ir.NewObject("test", "DisjunctionWithTwoScalarsAndConst", ir.NewDisjunction([]ir.Type{
			ir.String(),
			ir.Bool(),
			ir.String(ir.Value("foo")),
		})),

		ir.NewObject("test", "DisjunctionWithTwoScalars", ir.NewDisjunction([]ir.Type{
			ir.String(),
			ir.Bool(),
		})),

		ir.NewObject("test", "DisjunctionWithTwoScalarConsts", ir.NewDisjunction([]ir.Type{
			ir.String(ir.Value("bar")),
			ir.String(ir.Value("foo")),
		})),

		ir.NewObject("test", "ScalarObject", ir.String()),
	}

	// Prepare expected output
	expectedObjects := []ir.Object{
		ir.NewObject("test", "DisjunctionWithScalarAndConst", ir.String(ir.Default("foo"), ir.Trail("DisjunctionWithConstantToDefault"))),
		objects[1],
		objects[2],
		objects[3],
		objects[4],
		objects[5],
	}

	// Call the compiler pass
	runPassOnObjects(t, &DisjunctionWithConstantToDefault{}, objects, expectedObjects)
}
