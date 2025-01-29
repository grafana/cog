package compiler

import (
	"testing"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/testutils"
)

func TestReplaceReference(t *testing.T) {
	// Prepare test input
	schema := &ast.Schema{
		Package: "replace",
		Objects: testutils.ObjectsMap(
			ast.NewObject("replace", "SomeObject", ast.NewStruct(
				ast.NewStructField("ARef", ast.NewRef("common", "Bar")),
				ast.NewStructField("AString", ast.String()),
				ast.NewStructField("AReplacedRef", ast.NewRef("replace", "BadRef")),
			)),
		),
	}
	expected := &ast.Schema{
		Package: "replace",
		Objects: testutils.ObjectsMap(
			ast.NewObject("replace", "SomeObject", ast.NewStruct(
				ast.NewStructField("ARef", ast.NewRef("common", "Bar")),
				ast.NewStructField("AString", ast.String()),
				ast.NewStructField("AReplacedRef", ast.NewRef("common", "Ref", ast.Trail("ReplaceReference[replace.BadRef → common.Ref]"))),
			)),
		),
	}

	pass := &ReplaceReference{
		From: ObjectReference{Package: "replace", Object: "BadRef"},
		To:   ObjectReference{Package: "common", Object: "Ref"},
	}

	// Run the compiler pass
	runPassOnSchema(t, pass, schema, expected)
}
