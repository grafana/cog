package compiler

import (
	"testing"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/testutils"
)

func TestAnonymousEnumToExplicitType_withNoAnonymousEnum(t *testing.T) {
	// Prepare test input
	schema := &ast.Schema{
		Package: "without_enums",
		Objects: testutils.ObjectsMap(
			ast.NewObject("without_enums", "AString", ast.String()),
			ast.NewObject("without_enums", "AStruct", ast.NewStruct(
				ast.NewStructField("AString", ast.String()),
			)),
		),
	}

	// Run the compiler pass
	runPassOnSchema(t, &AnonymousEnumToExplicitType{}, schema, schema)
}

func TestAnonymousEnumToExplicitType_withAnonymousEnumInStruct(t *testing.T) {
	// Prepare test input
	schema := &ast.Schema{
		Package: "with_enums",
		Objects: testutils.ObjectsMap(
			ast.NewObject("with_enums", "Panel", ast.NewStruct(
				ast.NewStructField("title", ast.String()),
				ast.NewStructField("type", ast.NewEnum([]ast.EnumValue{
					{Name: "Foo", Value: "foo", Type: ast.String()},
					{Name: "Bar", Value: "bar", Type: ast.String()},
				}, ast.Nullable())),
			)),
			ast.NewObject("with_enums", "Mode", ast.NewEnum([]ast.EnumValue{
				{Name: "Auto", Value: "auto", Type: ast.String()},
				{Name: "Manual", Value: "manual", Type: ast.String()},
			})),
		),
	}

	// Prepare expected output
	expected := &ast.Schema{
		Package: "with_enums",
		Objects: testutils.ObjectsMap(
			ast.NewObject("with_enums", "Panel", ast.NewStruct(
				ast.NewStructField("title", ast.String()),
				ast.NewStructField("type", ast.NewRef("with_enums", "PanelType", ast.Nullable(), ast.Trail("AnonymousEnumToExplicitType"))),
			)),

			// this object is unchanged
			schema.Objects.Get("Mode"),

			// the anonymous enum, turned into an object
			ast.NewObject("with_enums", "PanelType", ast.NewEnum([]ast.EnumValue{
				{Name: "Foo", Value: "foo", Type: ast.String()},
				{Name: "Bar", Value: "bar", Type: ast.String()},
			}), "AnonymousEnumToExplicitType"),
		),
	}

	// Run the compiler pass
	runPassOnSchema(t, &AnonymousEnumToExplicitType{}, schema, expected)
}

func TestAnonymousEnumToExplicitType_withAnonymousEnumInArray(t *testing.T) {
	// Prepare test input
	schema := &ast.Schema{
		Package: "in_array",
		Objects: testutils.ObjectsMap(
			ast.NewObject("in_array", "TypesList", ast.NewArray(
				ast.NewEnum([]ast.EnumValue{
					{Name: "Foo", Value: "foo", Type: ast.String()},
					{Name: "Bar", Value: "bar", Type: ast.String()},
				}),
			)),
		),
	}

	// Prepare expected output
	expected := &ast.Schema{
		Package: "in_array",
		Objects: testutils.ObjectsMap(
			ast.NewObject("in_array", "TypesList", ast.NewArray(
				ast.NewRef("in_array", "TypesListEnum", ast.Trail("AnonymousEnumToExplicitType"))),
			),

			// the anonymous enum, turned into an object
			ast.NewObject("in_array", "TypesListEnum", ast.NewEnum([]ast.EnumValue{
				{Name: "Foo", Value: "foo", Type: ast.String()},
				{Name: "Bar", Value: "bar", Type: ast.String()},
			}), "AnonymousEnumToExplicitType"),
		),
	}

	// Run the compiler pass
	runPassOnSchema(t, &AnonymousEnumToExplicitType{}, schema, expected)
}

func TestAnonymousEnumToExplicitType_withAnonymousEnumInMap(t *testing.T) {
	// Prepare test input
	schema := &ast.Schema{
		Package: "in_map",
		Objects: testutils.ObjectsMap(
			ast.NewObject("in_map", "MapOfThings", ast.NewMap(
				ast.String(),
				ast.NewEnum([]ast.EnumValue{
					{Name: "Foo", Value: "foo", Type: ast.String()},
					{Name: "Bar", Value: "bar", Type: ast.String()},
				}),
			)),
		),
	}

	// Prepare expected output
	expected := &ast.Schema{
		Package: "in_map",
		Objects: testutils.ObjectsMap(
			ast.NewObject("in_map", "MapOfThings", ast.NewMap(
				ast.String(),
				ast.NewRef("in_map", "MapOfThingsEnum", ast.Trail("AnonymousEnumToExplicitType"))),
			),

			// the anonymous enum, turned into an object
			ast.NewObject("in_map", "MapOfThingsEnum", ast.NewEnum([]ast.EnumValue{
				{Name: "Foo", Value: "foo", Type: ast.String()},
				{Name: "Bar", Value: "bar", Type: ast.String()},
			}), "AnonymousEnumToExplicitType"),
		),
	}

	// Run the compiler pass
	runPassOnSchema(t, &AnonymousEnumToExplicitType{}, schema, expected)
}

