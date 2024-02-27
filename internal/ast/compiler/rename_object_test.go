package compiler

import (
	"testing"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/testutils"
)

func TestRenameObject(t *testing.T) {
	// Prepare test input
	schema := &ast.Schema{
		Package: "rename_object",
		Objects: testutils.ObjectsMap(
			ast.NewObject("rename_object", "SomeObject", ast.NewStruct(
				ast.NewStructField("foo", ast.String(ast.Nullable())),
				ast.NewStructField("ref_to_nice_object", ast.NewRef("rename_object", "NotANiceName")),
			)),
			ast.NewObject("rename_object", "NotANiceName", ast.NewStruct(
				ast.NewStructField("AString", ast.String(ast.Nullable())),
			)),
		),
	}
	expected := &ast.Schema{
		Package: "rename_object",
		Objects: testutils.ObjectsMap(
			ast.NewObject("rename_object", "SomeObject", ast.NewStruct(
				ast.NewStructField("foo", ast.String(ast.Nullable())),
				ast.NewStructField("ref_to_nice_object", ast.NewRef("rename_object", "ReallyNiceName")),
			)),
			ast.NewObject("rename_object", "ReallyNiceName", ast.NewStruct(
				ast.NewStructField("AString", ast.String(ast.Nullable())),
			), "RenameObject[NotANiceName → ReallyNiceName]"),
		),
	}

	pass := &RenameObject{
		From: ObjectReference{Package: schema.Package, Object: "NotANiceName"},
		To:   "ReallyNiceName",
	}

	// Run the compiler pass
	runPassOnSchema(t, pass, schema, expected)
}
