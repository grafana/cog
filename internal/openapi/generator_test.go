package openapi

import (
	"testing"

	"github.com/grafana/cog/internal/ast"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testFolder = "testdata/"

func TestDataTypes(t *testing.T) {
	f, err := GenerateAST(testFolder+"datatypes.json", Config{Package: "datatypes"})
	require.NoError(t, err)

	require.Len(t, f.Objects, 1)

	def := f.Objects[0]
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
			defValue:   float64(3),
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
		t.Run(tc.name, func(t *testing.T) {
			validateFields(t, structType, tc)
		})
	}
}

func TestEnums(t *testing.T) {
	f, err := GenerateAST(testFolder+"enums.json", Config{Package: "enums"})
	require.NoError(t, err)

	require.Len(t, f.Objects, 1)

	def := f.Objects[0]
	assert.Equal(t, def.Name, "Enums")
	assert.Equal(t, def.Type.Kind, ast.KindStruct)

	structType := def.Type.AsStruct()
	testCases := []field{
		{
			name: "enumString",
			kind: ast.KindEnum,
			enumValues: []ast.EnumValue{
				{
					Type:  ast.NewScalar(ast.KindString),
					Name:  "a",
					Value: "a",
				},
				{
					Type:  ast.NewScalar(ast.KindString),
					Name:  "b",
					Value: "b",
				},
			},
			required: true,
		},
		{
			name: "enumInt",
			kind: ast.KindEnum,
			enumValues: []ast.EnumValue{
				{
					Type:  ast.NewScalar(ast.KindInt64),
					Name:  "3",
					Value: float64(3),
				},
				{
					Type:  ast.NewScalar(ast.KindInt64),
					Name:  "4",
					Value: float64(4),
				},
			},
		},
		{
			name: "enumWithDefault",
			kind: ast.KindEnum,
			enumValues: []ast.EnumValue{
				{
					Type:  ast.NewScalar(ast.KindString),
					Name:  "a",
					Value: "a",
				},
				{
					Type:  ast.NewScalar(ast.KindString),
					Name:  "b",
					Value: "b",
				},
			},
			defValue: "a",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			validateFields(t, structType, tc)
		})
	}
}

func TestArrays(t *testing.T) {
	f, err := GenerateAST(testFolder+"arrays.json", Config{Package: "arrays"})
	require.NoError(t, err)

	require.Len(t, f.Objects, 2)
	def := f.LocateDefinition("Arrays")
	assert.Equal(t, def.Type.Kind, ast.KindStruct)

	structType := def.Type.AsStruct()
	testCases := []field{
		{
			name:          "arrayString",
			kind:          ast.KindArray,
			required:      true,
			arrayTypeKind: ast.KindScalar,
			scalarKind:    ast.KindString,
		},
		{
			name:          "arrayInt",
			kind:          ast.KindArray,
			arrayTypeKind: ast.KindScalar,
			scalarKind:    ast.KindInt64,
		},
		{
			name:          "arrayRef",
			kind:          ast.KindArray,
			arrayTypeKind: ast.KindRef,
			refType:       "Test",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			validateFields(t, structType, tc)
		})
	}
}

func TestRefs(t *testing.T) {
	f, err := GenerateAST(testFolder+"refs.json", Config{Package: "refs"})
	require.NoError(t, err)

	require.Len(t, f.Objects, 2)
	def := f.LocateDefinition("Refs")
	assert.Equal(t, def.Type.Kind, ast.KindStruct)

	structType := def.Type.AsStruct()
	testCases := []field{
		{
			name:     "ref",
			kind:     ast.KindRef,
			refType:  "Test",
			required: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			validateFields(t, structType, tc)
		})
	}
}

type field struct {
	name          string
	kind          ast.Kind
	scalarKind    ast.ScalarKind
	required      bool
	nullable      bool
	defValue      any
	enumValues    []ast.EnumValue
	refType       string
	arrayTypeKind ast.Kind
}

func validateFields(t *testing.T, s ast.StructType, f field) {
	t.Helper()
	field, ok := s.FieldByName(f.name)
	require.True(t, ok)

	assert.Equal(t, field.Name, f.name)
	assert.Equal(t, field.Type.Kind, f.kind)
	assert.Equal(t, field.Required, f.required)
	assert.Equal(t, field.Type.Default, f.defValue)
	assert.Equal(t, field.Type.Nullable, f.nullable)

	validateKind(t, field.Type, f)
}

func validateKind(t *testing.T, tp ast.Type, f field) {
	t.Helper()
	switch tp.Kind {
	case ast.KindScalar:
		assert.Equal(t, tp.AsScalar().ScalarKind, f.scalarKind)
	case ast.KindEnum:
		assert.Equal(t, tp.AsEnum().Values, f.enumValues)
	case ast.KindRef:
		assert.Equal(t, tp.AsRef().ReferredType, f.refType)
	case ast.KindArray:
		assert.Equal(t, tp.AsArray().ValueType.Kind, f.arrayTypeKind)
		validateKind(t, tp.AsArray().ValueType, f)
	}
}
