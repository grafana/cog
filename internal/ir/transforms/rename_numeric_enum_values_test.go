package transforms

import (
	"testing"

	"github.com/grafana/cog/internal/ir"
)

func TestRenameNumericEnumValues(t *testing.T) {
	// Prepare test input
	objects := []ir.Object{
		ir.NewObject("pkg", "NotAnEnumStruct", ir.String()),

		ir.NewObject("pkg", "AnEnumWithNumericValues", ir.NewEnum([]ir.EnumValue{
			{Type: ir.NewScalar(ir.KindInt64), Name: "-1", Value: -1},
			{Type: ir.NewScalar(ir.KindInt64), Name: "1", Value: 1},
			{Type: ir.NewScalar(ir.KindInt64), Name: "2", Value: 2},
		})),

		ir.NewObject("pkg", "AnEnumWithNoNumericValues", ir.NewEnum([]ir.EnumValue{
			{Type: ir.String(), Name: "Hide", Value: "hide"},
			{Type: ir.String(), Name: "DontHide", Value: "dont_hide"},
		})),
	}

	// Prepare expected output
	expected := []ir.Object{
		ir.NewObject("pkg", "NotAnEnumStruct", ir.String()),

		ir.NewObject("pkg", "AnEnumWithNumericValues", ir.NewEnum([]ir.EnumValue{
			{Type: ir.NewScalar(ir.KindInt64), Name: "Negative1", Value: -1},
			{Type: ir.NewScalar(ir.KindInt64), Name: "N1", Value: 1},
			{Type: ir.NewScalar(ir.KindInt64), Name: "N2", Value: 2},
		}), "RenameNumericEnumValues"),

		ir.NewObject("pkg", "AnEnumWithNoNumericValues", ir.NewEnum([]ir.EnumValue{
			{Type: ir.String(), Name: "Hide", Value: "hide"},
			{Type: ir.String(), Name: "DontHide", Value: "dont_hide"},
		})),
	}

	// Run the compiler pass
	runPassOnObjects(t, &RenameNumericEnumValues{}, objects, expected)
}
