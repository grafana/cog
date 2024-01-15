package compiler

import (
	"testing"

	"github.com/grafana/cog/internal/ast"
	"github.com/stretchr/testify/require"
)

func TestDisjunctionToType_WithNonDisjunctionObjects_HasNoImpact(t *testing.T) {
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
	runPassOnObjects(t, &DisjunctionToType{}, objects, objects)
}

func TestDisjunctionToType_WithDisjunctionOfScalars_AsAnObject(t *testing.T) {
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

	expectedObjects := []ast.Object{
		ast.NewObject("test", "ADisjunctionOfScalars", ast.NewRef("test", "StringOrBool")),
		ast.NewObject("test", "StringOrBool", disjunctionStructType),
	}

	// Call the compiler pass
	runPassOnObjects(t, &DisjunctionToType{}, objects, expectedObjects)
}

func TestDisjunctionToType_WithDisjunctionOfScalars_AsAMapValueType(t *testing.T) {
	// Prepare test input
	objects := []ast.Object{
		ast.NewObject("test", "ADisjunctionOfScalars", ast.NewMap(
			ast.String(),
			ast.NewDisjunction([]ast.Type{
				ast.String(),
				ast.Bool(),
			}),
		)),
	}

	// Prepare expected output
	disjunctionStructType := ast.NewStruct(
		ast.NewStructField("String", ast.String(ast.Nullable())),
		ast.NewStructField("Bool", ast.Bool(ast.Nullable())),
	)
	// The original disjunction definition is preserved as a hint
	disjunctionStructType.Hints[ast.HintDisjunctionOfScalars] = objects[0].Type.AsMap().ValueType.AsDisjunction()

	expectedObjects := []ast.Object{
		ast.NewObject("test", "ADisjunctionOfScalars", ast.NewMap(
			ast.String(),
			ast.NewRef("test", "StringOrBool"),
		)),
		ast.NewObject("test", "StringOrBool", disjunctionStructType),
	}

	// Call the compiler pass
	runPassOnObjects(t, &DisjunctionToType{}, objects, expectedObjects)
}

func TestDisjunctionToType_WithDisjunctionOfScalars_AsAStructField(t *testing.T) {
	// Prepare test input
	disjunctionType := ast.NewDisjunction([]ast.Type{
		ast.String(),
		ast.Bool(),
	})
	objects := []ast.Object{
		ast.NewObject("test", "AStructWithADisjunctionOfScalars", ast.NewStruct(
			ast.NewStructField("AFieldWithADisjunctionOfScalars", disjunctionType),
		)),
	}

	// Prepare expected output
	disjunctionStructType := ast.NewStruct(
		ast.NewStructField("String", ast.String(ast.Nullable())),
		ast.NewStructField("Bool", ast.Bool(ast.Nullable())),
	)
	// The original disjunction definition is preserved as a hint
	disjunctionStructType.Hints[ast.HintDisjunctionOfScalars] = disjunctionType.AsDisjunction()

	expectedObjects := []ast.Object{
		ast.NewObject("test", "AStructWithADisjunctionOfScalars", ast.NewStruct(
			ast.NewStructField("AFieldWithADisjunctionOfScalars", ast.NewRef("test", "StringOrBool")),
		)),
		ast.NewObject("test", "StringOrBool", disjunctionStructType),
	}

	// Call the compiler pass
	runPassOnObjects(t, &DisjunctionToType{}, objects, expectedObjects)
}

func TestDisjunctionToType_WithDisjunctionOfScalars_AsNullableAStructField(t *testing.T) {
	// Prepare test input
	disjunctionType := ast.NewDisjunction([]ast.Type{
		ast.String(),
		ast.Bool(),
	}, ast.Nullable())
	objects := []ast.Object{
		ast.NewObject("test", "AStructWithADisjunctionOfScalars", ast.NewStruct(
			ast.NewStructField("AFieldWithADisjunctionOfScalars", disjunctionType),
		)),
	}

	// Prepare expected output
	disjunctionStructType := ast.NewStruct(
		ast.NewStructField("String", ast.String(ast.Nullable())),
		ast.NewStructField("Bool", ast.Bool(ast.Nullable())),
	)
	// The original disjunction definition is preserved as a hint
	disjunctionStructType.Hints[ast.HintDisjunctionOfScalars] = disjunctionType.AsDisjunction()

	expectedObjects := []ast.Object{
		ast.NewObject("test", "AStructWithADisjunctionOfScalars", ast.NewStruct(
			ast.NewStructField("AFieldWithADisjunctionOfScalars", ast.NewRef("test", "StringOrBool", ast.Nullable())),
		)),
		ast.NewObject("test", "StringOrBool", disjunctionStructType),
	}

	// Call the compiler pass
	runPassOnObjects(t, &DisjunctionToType{}, objects, expectedObjects)
}

