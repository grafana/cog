package python

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
		),
	}
}

func TestEquality_Python_GeneratesEqMethods(t *testing.T) {
	req := require.New(t)

	config := Config{GenerateEqual: true}
	jenny := RawTypes{
		config:          config,
		tmpl:            initTemplates(config, common.NewAPIReferenceCollector()),
		apiRefCollector: common.NewAPIReferenceCollector(),
	}

	schema := equalitySchema()

	// Run Python compiler passes so nullable types are handled correctly
	processedSchemas, err := New(config).CompilerPasses().Process(ast.Schemas{schema})
	req.NoError(err)

	context := languages.Context{Schemas: processedSchemas}
	files, err := jenny.Generate(context)
	req.NoError(err)
	req.Len(files, 1)

	output := string(files[0].Data)

	// Variable struct
	req.Contains(output, "def __eq__(self, other: object) -> bool:")
	req.Contains(output, "if not isinstance(other, Variable):")
	req.Contains(output, "if self.name != other.name:")

	// Container struct
	req.Contains(output, "if not isinstance(other, Container):")
	req.Contains(output, "if self.string_field != other.string_field:")
	req.Contains(output, "if self.int_field != other.int_field:")
	req.Contains(output, "if self.ref_field != other.ref_field:")

	// Optionals struct
	req.Contains(output, "if not isinstance(other, Optionals):")
	req.Contains(output, "if self.string_field != other.string_field:")
	req.Contains(output, "if self.ref_field != other.ref_field:")
}

func TestEquality_Python_SkipsNonStructTypes(t *testing.T) {
	req := require.New(t)

	config := Config{GenerateEqual: true}
	jenny := RawTypes{
		config:          config,
		tmpl:            initTemplates(config, common.NewAPIReferenceCollector()),
		apiRefCollector: common.NewAPIReferenceCollector(),
	}

	// A scalar constant (non-struct type) should not generate __eq__
	schema := &ast.Schema{
		Package: "test",
		Objects: testutils.ObjectsMap(
			ast.NewObject("test", "MyConstant", ast.NewScalar(ast.KindString, ast.Value("hello"))),
		),
	}

	context := languages.Context{Schemas: ast.Schemas{schema}}
	files, err := jenny.Generate(context)
	req.NoError(err)
	req.Len(files, 1)

	output := string(files[0].Data)
	req.False(strings.Contains(output, "__eq__"), "scalar types should not generate __eq__ methods")
}

func TestEquality_Python_DisabledByDefault(t *testing.T) {
	req := require.New(t)

	config := Config{} // GenerateEqual defaults to false
	jenny := RawTypes{
		config:          config,
		tmpl:            initTemplates(config, common.NewAPIReferenceCollector()),
		apiRefCollector: common.NewAPIReferenceCollector(),
	}

	schema := equalitySchema()
	context := languages.Context{Schemas: ast.Schemas{schema}}

	files, err := jenny.Generate(context)
	req.NoError(err)
	req.Len(files, 1)

	output := string(files[0].Data)
	req.False(strings.Contains(output, "__eq__"), "equality methods should not be generated when GenerateEqual is false")
}
