package php

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

func TestEquality_PHP_GeneratesEqualsMethods(t *testing.T) {
	req := require.New(t)

	config := Config{GenerateEqual: true}
	jenny := RawTypes{
		config:          config,
		tmpl:            initTemplates(config, common.NewAPIReferenceCollector()),
		apiRefCollector: common.NewAPIReferenceCollector(),
	}

	schema := equalitySchema()

	// Run PHP compiler passes so nullable types are handled correctly
	processedSchemas, err := New(config).CompilerPasses().Process(ast.Schemas{schema})
	req.NoError(err)

	context := languages.Context{Schemas: processedSchemas}

	files, err := jenny.Generate(context)
	req.NoError(err)
	req.True(len(files) > 0)

	// Collect all output
	allOutput := make(map[string]string)
	for _, f := range files {
		allOutput[f.RelativePath] = string(f.Data)
	}

	// Variable class
	varOutput, found := allOutput["src/Equality/Variable.php"]
	req.True(found, "Variable.php should be generated")
	req.Contains(varOutput, "public function equals(mixed $other): bool")
	req.Contains(varOutput, "if (!($other instanceof self))")
	req.Contains(varOutput, "$this->name !== $other->name")

	// Container class with required struct ref
	containerOutput, found := allOutput["src/Equality/Container.php"]
	req.True(found, "Container.php should be generated")
	req.Contains(containerOutput, "public function equals(mixed $other): bool")
	req.Contains(containerOutput, "$this->stringField !== $other->stringField")
	req.Contains(containerOutput, "$this->intField !== $other->intField")
	req.Contains(containerOutput, "$this->refField->equals($other->refField)")

	// Optionals class with nullable fields (after compiler passes)
	optOutput, found := allOutput["src/Equality/Optionals.php"]
	req.True(found, "Optionals.php should be generated")
	req.Contains(optOutput, "public function equals(mixed $other): bool")
	req.Contains(optOutput, "$this->stringField === null") // null check for optional → nullable field
}

func TestEquality_PHP_SkipsNonStructTypes(t *testing.T) {
	req := require.New(t)

	config := Config{GenerateEqual: true}
	jenny := RawTypes{
		config:          config,
		tmpl:            initTemplates(config, common.NewAPIReferenceCollector()),
		apiRefCollector: common.NewAPIReferenceCollector(),
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
		req.False(strings.Contains(string(f.Data), "public function equals"), "scalar types should not generate equals methods")
	}
}

func TestEquality_PHP_DisabledByDefault(t *testing.T) {
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

	for _, f := range files {
		req.False(strings.Contains(string(f.Data), "public function equals"), "equality methods should not be generated when GenerateEqual is false")
	}
}
