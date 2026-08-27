package transforms

import (
	"testing"

	"github.com/grafana/cog/pkg/ir"
)

func TestFlattenDisjunctions_WithNestedDisjunctionOfRefs_AsAnObject(t *testing.T) {
	// Prepare test input
	objects := []ir.Object{
		ir.NewObject("test", "ADisjunctionOfRefs", ir.NewDisjunction([]ir.Type{
			ir.NewRef("test", "SomeOrOther"),
			ir.NewRef("test", "LastStruct"),
		})),

		ir.NewObject("test", "SomeOrOther", ir.NewDisjunction([]ir.Type{
			ir.NewRef("test", "SomeStruct"),
			ir.NewRef("test", "OtherStruct"),
		})),

		ir.NewObject("test", "SomeStruct", ir.NewStruct(
			ir.NewStructField("Type", ir.String(ir.Value("some-struct"))),
			ir.NewStructField("FieldFoo", ir.String()),
		)),
		ir.NewObject("test", "OtherStruct", ir.NewStruct(
			ir.NewStructField("FieldBar", ir.NewMap(ir.String(), ir.String())),
			ir.NewStructField("Type", ir.String(ir.Value("other-struct"))),
		)),
		ir.NewObject("test", "LastStruct", ir.NewStruct(
			ir.NewStructField("FieldLast", ir.NewMap(ir.String(), ir.String())),
			ir.NewStructField("Type", ir.String(ir.Value("last-struct"))),
		)),
	}

	// Prepare expected output
	expectedObjects := []ir.Object{
		ir.NewObject("test", "ADisjunctionOfRefs", ir.NewDisjunction([]ir.Type{
			ir.NewRef("test", "SomeStruct"),
			ir.NewRef("test", "OtherStruct"),
			ir.NewRef("test", "LastStruct"),
		})),

		objects[1],
		objects[2],
		objects[3],
		objects[4],
	}

	// Call the compiler pass
	runPassOnObjects(t, &FlattenDisjunctions{}, objects, expectedObjects)
}

func TestFlattenDisjunctions_WithDisjunctionOfStringAndConstants(t *testing.T) {
	// Prepare test input
	objects := []ir.Object{
		ir.NewObject("test", "ADisjunction", ir.NewDisjunction([]ir.Type{
			ir.String(),
			ir.String(ir.Value("*")),
			ir.String(ir.Value("none")),
		})),
	}

	// Call the compiler pass
	runPassOnObjects(t, &FlattenDisjunctions{}, objects, objects)
}

func TestFlattenDisjunctions_WithDisjunctionsOfAnonymousStructs(t *testing.T) {
	// Prepare test input
	objects := []ir.Object{
		ir.NewObject("test", "ADisjunctionOfStructs", ir.NewDisjunction([]ir.Type{
			ir.NewStruct(
				ir.NewStructField("Type", ir.String(ir.Value("root"))),
				ir.NewStructField("FieldRoot", ir.String()),
			),
			ir.NewRef("test", "SomeOrOther"),
			ir.NewStruct(
				ir.NewStructField("FieldLast", ir.Any()),
				ir.NewStructField("Type", ir.String(ir.Value("last-struct"))),
			),
		})),

		ir.NewObject("test", "SomeOrOther", ir.NewDisjunction([]ir.Type{
			ir.NewRef("test", "SomeStruct"),
			ir.NewRef("test", "OtherStruct"),
		})),

		ir.NewObject("test", "SomeStruct", ir.NewStruct(
			ir.NewStructField("Type", ir.String(ir.Value("some-struct"))),
			ir.NewStructField("FieldFoo", ir.String()),
		)),
		ir.NewObject("test", "OtherStruct", ir.NewStruct(
			ir.NewStructField("FieldBar", ir.NewMap(ir.String(), ir.String())),
			ir.NewStructField("Type", ir.String(ir.Value("other-struct"))),
		)),
	}

	// Prepare expected output
	expectedObjects := []ir.Object{
		ir.NewObject("test", "ADisjunctionOfStructs", ir.NewDisjunction([]ir.Type{
			ir.NewStruct(
				ir.NewStructField("Type", ir.String(ir.Value("root"))),
				ir.NewStructField("FieldRoot", ir.String()),
			),
			ir.NewRef("test", "SomeStruct"),
			ir.NewRef("test", "OtherStruct"),
			ir.NewStruct(
				ir.NewStructField("FieldLast", ir.Any()),
				ir.NewStructField("Type", ir.String(ir.Value("last-struct"))),
			),
		})),
		objects[1],
		objects[2],
		objects[3],
	}

	// Call the compiler pass
	runPassOnObjects(t, &FlattenDisjunctions{}, objects, expectedObjects)
}