func TestDisjunctionToType_WithDisjunctionOfScalars_AsAnArrayValueType(t *testing.T) {
	// Prepare test input
	disjunctionType := ast.NewDisjunction([]ast.Type{
		ast.String(),
		ast.Bool(),
	})
	objects := []ast.Object{
		ast.NewObject("test", "AnArrayWithADisjunctionOfScalars", ast.NewArray(disjunctionType)),
	}

	// Prepare expected output
	disjunctionStructType := ast.NewStruct(
		ast.NewStructField("String", ast.String(ast.Nullable())),
		ast.NewStructField("Bool", ast.Bool(ast.Nullable())),
	)
	// The original disjunction definition is preserved as a hint
	disjunctionStructType.Hints[ast.HintDisjunctionOfScalars] = disjunctionType.AsDisjunction()

	expectedObjects := []ast.Object{
		ast.NewObject("test", "AnArrayWithADisjunctionOfScalars", ast.NewArray(ast.NewRef("test", "StringOrBool"))),
		ast.NewObject("test", "StringOrBool", disjunctionStructType),
	}

	// Call the compiler pass
	runPassOnObjects(t, &DisjunctionToType{}, objects, expectedObjects)
}

func TestDisjunctionToType_WithDisjunctionOfRefs_AsAnObject_NoDiscriminatorMetadata_NoDiscriminatorFieldCandidate(t *testing.T) {
	req := require.New(t)

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

	compilerPass := &DisjunctionToType{}
	_, err := compilerPass.Process([]*ast.Schema{
		{Package: "test", Objects: objects},
	})
	req.Error(err)
	req.ErrorIs(err, ErrCanNotInferDiscriminator)
	req.ErrorContains(err, "discriminator field is empty")
}

func TestDisjunctionToType_WithDisjunctionOfRefs_AsAnObject_NoDiscriminatorMetadata_NonScalarDiscriminator(t *testing.T) {
	req := require.New(t)

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

	compilerPass := &DisjunctionToType{}
	_, err := compilerPass.Process([]*ast.Schema{
		{Package: "test", Objects: objects},
	})
	req.Error(err)
	req.ErrorIs(err, ErrCanNotInferDiscriminator)
	req.ErrorContains(err, "field is not a scalar")
}

func TestDisjunctionToType_WithDisjunctionOfRefs_AsAnObject_NoDiscriminatorMetadata_NonConcreteDiscriminator(t *testing.T) {
	req := require.New(t)

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

	compilerPass := &DisjunctionToType{}
	_, err := compilerPass.Process([]*ast.Schema{
		{Package: "test", Objects: objects},
	})
	req.Error(err)
	req.ErrorIs(err, ErrCanNotInferDiscriminator)
	req.ErrorContains(err, "field is not concrete")
}

func TestDisjunctionToType_WithDisjunctionOfRefs_AsAnObject_NoDiscriminatorMetadata_UnknownDiscriminatorField(t *testing.T) {
	req := require.New(t)

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

	compilerPass := &DisjunctionToType{}
	_, err := compilerPass.Process([]*ast.Schema{
		{Package: "test", Objects: objects},
	})
	req.Error(err)
	req.ErrorIs(err, ErrCanNotInferDiscriminator)
	req.ErrorContains(err, "discriminator field 'DoesNotExist' not found")
}

