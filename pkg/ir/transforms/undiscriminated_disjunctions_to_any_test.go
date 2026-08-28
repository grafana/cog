package transforms

import (
	"testing"

	"github.com/grafana/cog/pkg/ir"
)

func TestUndiscriminatedDisjunctionToAny(t *testing.T) {
	// Prepare test input
	disjunctionTypeNoMapping := ir.NewDisjunction([]ir.Type{
		ir.NewRef("test", "SomeStruct"),
		ir.NewRef("test", "OtherStruct"),
	})
	disjunctionTypeNoMapping.Disjunction.Discriminator = "Type"
	disjunctionTypeMapping := ir.NewDisjunction([]ir.Type{
		ir.NewRef("test", "SomeStruct"),
		ir.NewRef("test", "OtherStruct"),
	})
	disjunctionTypeMapping.Disjunction.Discriminator = "Type"
	disjunctionTypeMapping.Disjunction.DiscriminatorMapping = map[string]string{
		"some-struct":  "SomeStruct",
		"other-struct": "OtherStruct",
	}

	objects := []ir.Object{
		ir.NewObject("test", "ADisjunctionOfRefsNoMapping", disjunctionTypeNoMapping),
		ir.NewObject("test", "ADisjunctionOfRefsMapping", disjunctionTypeMapping),
		ir.NewObject("test", "ADisjunctionOfScalars", ir.NewDisjunction([]ir.Type{
			ir.String(),
			ir.Bool(),
		})),

		ir.NewObject("test", "SomeStruct", ir.NewStruct(
			ir.NewStructField("Type", ir.String()), // Not a concrete scalar
			ir.NewStructField("FieldFoo", ir.String()),
		)),
		ir.NewObject("test", "OtherStruct", ir.NewStruct(
			ir.NewStructField("Type", ir.String(ir.Value("other-struct"))),
			ir.NewStructField("FieldBar", ir.Bool()),
		)),
	}

	disjunctionOfRefNoMapping := objects[0].DeepCopy()
	disjunctionOfRefNoMapping.Type = ir.Any(ir.Trail("UndiscriminatedDisjunctionToAny"))
	expected := []ir.Object{
		disjunctionOfRefNoMapping,
		objects[1],
		objects[2],
		objects[3],
		objects[4],
	}

	runPassOnObjects(t, &UndiscriminatedDisjunctionToAny{}, objects, expected)
}

func TestUndiscriminatedDisjunctionToAny_WhenGenerateEnabled_HasNoImpact(t *testing.T) {
	disjunctionType := ir.NewDisjunction([]ir.Type{
		ir.NewRef("test", "SomeStruct"),
		ir.NewRef("test", "OtherStruct"),
	})

	objects := []ir.Object{
		ir.NewObject("test", "ADisjunctionOfRefs", disjunctionType),
		ir.NewObject("test", "SomeStruct", ir.NewStruct(
			ir.NewStructField("FieldFoo", ir.String()),
		)),
		ir.NewObject("test", "OtherStruct", ir.NewStruct(
			ir.NewStructField("FieldBar", ir.Bool()),
		)),
	}

	runPassOnObjects(t, &UndiscriminatedDisjunctionToAny{GenerateUndiscriminatedDisjunctions: true}, objects, objects)
}
