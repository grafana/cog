package transforms

import (
	"testing"

	"github.com/grafana/cog/internal/testutils"
	"github.com/grafana/cog/pkg/ir"
	"github.com/stretchr/testify/require"
)

func TestDisjunctionToType_WithNonDisjunctionObjects_HasNoImpact(t *testing.T) {
	// Prepare test input
	objects := []ir.Object{
		ir.NewObject("test", "AMap", ir.NewMap(ir.String(), ir.String())),
		ir.NewObject("test", "ARef", ir.NewRef("test", "AMap")),
		ir.NewObject("test", "AnEnum", ir.NewEnum([]ir.EnumValue{
			{
				Name:  "Foo",
				Type:  ir.String(),
				Value: "foo",
			},
			{
				Name:  "Bar",
				Type:  ir.String(),
				Value: "bar",
			},
		})),
		ir.NewObject("test", "AnArray", ir.NewArray(ir.String())),
		ir.NewObject("test", "AScalar", ir.NewScalar(ir.KindInt8)),
		ir.NewObject("test", "AStruct", ir.NewStruct(
			ir.NewStructField("SomeNonDisjunctionField", ir.NewScalar(ir.KindInt8)),
		)),
	}

	// Call the compiler pass
	runPassOnObjects(t, &DisjunctionToType{}, objects, objects)
}

func TestDisjunctionToType_WithDisjunctionOfScalars_AsAnObject(t *testing.T) {
	// Prepare test input
	objects := []ir.Object{
		ir.NewObject("test", "ADisjunctionOfScalars", ir.NewDisjunction([]ir.Type{
			ir.String(),
			ir.Bool(),
		})),
	}

	// Prepare expected output
	disjunctionStructType := ir.NewStruct(
		ir.NewStructField("String", ir.String(ir.Nullable())),
		ir.NewStructField("Bool", ir.Bool(ir.Nullable())),
	)
	// The original disjunction definition is preserved as a hint
	disjunctionStructType.Hints[ir.HintDisjunctionOfScalars] = objects[0].Type.AsDisjunction()

	expectedObjects := []ir.Object{
		ir.NewObject("test", "ADisjunctionOfScalars", ir.NewRef("test", "StringOrBool", ir.Trail("DisjunctionToType[disjunction → ref]"))),
		ir.NewObject("test", "StringOrBool", disjunctionStructType, "DisjunctionToType[created]"),
	}

	// Call the compiler pass
	runPassOnObjects(t, &DisjunctionToType{}, objects, expectedObjects)
}

func TestDisjunctionToType_WithDisjunctionOfScalars_AsAMapValueType(t *testing.T) {
	// Prepare test input
	objects := []ir.Object{
		ir.NewObject("test", "ADisjunctionOfScalars", ir.NewMap(
			ir.String(),
			ir.NewDisjunction([]ir.Type{
				ir.String(),
				ir.Bool(),
			}),
		)),
	}

	// Prepare expected output
	disjunctionStructType := ir.NewStruct(
		ir.NewStructField("String", ir.String(ir.Nullable())),
		ir.NewStructField("Bool", ir.Bool(ir.Nullable())),
	)
	// The original disjunction definition is preserved as a hint
	disjunctionStructType.Hints[ir.HintDisjunctionOfScalars] = objects[0].Type.AsMap().ValueType.AsDisjunction()

	expectedObjects := []ir.Object{
		ir.NewObject("test", "ADisjunctionOfScalars", ir.NewMap(
			ir.String(),
			ir.NewRef("test", "StringOrBool", ir.Trail("DisjunctionToType[disjunction → ref]")),
		)),
		ir.NewObject("test", "StringOrBool", disjunctionStructType, "DisjunctionToType[created]"),
	}

	// Call the compiler pass
	runPassOnObjects(t, &DisjunctionToType{}, objects, expectedObjects)
}

