package openapi

import (
	"github.com/grafana/cog/internal/ast"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

var testFolder = "testdata/"

func TestDataTypes(t *testing.T) {
	f, err := GenerateAST(testFolder+"datatypes.json", Config{Package: "datatypes"})
	require.NoError(t, err)

	require.Len(t, f.Definitions, 1)

	def := f.Definitions[0]
	assert.Equal(t, def.Name, "DataTypes")
	assert.Equal(t, def.Type.Kind, ast.KindStruct)

	structType := def.Type.AsStruct()
	testCases := []field{
		{
			name:       "string",
			kind:       ast.KindScalar,
			scalarKind: ast.KindString,
			required:   true,
			nullable:   true,
		},
		{
			name:       "int32",
			kind:       ast.KindScalar,
			scalarKind: ast.KindInt32,
			defValue:   3,
		},
		{
			name:       "int64",
			kind:       ast.KindScalar,
			scalarKind: ast.KindInt64,
		},
		{
			name:       "float32",
			kind:       ast.KindScalar,
			scalarKind: ast.KindFloat32,
			defValue:   3.5,
		},
		{
			name:       "float64",
			kind:       ast.KindScalar,
			scalarKind: ast.KindFloat64,
		},
		{
			name:       "boolean",
			kind:       ast.KindScalar,
			scalarKind: ast.KindBool,
		},
		{
			name:       "bytes",
			kind:       ast.KindScalar,
			scalarKind: ast.KindBytes,
			required:   true,
		},
	}
	for _, tc := range testCases {
		validateFields(t, structType, tc)
	}
}

type field struct {
	name       string
	kind       ast.Kind
	scalarKind ast.ScalarKind
	required   bool
	nullable   bool
	defValue   any
}

func validateFields(t *testing.T, s ast.StructType, f field) {
	field, ok := s.FieldByName(f.name)
	require.True(t, ok)

	assert.Equal(t, field.Name, f.name)
	assert.Equal(t, field.Type.Kind, f.kind)
	assert.Equal(t, field.Required, f.required)
	assert.Equal(t, field.Default, f.defValue)
	assert.Equal(t, field.Type.Nullable, f.nullable)
	assert.Equal(t, field.Type.AsScalar().ScalarKind, f.scalarKind)
}
