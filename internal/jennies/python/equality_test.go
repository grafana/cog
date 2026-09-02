package python

import (
	"strings"
	"testing"

	"github.com/grafana/cog/internal/testutils"
	"github.com/grafana/cog/pkg/apiref"
	"github.com/grafana/cog/pkg/ir"
	"github.com/grafana/cog/pkg/languages"
	"github.com/grafana/cog/pkg/logs"
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
		),
	}
}

func TestEquality_Python_GeneratesEqMethods(t *testing.T) {
	req := require.New(t)

	config := Config{GenerateEqual: true}
	jenny := RawTypes{
		config:          config,
		tmpl:            initTemplates(config, apiref.NewAPIReferenceCollector()),
		apiRefCollector: apiref.NewAPIReferenceCollector(),
	}

	schema := equalitySchema()

	// Run Python compiler passes so nullable types are handled correctly
	processedSchemas, err := New(logs.NoopLogger(), config).Transform(ir.Schemas{schema})
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
		tmpl:            initTemplates(config, apiref.NewAPIReferenceCollector()),
		apiRefCollector: apiref.NewAPIReferenceCollector(),
	}

	// A scalar constant (non-struct type) should not generate __eq__
	schema := &ir.Schema{
		Package: "test",
		Objects: testutils.ObjectsMap(
			ir.NewObject("test", "MyConstant", ir.NewScalar(ir.KindString, ir.Value("hello"))),
		),
	}

	context := languages.Context{Schemas: ir.Schemas{schema}}
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
		tmpl:            initTemplates(config, apiref.NewAPIReferenceCollector()),
		apiRefCollector: apiref.NewAPIReferenceCollector(),
	}

	schema := equalitySchema()
	context := languages.Context{Schemas: ir.Schemas{schema}}

	files, err := jenny.Generate(context)
	req.NoError(err)
	req.Len(files, 1)

	output := string(files[0].Data)
	req.False(strings.Contains(output, "__eq__"), "equality methods should not be generated when GenerateEqual is false")
}
