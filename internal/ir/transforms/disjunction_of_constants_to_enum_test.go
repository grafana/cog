package transforms

import (
	"testing"

	"github.com/grafana/cog/internal/ir"
)

func TestDisjunctionOfConstantsToEnum_withInvalidTypes(t *testing.T) {
	// Prepare test input
	objects := []ir.Object{
		ir.NewObject(testPkgName, "AString", ir.String()),
		ir.NewObject(testPkgName, "AStruct", ir.NewStruct(
			ir.NewStructField("AString", ir.String()),
		)),
		ir.NewObject(testPkgName, "ADisjunction", ir.NewDisjunction(ir.Types{
			ir.NewRef(testPkgName, "AString"),
			ir.NewRef(testPkgName, "AStruct"),
		})),
		ir.NewObject(testPkgName, "ADisjunctionWithConstAndNonConst", ir.NewDisjunction(ir.Types{
			ir.String(),
			ir.String(ir.Value("foo")),
		})),
		ir.NewObject(testPkgName, "ADisjunctionWithMixedTypes", ir.NewDisjunction(ir.Types{
			ir.String(ir.Value("foo")),
			ir.NewScalar(ir.KindFloat32, ir.Value(float32(42))),
		})),
	}

	// Run the compiler pass
	runPassOnObjects(t, &DisjunctionOfConstantsToEnum{}, objects, objects)
}

func TestDisjunctionOfConstantsToEnum(t *testing.T) {
	// Prepare test input
	objects := []ir.Object{
		ir.NewObject(testPkgName, "FirstDisjunction", ir.NewDisjunction(ir.Types{
			ir.String(ir.Value("first")),
			ir.String(ir.Value("second")),
		})),
		ir.NewObject(testPkgName, "SomeEnum", ir.NewEnum([]ir.EnumValue{
			{Type: ir.String(), Name: "foo", Value: "foo"},
			{Type: ir.String(), Name: "bar", Value: "bar"},
		})),
		ir.NewObject(testPkgName, "SecondDisjunction", ir.NewDisjunction(ir.Types{
			ir.String(ir.Value("third")),
			ir.NewRef(testPkgName, "FirstDisjunction"),
			ir.NewRef(testPkgName, "SomeEnum"),
		})),
	}

	expectedObjects := []ir.Object{
		ir.NewObject(testPkgName, "FirstDisjunction", ir.NewEnum([]ir.EnumValue{
			{Type: ir.String(), Name: "first", Value: "first"},
			{Type: ir.String(), Name: "second", Value: "second"},
		}, ir.Trail("DisjunctionOfConstantsToEnum"))),
		ir.NewObject(testPkgName, "SomeEnum", ir.NewEnum([]ir.EnumValue{
			{Type: ir.String(), Name: "foo", Value: "foo"},
			{Type: ir.String(), Name: "bar", Value: "bar"},
		})),
		ir.NewObject(testPkgName, "SecondDisjunction", ir.NewEnum([]ir.EnumValue{
			{Type: ir.String(), Name: "third", Value: "third"},
			{Type: ir.String(), Name: "first", Value: "first"},
			{Type: ir.String(), Name: "second", Value: "second"},
			{Type: ir.String(), Name: "foo", Value: "foo"},
			{Type: ir.String(), Name: "bar", Value: "bar"},
		}, ir.Trail("DisjunctionOfConstantsToEnum"))),
	}

	// Run the compiler pass
	runPassOnObjects(t, &DisjunctionOfConstantsToEnum{}, objects, expectedObjects)
}