func TestAnonymousEnumToExplicitType_withAnonymousEnumInDisjunction(t *testing.T) {
	// Prepare test input
	schema := &ast.Schema{
		Package: "in_disjunction",
		Objects: testutils.ObjectsMap(
			ast.NewObject("in_disjunction", "DisjunctionOfThings", ast.NewDisjunction([]ast.Type{
				ast.String(),
				ast.NewEnum([]ast.EnumValue{
					{Name: "Foo", Value: "foo", Type: ast.String()},
					{Name: "Bar", Value: "bar", Type: ast.String()},
				}),
			})),
		),
	}

	// Prepare expected output
	expected := &ast.Schema{
		Package: "in_disjunction",
		Objects: testutils.ObjectsMap(
			ast.NewObject("in_disjunction", "DisjunctionOfThings", ast.NewDisjunction([]ast.Type{
				ast.String(),
				ast.NewRef("in_disjunction", "DisjunctionOfThingsEnum", ast.Trail("AnonymousEnumToExplicitType")),
			})),

			// the anonymous enum, turned into an object
			ast.NewObject("in_disjunction", "DisjunctionOfThingsEnum", ast.NewEnum([]ast.EnumValue{
				{Name: "Foo", Value: "foo", Type: ast.String()},
				{Name: "Bar", Value: "bar", Type: ast.String()},
			}), "AnonymousEnumToExplicitType"),
		),
	}

	// Run the compiler pass
	runPassOnSchema(t, &AnonymousEnumToExplicitType{}, schema, expected)
}

func TestAnonymousEnumToExplicitType_withAnonymousEnumInIntersection(t *testing.T) {
	// Prepare test input
	schema := &ast.Schema{
		Package: "in_intersection",
		Objects: testutils.ObjectsMap(
			ast.NewObject("in_intersection", "IntersectionOfThings", ast.NewIntersection([]ast.Type{
				ast.String(),
				ast.NewEnum([]ast.EnumValue{
					{Name: "Foo", Value: "foo", Type: ast.String()},
					{Name: "Bar", Value: "bar", Type: ast.String()},
				}),
			})),
		),
	}

	// Prepare expected output
	expected := &ast.Schema{
		Package: "in_intersection",
		Objects: testutils.ObjectsMap(
			ast.NewObject("in_intersection", "IntersectionOfThings", ast.NewIntersection([]ast.Type{
				ast.String(),
				ast.NewRef("in_intersection", "IntersectionOfThingsEnum", ast.Trail("AnonymousEnumToExplicitType")),
			})),

			// the anonymous enum, turned into an object
			ast.NewObject("in_intersection", "IntersectionOfThingsEnum", ast.NewEnum([]ast.EnumValue{
				{Name: "Foo", Value: "foo", Type: ast.String()},
				{Name: "Bar", Value: "bar", Type: ast.String()},
			}), "AnonymousEnumToExplicitType"),
		),
	}

	// Run the compiler pass
	runPassOnSchema(t, &AnonymousEnumToExplicitType{}, schema, expected)
}

func TestAnonymousEnumToExplicitType_withFieldWithDefaultValue(t *testing.T) {
	// Prepare expected input
	schema := &ast.Schema{
		Package: "with_default",
		Objects: testutils.ObjectsMap(
			ast.NewObject("with_default", "Panel", ast.NewStruct(
				ast.NewStructField("title", ast.String()),
				ast.NewStructField("type", ast.NewEnum([]ast.EnumValue{
					{Name: "Foo", Value: "foo", Type: ast.String()},
					{Name: "Bar", Value: "bar", Type: ast.String()},
				}, ast.Nullable(), ast.Default("foo"))),
			)),
		),
	}

	// Prepare expected output
	expected := &ast.Schema{
		Package: "with_default",
		Objects: testutils.ObjectsMap(
			ast.NewObject("with_default", "Panel", ast.NewStruct(
				ast.NewStructField("title", ast.String()),
				ast.NewStructField("type", ast.NewRef("with_default", "PanelType", ast.Nullable(), ast.Default("foo"), ast.Trail("AnonymousEnumToExplicitType"))),
			)),

			// the anonymous enum, turned into an object
			ast.NewObject("with_default", "PanelType", ast.NewEnum([]ast.EnumValue{
				{Name: "Foo", Value: "foo", Type: ast.String()},
				{Name: "Bar", Value: "bar", Type: ast.String()},
			}), "AnonymousEnumToExplicitType"),
		),
	}

	// Run the compiler pass
	runPassOnSchema(t, &AnonymousEnumToExplicitType{}, schema, expected)
}
