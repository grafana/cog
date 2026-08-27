package transforms

import (
	"testing"

	"github.com/grafana/cog/internal/testutils"
	"github.com/grafana/cog/pkg/ir"
)

func TestAnonymousEnumToExplicitType_withNoAnonymousEnum(t *testing.T) {
	// Prepare test input
	schema := &ir.Schema{
		Package: "without_enums",
		Objects: testutils.ObjectsMap(
			ir.NewObject("without_enums", "AString", ir.String()),
			ir.NewObject("without_enums", "AStruct", ir.NewStruct(
				ir.NewStructField("AString", ir.String()),
			)),
		),
	}

	// Run the compiler pass
	runPassOnSchema(t, &AnonymousEnumToExplicitType{}, schema, schema)
}

func TestAnonymousEnumToExplicitType_withAnonymousEnumInStruct(t *testing.T) {
	// Prepare test input
	schema := &ir.Schema{
		Package: "with_enums",
		Objects: testutils.ObjectsMap(
			ir.NewObject("with_enums", "Panel", ir.NewStruct(
				ir.NewStructField("title", ir.String()),
				ir.NewStructField("type", ir.NewEnum([]ir.EnumValue{
					{Name: "Foo", Value: "foo", Type: ir.String()},
					{Name: "Bar", Value: "bar", Type: ir.String()},
				}, ir.Nullable())),
			)),
			ir.NewObject("with_enums", "Mode", ir.NewEnum([]ir.EnumValue{
				{Name: "Auto", Value: "auto", Type: ir.String()},
				{Name: "Manual", Value: "manual", Type: ir.String()},
			})),
		),
	}

	// Prepare expected output
	expected := &ir.Schema{
		Package: "with_enums",
		Objects: testutils.ObjectsMap(
			ir.NewObject("with_enums", "Panel", ir.NewStruct(
				ir.NewStructField("title", ir.String()),
				ir.NewStructField("type", ir.NewRef("with_enums", "PanelType", ir.Nullable(), ir.Trail("AnonymousEnumToExplicitType"))),
			)),

			// this object is unchanged
			schema.Objects.Get("Mode"),

			// the anonymous enum, turned into an object
			ir.NewObject("with_enums", "PanelType", ir.NewEnum([]ir.EnumValue{
				{Name: "Foo", Value: "foo", Type: ir.String()},
				{Name: "Bar", Value: "bar", Type: ir.String()},
			}), "AnonymousEnumToExplicitType"),
		),
	}

	// Run the compiler pass
	runPassOnSchema(t, &AnonymousEnumToExplicitType{}, schema, expected)
}

func TestAnonymousEnumToExplicitType_withAnonymousEnumInArray(t *testing.T) {
	// Prepare test input
	schema := &ir.Schema{
		Package: "in_array",
		Objects: testutils.ObjectsMap(
			ir.NewObject("in_array", "TypesList", ir.NewArray(
				ir.NewEnum([]ir.EnumValue{
					{Name: "Foo", Value: "foo", Type: ir.String()},
					{Name: "Bar", Value: "bar", Type: ir.String()},
				}),
			)),
		),
	}

	// Prepare expected output
	expected := &ir.Schema{
		Package: "in_array",
		Objects: testutils.ObjectsMap(
			ir.NewObject("in_array", "TypesList", ir.NewArray(
				ir.NewRef("in_array", "TypesListEnum", ir.Trail("AnonymousEnumToExplicitType"))),
			),

			// the anonymous enum, turned into an object
			ir.NewObject("in_array", "TypesListEnum", ir.NewEnum([]ir.EnumValue{
				{Name: "Foo", Value: "foo", Type: ir.String()},
				{Name: "Bar", Value: "bar", Type: ir.String()},
			}), "AnonymousEnumToExplicitType"),
		),
	}

	// Run the compiler pass
	runPassOnSchema(t, &AnonymousEnumToExplicitType{}, schema, expected)
}

func TestAnonymousEnumToExplicitType_withAnonymousEnumInMap(t *testing.T) {
	// Prepare test input
	schema := &ir.Schema{
		Package: "in_map",
		Objects: testutils.ObjectsMap(
			ir.NewObject("in_map", "MapOfThings", ir.NewMap(
				ir.String(),
				ir.NewEnum([]ir.EnumValue{
					{Name: "Foo", Value: "foo", Type: ir.String()},
					{Name: "Bar", Value: "bar", Type: ir.String()},
				}),
			)),
		),
	}

	// Prepare expected output
	expected := &ir.Schema{
		Package: "in_map",
		Objects: testutils.ObjectsMap(
			ir.NewObject("in_map", "MapOfThings", ir.NewMap(
				ir.String(),
				ir.NewRef("in_map", "MapOfThingsEnum", ir.Trail("AnonymousEnumToExplicitType"))),
			),

			// the anonymous enum, turned into an object
			ir.NewObject("in_map", "MapOfThingsEnum", ir.NewEnum([]ir.EnumValue{
				{Name: "Foo", Value: "foo", Type: ir.String()},
				{Name: "Bar", Value: "bar", Type: ir.String()},
			}), "AnonymousEnumToExplicitType"),
		),
	}

	// Run the compiler pass
	runPassOnSchema(t, &AnonymousEnumToExplicitType{}, schema, expected)
}