func TestDisjunctionToType_WithDisjunctionOfScalars_AsAStructField(t *testing.T) {
	// Prepare test input
	disjunctionType := ir.NewDisjunction([]ir.Type{
		ir.String(),
		ir.Bool(),
	})
	objects := []ir.Object{
		ir.NewObject("test", "AStructWithADisjunctionOfScalars", ir.NewStruct(
			ir.NewStructField("AFieldWithADisjunctionOfScalars", disjunctionType),
		)),
	}

	// Prepare expected output
	disjunctionStructType := ir.NewStruct(
		ir.NewStructField("String", ir.String(ir.Nullable())),
		ir.NewStructField("Bool", ir.Bool(ir.Nullable())),
	)
	// The original disjunction definition is preserved as a hint
	disjunctionStructType.Hints[ir.HintDisjunctionOfScalars] = disjunctionType.AsDisjunction()

	expectedObjects := []ir.Object{
		ir.NewObject("test", "AStructWithADisjunctionOfScalars", ir.NewStruct(
			ir.NewStructField("AFieldWithADisjunctionOfScalars", ir.NewRef("test", "StringOrBool", ir.Trail("DisjunctionToType[disjunction → ref]"))),
		)),
		ir.NewObject("test", "StringOrBool", disjunctionStructType, "DisjunctionToType[created]"),
	}

	// Call the compiler pass
	runPassOnObjects(t, &DisjunctionToType{}, objects, expectedObjects)
}

func TestDisjunctionToType_WithDisjunctionOfScalars_AsNullableAStructField(t *testing.T) {
	// Prepare test input
	disjunctionType := ir.NewDisjunction([]ir.Type{
		ir.String(),
		ir.Bool(),
	}, ir.Nullable())
	objects := []ir.Object{
		ir.NewObject("test", "AStructWithADisjunctionOfScalars", ir.NewStruct(
			ir.NewStructField("AFieldWithADisjunctionOfScalars", disjunctionType),
		)),
	}

	// Prepare expected output
	disjunctionStructType := ir.NewStruct(
		ir.NewStructField("String", ir.String(ir.Nullable())),
		ir.NewStructField("Bool", ir.Bool(ir.Nullable())),
	)
	// The original disjunction definition is preserved as a hint
	disjunctionStructType.Hints[ir.HintDisjunctionOfScalars] = disjunctionType.AsDisjunction()

	expectedObjects := []ir.Object{
		ir.NewObject("test", "AStructWithADisjunctionOfScalars", ir.NewStruct(
			ir.NewStructField("AFieldWithADisjunctionOfScalars", ir.NewRef("test", "StringOrBool", ir.Nullable(), ir.Trail("DisjunctionToType[disjunction → ref]"))),
		)),
		ir.NewObject("test", "StringOrBool", disjunctionStructType, "DisjunctionToType[created]"),
	}

	// Call the compiler pass
	runPassOnObjects(t, &DisjunctionToType{}, objects, expectedObjects)
}

func TestDisjunctionToType_WithDisjunctionOfScalars_AsAnArrayValueType(t *testing.T) {
	// Prepare test input
	disjunctionType := ir.NewDisjunction([]ir.Type{
		ir.String(),
		ir.Bool(),
	})
	objects := []ir.Object{
		ir.NewObject("test", "AnArrayWithADisjunctionOfScalars", ir.NewArray(disjunctionType)),
	}

	// Prepare expected output
	disjunctionStructType := ir.NewStruct(
		ir.NewStructField("String", ir.String(ir.Nullable())),
		ir.NewStructField("Bool", ir.Bool(ir.Nullable())),
	)
	// The original disjunction definition is preserved as a hint
	disjunctionStructType.Hints[ir.HintDisjunctionOfScalars] = disjunctionType.AsDisjunction()

	expectedObjects := []ir.Object{
		ir.NewObject("test", "AnArrayWithADisjunctionOfScalars", ir.NewArray(ir.NewRef("test", "StringOrBool", ir.Trail("DisjunctionToType[disjunction → ref]")))),
		ir.NewObject("test", "StringOrBool", disjunctionStructType, "DisjunctionToType[created]"),
	}

	// Call the compiler pass
	runPassOnObjects(t, &DisjunctionToType{}, objects, expectedObjects)
}

