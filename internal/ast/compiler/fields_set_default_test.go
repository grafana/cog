package compiler

import (
	"testing"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/testutils"
)

func TestFieldsSetDefault(t *testing.T) {
	// Prepare test input
	schema := &ast.Schema{
		Package: "set_required",
		Objects: testutils.ObjectsMap(
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
			ast.NewObject("set_required", "SomeObject", ast.NewStruct(
				ast.NewStructField("AString", ast.String(ast.Nullable(), ast.Default("default-foo")), ast.PassesTrail("FieldsSetDefault[default=default-foo]")),
				ast.NewStructField("AnotherString", ast.String(ast.Nullable())),
				ast.NewStructField("ABool", ast.String(ast.Nullable(), ast.Default(true)), ast.PassesTrail("FieldsSetDefault[default=true]")),
			)),
		),
	}

	pass := &FieldsSetDefault{
		DefaultValues: map[FieldReference]any{
			{Package: schema.Package, Object: "SomeObject", Field: "AString"}: "default-foo",
			{Package: schema.Package, Object: "SomeObject", Field: "ABool"}:   true,
		},
	}

	// Run the compiler pass
	runPassOnSchema(t, pass, schema, expected)
}
