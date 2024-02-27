package compiler

import (
	"testing"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/testutils"
)

func TestNameAnonymousStruct(t *testing.T) {
	// Prepare test input
	schema := &ast.Schema{
		Package: "name_anonymous_struct",
		Objects: testutils.ObjectsMap(
			ast.NewObject("name_anonymous_struct", "SomeObject", ast.NewStruct(
				ast.NewStructField("inner", ast.NewStruct(
					ast.NewStructField("title", ast.String()),
				)),
			)),
		),
	}
	expected := &ast.Schema{
		Package: "name_anonymous_struct",
		Objects: testutils.ObjectsMap(
			ast.NewObject("name_anonymous_struct", "SomeObject", ast.NewStruct(
				ast.NewStructField("inner", ast.NewRef(schema.Package, "Inner")),
			)),
			ast.NewObject("name_anonymous_struct", "Inner", ast.NewStruct(
				ast.NewStructField("title", ast.String()),
			), "NameAnonymousStruct"),
		),
	}

	pass := &NameAnonymousStruct{
		Field: FieldReference{Package: schema.Package, Object: "SomeObject", Field: "inner"},
		As:    "Inner",
	}

	// Run the compiler pass
	runPassOnSchema(t, pass, schema, expected)
}

func TestNameAnonymousStruct_onNonStructObject(t *testing.T) {
	// Prepare test input
	schema := &ast.Schema{
		Package: "name_anonymous_struct",
		Objects: testutils.ObjectsMap(
			ast.NewObject("name_anonymous_struct", "SomeObject", ast.NewStruct(
				ast.NewStructField("inner", ast.NewStruct(
					ast.NewStructField("title", ast.String()),
				)),
			)),
		),
	}

	pass := &NameAnonymousStruct{
		// no-op since `doesNotExist` does not exist
		Field: FieldReference{Package: schema.Package, Object: "SomeObject", Field: "doesNotExist"},
		As:    "Inner",
	}

	// Run the compiler pass
	runPassOnSchema(t, pass, schema, schema)
}

func TestNameAnonymousStruct_onNonStructField(t *testing.T) {
	// Prepare test input
	schema := &ast.Schema{
		Package: "name_anonymous_struct",
		Objects: testutils.ObjectsMap(
			ast.NewObject("name_anonymous_struct", "SomeObject", ast.NewStruct(
				ast.NewStructField("inner", ast.Bool()),
			)),
		),
	}

	pass := &NameAnonymousStruct{
		// no-op since `AString` is not a struct
		Field: FieldReference{Package: schema.Package, Object: "SomeObject", Field: "inner"},
		As:    "Inner",
	}

	// Run the compiler pass
	runPassOnSchema(t, pass, schema, schema)
}

func TestNameAnonymousStruct_onUnknownField(t *testing.T) {
	// Prepare test input
	schema := &ast.Schema{
		Package: "name_anonymous_struct",
		Objects: testutils.ObjectsMap(
			ast.NewObject("name_anonymous_struct", "AString", ast.String()),
		),
	}

	pass := &NameAnonymousStruct{
		// no-op since `AString` is not a struct
		Field: FieldReference{Package: schema.Package, Object: "AString", Field: "doesNotExist"},
		As:    "Inner",
	}

	// Run the compiler pass
	runPassOnSchema(t, pass, schema, schema)
}