func TestDisjunctionToType_WithSingleTypeScalarsOnly(t *testing.T) {
	disjunctionType := ir.NewDisjunction([]ir.Type{
		ir.String(ir.Value("foo")),
		ir.String(ir.Value("bar")),
		ir.String(),
	})
	objects := []ir.Object{
		ir.NewObject("test", "ADisjunctionOfSingleTypeScalars", disjunctionType),
	}

	expectedObjects := []ir.Object{
		ir.NewObject("test", "ADisjunctionOfSingleTypeScalars", ir.NewScalar(ir.KindString)),
	}

	runPassOnObjects(t, &DisjunctionToType{}, objects, expectedObjects)
}

func TestDisjunctionToType_WithMixedDisjunction(t *testing.T) {
	disjunctionType := ir.NewDisjunction([]ir.Type{
		ir.String(),
		ir.NewRef("test", "SomeStruct"),
	})

	objects := []ir.Object{
		ir.NewObject("test", "ADisjunctionOfScalarsAndRefs", disjunctionType),
		ir.NewObject("test", "SomeStruct", ir.NewStruct(
			ir.NewStructField("FieldFoo", ir.String()),
		)),
	}

	disjunctionStructType := ir.NewStruct(
		ir.NewStructField("String", ir.String(ir.Nullable())),
		ir.NewStructField("SomeStruct", ir.NewRef("test", "SomeStruct", ir.Nullable())),
	)
	disjunctionStructType.Hints[ir.HintDisjunctionOfScalarsAndRefs] = disjunctionType.AsDisjunction()

	expectedObjects := []ir.Object{
		ir.NewObject("test", "ADisjunctionOfScalarsAndRefs", ir.NewRef("test", "StringOrSomeStruct", ir.Trail("DisjunctionToType[disjunction → ref]"))),
		objects[1],
		ir.NewObject("test", "StringOrSomeStruct", disjunctionStructType, "DisjunctionToType[created]"),
	}

	runPassOnObjects(t, &DisjunctionToType{}, objects, expectedObjects)
}

func TestDisjunctionToType_WithUndiscriminatedDisjunctionOfRefs_GenerateEnabled(t *testing.T) {
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

	disjunctionStructType := ir.NewStruct(
		ir.NewStructField("SomeStruct", ir.NewRef("test", "SomeStruct", ir.Nullable())),
		ir.NewStructField("OtherStruct", ir.NewRef("test", "OtherStruct", ir.Nullable())),
	)
	disjunctionStructType.Hints[ir.HintUndiscriminatedDisjunctionOfRefs] = disjunctionType.AsDisjunction()

	expectedObjects := []ir.Object{
		ir.NewObject("test", "ADisjunctionOfRefs", ir.NewRef("test", "SomeStructOrOtherStruct", ir.Trail("DisjunctionToType[disjunction → ref]"))),
		objects[1],
		objects[2],
		ir.NewObject("test", "SomeStructOrOtherStruct", disjunctionStructType, "DisjunctionToType[created]"),
	}

	runPassOnObjects(t, &DisjunctionToType{GenerateUndiscriminatedDisjunctions: true}, objects, expectedObjects)
}

