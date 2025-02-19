package compiler

import (
	"testing"

	"github.com/grafana/cog/internal/ast"
)

func TestDisjunctionInferMapping_WithNonDisjunctionObjects_HasNoImpact(t *testing.T) {
	// Prepare test input
	objects := []ast.Object{
		ast.NewObject("test", "AMap", ast.NewMap(ast.String(), ast.String())),
		ast.NewObject("test", "ARef", ast.NewRef("test", "AMap")),
		ast.NewObject("test", "AnEnum", ast.NewEnum([]ast.EnumValue{
			{
				Name:  "Foo",
				Type:  ast.String(),
				Value: "foo",
			},
			{
				Name:  "Bar",
				Type:  ast.String(),
				Value: "bar",
			},
		})),
		ast.NewObject("test", "AnArray", ast.NewArray(ast.String())),
		ast.NewObject("test", "AScalar", ast.NewScalar(ast.KindInt8)),
		ast.NewObject("test", "AStruct", ast.NewStruct(
			ast.NewStructField("SomeNonDisjunctionField", ast.NewScalar(ast.KindInt8)),
		)),
	}

	// Call the compiler pass
	runPassOnObjects(t, &DisjunctionInferMapping{}, objects, objects)
}

func TestDisjunctionInferMapping_WithDisjunctionOfScalars_AsAnObject_hasNoImpact(t *testing.T) {
	// Prepare test input
	objects := []ast.Object{
		ast.NewObject("test", "ADisjunctionOfScalars", ast.NewDisjunction([]ast.Type{
			ast.String(),
			ast.Bool(),
		})),
	}

	// Prepare expected output
	disjunctionStructType := ast.NewStruct(
		ast.NewStructField("String", ast.String(ast.Nullable())),
		ast.NewStructField("Bool", ast.Bool(ast.Nullable())),
	)
	// The original disjunction definition is preserved as a hint
	disjunctionStructType.Hints[ast.HintDisjunctionOfScalars] = objects[0].Type.AsDisjunction()

	// Call the compiler pass
	runPassOnObjects(t, &DisjunctionInferMapping{}, objects, objects)
}

func TestDisjunctionInferMapping_WithDisjunctionOfRefs_AsAnObject_NoDiscriminatorMetadata_NoDiscriminatorFieldCandidate(t *testing.T) {
	// Prepare test input
	objects := []ast.Object{
		ast.NewObject("test", "ADisjunctionOfRefs", ast.NewDisjunction([]ast.Type{
			ast.NewRef("test", "SomeStruct"),
			ast.NewRef("test", "OtherStruct"),
		})),

		ast.NewObject("test", "SomeStruct", ast.NewStruct(
			ast.NewStructField("Kind", ast.String(ast.Value("some-struct"))), // No equivalent in OtherStruct
			ast.NewStructField("FieldFoo", ast.String()),
		)),
		ast.NewObject("test", "OtherStruct", ast.NewStruct(
			ast.NewStructField("Type", ast.String(ast.Value("other-struct"))),
			ast.NewStructField("FieldBar", ast.Bool()),
		)),
	}

	disjunctionOfRef := objects[0].DeepCopy()
	disjunctionOfRef.Type.PassesTrail = []string{"DisjunctionInferMapping[no_mapping_found:could not identify discriminator field]"}
	expected := []ast.Object{
		disjunctionOfRef,
		objects[1],
		objects[2],
	}

	runPassOnObjects(t, &DisjunctionInferMapping{}, objects, expected)
}

