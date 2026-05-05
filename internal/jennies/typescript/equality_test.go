package typescript

import (
	"strings"
	"testing"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/languages"
	"github.com/grafana/cog/internal/testutils"
	"github.com/stretchr/testify/require"
)

func equalitySchema() *ast.Schema {
	return &ast.Schema{
		Package: "equality",
		Objects: testutils.ObjectsMap(
			ast.NewObject("equality", "Variable", ast.NewStruct(
				ast.NewStructField("name", ast.NewScalar(ast.KindString), ast.Required()),
			)),
			ast.NewObject("equality", "Container", ast.NewStruct(
				ast.NewStructField("stringField", ast.NewScalar(ast.KindString), ast.Required()),
				ast.NewStructField("intField", ast.NewScalar(ast.KindInt64), ast.Required()),
				ast.NewStructField("refField", ast.NewRef("equality", "Variable"), ast.Required()),
			)),
			ast.NewObject("equality", "Optionals", ast.NewStruct(
				ast.NewStructField("stringField", ast.NewScalar(ast.KindString)),
				ast.NewStructField("refField", ast.NewRef("equality", "Variable")),
			)),
			ast.NewObject("equality", "Arrays", ast.NewStruct(
				ast.NewStructField("ints", ast.NewArray(ast.NewScalar(ast.KindInt64)), ast.Required()),
				ast.NewStructField("refs", ast.NewArray(ast.NewRef("equality", "Variable")), ast.Required()),
			)),
			ast.NewObject("equality", "Maps", ast.NewStruct(
				ast.NewStructField("ints", ast.NewMap(ast.String(), ast.NewScalar(ast.KindInt64)), ast.Required()),
				ast.NewStructField("refs", ast.NewMap(ast.String(), ast.NewRef("equality", "Variable")), ast.Required()),
			)),
		),
	}
}

func TestEquality_TypeScript_GeneratesEqualsFunctions(t *testing.T) {
	req := require.New(t)

	config := Config{GenerateEqual: true}
	config.applyDefaults()

	jenny := RawTypes{
		config: config,
		tmpl:   initTemplates(config, common.NewAPIReferenceCollector()),
	}

	schema := equalitySchema()
	context := languages.Context{Schemas: ast.Schemas{schema}}
	jenny.schemas = context.Schemas

	files, err := jenny.Generate(context)
	req.NoError(err)
	req.Len(files, 1)

	output := string(files[0].Data)

	// Struct with only a name field → equalsVariable
	req.Contains(output, "export const equalsVariable = (a: Variable, b: Variable): boolean => {")
	req.Contains(output, "if (a.name !== b.name) return false;")

	// Container with scalar and ref fields
	req.Contains(output, "export const equalsContainer = (a: Container, b: Container): boolean => {")
	req.Contains(output, "if (a.stringField !== b.stringField) return false;")
	req.Contains(output, "if (a.intField !== b.intField) return false;")
	req.Contains(output, "if (!equalsVariable(a.refField, b.refField)) return false;")

	// Optional fields get undefined checks
	req.Contains(output, "export const equalsOptionals = (a: Optionals, b: Optionals): boolean => {")
	req.Contains(output, "(a.stringField === undefined) !== (b.stringField === undefined)")
	req.Contains(output, "(a.refField === undefined) !== (b.refField === undefined)")

	// Arrays get length check + iteration
	req.Contains(output, "export const equalsArrays = (a: Arrays, b: Arrays): boolean => {")
	req.Contains(output, "a.ints.length !== b.ints.length")
	req.Contains(output, "for (let i")

	// Maps get Object.keys check + iteration
	req.Contains(output, "export const equalsMaps = (a: Maps, b: Maps): boolean => {")
	req.Contains(output, "Object.keys(a.ints).length !== Object.keys(b.ints).length")
	req.Contains(output, "for (const key")
}

func TestEquality_TypeScript_SkipsNonStructTypes(t *testing.T) {
	req := require.New(t)

	config := Config{GenerateEqual: true}
	config.applyDefaults()

	jenny := RawTypes{
		config: config,
		tmpl:   initTemplates(config, common.NewAPIReferenceCollector()),
	}

	// A scalar constant (non-struct type) should not generate an equality function
	schema := &ast.Schema{
		Package: "test",
		Objects: testutils.ObjectsMap(
			ast.NewObject("test", "MyConstant", ast.NewScalar(ast.KindString, ast.Value("hello"))),
		),
	}

	context := languages.Context{Schemas: ast.Schemas{schema}}
	jenny.schemas = context.Schemas

	files, err := jenny.Generate(context)
	req.NoError(err)
	req.Len(files, 1)

	output := string(files[0].Data)
	req.False(strings.Contains(output, "equalsMyConstant"), "scalar types should not generate equality functions")
}

func TestEquality_TypeScript_DisabledByDefault(t *testing.T) {
	req := require.New(t)

	config := Config{} // GenerateEqual defaults to false
	config.applyDefaults()

	jenny := RawTypes{
		config: config,
		tmpl:   initTemplates(config, common.NewAPIReferenceCollector()),
	}

	schema := equalitySchema()
	context := languages.Context{Schemas: ast.Schemas{schema}}
	jenny.schemas = context.Schemas

	files, err := jenny.Generate(context)
	req.NoError(err)
	req.Len(files, 1)

	output := string(files[0].Data)
	req.False(strings.Contains(output, "equalsContainer"), "equality functions should not be generated when GenerateEqual is false")
}
