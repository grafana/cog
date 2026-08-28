package transforms

import (
	"testing"

	"github.com/grafana/cog/pkg/ir"
)

func TestNotRequiredFieldAsNullableType(t *testing.T) {
	// Prepare test input
	objects := []ir.Object{
		ir.NewObject("pkg", "NotAStruct", ir.String()),

		ir.NewObject("pkg", "AStruct", ir.NewStruct(
			ir.NewStructField("RequiredString", ir.String(), ir.Required()),
			ir.NewStructField("RequiredNullableString", ir.String(ir.Nullable()), ir.Required()),
			ir.NewStructField("NotRequiredString", ir.String()),

			ir.NewStructField("RequiredRef", ir.NewRef("test", "SomeStruct"), ir.Required()),
			ir.NewStructField("RequiredNullableRef", ir.NewRef("test", "SomeStruct", ir.Nullable()), ir.Required()),
			ir.NewStructField("NotRequiredRef", ir.NewRef("test", "SomeStruct")),

			ir.NewStructField("NotRequiredArray", ir.NewArray(ir.String())),
			ir.NewStructField("RequiredArray", ir.NewArray(ir.String()), ir.Required()),

			ir.NewStructField("NotRequiredMap", ir.NewMap(
				ir.String(),
				ir.Bool(),
			)),
			ir.NewStructField("RequiredMap", ir.NewMap(
				ir.String(),
				ir.Bool(),
			), ir.Required()),
		)),
	}

	// Prepare expected output
	expected := []ir.Object{
		ir.NewObject("pkg", "NotAStruct", ir.String()),

		ir.NewObject("pkg", "AStruct", ir.NewStruct(
			ir.NewStructField("RequiredString", ir.String(), ir.Required()),
			ir.NewStructField("RequiredNullableString", ir.String(ir.Nullable()), ir.Required()),
			ir.NewStructField("NotRequiredString", ir.String(ir.Nullable()), ir.PassesTrail("NotRequiredFieldAsNullableType[nullable=true]")), // should become nullable

			ir.NewStructField("RequiredRef", ir.NewRef("test", "SomeStruct"), ir.Required()),
			ir.NewStructField("RequiredNullableRef", ir.NewRef("test", "SomeStruct", ir.Nullable()), ir.Required()),
			ir.NewStructField("NotRequiredRef", ir.NewRef("test", "SomeStruct", ir.Nullable()), ir.PassesTrail("NotRequiredFieldAsNullableType[nullable=true]")), // should become nullable

			ir.NewStructField("NotRequiredArray", ir.NewArray(ir.String(), ir.Nullable()), ir.PassesTrail("NotRequiredFieldAsNullableType[nullable=true]")), // should become nullable
			ir.NewStructField("RequiredArray", ir.NewArray(ir.String()), ir.Required()),

			ir.NewStructField("NotRequiredMap", ir.NewMap( // should become nullable
				ir.String(),
				ir.Bool(),
				ir.Nullable(),
			), ir.PassesTrail("NotRequiredFieldAsNullableType[nullable=true]")),
			ir.NewStructField("RequiredMap", ir.NewMap(
				ir.String(),
				ir.Bool(),
			), ir.Required()),
		)),
	}

	// Run the compiler pass
	runPassOnObjects(t, &NotRequiredFieldAsNullableType{}, objects, expected)
}
