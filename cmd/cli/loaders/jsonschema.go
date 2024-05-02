package loaders

import (
	"os"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jsonschema"
)

type JSONSchemaInput struct {
	// Path to a JSONSchema file.
	Path string `yaml:"path"`

	Package string `yaml:"package"`
}

func (input *JSONSchemaInput) InterpolateParameters(interpolator ParametersInterpolator) {
	input.Path = interpolator(input.Path)
	input.Package = interpolator(input.Package)
}

func (input *JSONSchemaInput) LoadSchemas() (ast.Schemas, error) {
	reader, err := os.Open(input.Path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = reader.Close() }()

	pkg := input.Package
	if pkg == "" {
		pkg = guessPackageFromFilename(input.Path)
	}

	schema, err := jsonschema.GenerateAST(reader, jsonschema.Config{
		Package:        pkg,
		SchemaMetadata: ast.SchemaMeta{}, // TODO: extract these from somewhere
	})
	if err != nil {
		return nil, err
	}

	return ast.Schemas{schema}, nil
}
