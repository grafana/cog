package compiler

import (
	"testing"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/testutils"
)

func TestOmitFields(t *testing.T) {
	// Prepare test input
	schema := &ast.Schema{
		Package: "omit",
		Objects: testutils.ObjectsMap(
			ast.NewObject("omit", "AString", ast.String()),
			ast.NewObject("omit", "SomeObject", ast.NewStruct(
				ast.NewStructField("AString", ast.String()),
				ast.NewStructField("AnotherString", ast.String()),
			)),
			ast.NewObject("omit", "OtherObject", ast.NewStruct(
				ast.NewStructField("AString", ast.String()),
			)),
		),
	}
	expected := &ast.Schema{
		Package: schema.Package,
		Objects: testutils.ObjectsMap(
			ast.NewObject("omit", "AString", ast.String()),
			ast.NewObject("omit", "SomeObject", ast.NewStruct(
				ast.NewStructField("AnotherString", ast.String()),
			)),
			ast.NewObject("omit", "OtherObject", ast.NewStruct(
				ast.NewStructField("AString", ast.String()),
			)),
		),
	}

	pass := &OmitFields{
		Fields: []FieldReference{
			{Package: schema.Package, Object: "SomeObject", Field: "AString"},
		},
	}

	// Run the compiler pass
	runPassOnSchema(t, pass, schema, expected)
}
