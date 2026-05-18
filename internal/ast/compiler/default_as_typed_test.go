package compiler

import (
	"testing"

	"github.com/grafana/cog/internal/ast"
)

func TestDefaultAsTyped_scalar(t *testing.T) {
	scalar := ast.String(ast.Default("hello"))
	scalar.TypedDefault = ast.NewScalarDefault("hello")

	input := []ast.Object{
		ast.NewObject(testPkgName, "AString", ast.String(ast.Default("hello"))),
	}
	expected := []ast.Object{
		ast.NewObject(testPkgName, "AString", scalar),
	}

	runPassOnObjects(t, &DefaultAsTyped{}, input, expected)
}

func TestDefaultAsTyped_enum(t *testing.T) {
	enumType := ast.NewEnum([]ast.EnumValue{
		{Type: ast.String(), Name: "A", Value: "a"},
		{Type: ast.String(), Name: "B", Value: "b"},
	}, ast.Default("b"))
	enumType.TypedDefault = ast.NewEnumDefault("b")

	input := []ast.Object{
		ast.NewObject(testPkgName, "Letter", ast.NewEnum([]ast.EnumValue{
			{Type: ast.String(), Name: "A", Value: "a"},
			{Type: ast.String(), Name: "B", Value: "b"},
		}, ast.Default("b"))),
	}
	expected := []ast.Object{
		ast.NewObject(testPkgName, "Letter", enumType),
	}

	runPassOnObjects(t, &DefaultAsTyped{}, input, expected)
}

func TestDefaultAsTyped_array(t *testing.T) {
	arrayType := ast.NewArray(ast.String(), ast.Default([]any{"a", "b"}))
	arrayType.TypedDefault = ast.NewArrayDefault([]ast.TypedDefault{
		{Kind: ast.KindScalar, Scalar: &ast.ScalarDefault{Value: "a"}},
		{Kind: ast.KindScalar, Scalar: &ast.ScalarDefault{Value: "b"}},
	})

	input := []ast.Object{
		ast.NewObject(testPkgName, "Words", ast.NewArray(ast.String(), ast.Default([]any{"a", "b"}))),
	}
	expected := []ast.Object{
		ast.NewObject(testPkgName, "Words", arrayType),
	}

	runPassOnObjects(t, &DefaultAsTyped{}, input, expected)
}

func TestDefaultAsTyped_map(t *testing.T) {
	mapType := ast.NewMap(ast.String(), ast.NewScalar(ast.KindInt32), ast.Default(map[string]any{"x": 1, "y": 2}))
	mapType.TypedDefault = ast.NewMapDefault(map[string]ast.TypedDefault{
		"x": {Kind: ast.KindScalar, Scalar: &ast.ScalarDefault{Value: 1}},
		"y": {Kind: ast.KindScalar, Scalar: &ast.ScalarDefault{Value: 2}},
	})

	input := []ast.Object{
		ast.NewObject(testPkgName, "Counts", ast.NewMap(ast.String(), ast.NewScalar(ast.KindInt32), ast.Default(map[string]any{"x": 1, "y": 2}))),
	}
	expected := []ast.Object{
		ast.NewObject(testPkgName, "Counts", mapType),
	}

	runPassOnObjects(t, &DefaultAsTyped{}, input, expected)
}

func TestDefaultAsTyped_struct(t *testing.T) {
	structType := ast.NewStruct(
		ast.NewStructField("Name", ast.String()),
		ast.NewStructField("Count", ast.NewScalar(ast.KindInt32)),
	)
	structType.Default = map[string]any{"Name": "foo", "Count": 7}
	structType.TypedDefault = ast.NewStructDefault(map[string]ast.TypedDefault{
		"Name":  {Kind: ast.KindScalar, Scalar: &ast.ScalarDefault{Value: "foo"}},
		"Count": {Kind: ast.KindScalar, Scalar: &ast.ScalarDefault{Value: 7}},
	})

	inputStruct := ast.NewStruct(
		ast.NewStructField("Name", ast.String()),
		ast.NewStructField("Count", ast.NewScalar(ast.KindInt32)),
	)
	inputStruct.Default = map[string]any{"Name": "foo", "Count": 7}

	input := []ast.Object{
		ast.NewObject(testPkgName, "Item", inputStruct),
	}
	expected := []ast.Object{
		ast.NewObject(testPkgName, "Item", structType),
	}

	runPassOnObjects(t, &DefaultAsTyped{}, input, expected)
}

