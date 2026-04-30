package java

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

func TestEquality_Java_GeneratesEqualsMethods(t *testing.T) {
	req := require.New(t)

	config := Config{GenerateEqual: true}
	jenny := RawTypes{
		config: config,
		tmpl:   initTemplates(config, common.NewAPIReferenceCollector()),
	}

	schema := equalitySchema()

	// Run Java compiler passes so nullable types are handled correctly
	processedSchemas, err := New(config).CompilerPasses().Process(ast.Schemas{schema})
	req.NoError(err)

	context := languages.Context{Schemas: processedSchemas}

	files, err := jenny.Generate(context)
	req.NoError(err)
	req.True(len(files) > 0)

	// Collect all output by filename
	allOutput := make(map[string]string)
	for _, f := range files {
		allOutput[f.RelativePath] = string(f.Data)
	}

	// Variable class - path depends on config.ProjectPath (empty by default)
	varOutput, found := allOutput["equality/Variable.java"]
	req.True(found, "Variable.java should be generated, got: %v", func() []string {
		keys := make([]string, 0, len(allOutput))
		for k := range allOutput {
			keys = append(keys, k)
		}
		return keys
	}())
	req.Contains(varOutput, "@Override")
	req.Contains(varOutput, "public boolean equals(Object other)")
	req.Contains(varOutput, "!(other instanceof Variable)")
	req.Contains(varOutput, "Objects.equals(this.name, o.name)")
	req.Contains(varOutput, "public int hashCode()")
	req.Contains(varOutput, "Objects.hash(")
	req.Contains(varOutput, "import java.util.Objects;")

	// Container class with scalar + struct ref fields
	containerOutput, found := allOutput["equality/Container.java"]
	req.True(found, "Container.java should be generated")
	req.Contains(containerOutput, "public boolean equals(Object other)")
	req.Contains(containerOutput, "!(other instanceof Container)")
	req.Contains(containerOutput, "Objects.equals(this.stringField, o.stringField)")
	req.Contains(containerOutput, "Objects.equals(this.intField, o.intField)")
	req.Contains(containerOutput, "Objects.equals(this.refField, o.refField)")
}

func TestEquality_Java_SkipsNonStructTypes(t *testing.T) {
	req := require.New(t)

	config := Config{GenerateEqual: true}
	jenny := RawTypes{
		config: config,
		tmpl:   initTemplates(config, common.NewAPIReferenceCollector()),
	}

	// A scalar constant (non-struct type) should not generate equals()
	schema := &ast.Schema{
		Package: "test",
		Objects: testutils.ObjectsMap(
			ast.NewObject("test", "MyConstant", ast.NewScalar(ast.KindString, ast.Value("hello"))),
		),
	}

	context := languages.Context{Schemas: ast.Schemas{schema}}
	files, err := jenny.Generate(context)
	req.NoError(err)

	for _, f := range files {
		req.False(strings.Contains(string(f.Data), "public boolean equals"), "scalar types should not generate equals methods")
	}
}

func TestEquality_Java_DisabledByDefault(t *testing.T) {
	req := require.New(t)

	config := Config{} // GenerateEqual defaults to false
	jenny := RawTypes{
		config: config,
		tmpl:   initTemplates(config, common.NewAPIReferenceCollector()),
	}

	schema := equalitySchema()
	context := languages.Context{Schemas: ast.Schemas{schema}}

	files, err := jenny.Generate(context)
	req.NoError(err)

	for _, f := range files {
		req.False(strings.Contains(string(f.Data), "public boolean equals"), "equality methods should not be generated when GenerateEqual is false")
	}
}
