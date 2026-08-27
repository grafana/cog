package transforms

import (
	"testing"

	"github.com/grafana/cog/internal/testutils"
	"github.com/grafana/cog/pkg/ir"
)

func TestFieldsSetRequired(t *testing.T) {
	// Prepare test input
	schema := &ir.Schema{
		Package: "set_required",
		Objects: testutils.ObjectsMap(
			ir.NewObject("set_required", "AString", ir.String()),
			ir.NewObject("set_required", "SomeObject", ir.NewStruct(
				ir.NewStructField("AString", ir.String(ir.Nullable())),
				ir.NewStructField("AnotherString", ir.String(ir.Nullable())),
				ir.NewStructField("ABool", ir.String(ir.Nullable())),
			)),
		),
	}
	expected := &ir.Schema{
		Package: "set_required",
		Objects: testutils.ObjectsMap(
			ir.NewObject("set_required", "AString", ir.String()),
			ir.NewObject("set_required", "SomeObject", ir.NewStruct(
				ir.NewStructField("AString", ir.String(), ir.Required(), ir.PassesTrail("FieldsSetRequired[nullable=false, required=true]")),
				ir.NewStructField("AnotherString", ir.String(ir.Nullable())),
				ir.NewStructField("ABool", ir.String(), ir.Required(), ir.PassesTrail("FieldsSetRequired[nullable=false, required=true]")),
			)),
		),
	}

	pass := &FieldsSetRequired{
		Fields: []FieldReference{
			// no-op: `AString` isn't a struct
			{Package: schema.Package, Object: "AString", Field: "Foo"},

			{Package: schema.Package, Object: "SomeObject", Field: "AString"},
			{Package: schema.Package, Object: "SomeObject", Field: "ABool"},
		},
	}

	// Run the compiler pass
	runPassOnSchema(t, pass, schema, expected)
}
