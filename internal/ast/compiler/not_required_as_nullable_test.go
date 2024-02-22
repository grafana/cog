package compiler

import (
	"testing"

	"github.com/grafana/cog/internal/ast"
)

func TestNotRequiredFieldAsNullableType(t *testing.T) {
	// Prepare test input
	objects := []ast.Object{
		ast.NewObject("pkg", "NotAStruct", ast.String()),

		ast.NewObject("pkg", "AStruct", ast.NewStruct(
			ast.NewStructField("RequiredString", ast.String(), ast.Required()),
			ast.NewStructField("RequiredNullableString", ast.String(ast.Nullable()), ast.Required()),
			ast.NewStructField("NotRequiredString", ast.String()),

			ast.NewStructField("RequiredRef", ast.NewRef("test", "SomeStruct"), ast.Required()),
			ast.NewStructField("RequiredNullableRef", ast.NewRef("test", "SomeStruct", ast.Nullable()), ast.Required()),
			ast.NewStructField("NotRequiredRef", ast.NewRef("test", "SomeStruct")),

			ast.NewStructField("NotRequiredArray", ast.NewArray(ast.String())),
			ast.NewStructField("RequiredArray", ast.NewArray(ast.String()), ast.Required()),

			ast.NewStructField("NotRequiredMap", ast.NewMap(
				ast.String(),
				ast.Bool(),
			)),
			ast.NewStructField("RequiredMap", ast.NewMap(
				ast.String(),
				ast.Bool(),
			), ast.Required()),
		)),
	}

	// Prepare expected output
	expected := []ast.Object{
		ast.NewObject("pkg", "NotAStruct", ast.String()),

		ast.NewObject("pkg", "AStruct", ast.NewStruct(
			ast.NewStructField("RequiredString", ast.String(), ast.Required()),
			ast.NewStructField("RequiredNullableString", ast.String(ast.Nullable()), ast.Required()),
			ast.NewStructField("NotRequiredString", ast.String(ast.Nullable()), ast.PassesTrail("NotRequiredFieldAsNullableType[nullable=true]")), // should become nullable

			ast.NewStructField("RequiredRef", ast.NewRef("test", "SomeStruct"), ast.Required()),
			ast.NewStructField("RequiredNullableRef", ast.NewRef("test", "SomeStruct", ast.Nullable()), ast.Required()),
			ast.NewStructField("NotRequiredRef", ast.NewRef("test", "SomeStruct", ast.Nullable()), ast.PassesTrail("NotRequiredFieldAsNullableType[nullable=true]")), // should become nullable

			ast.NewStructField("NotRequiredArray", ast.NewArray(ast.String(), ast.Nullable()), ast.PassesTrail("NotRequiredFieldAsNullableType[nullable=true]")), // should become nullable
			ast.NewStructField("RequiredArray", ast.NewArray(ast.String()), ast.Required()),

			ast.NewStructField("NotRequiredMap", ast.NewMap( // should become nullable
				ast.String(),
				ast.Bool(),
				ast.Nullable(),
			), ast.PassesTrail("NotRequiredFieldAsNullableType[nullable=true]")),
			ast.NewStructField("RequiredMap", ast.NewMap(
				ast.String(),
				ast.Bool(),
			), ast.Required()),
		)),
	}

	// Run the compiler pass
	runPassOnObjects(t, &NotRequiredFieldAsNullableType{}, objects, expected)
}
