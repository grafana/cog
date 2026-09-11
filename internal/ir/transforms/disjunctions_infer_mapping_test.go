package transforms

import (
	"testing"

	"github.com/grafana/cog/internal/ir"
)

func TestDisjunctionInferMapping_WithNonDisjunctionObjects_HasNoImpact(t *testing.T) {
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
	runPassOnObjects(t, &DisjunctionInferMapping{}, objects, objects)
}

func TestDisjunctionInferMapping_WithDisjunctionOfScalars_AsAnObject_hasNoImpact(t *testing.T) {
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

	// Call the compiler pass
	runPassOnObjects(t, &DisjunctionInferMapping{}, objects, objects)
}

func TestDisjunctionInferMapping_WithDisjunctionOfRefs_AsAnObject_NoDiscriminatorMetadata_NoDiscriminatorFieldCandidate(t *testing.T) {
	// Prepare test input
	objects := []ir.Object{
		ir.NewObject("test", "ADisjunctionOfRefs", ir.NewDisjunction([]ir.Type{
			ir.NewRef("test", "SomeStruct"),
			ir.NewRef("test", "OtherStruct"),
		})),

		ir.NewObject("test", "SomeStruct", ir.NewStruct(
			ir.NewStructField("Kind", ir.String(ir.Value("some-struct"))), // No equivalent in OtherStruct
			ir.NewStructField("FieldFoo", ir.String()),
		)),
		ir.NewObject("test", "OtherStruct", ir.NewStruct(
			ir.NewStructField("Type", ir.String(ir.Value("other-struct"))),
			ir.NewStructField("FieldBar", ir.Bool()),
		)),
	}

	disjunctionOfRef := objects[0].DeepCopy()
	disjunctionOfRef.Type.PassesTrail = []string{"DisjunctionInferMapping[no_mapping_found:could not identify discriminator field]"}
	expected := []ir.Object{
		disjunctionOfRef,
		objects[1],
		objects[2],
	}

	runPassOnObjects(t, &DisjunctionInferMapping{}, objects, expected)
}

func TestDisjunctionInferMapping_WithDisjunctionOfRefs_AsAnObject_NoDiscriminatorMetadata_NonScalarDiscriminator_NonConstantReference(t *testing.T) {
	// Prepare test input
	disjunctionType := ir.NewDisjunction([]ir.Type{
		ir.NewRef("test", "SomeStruct"),
		ir.NewRef("test", "OtherStruct"),
	})
	disjunctionType.Disjunction.Discriminator = "MapOfString"

	objects := []ir.Object{
		ir.NewObject("test", "ADisjunctionOfRefs", disjunctionType),

		ir.NewObject("test", "SomeStruct", ir.NewStruct(
			ir.NewStructField("FieldFoo", ir.String()),
			ir.NewStructField("MapOfString", ir.NewMap(ir.String(), ir.String())),
		)),
		ir.NewObject("test", "OtherStruct", ir.NewStruct(
			ir.NewStructField("FieldBar", ir.Bool()),
			ir.NewStructField("MapOfString", ir.NewMap(ir.String(), ir.String())),
		)),
	}

	disjunctionOfRef := objects[0].DeepCopy()
	disjunctionOfRef.Type.PassesTrail = []string{"DisjunctionInferMapping[no_mapping_found:discriminator field 'MapOfString' is not a scalar or constant reference]"}
	expected := []ir.Object{
		disjunctionOfRef,
		objects[1],
		objects[2],
	}

	runPassOnObjects(t, &DisjunctionInferMapping{}, objects, expected)
}

