package compiler

import (
	"testing"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/testutils"
)

func TestHintObject(t *testing.T) {
	// Prepare test input
	schema := &ast.Schema{
		Package: "hint_object",
		Objects: testutils.ObjectsMap(
			ast.NewObject("hint_object", "IWantHintsPlz", ast.String()),
		),
	}
	expected := &ast.Schema{
		Package: "hint_object",
		Objects: testutils.ObjectsMap(
			ast.NewObject(
				"hint_object",
				"IWantHintsPlz",
				ast.String(ast.Hints(ast.JenniesHints{"foo": "hint_value"})),
				"HintObject[foo=hint_value]",
			),
		),
	}

	pass := &HintObject{
		Object: ObjectReference{Package: "hint_object", Object: "IWantHintsPlz"},
		Hints:  map[string]any{"foo": "hint_value"},
	}

	// Run the compiler pass
	runPassOnSchema(t, pass, schema, expected)
}
