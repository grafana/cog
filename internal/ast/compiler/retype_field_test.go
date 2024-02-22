package compiler

import (
	"testing"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/testutils"
)

func TestRetypeField(t *testing.T) {
	// Prepare test input
	schema := &ast.Schema{
		Package: "retype_field",
		Objects: testutils.ObjectsMap(
			ast.NewObject("retype_field", "SomeObject", ast.NewStruct(
				ast.NewStructField("AString", ast.String()),
			)),
		),
	}
	expected := &ast.Schema{
		Package: "retype_field",
		Objects: testutils.ObjectsMap(
			ast.NewObject("retype_field", "SomeObject", ast.NewStruct(
				ast.NewStructField("AString", ast.Bool(), ast.PassesTrail("RetypeField[String → Bool]")),
			)),
		),
	}

	pass := &RetypeField{
		Field: FieldReference{Package: schema.Package, Object: "SomeObject", Field: "AString"},
		As:    ast.Bool(),
	}

	// Run the compiler pass
	runPassOnSchema(t, pass, schema, expected)
}

func TestRetypeField_notFoundFieldRef(t *testing.T) {
	// Prepare test input
	schema := &ast.Schema{
		Package: "retype_field",
		Objects: testutils.ObjectsMap(
			ast.NewObject("retype_field", "SomeObject", ast.NewStruct(
				ast.NewStructField("AString", ast.String()),
			)),
		),
	}

	pass := &RetypeField{
		// no-op since `SomeObject.NotFound` does not exist
		Field: FieldReference{Package: schema.Package, Object: "SomeObject", Field: "NotFound"},
		As:    ast.Bool(),
	}

	// Run the compiler pass
	runPassOnSchema(t, pass, schema, schema)
}

func TestRetypeField_onNonStruct(t *testing.T) {
	// Prepare test input
	schema := &ast.Schema{
		Package: "retype_field",
		Objects: testutils.ObjectsMap(
			ast.NewObject("retype_field", "AString", ast.String()),
		),
	}

	pass := &RetypeField{
		// no-op since `AString` is not a struct
		Field: FieldReference{Package: schema.Package, Object: "AString", Field: "NotAField"},
		As:    ast.Bool(),
	}

	// Run the compiler pass
	runPassOnSchema(t, pass, schema, schema)
}
