package compiler

import (
	"testing"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/testutils"
)

func TestOmit(t *testing.T) {
	// Prepare test input
	schema := &ast.Schema{
		Package: "omit",
		Objects: testutils.ObjectsMap(
			ast.NewObject("omit", "AString", ast.String()),
			ast.NewObject("omit", "SomeObject", ast.NewStruct(
				ast.NewStructField("AString", ast.String()),
			)),
			ast.NewObject("omit", "OtherObject", ast.NewStruct(
				ast.NewStructField("Foo", ast.String()),
			)),
		),
	}
	expected := &ast.Schema{
		Package: "omit",
		Objects: testutils.ObjectsMap(
			ast.NewObject("omit", "OtherObject", ast.NewStruct(
				ast.NewStructField("Foo", ast.String()),
			)),
		),
	}

	pass := &Omit{
		Objects: []ObjectReference{
			{Package: schema.Package, Object: "AString"},
			{Package: schema.Package, Object: "SomeObject"},
			{Package: schema.Package, Object: "DoesNotExist"}, // no-op since it's not defined in the schema
		},
	}

	// Run the compiler pass
	runPassOnSchema(t, pass, schema, expected)
}