func TestDisjunctionInferMapping_WithDisjunctionOfRefs_AsAnObject_NoDiscriminatorMetadata_NonConcreteDiscriminator(t *testing.T) {
	// Prepare test input
	disjunctionType := ir.NewDisjunction([]ir.Type{
		ir.NewRef("test", "SomeStruct"),
		ir.NewRef("test", "OtherStruct"),
	})
	disjunctionType.Disjunction.Discriminator = "Type"

	objects := []ir.Object{
		ir.NewObject("test", "ADisjunctionOfRefs", disjunctionType),

		ir.NewObject("test", "SomeStruct", ir.NewStruct(
			ir.NewStructField("Type", ir.String()), // Not a concrete scalar
			ir.NewStructField("FieldFoo", ir.String()),
		)),
		ir.NewObject("test", "OtherStruct", ir.NewStruct(
			ir.NewStructField("Type", ir.String(ir.Value("other-struct"))),
			ir.NewStructField("FieldBar", ir.Bool()),
		)),
	}

	disjunctionOfRef := objects[0].DeepCopy()
	disjunctionOfRef.Type.PassesTrail = []string{"DisjunctionInferMapping[no_mapping_found:discriminator field 'Type' is not a scalar or constant reference]"}
	expected := []ir.Object{
		disjunctionOfRef,
		objects[1],
		objects[2],
	}

	runPassOnObjects(t, &DisjunctionInferMapping{}, objects, expected)
}

func TestDisjunctionInferMapping_WithDisjunctionOfRefs_AsAnObject_NoDiscriminatorMetadata_UnknownDiscriminatorField(t *testing.T) {
	// Prepare test input
	disjunctionType := ir.NewDisjunction([]ir.Type{
		ir.NewRef("test", "SomeStruct"),
		ir.NewRef("test", "OtherStruct"),
	})
	disjunctionType.Disjunction.Discriminator = "DoesNotExist"

	objects := []ir.Object{
		ir.NewObject("test", "ADisjunctionOfRefs", disjunctionType),

		ir.NewObject("test", "SomeStruct", ir.NewStruct(
			ir.NewStructField("Type", ir.String(ir.Value("some-struct"))),
			ir.NewStructField("FieldFoo", ir.String()),
		)),
		ir.NewObject("test", "OtherStruct", ir.NewStruct(
			ir.NewStructField("Type", ir.String(ir.Value("other-struct"))),
			ir.NewStructField("FieldBar", ir.Bool()),
		)),
	}

	disjunctionOfRef := objects[0].DeepCopy()
	disjunctionOfRef.Type.PassesTrail = []string{"DisjunctionInferMapping[no_mapping_found:discriminator field 'DoesNotExist' not found]"}
	expected := []ir.Object{
		disjunctionOfRef,
		objects[1],
		objects[2],
	}

	runPassOnObjects(t, &DisjunctionInferMapping{}, objects, expected)
}

func TestDisjunctionInferMapping_WithDisjunctionOfRefs_AsAnObject_NoDiscriminatorMetadata(t *testing.T) {
	// Prepare test input
	objects := []ir.Object{
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
	}

	// Prepare expected output
	newDisjunction := objects[0].DeepCopy()
	newDisjunction.Type.Disjunction.Discriminator = "Type"
	newDisjunction.Type.Disjunction.DiscriminatorMapping = map[string]string{
		"other-struct": "OtherStruct",
		"some-struct":  "SomeStruct",
	}

	expectedObjects := []ir.Object{
		newDisjunction,
		objects[1],
		objects[2],
	}

	// Call the compiler pass
	runPassOnObjects(t, &DisjunctionInferMapping{}, objects, expectedObjects)
}

func TestDisjunctionInferMapping_WithDisjunctionOfRefs_AsAnObject_Scalar_WithDiscriminatorFieldSet(t *testing.T) {
	// Prepare test input
	disjunctionType := ir.NewDisjunction([]ir.Type{
		ir.NewRef("test", "SomeStruct"),
		ir.NewRef("test", "OtherStruct"),
	})
	// Add discriminator-related metadata to the disjunction
	// Mapping omitted: it will be inferred
	disjunctionType.Disjunction.Discriminator = "Kind"

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
	newDisjunction := objects[0].DeepCopy()
	newDisjunction.Type.Disjunction.DiscriminatorMapping = map[string]string{
		"other-kind": "OtherStruct",
		"some-kind":  "SomeStruct",
	}

	expectedObjects := []ir.Object{
		newDisjunction,
		objects[1],
		objects[2],
	}

	// Call the compiler pass
	runPassOnObjects(t, &DisjunctionInferMapping{}, objects, expectedObjects)
}

