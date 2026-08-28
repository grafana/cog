package transforms

import (
	"testing"

	"github.com/grafana/cog/internal/testutils"
	"github.com/grafana/cog/pkg/ir"
)

func TestFieldsSetDefault(t *testing.T) {
	// Prepare test input
	schema := &ir.Schema{
		Package: "set_required",
		Objects: testutils.ObjectsMap(
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
			ir.NewObject("set_required", "SomeObject", ir.NewStruct(
				ir.NewStructField("AString", ir.String(ir.Nullable(), ir.Default("default-foo")), ir.PassesTrail("FieldsSetDefault[default=default-foo]")),
				ir.NewStructField("AnotherString", ir.String(ir.Nullable())),
				ir.NewStructField("ABool", ir.String(ir.Nullable(), ir.Default(true)), ir.PassesTrail("FieldsSetDefault[default=true]")),
			)),
		),
	}

	pass := &FieldsSetDefault{
		DefaultValues: map[FieldReference]any{
			{Package: schema.Package, Object: "SomeObject", Field: "AString"}: "default-foo",
			{Package: schema.Package, Object: "SomeObject", Field: "ABool"}:   true,
		},
	}

	// Run the compiler pass
	runPassOnSchema(t, pass, schema, expected)
}
