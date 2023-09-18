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
		ast.NewObject("AMap", ast.NewMap(ast.NewScalar(ast.KindString), ast.NewScalar(ast.KindString))),
		ast.NewObject("ARef", ast.NewRef("AMap")),
		ast.NewObject("AnEnum", ast.NewEnum([]ast.EnumValue{
			{
				Name:  "Foo",
				Type:  ast.NewScalar(ast.KindString),
				Value: "foo",
			},
			{
				Name:  "Bar",
				Type:  ast.NewScalar(ast.KindString),
				Value: "bar",
			},
		})),
		ast.NewObject("AnArray", ast.NewArray(ast.NewScalar(ast.KindString))),
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
			ast.NewScalar(ast.KindString),
			ast.NewScalar(ast.KindBool),
		})),
	}

	// Prepare expected output
	disjunctionStructType := ast.NewStruct(
		ast.NewStructField("ValString", ast.NewScalar(ast.KindString)),
		ast.NewStructField("ValBool", ast.NewScalar(ast.KindBool)),
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
		ast.NewScalar(ast.KindString),
		ast.NewScalar(ast.KindBool),
	})
	objects := []ast.Object{
		ast.NewObject("AStructWithADisjunctionOfScalars", ast.NewStruct(
			ast.NewStructField("AFieldWithADisjunctionOfScalars", disjunctionType),
		)),
	}

	// Prepare expected output
	disjunctionStructType := ast.NewStruct(
		ast.NewStructField("ValString", ast.NewScalar(ast.KindString)),
		ast.NewStructField("ValBool", ast.NewScalar(ast.KindBool)),
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
		ast.NewScalar(ast.KindString),
		ast.NewScalar(ast.KindBool),
	})
	objects := []ast.Object{
		ast.NewObject("AnArrayWithADisjunctionOfScalars", ast.NewArray(disjunctionType)),
	}

	// Prepare expected output
	disjunctionStructType := ast.NewStruct(
		ast.NewStructField("ValString", ast.NewScalar(ast.KindString)),
		ast.NewStructField("ValBool", ast.NewScalar(ast.KindBool)),
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
