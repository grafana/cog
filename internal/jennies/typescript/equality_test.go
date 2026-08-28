package typescript

import (
	"strings"
	"testing"

	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/testutils"
	"github.com/grafana/cog/pkg/ir"
	"github.com/grafana/cog/pkg/languages"
	"github.com/stretchr/testify/require"
)

func equalitySchema() *ir.Schema {
	return &ir.Schema{
		Package: "equality",
		Objects: testutils.ObjectsMap(
			ir.NewObject("equality", "Variable", ir.NewStruct(
				ir.NewStructField("name", ir.NewScalar(ir.KindString), ir.Required()),
			)),
			ir.NewObject("equality", "Container", ir.NewStruct(
				ir.NewStructField("stringField", ir.NewScalar(ir.KindString), ir.Required()),
				ir.NewStructField("intField", ir.NewScalar(ir.KindInt64), ir.Required()),
				ir.NewStructField("refField", ir.NewRef("equality", "Variable"), ir.Required()),
			)),
			ir.NewObject("equality", "Optionals", ir.NewStruct(
				ir.NewStructField("stringField", ir.NewScalar(ir.KindString)),
				ir.NewStructField("refField", ir.NewRef("equality", "Variable")),
			)),
			ir.NewObject("equality", "Arrays", ir.NewStruct(
				ir.NewStructField("ints", ir.NewArray(ir.NewScalar(ir.KindInt64)), ir.Required()),
				ir.NewStructField("refs", ir.NewArray(ir.NewRef("equality", "Variable")), ir.Required()),
			)),
			ir.NewObject("equality", "Maps", ir.NewStruct(
				ir.NewStructField("ints", ir.NewMap(ir.String(), ir.NewScalar(ir.KindInt64)), ir.Required()),
				ir.NewStructField("refs", ir.NewMap(ir.String(), ir.NewRef("equality", "Variable")), ir.Required()),
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
	context := languages.Context{Schemas: ir.Schemas{schema}}
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
	schema := &ir.Schema{
		Package: "test",
		Objects: testutils.ObjectsMap(
			ir.NewObject("test", "MyConstant", ir.NewScalar(ir.KindString, ir.Value("hello"))),
		),
	}

	context := languages.Context{Schemas: ir.Schemas{schema}}
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
	context := languages.Context{Schemas: ir.Schemas{schema}}
	jenny.schemas = context.Schemas

	files, err := jenny.Generate(context)
	req.NoError(err)
	req.Len(files, 1)

	output := string(files[0].Data)
	req.False(strings.Contains(output, "equalsContainer"), "equality functions should not be generated when GenerateEqual is false")
}