func TestDefaultAsTyped_ref_withStructDefault(t *testing.T) {
	refType := ast.NewRef("test", "Item")
	refType.Default = map[string]any{"Name": "foo"}
	refType.TypedDefault = ast.NewRefDefault(ast.TypedDefault{
		Kind: ast.KindStruct,
		Struct: &ast.StructDefault{
			Fields: map[string]ast.TypedDefault{
				"Name": {Kind: ast.KindScalar, Scalar: &ast.ScalarDefault{Value: "foo"}},
			},
		},
	})

	inputRef := ast.NewRef("test", "Item")
	inputRef.Default = map[string]any{"Name": "foo"}

	input := []ast.Object{
		ast.NewObject(testPkgName, "Holder", ast.NewStruct(ast.NewStructField("It", inputRef))),
	}
	expected := []ast.Object{
		ast.NewObject(testPkgName, "Holder", ast.NewStruct(ast.NewStructField("It", refType))),
	}

	runPassOnObjects(t, &DefaultAsTyped{}, input, expected)
}

func TestDefaultAsTyped_ref_withScalarDefault(t *testing.T) {
	refType := ast.NewRef("test", "Letter")
	refType.Default = "b"
	refType.TypedDefault = ast.NewRefDefault(ast.TypedDefault{
		Kind:   ast.KindScalar,
		Scalar: &ast.ScalarDefault{Value: "b"},
	})

	inputRef := ast.NewRef("test", "Letter")
	inputRef.Default = "b"

	input := []ast.Object{
		ast.NewObject(testPkgName, "Holder", ast.NewStruct(ast.NewStructField("It", inputRef))),
	}
	expected := []ast.Object{
		ast.NewObject(testPkgName, "Holder", ast.NewStruct(ast.NewStructField("It", refType))),
	}

	runPassOnObjects(t, &DefaultAsTyped{}, input, expected)
}

func TestDefaultAsTyped_struct_withNestedStruct(t *testing.T) {
	// Nested struct field with a map default — exactly the case that motivates
	// this pass: avoid def.Default.(map[string]any) in jennies.
	innerStruct := ast.NewStruct(ast.NewStructField("Name", ast.String()))
	outerStruct := ast.NewStruct(ast.NewStructField("Inner", innerStruct))
	outerStruct.Default = map[string]any{"Inner": map[string]any{"Name": "foo"}}

	expectedOuter := outerStruct.DeepCopy()
	expectedOuter.TypedDefault = ast.NewStructDefault(map[string]ast.TypedDefault{
		"Inner": {
			Kind: ast.KindStruct,
			Struct: &ast.StructDefault{Fields: map[string]ast.TypedDefault{
				"Name": {Kind: ast.KindScalar, Scalar: &ast.ScalarDefault{Value: "foo"}},
			}},
		},
	})

	input := []ast.Object{ast.NewObject(testPkgName, "Outer", outerStruct)}
	expected := []ast.Object{ast.NewObject(testPkgName, "Outer", expectedOuter)}

	runPassOnObjects(t, &DefaultAsTyped{}, input, expected)
}

func TestDefaultAsTyped_noDefault(t *testing.T) {
	input := []ast.Object{
		ast.NewObject(testPkgName, "Item", ast.NewStruct(
			ast.NewStructField("Name", ast.String()),
		)),
	}
	// No default => no TypedDefault.
	runPassOnObjects(t, &DefaultAsTyped{}, input, input)
}
