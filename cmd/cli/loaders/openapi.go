package loaders

import (
	"context"
	"net/url"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/openapi"
	"github.com/grafana/cog/internal/tools"
)

type OpenAPIInput struct {
	// Path to an OpenAPI file.
	Path string `yaml:"path"`

	// URL to an OpenAPI file.
	URL string `yaml:"url"`

	// Package name to use for the input schema. If empty, it will be guessed
	// from the input file name.
	Package string `yaml:"package"`

	// NoValidate disables validation of the OpenAPI spec.
	NoValidate bool `yaml:"no_validate"`

	// AllowedObjects is a list of object names that will be allowed when
	// parsing the input schema.
	// Note: if AllowedObjects is empty, no filter is applied.
	AllowedObjects []string `yaml:"allowed_objects"`
}

func (input *OpenAPIInput) loadSchema(ctx context.Context) (*openapi3.T, error) {
	loader := openapi3.NewLoader()
	loader.Context = ctx
	loader.IsExternalRefsAllowed = true

	if input.Path != "" {
		return loader.LoadFromFile(input.Path)
	}

	parsedURL, err := url.Parse(input.URL)
	if err != nil {
		return nil, err
	}

	return loader.LoadFromURI(parsedURL)
}

func (input *OpenAPIInput) packageName() string {
	if input.Package != "" {
		return input.Package
	}

	return guessPackageFromFilename(input.Path)
}

func (input *OpenAPIInput) InterpolateParameters(interpolator ParametersInterpolator) {
	input.Path = interpolator(input.Path)
	input.URL = interpolator(input.URL)
	input.Package = interpolator(input.Package)
	input.AllowedObjects = tools.Map(input.AllowedObjects, interpolator)
}

func (input OpenAPIInput) LoadSchemas(ctx context.Context) (ast.Schemas, error) {
	oapiSchema, err := input.loadSchema(ctx)
	if err != nil {
		return nil, err
	}

	schema, err := openapi.GenerateAST(ctx, oapiSchema, openapi.Config{
		Package:        input.packageName(),
		SchemaMetadata: ast.SchemaMeta{}, // TODO: extract these from somewhere
		Validate:       !input.NoValidate,
	})
	if err != nil {
		return nil, err
	}

	return filterSchema(schema, input.AllowedObjects)
}
