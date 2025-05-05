package compiler

import (
	"testing"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/testutils"
	"github.com/stretchr/testify/require"
)

func TestAddFormatInt32ForUint8(t *testing.T) {
	// Prepare test input
	schema := &ast.Schema{
		Package: "test",
		Objects: testutils.ObjectsMap(
			ast.NewObject("test", "StructWithBytes", ast.NewStruct(
				ast.NewStructField("FieldUint8", ast.NewScalar(ast.KindUint8)),
				ast.NewStructField("FieldUint8Array", ast.NewArray(ast.NewScalar(ast.KindUint8))),
				ast.NewStructField("FieldInt32", ast.NewScalar(ast.KindInt32)),
				ast.NewStructField("FieldBytes", ast.NewScalar(ast.KindBytes)), // KindBytes should not be modified
			)),
			ast.NewObject("test", "NotAStruct", ast.NewScalar(ast.KindString)), // Non-struct object
		),
	}

	// Prepare expected output
	expected := schema.DeepCopy()
	structWithBytes := expected.Objects.Get("StructWithBytes")
	structWithBytes.Type.Struct.Fields[0].Comments = []string{"+format=int32"}
	structWithBytes.Type.Struct.Fields[1].Comments = []string{"+format=int32"}
	structWithBytes.AddToPassesTrail("AddFormatInt32ForUint8")

	// Run the compiler pass
	processedSchemas, err := (&AddFormatUint8{}).Process([]*ast.Schema{schema})

	// Verify results
	require.NoError(t, err)
	require.Len(t, processedSchemas, 1)
	require.Equal(t, expected.Objects, processedSchemas[0].Objects)
}