func TestDisjunctionInferMapping_WithDisjunctionOfRefs_AsAnObject_Scalar_WithDiscriminatorFieldAndMappingSet(t *testing.T) {
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

	// Call the compiler pass
	runPassOnObjects(t, &DisjunctionInferMapping{}, objects, objects)
}

func TestDisjunctionInferMapping_WithDisjunctionOfRefs_AsAnObject_ConcreteReference_WithDiscriminatorFieldSet(t *testing.T) {
	// Prepare test input
	disjunctionType := ir.NewDisjunction([]ir.Type{
		ir.NewRef("test", "SomeStruct"),
		ir.NewRef("test", "OtherStruct"),
	})
	// Add discriminator-related metadata to the disjunction
	// Mapping omitted: it will be inferred
	disjunctionType.Disjunction.Discriminator = "Kind"

	objects := []ir.Object{
		ir.NewObject("test", "ADisjunctionOfRefs", disjunctionType),
		ir.NewObject("test", "AnEnum", ir.NewEnum([]ir.EnumValue{
			{Type: ir.String(), Name: "ValueA", Value: "a"},
			{Type: ir.String(), Name: "ValueB", Value: "b"},
		})),
		ir.NewObject("test", "SomeStruct", ir.NewStruct(
			ir.NewStructField("Type", ir.String(ir.Value("some-struct"))),
			ir.NewStructField("Kind", ir.NewConstantReferenceType("test", "AnEnum", "a")),
			ir.NewStructField("FieldFoo", ir.String()),
		)),
		ir.NewObject("test", "OtherStruct", ir.NewStruct(
			ir.NewStructField("Type", ir.String(ir.Value("other-struct"))),
			ir.NewStructField("Kind", ir.NewConstantReferenceType("test", "AnEnum", "b")),
			ir.NewStructField("FieldBar", ir.Bool()),
		)),
	}

	// Prepare expected output
	newDisjunction := objects[0].DeepCopy()
	newDisjunction.Type.Disjunction.DiscriminatorMapping = map[string]string{
		"b": "OtherStruct",
		"a": "SomeStruct",
	}

	expectedObjects := []ir.Object{
		newDisjunction,
		objects[1],
		objects[2],
		objects[3],
	}

	// Call the compiler pass
	runPassOnObjects(t, &DisjunctionInferMapping{}, objects, expectedObjects)
}

func TestDisjunctionInferMapping_WithDisjunctionOfRefs_AsAnObject_ConcreteReference_WithDiscriminatorFieldAndMappingSet(t *testing.T) {
	// Prepare test input
	disjunctionType := ir.NewDisjunction([]ir.Type{
		ir.NewRef("test", "SomeStruct"),
		ir.NewRef("test", "OtherStruct"),
	})
	// Add discriminator-related metadata to the disjunction
	disjunctionType.Disjunction.Discriminator = "Kind"
	disjunctionType.Disjunction.DiscriminatorMapping = map[string]string{
		"b": "OtherStruct",
		"a": "SomeStruct",
	}

	objects := []ir.Object{
		ir.NewObject("test", "ADisjunctionOfRefs", disjunctionType),
		ir.NewObject("test", "AnEnum", ir.NewEnum([]ir.EnumValue{
			{Type: ir.String(), Name: "ValueA", Value: "a"},
			{Type: ir.String(), Name: "ValueB", Value: "b"},
		})),
		ir.NewObject("test", "SomeStruct", ir.NewStruct(
			ir.NewStructField("Type", ir.String(ir.Value("some-struct"))),
			ir.NewStructField("Kind", ir.NewConstantReferenceType("test", "AnEnum", "a")),
			ir.NewStructField("FieldFoo", ir.String()),
		)),
		ir.NewObject("test", "OtherStruct", ir.NewStruct(
			ir.NewStructField("Type", ir.String(ir.Value("other-struct"))),
			ir.NewStructField("Kind", ir.NewConstantReferenceType("test", "AnEnum", "b")),
			ir.NewStructField("FieldBar", ir.Bool()),
		)),
	}

	// Call the compiler pass
	runPassOnObjects(t, &DisjunctionInferMapping{}, objects, objects)
}
