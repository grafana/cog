package transforms

import (
	"testing"

	"github.com/grafana/cog/internal/testutils"
	"github.com/grafana/cog/pkg/ir"
)

func TestReplaceReference(t *testing.T) {
	// Prepare test input
	schema := &ir.Schema{
		Package: "replace",
		Objects: testutils.ObjectsMap(
			ir.NewObject("replace", "SomeObject", ir.NewStruct(
				ir.NewStructField("ARef", ir.NewRef("common", "Bar")),
				ir.NewStructField("AString", ir.String()),
				ir.NewStructField("AReplacedRef", ir.NewRef("replace", "BadRef")),
			)),
		),
	}
	expected := &ir.Schema{
		Package: "replace",
		Objects: testutils.ObjectsMap(
			ir.NewObject("replace", "SomeObject", ir.NewStruct(
				ir.NewStructField("ARef", ir.NewRef("common", "Bar")),
				ir.NewStructField("AString", ir.String()),
				ir.NewStructField("AReplacedRef", ir.NewRef("common", "Ref", ir.Trail("ReplaceReference[replace.BadRef → common.Ref]"))),
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
