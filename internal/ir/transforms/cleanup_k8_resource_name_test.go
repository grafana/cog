package transforms

import (
	"testing"

	"github.com/grafana/cog/internal/ir"
)

func TestCleanupK8ResourceNames(t *testing.T) {
	objects := []ir.Object{
		ir.NewObject("test", "my.super.large.name.Resource",
			ir.NewStruct(
				ir.NewStructField("aRef", ir.NewRef("test", "my.super.large.name.ARef")),
				ir.NewStructField("aConstantRef", ir.NewConstantReferenceType("test", "other.name.with.ugly.name.AConstantRef", "a")),
				ir.NewStructField("aDisjunction", ir.NewDisjunction([]ir.Type{
					ir.NewRef("test", "my.super.large.name.ARef"),
					ir.NewRef("test", "other.name.with.ugly.name.AConstantRef"),
				})),
			),
		),
		ir.NewObject("test", "my.super.large.name.ARef", ir.NewStruct(
			ir.NewStructField("aString", ir.String()),
		)),
		ir.NewObject("test", "other.name.with.ugly.name.AConstantRef", ir.NewEnum([]ir.EnumValue{
			{Type: ir.String(), Name: "A", Value: "a"},
			{Type: ir.String(), Name: "B", Value: "b"},
		})),
	}

	expected := []ir.Object{
		ir.NewObject("test", "Resource",
			ir.NewStruct(
				ir.NewStructField("aRef", ir.NewRef("test", "ARef")),
				ir.NewStructField("aConstantRef", ir.NewConstantReferenceType("test", "AConstantRef", "a")),
				ir.NewStructField("aDisjunction", ir.NewDisjunction([]ir.Type{
					ir.NewRef("test", "ARef"),
					ir.NewRef("test", "AConstantRef"),
				})),
			),
		),
		ir.NewObject("test", "ARef", ir.NewStruct(
			ir.NewStructField("aString", ir.String()),
		)),
		ir.NewObject("test", "AConstantRef", ir.NewEnum([]ir.EnumValue{
			{Type: ir.String(), Name: "A", Value: "a"},
			{Type: ir.String(), Name: "B", Value: "b"},
		})),
	}

	runPassOnObjects(t, &CleanupK8ResourceNames{}, objects, expected)
}

func TestCleanupK8ResourceNamesPrefix(t *testing.T) {
	objects := []ir.Object{
		ir.NewObject("test", "my.super.large.name.Resource",
			ir.NewStruct(
				ir.NewStructField("aRef", ir.NewRef("test", "my.super.large.name.HelloARef")),
				ir.NewStructField("aConstantRef", ir.NewConstantReferenceType("test", "other.name.with.ugly.name.HelloAConstantRef", "a")),
				ir.NewStructField("aDisjunction", ir.NewDisjunction([]ir.Type{
					ir.NewRef("test", "my.super.large.name.HelloARef"),
					ir.NewRef("test", "other.name.with.ugly.name.HelloAConstantRef"),
				})),
			),
		),
		ir.NewObject("test", "my.super.large.name.HelloARef", ir.NewStruct(
			ir.NewStructField("aString", ir.String()),
		)),
		ir.NewObject("test", "other.name.with.ugly.name.HelloAConstantRef", ir.NewEnum([]ir.EnumValue{
			{Type: ir.String(), Name: "A", Value: "a"},
			{Type: ir.String(), Name: "B", Value: "b"},
		})),
	}

	expected := []ir.Object{
		ir.NewObject("test", "Resource",
			ir.NewStruct(
				ir.NewStructField("aRef", ir.NewRef("test", "ARef")),
				ir.NewStructField("aConstantRef", ir.NewConstantReferenceType("test", "AConstantRef", "a")),
				ir.NewStructField("aDisjunction", ir.NewDisjunction([]ir.Type{
					ir.NewRef("test", "ARef"),
					ir.NewRef("test", "AConstantRef"),
				})),
			),
		),
		ir.NewObject("test", "ARef", ir.NewStruct(
			ir.NewStructField("aString", ir.String()),
		)),
		ir.NewObject("test", "AConstantRef", ir.NewEnum([]ir.EnumValue{
			{Type: ir.String(), Name: "A", Value: "a"},
			{Type: ir.String(), Name: "B", Value: "b"},
		})),
	}

	runPassOnObjects(t, &CleanupK8ResourceNames{PrefixToRemove: "Hello"}, objects, expected)
}
