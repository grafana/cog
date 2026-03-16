package compiler

import (
	"testing"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/testutils"
)

func TestDefaultAsTyped_ScalarStringDefault(t *testing.T) {
	schema := &ast.Schema{
		Package: "test",
		Objects: testutils.ObjectsMap(
			ast.NewObject("test", "SomeObject", ast.NewStruct(
				ast.NewStructField("AString", ast.String(ast.Default("hello"))),
			)),
		),
	}

	expected := &ast.Schema{
		Package: "test",
		Objects: testutils.ObjectsMap(
			ast.NewObject("test", "SomeObject", ast.NewStruct(
				ast.NewStructField("AString", ast.String(
					ast.Default("hello"),
					ast.TypedDefaultOpt(&ast.TypeDefault{Scalar: &ast.ScalarType{ScalarKind: ast.KindString, Value: "hello"}}),
				)),
			)),
		),
	}

	pass := &DefaultAsTyped{}
	runPassOnSchema(t, pass, schema, expected)
}

func TestDefaultAsTyped_ScalarBoolDefault(t *testing.T) {
	schema := &ast.Schema{
		Package: "test",
		Objects: testutils.ObjectsMap(
			ast.NewObject("test", "SomeObject", ast.NewStruct(
				ast.NewStructField("ABool", ast.Bool(ast.Default(true))),
			)),
		),
	}

	expected := &ast.Schema{
		Package: "test",
		Objects: testutils.ObjectsMap(
			ast.NewObject("test", "SomeObject", ast.NewStruct(
				ast.NewStructField("ABool", ast.Bool(
					ast.Default(true),
					ast.TypedDefaultOpt(&ast.TypeDefault{Scalar: &ast.ScalarType{ScalarKind: ast.KindBool, Value: true}}),
				)),
			)),
		),
	}

	pass := &DefaultAsTyped{}
	runPassOnSchema(t, pass, schema, expected)
}

func TestDefaultAsTyped_Float64Default(t *testing.T) {
	schema := &ast.Schema{
		Package: "test",
		Objects: testutils.ObjectsMap(
			ast.NewObject("test", "SomeObject", ast.NewStruct(
				ast.NewStructField("AFloat", ast.NewScalar(ast.KindFloat64, ast.Default(float64(3.14)))),
			)),
		),
	}

	expected := &ast.Schema{
		Package: "test",
		Objects: testutils.ObjectsMap(
			ast.NewObject("test", "SomeObject", ast.NewStruct(
				ast.NewStructField("AFloat", ast.NewScalar(ast.KindFloat64,
					ast.Default(3.14),
					ast.TypedDefaultOpt(&ast.TypeDefault{Scalar: &ast.ScalarType{ScalarKind: ast.KindFloat64, Value: float64(3.14)}}),
				)),
			)),
		),
	}

	pass := &DefaultAsTyped{}
	runPassOnSchema(t, pass, schema, expected)
}

func TestDefaultAsTyped_NoDefault(t *testing.T) {
	schema := &ast.Schema{
		Package: "test",
		Objects: testutils.ObjectsMap(
			ast.NewObject("test", "SomeObject", ast.NewStruct(
				ast.NewStructField("AString", ast.String()),
			)),
		),
	}

	pass := &DefaultAsTyped{}
	runPassOnSchema(t, pass, schema, schema)
}

func TestDefaultAsTyped_StructDefault(t *testing.T) {
	structDefault := map[string]any{
		"name": "alice",
	}

	// Build schema with a field whose type has a struct default set directly
	innerType := ast.NewStruct(
		ast.NewStructField("name", ast.String()),
	)
	innerType.Default = structDefault

	schema := &ast.Schema{
		Package: "test",
		Objects: testutils.ObjectsMap(
			ast.NewObject("test", "SomeObject", ast.NewStruct(
				ast.NewStructField("inner", innerType),
			)),
		),
	}

	expectedInnerType := innerType
	expectedInnerType.TypedDefault = &ast.TypeDefault{
		Struct: map[string]*ast.TypeDefault{
			"name": {Scalar: &ast.ScalarType{ScalarKind: ast.KindString, Value: "alice"}},
		},
	}

	expected := &ast.Schema{
		Package: "test",
		Objects: testutils.ObjectsMap(
			ast.NewObject("test", "SomeObject", ast.NewStruct(
				ast.NewStructField("inner", expectedInnerType),
			)),
		),
	}

	pass := &DefaultAsTyped{}
	runPassOnSchema(t, pass, schema, expected)
}
