package codegen

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/openapi"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type K8APIInput struct {
	InputBase `yaml:",inline"`

	URL string `yaml:"url"`
}

func (k8 *K8APIInput) interpolateParameters(interpolator ParametersInterpolator) {
	k8.InputBase.interpolateParameters(interpolator)

	k8.URL = interpolator(k8.URL)
}

func (k8 *K8APIInput) LoadSchemas(ctx context.Context) (ast.Schemas, error) {
	specs, err := k8.extractGroups()
	if err != nil {
		return nil, err
	}

	var schemas ast.Schemas

	for _, spec := range specs {
		oapiSchema, err := k8.loadSchema(ctx, spec.url)
		if err != nil {
			return nil, err
		}

		name := strings.Split(spec.name, ".")[0]

		schema, err := openapi.GenerateAST(ctx, oapiSchema, openapi.Config{
			Package: strings.Split(spec.name, ".")[0],
			SchemaMetadata: ast.SchemaMeta{
				Kind:       ast.SchemaKindCore,
				Identifier: name,
			},
			Validate: true,
		})
		if err != nil {
			return nil, err
		}

		schemas = append(schemas, schema)
	}

	return schemas, nil
}

func (k8 *K8APIInput) extractGroups() ([]schemaSpecs, error) {
	res, err := http.Get(k8.URL + "/apis")
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var specs groupSpec
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(b, &specs); err != nil {
		return nil, err
	}

	var paths []schemaSpecs
	for _, g := range specs.Groups {
		paths = append(paths, schemaSpecs{
			name:    g.Name,
			version: g.PreferredVersion.Version,
			url:     fmt.Sprintf("%s/%s/%s", k8.URL, "openapi/v3/apis", g.PreferredVersion.GroupVersion),
		})
	}

	return paths, nil
}

func (k8 *K8APIInput) loadSchema(ctx context.Context, u string) (*openapi3.T, error) {
	loader := openapi3.NewLoader()
	loader.Context = ctx
	loader.IsExternalRefsAllowed = true

	parsedURL, err := url.Parse(u)
	if err != nil {
		return nil, err
	}

	return loader.LoadFromURI(parsedURL)
}

type schemaSpecs struct {
	name    string
	version string
	url     string
}

type groupSpec struct {
	Kind       string  `yaml:"kind"`
	APIVersion string  `yaml:"apiVersion"`
	Groups     []group `yaml:"groups"`
}

type group struct {
	Name             string `yaml:"name"`
	PreferredVersion struct {
		GroupVersion string `yaml:"groupVersion"`
		Version      string `yaml:"version"`
	} `yaml:"preferredVersion"`
}
