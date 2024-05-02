package loaders

import (
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/openapi"
)

type OpenAPIInput struct {
	// Path to an OpenAPI file.
	Path string `yaml:"path"`

	Package string `yaml:"package"`
}

func (input *OpenAPIInput) InterpolateParameters(interpolator ParametersInterpolator) {
	input.Path = interpolator(input.Path)
	input.Package = interpolator(input.Package)
}

func (input OpenAPIInput) LoadSchemas() (ast.Schemas, error) {
	pkg := input.Package
	if pkg == "" {
		pkg = guessPackageFromFilename(input.Path)
	}

	schema, err := openapi.GenerateAST(input.Path, openapi.Config{
		Package:        pkg,
		SchemaMetadata: ast.SchemaMeta{}, // TODO: extract these from somewhere
	})
	if err != nil {
		return nil, err
	}

	return ast.Schemas{schema}, nil
}