func TestDisjunctionInferMapping_WithDisjunctionOfRefs_AsAnObject_NoDiscriminatorMetadata_NonScalarDiscriminator_NonConstantReference(t *testing.T) {
	// Prepare test input
	disjunctionType := ast.NewDisjunction([]ast.Type{
		ast.NewRef("test", "SomeStruct"),
		ast.NewRef("test", "OtherStruct"),
	})
	disjunctionType.Disjunction.Discriminator = "MapOfString"

	objects := []ast.Object{
		ast.NewObject("test", "ADisjunctionOfRefs", disjunctionType),

		ast.NewObject("test", "SomeStruct", ast.NewStruct(
			ast.NewStructField("FieldFoo", ast.String()),
			ast.NewStructField("MapOfString", ast.NewMap(ast.String(), ast.String())),
		)),
		ast.NewObject("test", "OtherStruct", ast.NewStruct(
			ast.NewStructField("FieldBar", ast.Bool()),
			ast.NewStructField("MapOfString", ast.NewMap(ast.String(), ast.String())),
		)),
	}

	disjunctionOfRef := objects[0].DeepCopy()
	disjunctionOfRef.Type.PassesTrail = []string{"DisjunctionInferMapping[no_mapping_found:discriminator field 'MapOfString' is not a scalar or constant reference]"}
	expected := []ast.Object{
		disjunctionOfRef,
		objects[1],
		objects[2],
	}

	runPassOnObjects(t, &DisjunctionInferMapping{}, objects, expected)
}

func TestDisjunctionInferMapping_WithDisjunctionOfRefs_AsAnObject_NoDiscriminatorMetadata_NonConcreteDiscriminator(t *testing.T) {
	// Prepare test input
	disjunctionType := ast.NewDisjunction([]ast.Type{
		ast.NewRef("test", "SomeStruct"),
		ast.NewRef("test", "OtherStruct"),
	})
	disjunctionType.Disjunction.Discriminator = "Type"

	objects := []ast.Object{
		ast.NewObject("test", "ADisjunctionOfRefs", disjunctionType),

		ast.NewObject("test", "SomeStruct", ast.NewStruct(
			ast.NewStructField("Type", ast.String()), // Not a concrete scalar
			ast.NewStructField("FieldFoo", ast.String()),
		)),
		ast.NewObject("test", "OtherStruct", ast.NewStruct(
			ast.NewStructField("Type", ast.String(ast.Value("other-struct"))),
			ast.NewStructField("FieldBar", ast.Bool()),
		)),
	}

	disjunctionOfRef := objects[0].DeepCopy()
	disjunctionOfRef.Type.PassesTrail = []string{"DisjunctionInferMapping[no_mapping_found:discriminator field 'Type' is not a scalar or constant reference]"}
	expected := []ast.Object{
		disjunctionOfRef,
		objects[1],
		objects[2],
	}

	runPassOnObjects(t, &DisjunctionInferMapping{}, objects, expected)
}

func TestDisjunctionInferMapping_WithDisjunctionOfRefs_AsAnObject_NoDiscriminatorMetadata_UnknownDiscriminatorField(t *testing.T) {
	// Prepare test input
	disjunctionType := ast.NewDisjunction([]ast.Type{
		ast.NewRef("test", "SomeStruct"),
		ast.NewRef("test", "OtherStruct"),
	})
	disjunctionType.Disjunction.Discriminator = "DoesNotExist"

	objects := []ast.Object{
		ast.NewObject("test", "ADisjunctionOfRefs", disjunctionType),

		ast.NewObject("test", "SomeStruct", ast.NewStruct(
			ast.NewStructField("Type", ast.String(ast.Value("some-struct"))),
			ast.NewStructField("FieldFoo", ast.String()),
		)),
		ast.NewObject("test", "OtherStruct", ast.NewStruct(
			ast.NewStructField("Type", ast.String(ast.Value("other-struct"))),
			ast.NewStructField("FieldBar", ast.Bool()),
		)),
	}

	disjunctionOfRef := objects[0].DeepCopy()
	disjunctionOfRef.Type.PassesTrail = []string{"DisjunctionInferMapping[no_mapping_found:discriminator field 'DoesNotExist' not found]"}
	expected := []ast.Object{
		disjunctionOfRef,
		objects[1],
		objects[2],
	}

	runPassOnObjects(t, &DisjunctionInferMapping{}, objects, expected)
}

