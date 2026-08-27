package transforms

import (
	"testing"

	"github.com/grafana/cog/pkg/ir"
)

func TestTrimEnumValues(t *testing.T) {
	objects := []ir.Object{
		ir.NewObject("enum_with_leading_and_trailing_spaces", "MyEnum", ir.NewEnum([]ir.EnumValue{
			{Type: ir.String(), Name: "Leading", Value: " Leading"},
			{Type: ir.String(), Name: "Trailing", Value: "Trailing "},
			{Type: ir.String(), Name: "Both", Value: " Both "},
			{Type: ir.String(), Name: "SpacesInMiddle", Value: "Spaces in middle"},
		})),
	}

	expected := []ir.Object{
		ir.NewObject("enum_with_leading_and_trailing_spaces", "MyEnum", ir.NewEnum([]ir.EnumValue{
			{Type: ir.String(), Name: "Leading", Value: "Leading"},
			{Type: ir.String(), Name: "Trailing", Value: "Trailing"},
			{Type: ir.String(), Name: "Both", Value: "Both"},
			{Type: ir.String(), Name: "SpacesInMiddle", Value: "Spaces in middle"},
		})),
	}

	runPassOnObjects(t, &TrimEnumValues{}, objects, expected)
}
