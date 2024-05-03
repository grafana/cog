package loaders

import (
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/openapi"
	"github.com/grafana/cog/internal/tools"
)

type OpenAPIInput struct {
	// Path to an OpenAPI file.
	Path string `yaml:"path"`

	// Package name to use for the input schema. If empty, it will be guessed
	// from the input file name.
	Package string `yaml:"package"`

	// AllowedObjects is a list of object names that will be allowed when
	// parsing the input schema.
	// Note: if AllowedObjects is empty, no filter is applied.
	AllowedObjects []string `yaml:"allowed_objects"`
}

func (input *OpenAPIInput) InterpolateParameters(interpolator ParametersInterpolator) {
	input.Path = interpolator(input.Path)
	input.Package = interpolator(input.Package)
	input.AllowedObjects = tools.Map(input.AllowedObjects, interpolator)
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

	return filterSchema(schema, input.AllowedObjects)
}
