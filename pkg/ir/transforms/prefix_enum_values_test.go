package transforms

import (
	"testing"

	"github.com/grafana/cog/pkg/ir"
)

func TestPrefixEnumValues(t *testing.T) {
	// Prepare test input
	objects := []ir.Object{
		ir.NewObject("pkg", "VariableRefresh", ir.NewEnum([]ir.EnumValue{
			{Name: "Never", Value: "never", Type: ir.String()},
			{Name: "Always", Value: "always", Type: ir.String()},
		})),

		ir.NewObject("pkg", "SomeType", ir.String()),
	}

	// Prepare expected output
	expected := []ir.Object{
		ir.NewObject("pkg", "VariableRefresh", ir.NewEnum([]ir.EnumValue{
			{Name: "VariableRefreshNever", Value: "never", Type: ir.String()},
			{Name: "VariableRefreshAlways", Value: "always", Type: ir.String()},
		}), "PrefixEnumValues"),

		ir.NewObject("pkg", "SomeType", ir.String()),
	}

	// Run the compiler pass
	runPassOnObjects(t, &PrefixEnumValues{}, objects, expected)
}

func TestPrefixEnumValuesWithNegativeIntegerName(t *testing.T) {
	// Prepare test input
	objects := []ir.Object{
		ir.NewObject("pkg", "BarAlignment", ir.NewEnum([]ir.EnumValue{
			{Name: "1", Value: 1, Type: ir.NewScalar(ir.KindInt64)},
			{Name: "-1", Value: -1, Type: ir.NewScalar(ir.KindInt64)},
		})),
	}

	// Prepare expected output
	expected := []ir.Object{
		ir.NewObject("pkg", "BarAlignment", ir.NewEnum([]ir.EnumValue{
			{Name: "BarAlignment1", Value: 1, Type: ir.NewScalar(ir.KindInt64)},
			{Name: "BarAlignmentNegative1", Value: -1, Type: ir.NewScalar(ir.KindInt64)},
		}), "PrefixEnumValues"),
	}

	// Run the compiler pass
	runPassOnObjects(t, &PrefixEnumValues{}, objects, expected)
}

func TestPrefixEnumValuesWithEmptyStringMember(t *testing.T) {
	// Prepare test input
	objects := []ir.Object{
		ir.NewObject("pkg", "BarAlignment", ir.NewEnum([]ir.EnumValue{
			{Name: "", Value: "", Type: ir.String()},
			{Name: "foo", Value: "foo", Type: ir.String()},
		})),
	}

	// Prepare expected output
	expected := []ir.Object{
		ir.NewObject("pkg", "BarAlignment", ir.NewEnum([]ir.EnumValue{
			{Name: "BarAlignmentNone", Value: "", Type: ir.String()},
			{Name: "BarAlignmentFoo", Value: "foo", Type: ir.String()},
		}), "PrefixEnumValues"),
	}

	// Run the compiler pass
	runPassOnObjects(t, &PrefixEnumValues{}, objects, expected)
}
