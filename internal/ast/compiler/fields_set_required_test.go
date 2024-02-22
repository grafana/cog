package compiler

import (
	"testing"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/testutils"
)

func TestFieldsSetRequired(t *testing.T) {
	// Prepare test input
	schema := &ast.Schema{
		Package: "set_required",
		Objects: testutils.ObjectsMap(
			ast.NewObject("set_required", "AString", ast.String()),
			ast.NewObject("set_required", "SomeObject", ast.NewStruct(
				ast.NewStructField("AString", ast.String(ast.Nullable())),
				ast.NewStructField("AnotherString", ast.String(ast.Nullable())),
				ast.NewStructField("ABool", ast.String(ast.Nullable())),
			)),
		),
	}
	expected := &ast.Schema{
		Package: "set_required",
		Objects: testutils.ObjectsMap(
			ast.NewObject("set_required", "AString", ast.String()),
			ast.NewObject("set_required", "SomeObject", ast.NewStruct(
				ast.NewStructField("AString", ast.String(), ast.Required(), ast.PassesTrail("FieldsSetRequired[nullable=false, required=true]")),
				ast.NewStructField("AnotherString", ast.String(ast.Nullable())),
				ast.NewStructField("ABool", ast.String(), ast.Required(), ast.PassesTrail("FieldsSetRequired[nullable=false, required=true]")),
			)),
		),
	}

	pass := &FieldsSetRequired{
		Fields: []FieldReference{
			// no-op: `AString` isn't a struct
			{Package: schema.Package, Object: "AString", Field: "Foo"},

			{Package: schema.Package, Object: "SomeObject", Field: "AString"},
			{Package: schema.Package, Object: "SomeObject", Field: "ABool"},
		},
	}

	// Run the compiler pass
	runPassOnSchema(t, pass, schema, expected)
}