func TestAnonymousEnumToExplicitType_withAnonymousEnumInDisjunction(t *testing.T) {
	// Prepare test input
	schema := &ir.Schema{
		Package: "in_disjunction",
		Objects: testutils.ObjectsMap(
			ir.NewObject("in_disjunction", "DisjunctionOfThings", ir.NewDisjunction([]ir.Type{
				ir.String(),
				ir.NewEnum([]ir.EnumValue{
					{Name: "Foo", Value: "foo", Type: ir.String()},
					{Name: "Bar", Value: "bar", Type: ir.String()},
				}),
			})),
		),
	}

	// Prepare expected output
	expected := &ir.Schema{
		Package: "in_disjunction",
		Objects: testutils.ObjectsMap(
			ir.NewObject("in_disjunction", "DisjunctionOfThings", ir.NewDisjunction([]ir.Type{
				ir.String(),
				ir.NewRef("in_disjunction", "DisjunctionOfThingsEnum", ir.Trail("AnonymousEnumToExplicitType")),
			})),

			// the anonymous enum, turned into an object
			ir.NewObject("in_disjunction", "DisjunctionOfThingsEnum", ir.NewEnum([]ir.EnumValue{
				{Name: "Foo", Value: "foo", Type: ir.String()},
				{Name: "Bar", Value: "bar", Type: ir.String()},
			}), "AnonymousEnumToExplicitType"),
		),
	}

	// Run the compiler pass
	runPassOnSchema(t, &AnonymousEnumToExplicitType{}, schema, expected)
}

func TestAnonymousEnumToExplicitType_withAnonymousEnumInIntersection(t *testing.T) {
	// Prepare test input
	schema := &ir.Schema{
		Package: "in_intersection",
		Objects: testutils.ObjectsMap(
			ir.NewObject("in_intersection", "IntersectionOfThings", ir.NewIntersection([]ir.Type{
				ir.String(),
				ir.NewEnum([]ir.EnumValue{
					{Name: "Foo", Value: "foo", Type: ir.String()},
					{Name: "Bar", Value: "bar", Type: ir.String()},
				}),
			})),
		),
	}

	// Prepare expected output
	expected := &ir.Schema{
		Package: "in_intersection",
		Objects: testutils.ObjectsMap(
			ir.NewObject("in_intersection", "IntersectionOfThings", ir.NewIntersection([]ir.Type{
				ir.String(),
				ir.NewRef("in_intersection", "IntersectionOfThingsEnum", ir.Trail("AnonymousEnumToExplicitType")),
			})),

			// the anonymous enum, turned into an object
			ir.NewObject("in_intersection", "IntersectionOfThingsEnum", ir.NewEnum([]ir.EnumValue{
				{Name: "Foo", Value: "foo", Type: ir.String()},
				{Name: "Bar", Value: "bar", Type: ir.String()},
			}), "AnonymousEnumToExplicitType"),
		),
	}

	// Run the compiler pass
	runPassOnSchema(t, &AnonymousEnumToExplicitType{}, schema, expected)
}

func TestAnonymousEnumToExplicitType_withFieldWithDefaultValue(t *testing.T) {
	// Prepare expected input
	schema := &ir.Schema{
		Package: "with_default",
		Objects: testutils.ObjectsMap(
			ir.NewObject("with_default", "Panel", ir.NewStruct(
				ir.NewStructField("title", ir.String()),
				ir.NewStructField("type", ir.NewEnum([]ir.EnumValue{
					{Name: "Foo", Value: "foo", Type: ir.String()},
					{Name: "Bar", Value: "bar", Type: ir.String()},
				}, ir.Nullable(), ir.Default("foo"))),
			)),
		),
	}

	// Prepare expected output
	expected := &ir.Schema{
		Package: "with_default",
		Objects: testutils.ObjectsMap(
			ir.NewObject("with_default", "Panel", ir.NewStruct(
				ir.NewStructField("title", ir.String()),
				ir.NewStructField("type", ir.NewRef("with_default", "PanelType", ir.Nullable(), ir.Default("foo"), ir.Trail("AnonymousEnumToExplicitType"))),
			)),

			// the anonymous enum, turned into an object
			ir.NewObject("with_default", "PanelType", ir.NewEnum([]ir.EnumValue{
				{Name: "Foo", Value: "foo", Type: ir.String()},
				{Name: "Bar", Value: "bar", Type: ir.String()},
			}), "AnonymousEnumToExplicitType"),
		),
	}

	// Run the compiler pass
	runPassOnSchema(t, &AnonymousEnumToExplicitType{}, schema, expected)
}
