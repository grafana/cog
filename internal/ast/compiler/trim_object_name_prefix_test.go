package compiler

import (
	"testing"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/testutils"
)

func TestTrimObjectNamePrefix(t *testing.T) {
	// Prepare test input
	schema := &ast.Schema{
		Package: "prefix_names",
		Objects: testutils.ObjectsMap(
			ast.NewObject("prefix_names", "MyPrefixSomeObject", ast.NewStruct(
				ast.NewStructField("foo", ast.String(ast.Nullable())),
				ast.NewStructField("ref_to_nice_object", ast.NewRef("prefix_names", "MyPrefixNotANiceName")),
			)),
			ast.NewObject("prefix_names", "MyPrefixNotANiceName", ast.NewStruct(
				ast.NewStructField("AString", ast.String(ast.Nullable())),
			)),
			ast.NewObject("prefix_names", "VariableRefresh", ast.NewEnum([]ast.EnumValue{
				{Name: "MyPrefixNever", Value: "never", Type: ast.String()},
				{Name: "MyPrefixAlways", Value: "always", Type: ast.String()},
			})),
		),
	}
	expected := &ast.Schema{
		Package: "prefix_names",
		Objects: testutils.ObjectsMap(
			ast.NewObject("prefix_names", "SomeObject", ast.NewStruct(
				ast.NewStructField("foo", ast.String(ast.Nullable())),
				ast.NewStructField("ref_to_nice_object", ast.NewRef("prefix_names", "NotANiceName", ast.Trail("TrimObjectNamePrefix[MyPrefixNotANiceName → NotANiceName]"))),
			), "TrimObjectNamePrefix[MyPrefixSomeObject → SomeObject]"),
			ast.NewObject("prefix_names", "NotANiceName", ast.NewStruct(
				ast.NewStructField("AString", ast.String(ast.Nullable())),
			), "TrimObjectNamePrefix[MyPrefixNotANiceName → NotANiceName]"),
			ast.NewObject("prefix_names", "VariableRefresh", ast.NewEnum([]ast.EnumValue{
				{Name: "Never", Value: "never", Type: ast.String()},
				{Name: "Always", Value: "always", Type: ast.String()},
			})),
		),
	}

	pass := &TrimObjectNamePrefix{
		Prefix: "MyPrefix",
	}

	// Run the compiler pass
	runPassOnSchema(t, pass, schema, expected)
}
