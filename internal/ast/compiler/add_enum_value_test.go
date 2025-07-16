package compiler

import (
	"testing"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/testutils"
	"github.com/stretchr/testify/require"
)

func TestAddEnumFieldValueReference(t *testing.T) {
	schema := &ast.Schema{
		Package: "add_enum_value",
		Objects: testutils.ObjectsMap(
			ast.NewObject("test", "MyEnum", ast.NewEnum([]ast.EnumValue{
				{Value: "A", Name: "A", Type: ast.String()},
				{Value: "B", Name: "B", Type: ast.String()},
			})),
			ast.NewObject("test", "MyStruct", ast.NewStruct(
				ast.NewStructField("enum", ast.NewRef("test", "MyEnum")),
			)),
		),
	}

	expected := &ast.Schema{
		Package: "add_enum_value",
		Objects: testutils.ObjectsMap(
			ast.NewObject("test", "MyEnum", ast.NewEnum([]ast.EnumValue{
				{Value: "A", Name: "A", Type: ast.String()},
				{Value: "B", Name: "B", Type: ast.String()},
				{Value: "C", Name: "C", Type: ast.String()},
			})),
			ast.NewObject("test", "MyStruct", ast.NewStruct(
				ast.NewStructField("enum", ast.NewRef("test", "MyEnum")),
			)),
		),
	}

	pass := &AddEnumValue{
		FieldRef: FieldReference{
			Package: "test",
			Object:  "MyStruct",
			Field:   "enum",
		},
		Name:  "C",
		Value: "C",
	}

	runPassOnSchema(t, pass, schema, expected)
}

func TestAddEnumFieldValueDirectEnum(t *testing.T) {
	schema := &ast.Schema{
		Package: "add_enum_value",
		Objects: testutils.ObjectsMap(
			ast.NewObject("test", "MyStruct", ast.NewStruct(
				ast.NewStructField("enum", ast.NewEnum([]ast.EnumValue{
					{Value: 1, Name: "A", Type: ast.NewScalar(ast.KindInt64)},
					{Value: 2, Name: "B", Type: ast.NewScalar(ast.KindInt64)},
				})),
			)),
		),
	}

	expected := &ast.Schema{
		Package: "add_enum_value",
		Objects: testutils.ObjectsMap(
			ast.NewObject("test", "MyStruct", ast.NewStruct(
				ast.NewStructField("enum", ast.NewEnum([]ast.EnumValue{
					{Value: 1, Name: "A", Type: ast.NewScalar(ast.KindInt64)},
					{Value: 2, Name: "B", Type: ast.NewScalar(ast.KindInt64)},
					{Value: 3, Name: "C", Type: ast.NewScalar(ast.KindInt64)},
				})),
			)),
		),
	}

	pass := &AddEnumValue{
		FieldRef: FieldReference{
			Package: "test",
			Object:  "MyStruct",
			Field:   "enum",
		},
		Name:  "C",
		Value: 3,
	}

	runPassOnSchema(t, pass, schema, expected)
}

func TestAddEnumValueEnum(t *testing.T) {
	schema := &ast.Schema{
		Package: "add_enum_value",
		Objects: testutils.ObjectsMap(ast.NewObject("test", "MyEnum", ast.NewEnum([]ast.EnumValue{
			{Value: "A", Name: "A", Type: ast.String()},
			{Value: "B", Name: "B", Type: ast.String()},
		})),
		),
	}

	expected := &ast.Schema{
		Package: "add_enum_value",
		Objects: testutils.ObjectsMap(ast.NewObject("test", "MyEnum", ast.NewEnum([]ast.EnumValue{
			{Value: "A", Name: "A", Type: ast.String()},
			{Value: "B", Name: "B", Type: ast.String()},
			{Value: "C", Name: "C", Type: ast.String()},
		})),
		),
	}

	pass := &AddEnumValue{
		ObjectRef: ObjectReference{
			Package: "test",
			Object:  "MyEnum",
		},
		Name:  "C",
		Value: "C",
	}

	runPassOnSchema(t, pass, schema, expected)
}

func TestAddEnumValueInvalidValueKind(t *testing.T) {
	schema := &ast.Schema{
		Package: "add_enum_value",
		Objects: testutils.ObjectsMap(ast.NewObject("test", "MyEnum", ast.NewEnum([]ast.EnumValue{
			{Value: "A", Name: "A", Type: ast.String()},
			{Value: "B", Name: "B", Type: ast.String()},
		})),
		),
	}

	pass := &AddEnumValue{
		ObjectRef: ObjectReference{
			Package: "test",
			Object:  "MyEnum",
		},
		Name:  "C",
		Value: 1,
	}

	_, err := pass.Process(ast.Schemas{schema})
	require.Error(t, err)
}
