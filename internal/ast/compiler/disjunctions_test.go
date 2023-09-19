package compiler

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/grafana/cog/internal/ast"
	"github.com/stretchr/testify/require"
)

func TestDisjunctionToType_WithNonDisjunctionObjects_HasNoImpact(t *testing.T) {
	// Prepare test input
	objects := []ast.Object{
		ast.NewObject("AMap", ast.NewMap(ast.String(), ast.String())),
		ast.NewObject("ARef", ast.NewRef("AMap")),
		ast.NewObject("AnEnum", ast.NewEnum([]ast.EnumValue{
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
		ast.NewObject("AnArray", ast.NewArray(ast.String())),
		ast.NewObject("AScalar", ast.NewScalar(ast.KindInt8)),
		ast.NewObject("AStruct", ast.NewStruct(
			ast.NewStructField("SomeNonDisjunctionField", ast.NewScalar(ast.KindInt8)),
		)),
	}

	// Call the compiler pass
	runDisjunctionPass(t, objects, objects)
}

func TestDisjunctionToType_WithDisjunctionOfScalars_AsAnObject(t *testing.T) {
	// Prepare test input
	objects := []ast.Object{
		ast.NewObject("ADisjunctionOfScalars", ast.NewDisjunction([]ast.Type{
			ast.String(),
			ast.Bool(),
		})),
	}

	// Prepare expected output
	disjunctionStructType := ast.NewStruct(
		ast.NewStructField("ValString", ast.String(ast.Nullable())),
		ast.NewStructField("ValBool", ast.Bool(ast.Nullable())),
	)
	// The original disjunction definition is preserved as a hint
	disjunctionStructType.Struct.Hint[ast.HintDisjunctionOfScalars] = objects[0].Type.AsDisjunction()

	expectedObjects := []ast.Object{
		ast.NewObject("ADisjunctionOfScalars", ast.NewRef("StringOrBool")),
		ast.NewObject("StringOrBool", disjunctionStructType),
	}

	// Call the compiler pass
	runDisjunctionPass(t, objects, expectedObjects)
}

func TestDisjunctionToType_WithDisjunctionOfScalars_AsAStructField(t *testing.T) {
	// Prepare test input
	disjunctionType := ast.NewDisjunction([]ast.Type{
		ast.String(),
		ast.Bool(),
	})
	objects := []ast.Object{
		ast.NewObject("AStructWithADisjunctionOfScalars", ast.NewStruct(
			ast.NewStructField("AFieldWithADisjunctionOfScalars", disjunctionType),
		)),
	}

	// Prepare expected output
	disjunctionStructType := ast.NewStruct(
		ast.NewStructField("ValString", ast.String(ast.Nullable())),
		ast.NewStructField("ValBool", ast.Bool(ast.Nullable())),
	)
	// The original disjunction definition is preserved as a hint
	disjunctionStructType.Struct.Hint[ast.HintDisjunctionOfScalars] = disjunctionType.AsDisjunction()

	expectedObjects := []ast.Object{
		ast.NewObject("AStructWithADisjunctionOfScalars", ast.NewStruct(
			ast.NewStructField("AFieldWithADisjunctionOfScalars", ast.NewRef("StringOrBool")),
		)),
		ast.NewObject("StringOrBool", disjunctionStructType),
	}

	// Call the compiler pass
	runDisjunctionPass(t, objects, expectedObjects)
}

func TestDisjunctionToType_WithDisjunctionOfScalars_AsAnArrayValueType(t *testing.T) {
	// Prepare test input
	disjunctionType := ast.NewDisjunction([]ast.Type{
		ast.String(),
		ast.Bool(),
	})
	objects := []ast.Object{
		ast.NewObject("AnArrayWithADisjunctionOfScalars", ast.NewArray(disjunctionType)),
	}

	// Prepare expected output
	disjunctionStructType := ast.NewStruct(
		ast.NewStructField("ValString", ast.String(ast.Nullable())),
		ast.NewStructField("ValBool", ast.Bool(ast.Nullable())),
	)
	// The original disjunction definition is preserved as a hint
	disjunctionStructType.Struct.Hint[ast.HintDisjunctionOfScalars] = disjunctionType.AsDisjunction()

	expectedObjects := []ast.Object{
		ast.NewObject("AnArrayWithADisjunctionOfScalars", ast.NewArray(ast.NewRef("StringOrBool"))),
		ast.NewObject("StringOrBool", disjunctionStructType),
	}

	// Call the compiler pass
	runDisjunctionPass(t, objects, expectedObjects)
}

func TestDisjunctionToType_WithDisjunctionOfRefs_AsAnObject_NoDiscriminatorMetadata_NoDiscriminatorFieldCandidate(t *testing.T) {
	req := require.New(t)

	// Prepare test input
	objects := []ast.Object{
		ast.NewObject("ADisjunctionOfRefs", ast.NewDisjunction([]ast.Type{
			ast.NewRef("SomeStruct"),
			ast.NewRef("OtherStruct"),
		})),

		ast.NewObject("SomeStruct", ast.NewStruct(
			ast.NewStructField("Kind", ast.String(ast.Value("some-struct"))), // No equivalent in OtherStruct
			ast.NewStructField("FieldFoo", ast.String()),
		)),
		ast.NewObject("OtherStruct", ast.NewStruct(
			ast.NewStructField("Type", ast.String(ast.Value("other-struct"))),
			ast.NewStructField("FieldBar", ast.Bool()),
		)),
	}

	compilerPass := &DisjunctionToType{}
	_, err := compilerPass.Process([]*ast.File{
		{Package: "test", Definitions: objects},
	})
	req.Error(err)
	req.ErrorIs(err, ErrCanNotInferDiscriminator)
	req.ErrorContains(err, "discriminator field is empty")
}

func TestDisjunctionToType_WithDisjunctionOfRefs_AsAnObject_NoDiscriminatorMetadata_NonScalarDiscriminator(t *testing.T) {
	req := require.New(t)

	// Prepare test input
	disjunctionType := ast.NewDisjunction([]ast.Type{
		ast.NewRef("SomeStruct"),
		ast.NewRef("OtherStruct"),
	})
	disjunctionType.Disjunction.Discriminator = "MapOfString"

	objects := []ast.Object{
		ast.NewObject("ADisjunctionOfRefs", disjunctionType),

		ast.NewObject("SomeStruct", ast.NewStruct(
			ast.NewStructField("FieldFoo", ast.String()),
			ast.NewStructField("MapOfString", ast.NewMap(ast.String(), ast.String())),
		)),
		ast.NewObject("OtherStruct", ast.NewStruct(
			ast.NewStructField("FieldBar", ast.Bool()),
			ast.NewStructField("MapOfString", ast.NewMap(ast.String(), ast.String())),
		)),
	}

	compilerPass := &DisjunctionToType{}
	_, err := compilerPass.Process([]*ast.File{
		{Package: "test", Definitions: objects},
	})
	req.Error(err)
	req.ErrorIs(err, ErrCanNotInferDiscriminator)
	req.ErrorContains(err, "field is not a scalar")
}

func TestDisjunctionToType_WithDisjunctionOfRefs_AsAnObject_NoDiscriminatorMetadata_NonConcreteDiscriminator(t *testing.T) {
	req := require.New(t)

	// Prepare test input
	disjunctionType := ast.NewDisjunction([]ast.Type{
		ast.NewRef("SomeStruct"),
		ast.NewRef("OtherStruct"),
	})
	disjunctionType.Disjunction.Discriminator = "Type"

	objects := []ast.Object{
		ast.NewObject("ADisjunctionOfRefs", disjunctionType),

		ast.NewObject("SomeStruct", ast.NewStruct(
			ast.NewStructField("Type", ast.String()), // Not a concrete scalar
			ast.NewStructField("FieldFoo", ast.String()),
		)),
		ast.NewObject("OtherStruct", ast.NewStruct(
			ast.NewStructField("Type", ast.String(ast.Value("other-struct"))),
			ast.NewStructField("FieldBar", ast.Bool()),
		)),
	}

	compilerPass := &DisjunctionToType{}
	_, err := compilerPass.Process([]*ast.File{
		{Package: "test", Definitions: objects},
	})
	req.Error(err)
	req.ErrorIs(err, ErrCanNotInferDiscriminator)
	req.ErrorContains(err, "field is not concrete")
}

func TestDisjunctionToType_WithDisjunctionOfRefs_AsAnObject_NoDiscriminatorMetadata_UnknownDiscriminatorField(t *testing.T) {
	req := require.New(t)

	// Prepare test input
	disjunctionType := ast.NewDisjunction([]ast.Type{
		ast.NewRef("SomeStruct"),
		ast.NewRef("OtherStruct"),
	})
	disjunctionType.Disjunction.Discriminator = "DoesNotExist"

	objects := []ast.Object{
		ast.NewObject("ADisjunctionOfRefs", disjunctionType),

		ast.NewObject("SomeStruct", ast.NewStruct(
			ast.NewStructField("Type", ast.String(ast.Value("some-struct"))),
			ast.NewStructField("FieldFoo", ast.String()),
		)),
		ast.NewObject("OtherStruct", ast.NewStruct(
			ast.NewStructField("Type", ast.String(ast.Value("other-struct"))),
			ast.NewStructField("FieldBar", ast.Bool()),
		)),
	}

	compilerPass := &DisjunctionToType{}
	_, err := compilerPass.Process([]*ast.File{
		{Package: "test", Definitions: objects},
	})
	req.Error(err)
	req.ErrorIs(err, ErrCanNotInferDiscriminator)
	req.ErrorContains(err, "discriminator field 'DoesNotExist' not found")
}

func TestDisjunctionToType_WithDisjunctionOfRefs_AsAnObject_NoDiscriminatorMetadata(t *testing.T) {
	// Prepare test input
	objects := []ast.Object{
		ast.NewObject("ADisjunctionOfRefs", ast.NewDisjunction([]ast.Type{
			ast.NewRef("SomeStruct"),
			ast.NewRef("OtherStruct"),
		})),

		ast.NewObject("SomeStruct", ast.NewStruct(
			ast.NewStructField("Type", ast.String(ast.Value("some-struct"))),
			ast.NewStructField("FieldFoo", ast.String()),
		)),
		ast.NewObject("OtherStruct", ast.NewStruct(
			ast.NewStructField("FieldBar", ast.NewMap(ast.String(), ast.String())),
			ast.NewStructField("Type", ast.String(ast.Value("other-struct"))),
		)),
	}

	// Prepare expected output
	disjunctionStructType := ast.NewStruct(
		ast.NewStructField("ValSomeStruct", ast.NewRef("SomeStruct", ast.Nullable())),
		ast.NewStructField("ValOtherStruct", ast.NewRef("OtherStruct", ast.Nullable())),
	)
	// The original disjunction definition is preserved as a hint
	disjunctionTypeWithDiscriminatorMeta := objects[0].Type.AsDisjunction()

	// Metadata should be inferred
	disjunctionTypeWithDiscriminatorMeta.Discriminator = "Type"
	disjunctionTypeWithDiscriminatorMeta.DiscriminatorMapping = map[string]any{
		"OtherStruct": "other-struct",
		"SomeStruct":  "some-struct",
	}
	disjunctionStructType.Struct.Hint[ast.HintDiscriminatedDisjunctionOfRefs] = disjunctionTypeWithDiscriminatorMeta

	expectedObjects := []ast.Object{
		ast.NewObject("ADisjunctionOfRefs", ast.NewRef("SomeStructOrOtherStruct")),
		objects[1],
		objects[2],
		ast.NewObject("SomeStructOrOtherStruct", disjunctionStructType),
	}

	// Call the compiler pass
	runDisjunctionPass(t, objects, expectedObjects)
}

func TestDisjunctionToType_WithDisjunctionOfRefs_AsAnObject_WithDiscriminatorFieldSet(t *testing.T) {
	// Prepare test input

	disjunctionType := ast.NewDisjunction([]ast.Type{
		ast.NewRef("SomeStruct"),
		ast.NewRef("OtherStruct"),
	})
	// Add discriminator-related metadata to the disjunction
	// Mapping omitted: it will be inferred
	disjunctionType.Disjunction.Discriminator = "Kind"

	objects := []ast.Object{
		ast.NewObject("ADisjunctionOfRefs", disjunctionType),

		ast.NewObject("SomeStruct", ast.NewStruct(
			ast.NewStructField("Type", ast.String(ast.Value("some-struct"))),
			ast.NewStructField("Kind", ast.String(ast.Value("some-kind"))),
			ast.NewStructField("FieldFoo", ast.String()),
		)),
		ast.NewObject("OtherStruct", ast.NewStruct(
			ast.NewStructField("Type", ast.String(ast.Value("other-struct"))),
			ast.NewStructField("Kind", ast.String(ast.Value("other-kind"))),
			ast.NewStructField("FieldBar", ast.Bool()),
		)),
	}

	// Prepare expected output
	disjunctionStructType := ast.NewStruct(
		ast.NewStructField("ValSomeStruct", ast.NewRef("SomeStruct", ast.Nullable())),
		ast.NewStructField("ValOtherStruct", ast.NewRef("OtherStruct", ast.Nullable())),
	)
	// The original disjunction definition is preserved as a hint
	disjunctionTypeWithDiscriminatorMeta := objects[0].Type.AsDisjunction()

	// Metadata should be inferred
	disjunctionTypeWithDiscriminatorMeta.Discriminator = "Kind"
	disjunctionTypeWithDiscriminatorMeta.DiscriminatorMapping = map[string]any{
		"OtherStruct": "other-kind",
		"SomeStruct":  "some-kind",
	}
	disjunctionStructType.Struct.Hint[ast.HintDiscriminatedDisjunctionOfRefs] = disjunctionTypeWithDiscriminatorMeta

	expectedObjects := []ast.Object{
		ast.NewObject("ADisjunctionOfRefs", ast.NewRef("SomeStructOrOtherStruct")),
		objects[1],
		objects[2],
		ast.NewObject("SomeStructOrOtherStruct", disjunctionStructType),
	}

	// Call the compiler pass
	runDisjunctionPass(t, objects, expectedObjects)
}

func TestDisjunctionToType_WithDisjunctionOfRefs_AsAnObject_WithDiscriminatorFieldAndMappingSet(t *testing.T) {
	// Prepare test input
	disjunctionType := ast.NewDisjunction([]ast.Type{
		ast.NewRef("SomeStruct"),
		ast.NewRef("OtherStruct"),
	})
	// Add discriminator-related metadata to the disjunction
	disjunctionType.Disjunction.Discriminator = "Kind"
	disjunctionType.Disjunction.DiscriminatorMapping = map[string]any{
		"OtherStruct": "other-kind",
		"SomeStruct":  "some-kind",
	}

	objects := []ast.Object{
		ast.NewObject("ADisjunctionOfRefs", disjunctionType),

		ast.NewObject("SomeStruct", ast.NewStruct(
			ast.NewStructField("Type", ast.String(ast.Value("some-struct"))),
			ast.NewStructField("Kind", ast.String(ast.Value("some-kind"))),
			ast.NewStructField("FieldFoo", ast.String()),
		)),
		ast.NewObject("OtherStruct", ast.NewStruct(
			ast.NewStructField("Type", ast.String(ast.Value("other-struct"))),
			ast.NewStructField("Kind", ast.String(ast.Value("other-kind"))),
			ast.NewStructField("FieldBar", ast.Bool()),
		)),
	}

	// Prepare expected output
	disjunctionStructType := ast.NewStruct(
		ast.NewStructField("ValSomeStruct", ast.NewRef("SomeStruct", ast.Nullable())),
		ast.NewStructField("ValOtherStruct", ast.NewRef("OtherStruct", ast.Nullable())),
	)
	// The original disjunction definition is preserved as a hint
	disjunctionTypeWithDiscriminatorMeta := objects[0].Type.AsDisjunction()

	// Metadata should be inferred
	disjunctionTypeWithDiscriminatorMeta.Discriminator = "Kind"
	disjunctionTypeWithDiscriminatorMeta.DiscriminatorMapping = map[string]any{
		"OtherStruct": "other-kind",
		"SomeStruct":  "some-kind",
	}
	disjunctionStructType.Struct.Hint[ast.HintDiscriminatedDisjunctionOfRefs] = disjunctionTypeWithDiscriminatorMeta

	expectedObjects := []ast.Object{
		ast.NewObject("ADisjunctionOfRefs", ast.NewRef("SomeStructOrOtherStruct")),
		objects[1],
		objects[2],
		ast.NewObject("SomeStructOrOtherStruct", disjunctionStructType),
	}

	// Call the compiler pass
	runDisjunctionPass(t, objects, expectedObjects)
}

func runDisjunctionPass(t *testing.T, input []ast.Object, expectedOutput []ast.Object) {
	t.Helper()

	req := require.New(t)

	compilerPass := &DisjunctionToType{}
	processedFiles, err := compilerPass.Process([]*ast.File{
		{
			Package:     "test",
			Definitions: input,
		},
	})
	req.NoError(err)
	req.Len(processedFiles, 1)
	req.Empty(cmp.Diff(expectedOutput, processedFiles[0].Definitions))
}
