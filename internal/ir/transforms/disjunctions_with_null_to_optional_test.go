package transforms

import (
	"testing"

	"github.com/grafana/cog/internal/ir"
)

func TestDisjunctionWithNullToOptional_WithDisjunctionOfTypeAndNull_AsAnObject(t *testing.T) {
	// Prepare test input
	objects := []ir.Object{
		ir.NewObject("test", "ScalarWithNull", ir.NewDisjunction([]ir.Type{
			ir.String(),
			ir.Null(),
		})),
		ir.NewObject("test", "RefWithNull", ir.NewDisjunction([]ir.Type{
			ir.NewRef("test", "SomeType"),
			ir.Null(),
		})),
	}

	expectedObjects := []ir.Object{
		ir.NewObject("test", "ScalarWithNull", ir.String(ir.Nullable(), ir.Trail("DisjunctionWithNullToOptional[String|null → String?]"))),
		ir.NewObject("test", "RefWithNull", ir.NewRef("test", "SomeType", ir.Nullable(), ir.Trail("DisjunctionWithNullToOptional[SomeType|null → SomeType?]"))),
	}

	// Call the compiler pass
	runPassOnObjects(t, &DisjunctionWithNullToOptional{}, objects, expectedObjects)
}

func TestDisjunctionWithNullToOptional_WithDisjunctionOfTypeAndNull_AsAStructField(t *testing.T) {
	// Prepare test input
	objects := []ir.Object{
		ir.NewObject("test", "StructWithScalarWithNull", ir.NewStruct(
			ir.NewStructField("Field", ir.NewDisjunction([]ir.Type{
				ir.String(),
				ir.Null(),
			})),
		)),
		ir.NewObject("test", "StructWithRefWithNull", ir.NewStruct(
			ir.NewStructField("Field", ir.NewDisjunction([]ir.Type{
				ir.NewRef("test", "SomeType"),
				ir.Null(),
			})),
		)),
	}

	expectedObjects := []ir.Object{
		ir.NewObject("test", "StructWithScalarWithNull", ir.NewStruct(
			ir.NewStructField("Field", ir.String(ir.Nullable(), ir.Trail("DisjunctionWithNullToOptional[String|null → String?]"))),
		)),
		ir.NewObject("test", "StructWithRefWithNull", ir.NewStruct(
			ir.NewStructField("Field", ir.NewRef("test", "SomeType", ir.Nullable(), ir.Trail("DisjunctionWithNullToOptional[SomeType|null → SomeType?]"))),
		)),
	}

	// Call the compiler pass
	runPassOnObjects(t, &DisjunctionWithNullToOptional{}, objects, expectedObjects)
}