func TestDisjunctionInferMapping_WithDisjunctionOfRefs_AsAnObject_NoDiscriminatorMetadata(t *testing.T) {
	// Prepare test input
	objects := []ast.Object{
		ast.NewObject("test", "ADisjunctionOfRefs", ast.NewDisjunction([]ast.Type{
			ast.NewRef("test", "SomeStruct"),
			ast.NewRef("test", "OtherStruct"),
		})),

		ast.NewObject("test", "SomeStruct", ast.NewStruct(
			ast.NewStructField("Type", ast.String(ast.Value("some-struct"))),
			ast.NewStructField("FieldFoo", ast.String()),
		)),
		ast.NewObject("test", "OtherStruct", ast.NewStruct(
			ast.NewStructField("FieldBar", ast.NewMap(ast.String(), ast.String())),
			ast.NewStructField("Type", ast.String(ast.Value("other-struct"))),
		)),
	}

	// Prepare expected output
	newDisjunction := objects[0].DeepCopy()
	newDisjunction.Type.Disjunction.Discriminator = "Type"
	newDisjunction.Type.Disjunction.DiscriminatorMapping = map[string]string{
		"other-struct": "OtherStruct",
		"some-struct":  "SomeStruct",
	}

	expectedObjects := []ast.Object{
		newDisjunction,
		objects[1],
		objects[2],
	}

	// Call the compiler pass
	runPassOnObjects(t, &DisjunctionInferMapping{}, objects, expectedObjects)
}

func TestDisjunctionInferMapping_WithDisjunctionOfRefs_AsAnObject_Scalar_WithDiscriminatorFieldSet(t *testing.T) {
	// Prepare test input
	disjunctionType := ast.NewDisjunction([]ast.Type{
		ast.NewRef("test", "SomeStruct"),
		ast.NewRef("test", "OtherStruct"),
	})
	// Add discriminator-related metadata to the disjunction
	// Mapping omitted: it will be inferred
	disjunctionType.Disjunction.Discriminator = "Kind"

	objects := []ast.Object{
		ast.NewObject("test", "ADisjunctionOfRefs", disjunctionType),

		ast.NewObject("test", "SomeStruct", ast.NewStruct(
			ast.NewStructField("Type", ast.String(ast.Value("some-struct"))),
			ast.NewStructField("Kind", ast.String(ast.Value("some-kind"))),
			ast.NewStructField("FieldFoo", ast.String()),
		)),
		ast.NewObject("test", "OtherStruct", ast.NewStruct(
			ast.NewStructField("Type", ast.String(ast.Value("other-struct"))),
			ast.NewStructField("Kind", ast.String(ast.Value("other-kind"))),
			ast.NewStructField("FieldBar", ast.Bool()),
		)),
	}

	// Prepare expected output
	newDisjunction := objects[0].DeepCopy()
	newDisjunction.Type.Disjunction.DiscriminatorMapping = map[string]string{
		"other-kind": "OtherStruct",
		"some-kind":  "SomeStruct",
	}

	expectedObjects := []ast.Object{
		newDisjunction,
		objects[1],
		objects[2],
	}

	// Call the compiler pass
	runPassOnObjects(t, &DisjunctionInferMapping{}, objects, expectedObjects)
}

func TestDisjunctionInferMapping_WithDisjunctionOfRefs_AsAnObject_Scalar_WithDiscriminatorFieldAndMappingSet(t *testing.T) {
	// Prepare test input
	disjunctionType := ast.NewDisjunction([]ast.Type{
		ast.NewRef("test", "SomeStruct"),
		ast.NewRef("test", "OtherStruct"),
	})
	// Add discriminator-related metadata to the disjunction
	disjunctionType.Disjunction.Discriminator = "Kind"
	disjunctionType.Disjunction.DiscriminatorMapping = map[string]string{
		"other-kind": "OtherStruct",
		"some-kind":  "SomeStruct",
	}

	objects := []ast.Object{
		ast.NewObject("test", "ADisjunctionOfRefs", disjunctionType),

		ast.NewObject("test", "SomeStruct", ast.NewStruct(
			ast.NewStructField("Type", ast.String(ast.Value("some-struct"))),
			ast.NewStructField("Kind", ast.String(ast.Value("some-kind"))),
			ast.NewStructField("FieldFoo", ast.String()),
		)),
		ast.NewObject("test", "OtherStruct", ast.NewStruct(
			ast.NewStructField("Type", ast.String(ast.Value("other-struct"))),
			ast.NewStructField("Kind", ast.String(ast.Value("other-kind"))),
			ast.NewStructField("FieldBar", ast.Bool()),
		)),
	}

	// Call the compiler pass
	runPassOnObjects(t, &DisjunctionInferMapping{}, objects, objects)
}

