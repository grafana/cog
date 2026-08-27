package transforms

import (
	"testing"

	"github.com/grafana/cog/internal/testutils"
	"github.com/grafana/cog/pkg/ir"
)

func TestPrefixObjectNames(t *testing.T) {
	// Prepare test input
	schema := &ir.Schema{
		Package: "prefix_names",
		Objects: testutils.ObjectsMap(
			ir.NewObject("prefix_names", "SomeObject", ir.NewStruct(
				ir.NewStructField("foo", ir.String(ir.Nullable())),
				ir.NewStructField("ref_to_nice_object", ir.NewRef("prefix_names", "NotANiceName")),
			)),
			ir.NewObject("prefix_names", "NotANiceName", ir.NewStruct(
				ir.NewStructField("AString", ir.String(ir.Nullable())),
			)),
			ir.NewObject("prefix_names", "VariableRefresh", ir.NewEnum([]ir.EnumValue{
				{Name: "Never", Value: "never", Type: ir.String()},
				{Name: "Always", Value: "always", Type: ir.String()},
			})),
		),
	}
	expected := &ir.Schema{
		Package: "prefix_names",
		Objects: testutils.ObjectsMap(
			ir.NewObject("prefix_names", "PreSomeObject", ir.NewStruct(
				ir.NewStructField("foo", ir.String(ir.Nullable())),
				ir.NewStructField("ref_to_nice_object", ir.NewRef("prefix_names", "PreNotANiceName", ir.Trail("PrefixObjectNames[NotANiceName → PreNotANiceName]"))),
			), "PrefixObjectNames[SomeObject → PreSomeObject]"),
			ir.NewObject("prefix_names", "PreNotANiceName", ir.NewStruct(
				ir.NewStructField("AString", ir.String(ir.Nullable())),
			), "PrefixObjectNames[NotANiceName → PreNotANiceName]"),
			ir.NewObject("prefix_names", "PreVariableRefresh", ir.NewEnum([]ir.EnumValue{
				{Name: "PreNever", Value: "never", Type: ir.String()},
				{Name: "PreAlways", Value: "always", Type: ir.String()},
			}), "PrefixObjectNames[VariableRefresh → PreVariableRefresh]"),
		),
	}

	pass := &PrefixObjectNames{
		Prefix: "Pre",
	}

	// Run the compiler pass
	runPassOnSchema(t, pass, schema, expected)
}
