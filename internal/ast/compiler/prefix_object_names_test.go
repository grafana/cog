package compiler

import (
	"testing"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/testutils"
)

func TestPrefixObjectNames(t *testing.T) {
	// Prepare test input
	schema := &ast.Schema{
		Package: "prefix_names",
		Objects: testutils.ObjectsMap(
			ast.NewObject("prefix_names", "SomeObject", ast.NewStruct(
				ast.NewStructField("foo", ast.String(ast.Nullable())),
				ast.NewStructField("ref_to_nice_object", ast.NewRef("prefix_names", "NotANiceName")),
			)),
			ast.NewObject("prefix_names", "NotANiceName", ast.NewStruct(
				ast.NewStructField("AString", ast.String(ast.Nullable())),
			)),
			ast.NewObject("prefix_names", "VariableRefresh", ast.NewEnum([]ast.EnumValue{
				{Name: "Never", Value: "never", Type: ast.String()},
				{Name: "Always", Value: "always", Type: ast.String()},
			})),
		),
	}
	expected := &ast.Schema{
		Package: "prefix_names",
		Objects: testutils.ObjectsMap(
			ast.NewObject("prefix_names", "PreSomeObject", ast.NewStruct(
				ast.NewStructField("foo", ast.String(ast.Nullable())),
				ast.NewStructField("ref_to_nice_object", ast.NewRef("prefix_names", "PreNotANiceName", ast.Trail("PrefixObjectNames[NotANiceName → PreNotANiceName]"))),
			), "PrefixObjectNames[SomeObject → PreSomeObject]"),
			ast.NewObject("prefix_names", "PreNotANiceName", ast.NewStruct(
				ast.NewStructField("AString", ast.String(ast.Nullable())),
			), "PrefixObjectNames[NotANiceName → PreNotANiceName]"),
			ast.NewObject("prefix_names", "PreVariableRefresh", ast.NewEnum([]ast.EnumValue{
				{Name: "PreNever", Value: "never", Type: ast.String()},
				{Name: "PreAlways", Value: "always", Type: ast.String()},
			}), "PrefixObjectNames[VariableRefresh → PreVariableRefresh]"),
		),
	}

	pass := &PrefixObjectNames{
		Prefix: "Pre",
	}

	// Run the compiler pass
	runPassOnSchema(t, pass, schema, expected)
}
