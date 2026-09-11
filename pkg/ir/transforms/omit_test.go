package transforms

import (
	"testing"

	"github.com/grafana/cog/internal/testutils"
	"github.com/grafana/cog/pkg/ir"
)

func TestOmit(t *testing.T) {
	// Prepare test input
	schema := &ir.Schema{
		Package: "omit",
		Objects: testutils.ObjectsMap(
			ir.NewObject("omit", "AString", ir.String()),
			ir.NewObject("omit", "SomeObject", ir.NewStruct(
				ir.NewStructField("AString", ir.String()),
			)),
			ir.NewObject("omit", "OtherObject", ir.NewStruct(
				ir.NewStructField("Foo", ir.String()),
			)),
		),
	}
	expected := &ir.Schema{
		Package: "omit",
		Objects: testutils.ObjectsMap(
			ir.NewObject("omit", "OtherObject", ir.NewStruct(
				ir.NewStructField("Foo", ir.String()),
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
