package transforms

import (
	"testing"

	"github.com/grafana/cog/internal/testutils"
	"github.com/grafana/cog/pkg/ir"
)

func TestTrimObjectNamePrefix(t *testing.T) {
	// Prepare test input
	schema := &ir.Schema{
		Package: "prefix_names",
		Objects: testutils.ObjectsMap(
			ir.NewObject("prefix_names", "MyPrefixSomeObject", ir.NewStruct(
				ir.NewStructField("foo", ir.String(ir.Nullable())),
				ir.NewStructField("ref_to_nice_object", ir.NewRef("prefix_names", "MyPrefixNotANiceName")),
			)),
			ir.NewObject("prefix_names", "MyPrefixNotANiceName", ir.NewStruct(
				ir.NewStructField("AString", ir.String(ir.Nullable())),
			)),
			ir.NewObject("prefix_names", "VariableRefresh", ir.NewEnum([]ir.EnumValue{
				{Name: "MyPrefixNever", Value: "never", Type: ir.String()},
				{Name: "MyPrefixAlways", Value: "always", Type: ir.String()},
			})),
		),
	}
	expected := &ir.Schema{
		Package: "prefix_names",
		Objects: testutils.ObjectsMap(
			ir.NewObject("prefix_names", "SomeObject", ir.NewStruct(
				ir.NewStructField("foo", ir.String(ir.Nullable())),
				ir.NewStructField("ref_to_nice_object", ir.NewRef("prefix_names", "NotANiceName", ir.Trail("TrimObjectNamePrefix[MyPrefixNotANiceName → NotANiceName]"))),
			), "TrimObjectNamePrefix[MyPrefixSomeObject → SomeObject]"),
			ir.NewObject("prefix_names", "NotANiceName", ir.NewStruct(
				ir.NewStructField("AString", ir.String(ir.Nullable())),
			), "TrimObjectNamePrefix[MyPrefixNotANiceName → NotANiceName]"),
			ir.NewObject("prefix_names", "VariableRefresh", ir.NewEnum([]ir.EnumValue{
				{Name: "Never", Value: "never", Type: ir.String()},
				{Name: "Always", Value: "always", Type: ir.String()},
			})),
		),
	}

	pass := &TrimObjectNamePrefix{
		Prefix: "MyPrefix",
	}

	// Run the compiler pass
	runPassOnSchema(t, pass, schema, expected)
}
