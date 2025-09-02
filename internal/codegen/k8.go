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

	URL            string                   `yaml:"url"`
	AllowedSchemas map[string]AllowedSchema `yaml:"allowed_schemas"`
}

type AllowedSchema struct {
	Version        string   `json:"version"`
	AllowedObjects []string `yaml:"allowed_objects"`
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
		name := strings.Split(spec.name, ".")[0]
		if !k8.shouldParseSchema(name) {
			continue
		}

		oapiSchema, err := k8.loadSchema(ctx, spec.url)
		if err != nil {
			return nil, err
		}

		schema, err := openapi.GenerateAST(ctx, oapiSchema, openapi.Config{
			Package: name + spec.version,
			SchemaMetadata: ast.SchemaMeta{
				Kind:       ast.SchemaKindCore,
				Identifier: name + spec.version,
			},
			Validate: false,
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
		v := k8.getSchemaVersion(g)
		paths = append(paths, schemaSpecs{
			name:    g.Name,
			version: v.Version,
			url:     fmt.Sprintf("%s/%s/%s", k8.URL, "openapi/v3/apis", v.GroupVersion),
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

func (k8 *K8APIInput) getSchemaVersion(g group) version {
	if k8.AllowedSchemas == nil {
		return g.PreferredVersion
	}

	name := strings.Split(g.Name, ".")[0]
	schema, ok := k8.AllowedSchemas[name]
	if !ok {
		return g.PreferredVersion
	}

	for _, v := range g.Versions {
		if v.Version == schema.Version {
			return v
		}
	}

	return g.PreferredVersion
}

func (k8 *K8APIInput) shouldParseSchema(name string) bool {
	if len(k8.AllowedSchemas) == 0 {
		return true
	}

	for allowed := range k8.AllowedSchemas {
		if allowed == name {
			return true
		}
	}

	return false
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
	Name             string    `yaml:"name"`
	Versions         []version `yaml:"versions"`
	PreferredVersion version   `yaml:"preferredVersion"`
}

type version struct {
	GroupVersion string `yaml:"groupVersion"`
	Version      string `yaml:"version"`
}
