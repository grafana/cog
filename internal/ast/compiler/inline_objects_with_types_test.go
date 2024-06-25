package compiler

import (
	"testing"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/testutils"
)

func TestInlineScalarAliases(t *testing.T) {
	// Prepare test input
	schema := &ast.Schema{
		Package: "inline_scalar_aliases",
		Objects: testutils.ObjectsMap(
			ast.NewObject("inline_scalar_aliases", "AliasToString", ast.String()),
			ast.NewObject("inline_scalar_aliases", "AliasToMap", ast.NewMap(ast.String(), ast.Any())),
			ast.NewObject("inline_scalar_aliases", "AliasToArray", ast.NewArray(ast.String())),
			ast.NewObject("inline_scalar_aliases", "Constant", ast.String(ast.Value("foo"))),
			ast.NewObject("inline_scalar_aliases", "SomeObject", ast.NewStruct(
				ast.NewStructField("aliasToString", ast.NewRef("inline_scalar_aliases", "AliasToString")),
				ast.NewStructField("aliasToMap", ast.NewRef("inline_scalar_aliases", "AliasToMap")),
				ast.NewStructField("aliasToArray", ast.NewRef("inline_scalar_aliases", "AliasToArray")),
			)),
		),
	}
	expected := &ast.Schema{
		Package: "inline_scalar_aliases",
		Objects: testutils.ObjectsMap(
			ast.NewObject("inline_scalar_aliases", "Constant", ast.String(ast.Value("foo"))),
			ast.NewObject("inline_scalar_aliases", "SomeObject", ast.NewStruct(
				ast.NewStructField("aliasToString", ast.String(ast.Trail("InlineObjectsWithTypes[original=inline_scalar_aliases.AliasToString]"))),
				ast.NewStructField("aliasToMap", ast.NewMap(ast.String(), ast.Any(), ast.Trail("InlineObjectsWithTypes[original=inline_scalar_aliases.AliasToMap]"))),
				ast.NewStructField("aliasToArray", ast.NewArray(ast.String(), ast.Trail("InlineObjectsWithTypes[original=inline_scalar_aliases.AliasToArray]"))),
			)),
		),
	}

	pass := &InlineObjectsWithTypes{
		InlineTypes: []ast.Kind{ast.KindScalar, ast.KindArray, ast.KindMap},
	}

	// Run the compiler pass
	runPassOnSchema(t, pass, schema, expected)
}
