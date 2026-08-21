package compiler

import (
	"testing"

	"github.com/grafana/cog/internal/testutils"
	"github.com/grafana/cog/pkg/ast"
)

func TestRetypeObject(t *testing.T) {
	// Prepare test input
	schema := &ast.Schema{
		Package: "retype_object",
		Objects: testutils.ObjectsMap(
			ast.NewObject("retype_object", "SomeObject", ast.NewStruct(
				ast.NewStructField("AString", ast.String()),
			)),
		),
	}
	expected := &ast.Schema{
		Package: "retype_object",
		Objects: testutils.ObjectsMap(
			ast.NewObject("retype_object", "SomeObject", ast.Bool(), "RetypeObject[Struct → Bool]"),
		),
	}

	pass := &RetypeObject{
		Object: ObjectReference{Package: schema.Package, Object: "SomeObject"},
		As:     ast.Bool(),
	}

	// Run the compiler pass
	runPassOnSchema(t, pass, schema, expected)
}
