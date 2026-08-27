package transforms

import (
	"testing"

	"github.com/grafana/cog/internal/testutils"
	"github.com/grafana/cog/pkg/ir"
)

func TestRetypeObject(t *testing.T) {
	// Prepare test input
	schema := &ir.Schema{
		Package: "retype_object",
		Objects: testutils.ObjectsMap(
			ir.NewObject("retype_object", "SomeObject", ir.NewStruct(
				ir.NewStructField("AString", ir.String()),
			)),
		),
	}
	expected := &ir.Schema{
		Package: "retype_object",
		Objects: testutils.ObjectsMap(
			ir.NewObject("retype_object", "SomeObject", ir.Bool(), "RetypeObject[Struct → Bool]"),
		),
	}

	pass := &RetypeObject{
		Object: ObjectReference{Package: schema.Package, Object: "SomeObject"},
		As:     ir.Bool(),
	}

	// Run the compiler pass
	runPassOnSchema(t, pass, schema, expected)
}