func TestDisjunctionToType_WithDisjunctionOfRefs_AsAnObject_NoDiscriminatorMetadata(t *testing.T) {
	req := require.New(t)

	// Prepare test input
	objects := testutils.ObjectsMap(
		ir.NewObject("test", "ADisjunctionOfRefs", ir.NewDisjunction([]ir.Type{
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
	)

	compilerPass := &DisjunctionToType{}
	_, err := compilerPass.Process(ir.Schemas{
		{Package: "test", Objects: objects},
	})
	req.Error(err)
	req.ErrorContains(err, "discriminator not set")
}

func TestDisjunctionToType_WithDisjunctionOfRefs_AsAnObject_WithDiscriminatorFieldSet(t *testing.T) {
	req := require.New(t)

	// Prepare test input
	disjunctionType := ir.NewDisjunction([]ir.Type{
		ir.NewRef("test", "SomeStruct"),
		ir.NewRef("test", "OtherStruct"),
	})
	// Add discriminator-related metadata to the disjunction
	// Mapping omitted: it will be inferred
	disjunctionType.Disjunction.Discriminator = "Kind"

	objects := testutils.ObjectsMap(
		ir.NewObject("test", "ADisjunctionOfRefs", disjunctionType),

		ir.NewObject("test", "SomeStruct", ir.NewStruct(
			ir.NewStructField("Type", ir.String(ir.Value("some-struct"))),
			ir.NewStructField("Kind", ir.String(ir.Value("some-kind"))),
			ir.NewStructField("FieldFoo", ir.String()),
		)),
		ir.NewObject("test", "OtherStruct", ir.NewStruct(
			ir.NewStructField("Type", ir.String(ir.Value("other-struct"))),
			ir.NewStructField("Kind", ir.String(ir.Value("other-kind"))),
			ir.NewStructField("FieldBar", ir.Bool()),
		)),
	)

	compilerPass := &DisjunctionToType{}
	_, err := compilerPass.Process(ir.Schemas{
		{Package: "test", Objects: objects},
	})
	req.Error(err)
	req.ErrorContains(err, "discriminator mapping not set")
}

func TestDisjunctionToType_WithDisjunctionOfRefs_AsAnObject_WithDiscriminatorFieldAndMappingSet(t *testing.T) {
	// Prepare test input
	disjunctionType := ir.NewDisjunction([]ir.Type{
		ir.NewRef("test", "SomeStruct"),
		ir.NewRef("test", "OtherStruct"),
	})
	// Add discriminator-related metadata to the disjunction
	disjunctionType.Disjunction.Discriminator = "Kind"
	disjunctionType.Disjunction.DiscriminatorMapping = map[string]string{
		"other-kind": "OtherStruct",
		"some-kind":  "SomeStruct",
	}

	objects := []ir.Object{
		ir.NewObject("test", "ADisjunctionOfRefs", disjunctionType),

		ir.NewObject("test", "SomeStruct", ir.NewStruct(
			ir.NewStructField("Type", ir.String(ir.Value("some-struct"))),
			ir.NewStructField("Kind", ir.String(ir.Value("some-kind"))),
			ir.NewStructField("FieldFoo", ir.String()),
		)),
		ir.NewObject("test", "OtherStruct", ir.NewStruct(
			ir.NewStructField("Type", ir.String(ir.Value("other-struct"))),
			ir.NewStructField("Kind", ir.String(ir.Value("other-kind"))),
			ir.NewStructField("FieldBar", ir.Bool()),
		)),
	}

	// Prepare expected output
	disjunctionStructType := ir.NewStruct(
		ir.NewStructField("SomeStruct", ir.NewRef("test", "SomeStruct", ir.Nullable())),
		ir.NewStructField("OtherStruct", ir.NewRef("test", "OtherStruct", ir.Nullable())),
	)
	// The original disjunction definition is preserved as a hint
	disjunctionTypeWithDiscriminatorMeta := objects[0].Type.AsDisjunction()

	// Metadata should be inferred
	disjunctionTypeWithDiscriminatorMeta.Discriminator = "Kind"
	disjunctionTypeWithDiscriminatorMeta.DiscriminatorMapping = map[string]string{
		"other-kind": "OtherStruct",
		"some-kind":  "SomeStruct",
	}
	disjunctionStructType.Hints[ir.HintDiscriminatedDisjunctionOfRefs] = disjunctionTypeWithDiscriminatorMeta

	expectedObjects := []ir.Object{
		ir.NewObject("test", "ADisjunctionOfRefs", ir.NewRef("test", "SomeStructOrOtherStruct", ir.Trail("DisjunctionToType[disjunction → ref]"))),
		objects[1],
		objects[2],
		ir.NewObject("test", "SomeStructOrOtherStruct", disjunctionStructType, "DisjunctionToType[created]"),
	}

	// Call the compiler pass
	runPassOnObjects(t, &DisjunctionToType{}, objects, expectedObjects)
}
