package compiler

import (
	"testing"

	"github.com/grafana/cog/internal/ast"
)

func TestDisjunctionOfConstantsToEnum_withInvalidTypes(t *testing.T) {
	// Prepare test input
	objects := []ast.Object{
		ast.NewObject(testPkgName, "AString", ast.String()),
		ast.NewObject(testPkgName, "AStruct", ast.NewStruct(
			ast.NewStructField("AString", ast.String()),
		)),
		ast.NewObject(testPkgName, "ADisjunction", ast.NewDisjunction(ast.Types{
			ast.NewRef(testPkgName, "AString"),
			ast.NewRef(testPkgName, "AStruct"),
		})),
		ast.NewObject(testPkgName, "ADisjunctionWithConstAndNonConst", ast.NewDisjunction(ast.Types{
			ast.String(),
			ast.String(ast.Value("foo")),
		})),
		ast.NewObject(testPkgName, "ADisjunctionWithMixedTypes", ast.NewDisjunction(ast.Types{
			ast.String(ast.Value("foo")),
			ast.NewScalar(ast.KindFloat32, ast.Value(float32(42))),
		})),
	}

	// Run the compiler pass
	runPassOnObjects(t, &DisjunctionOfConstantsToEnum{}, objects, objects)
}

func TestDisjunctionOfConstantsToEnum(t *testing.T) {
	// Prepare test input
	objects := []ast.Object{
		ast.NewObject(testPkgName, "FirstDisjunction", ast.NewDisjunction(ast.Types{
			ast.String(ast.Value("first")),
			ast.String(ast.Value("second")),
		})),
		ast.NewObject(testPkgName, "SomeEnum", ast.NewEnum([]ast.EnumValue{
			{Type: ast.String(), Name: "foo", Value: "foo"},
			{Type: ast.String(), Name: "bar", Value: "bar"},
		})),
		ast.NewObject(testPkgName, "SecondDisjunction", ast.NewDisjunction(ast.Types{
			ast.String(ast.Value("third")),
			ast.NewRef(testPkgName, "FirstDisjunction"),
			ast.NewRef(testPkgName, "SomeEnum"),
		})),
	}

	expectedObjects := []ast.Object{
		ast.NewObject(testPkgName, "FirstDisjunction", ast.NewEnum([]ast.EnumValue{
			{Type: ast.String(), Name: "first", Value: "first"},
			{Type: ast.String(), Name: "second", Value: "second"},
		}, ast.Trail("DisjunctionOfConstantsToEnum"))),
		ast.NewObject(testPkgName, "SomeEnum", ast.NewEnum([]ast.EnumValue{
			{Type: ast.String(), Name: "foo", Value: "foo"},
			{Type: ast.String(), Name: "bar", Value: "bar"},
		})),
		ast.NewObject(testPkgName, "SecondDisjunction", ast.NewEnum([]ast.EnumValue{
			{Type: ast.String(), Name: "third", Value: "third"},
			{Type: ast.String(), Name: "first", Value: "first"},
			{Type: ast.String(), Name: "second", Value: "second"},
			{Type: ast.String(), Name: "foo", Value: "foo"},
			{Type: ast.String(), Name: "bar", Value: "bar"},
		}, ast.Trail("DisjunctionOfConstantsToEnum"))),
	}

	// Run the compiler pass
	runPassOnObjects(t, &DisjunctionOfConstantsToEnum{}, objects, expectedObjects)
}