func TestDisjunctionToType_WithDisjunctionOfRefs_AsAnObject_NoDiscriminatorMetadata(t *testing.T) {
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
	disjunctionStructType := ast.NewStruct(
		ast.NewStructField("SomeStruct", ast.NewRef("test", "SomeStruct", ast.Nullable())),
		ast.NewStructField("OtherStruct", ast.NewRef("test", "OtherStruct", ast.Nullable())),
	)
	// The original disjunction definition is preserved as a hint
	disjunctionTypeWithDiscriminatorMeta := objects[0].Type.AsDisjunction()

	// Metadata should be inferred
	disjunctionTypeWithDiscriminatorMeta.Discriminator = "Type"
	disjunctionTypeWithDiscriminatorMeta.DiscriminatorMapping = map[string]string{
		"other-struct": "OtherStruct",
		"some-struct":  "SomeStruct",
	}
	disjunctionStructType.Hints[ast.HintDiscriminatedDisjunctionOfRefs] = disjunctionTypeWithDiscriminatorMeta

	expectedObjects := []ast.Object{
		ast.NewObject("test", "ADisjunctionOfRefs", ast.NewRef("test", "SomeStructOrOtherStruct")),
		objects[1],
		objects[2],
		ast.NewObject("test", "SomeStructOrOtherStruct", disjunctionStructType),
	}

	// Call the compiler pass
	runPassOnObjects(t, &DisjunctionToType{}, objects, expectedObjects)
}

func TestDisjunctionToType_WithDisjunctionOfRefs_AsAnObject_WithDiscriminatorFieldSet(t *testing.T) {
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
	disjunctionStructType := ast.NewStruct(
		ast.NewStructField("SomeStruct", ast.NewRef("test", "SomeStruct", ast.Nullable())),
		ast.NewStructField("OtherStruct", ast.NewRef("test", "OtherStruct", ast.Nullable())),
	)
	// The original disjunction definition is preserved as a hint
	disjunctionTypeWithDiscriminatorMeta := objects[0].Type.AsDisjunction()

	// Metadata should be inferred
	disjunctionTypeWithDiscriminatorMeta.Discriminator = "Kind"
	disjunctionTypeWithDiscriminatorMeta.DiscriminatorMapping = map[string]string{
		"other-kind": "OtherStruct",
		"some-kind":  "SomeStruct",
	}
	disjunctionStructType.Hints[ast.HintDiscriminatedDisjunctionOfRefs] = disjunctionTypeWithDiscriminatorMeta

	expectedObjects := []ast.Object{
		ast.NewObject("test", "ADisjunctionOfRefs", ast.NewRef("test", "SomeStructOrOtherStruct")),
		objects[1],
		objects[2],
		ast.NewObject("test", "SomeStructOrOtherStruct", disjunctionStructType),
	}

	// Call the compiler pass
	runPassOnObjects(t, &DisjunctionToType{}, objects, expectedObjects)
}

func TestDisjunctionToType_WithDisjunctionOfRefs_AsAnObject_WithDiscriminatorFieldAndMappingSet(t *testing.T) {
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

	// Prepare expected output
	disjunctionStructType := ast.NewStruct(
		ast.NewStructField("SomeStruct", ast.NewRef("test", "SomeStruct", ast.Nullable())),
		ast.NewStructField("OtherStruct", ast.NewRef("test", "OtherStruct", ast.Nullable())),
	)
	// The original disjunction definition is preserved as a hint
	disjunctionTypeWithDiscriminatorMeta := objects[0].Type.AsDisjunction()

	// Metadata should be inferred
	disjunctionTypeWithDiscriminatorMeta.Discriminator = "Kind"
	disjunctionTypeWithDiscriminatorMeta.DiscriminatorMapping = map[string]string{
		"other-kind": "OtherStruct",
		"some-kind":  "SomeStruct",
	}
	disjunctionStructType.Hints[ast.HintDiscriminatedDisjunctionOfRefs] = disjunctionTypeWithDiscriminatorMeta

	expectedObjects := []ast.Object{
		ast.NewObject("test", "ADisjunctionOfRefs", ast.NewRef("test", "SomeStructOrOtherStruct")),
		objects[1],
		objects[2],
		ast.NewObject("test", "SomeStructOrOtherStruct", disjunctionStructType),
	}

	// Call the compiler pass
	runPassOnObjects(t, &DisjunctionToType{}, objects, expectedObjects)
}
