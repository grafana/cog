package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/grafana/cog/internal/codegen"
	"github.com/grafana/cog/internal/tools"
	"github.com/grafana/cog/internal/yaml"
	"github.com/invopop/jsonschema"
)

func main() {
	types := []struct {
		name  string
		input any
	}{
		{name: "compiler_passes", input: &yaml.Compiler{}},
		{name: "veneers", input: &yaml.Veneers{}},
		{name: "pipeline", input: &codegen.Pipeline{}},
	}

	reflector := &jsonschema.Reflector{
		RequiredFromJSONSchemaTags: true, // we don't have good information as to what is required :/
		FieldNameTag:               "yaml",
		KeyNamer: func(key string) string {
			if strings.ToUpper(string(key[0])) == string(key[0]) {
				return strings.ToLower(key)
			}

			return key
		},
		Namer: func(reflected reflect.Type) string {
			name := reflected.Name()

			if reflected.PkgPath() != "" {
				parts := strings.Split(reflected.PkgPath(), "/")
				name = tools.UpperCamelCase(parts[len(parts)-1]) + tools.UpperCamelCase(name)
			}

			return name
		},
	}

	for _, t := range types {
		fmt.Printf("Generating schema for type '%s'\n", t.name)

		if err := reflector.AddGoComments("github.com/grafana/cog", "./internal"); err != nil {
			panic(fmt.Errorf("could not add Go comments to reflector: %w", err))
		}

		schema := reflector.Reflect(t.input)
		schema.ID = jsonschema.ID(fmt.Sprintf("https://raw.githubusercontent.com/grafana/cog/main/schemas/%s.json", t.name))

		schemaJSON, err := json.MarshalIndent(schema, "", "  ")
		if err != nil {
			panic(fmt.Errorf("could not marshal schema to JSON: %w", err))
		}

		if err := os.WriteFile(fmt.Sprintf("./schemas/%s.json", t.name), schemaJSON, 0600); err != nil {
			panic(fmt.Errorf("could not write schema: %w", err))
		}
	}
}