func TestDisjunctionInferMapping_WithDisjunctionOfRefs_AsAnObject_ConcreteReference_WithDiscriminatorFieldSet(t *testing.T) {
	// Prepare test input
	disjunctionType := ast.NewDisjunction([]ast.Type{
		ast.NewRef("test", "SomeStruct"),
		ast.NewRef("test", "OtherStruct"),
	})
	// Add discriminator-related metadata to the disjunction
	// Mapping omitted: it will be inferred
	disjunctionType.Disjunction.Discriminator = "Kind"

	objects := []ast.Object{
		ast.NewObject("test", "ADisjunctionOfRefs", disjunctionType),
		ast.NewObject("test", "AnEnum", ast.NewEnum([]ast.EnumValue{
			{ast.String(), "ValueA", "a"},
			{ast.String(), "ValueB", "b"},
		})),
		ast.NewObject("test", "SomeStruct", ast.NewStruct(
			ast.NewStructField("Type", ast.String(ast.Value("some-struct"))),
			ast.NewStructField("Kind", ast.NewConstantReferenceType("test", "AnEnum", "a")),
			ast.NewStructField("FieldFoo", ast.String()),
		)),
		ast.NewObject("test", "OtherStruct", ast.NewStruct(
			ast.NewStructField("Type", ast.String(ast.Value("other-struct"))),
			ast.NewStructField("Kind", ast.NewConstantReferenceType("test", "AnEnum", "b")),
			ast.NewStructField("FieldBar", ast.Bool()),
		)),
	}

	// Prepare expected output
	newDisjunction := objects[0].DeepCopy()
	newDisjunction.Type.Disjunction.DiscriminatorMapping = map[string]string{
		"b": "OtherStruct",
		"a": "SomeStruct",
	}

	expectedObjects := []ast.Object{
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
	disjunctionType := ast.NewDisjunction([]ast.Type{
		ast.NewRef("test", "SomeStruct"),
		ast.NewRef("test", "OtherStruct"),
	})
	// Add discriminator-related metadata to the disjunction
	disjunctionType.Disjunction.Discriminator = "Kind"
	disjunctionType.Disjunction.DiscriminatorMapping = map[string]string{
		"b": "OtherStruct",
		"a": "SomeStruct",
	}

	objects := []ast.Object{
		ast.NewObject("test", "ADisjunctionOfRefs", disjunctionType),
		ast.NewObject("test", "AnEnum", ast.NewEnum([]ast.EnumValue{
			{ast.String(), "ValueA", "a"},
			{ast.String(), "ValueB", "b"},
		})),
		ast.NewObject("test", "SomeStruct", ast.NewStruct(
			ast.NewStructField("Type", ast.String(ast.Value("some-struct"))),
			ast.NewStructField("Kind", ast.NewConstantReferenceType("test", "AnEnum", "a")),
			ast.NewStructField("FieldFoo", ast.String()),
		)),
		ast.NewObject("test", "OtherStruct", ast.NewStruct(
			ast.NewStructField("Type", ast.String(ast.Value("other-struct"))),
			ast.NewStructField("Kind", ast.NewConstantReferenceType("test", "AnEnum", "b")),
			ast.NewStructField("FieldBar", ast.Bool()),
		)),
	}

	// Call the compiler pass
	runPassOnObjects(t, &DisjunctionInferMapping{}, objects, objects)
}
