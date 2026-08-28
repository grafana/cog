package transforms

import (
	"testing"

	"github.com/grafana/cog/internal/testutils"
	"github.com/grafana/cog/pkg/ir"
)

func TestOmitFields(t *testing.T) {
	// Prepare test input
	schema := &ir.Schema{
		Package: "omit",
		Objects: testutils.ObjectsMap(
			ir.NewObject("omit", "AString", ir.String()),
			ir.NewObject("omit", "SomeObject", ir.NewStruct(
				ir.NewStructField("AString", ir.String()),
				ir.NewStructField("AnotherString", ir.String()),
			)),
			ir.NewObject("omit", "OtherObject", ir.NewStruct(
				ir.NewStructField("AString", ir.String()),
			)),
		),
	}
	expected := &ir.Schema{
		Package: schema.Package,
		Objects: testutils.ObjectsMap(
			ir.NewObject("omit", "AString", ir.String()),
			ir.NewObject("omit", "SomeObject", ir.NewStruct(
				ir.NewStructField("AnotherString", ir.String()),
			)),
			ir.NewObject("omit", "OtherObject", ir.NewStruct(
				ir.NewStructField("AString", ir.String()),
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
