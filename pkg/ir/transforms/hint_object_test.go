package transforms

import (
	"testing"

	"github.com/grafana/cog/internal/testutils"
	"github.com/grafana/cog/pkg/ir"
)

func TestHintObject(t *testing.T) {
	// Prepare test input
	schema := &ir.Schema{
		Package: "hint_object",
		Objects: testutils.ObjectsMap(
			ir.NewObject("hint_object", "IWantHintsPlz", ir.String()),
		),
	}
	expected := &ir.Schema{
		Package: "hint_object",
		Objects: testutils.ObjectsMap(
			ir.NewObject(
				"hint_object",
				"IWantHintsPlz",
				ir.String(ir.Hints(ir.JenniesHints{"foo": "hint_value"})),
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
