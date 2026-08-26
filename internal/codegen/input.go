package codegen

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/grafana/cog/internal/ir"
	"github.com/grafana/cog/internal/ir/transforms"
	"github.com/grafana/cog/internal/tools"
	cogyaml "github.com/grafana/cog/internal/yaml"
)

type interpolable interface {
	interpolateParameters(interpolator ParametersInterpolator)
}

type transformable interface {
	commonPasses() (transforms.Transforms, error)
}

type schemaLoader interface {
	LoadSchemas(ctx context.Context) (ir.Schemas, error)
}

// InputBase provides common options and behavior, meant to be re-used across
// all input types.
type InputBase struct {
	// AllowedObjects is a list of object names that will be allowed when
	// parsing the input schema.
	// Note: if AllowedObjects is empty, no filter is applied.
	AllowedObjects []string `yaml:"allowed_objects"`

	// Transforms holds a list of paths to files containing compiler passes
	// to apply to the input.
	Transforms []string `yaml:"transformations"`

	// Metadata to add to the schema, this can be used to set Kind and Variant
	Metadata *ir.SchemaMeta `yaml:"metadata"`
}

func (input *InputBase) schemaMetadata() ir.SchemaMeta {
	if input.Metadata != nil {
		return *input.Metadata
	}

	return ir.SchemaMeta{}
}

func (input *InputBase) commonPasses() (transforms.Transforms, error) {
	return cogyaml.NewTransformsLoader().LoadFiles(input.Transforms)
}

func (input *InputBase) interpolateParameters(interpolator ParametersInterpolator) {
	input.AllowedObjects = tools.Map(input.AllowedObjects, interpolator)
	input.Transforms = tools.Map(input.Transforms, interpolator)
}

func (input *InputBase) filterSchema(schema *ir.Schema) (ir.Schemas, error) {
	if len(input.AllowedObjects) == 0 {
		return ir.Schemas{schema}, nil
	}

	filterPass := transforms.FilterSchemas{
		AllowedObjects: tools.Map(input.AllowedObjects, func(objectName string) transforms.ObjectReference {
			return transforms.ObjectReference{Package: schema.Package, Object: objectName}
		}),
	}

	return filterPass.Process(ir.Schemas{schema})
}

type Input struct {
	JSONSchema *JSONSchemaInput `yaml:"jsonschema"`
	OpenAPI    *OpenAPIInput    `yaml:"openapi"`

	KindRegistry      *KindRegistryInput `yaml:"kind_registry"`
	KindsysCore       *CueInput          `yaml:"kindsys_core"`
	KindsysComposable *CueInput          `yaml:"kindsys_composable"`
	Cue               *CueInput          `yaml:"cue"`
}

func (input *Input) InterpolateParameters(interpolator ParametersInterpolator) error {
	loader, err := input.loader()
	if err != nil {
		return err
	}

	if interpolableLoader, ok := loader.(interpolable); ok {
		interpolableLoader.interpolateParameters(interpolator)
	}

	return nil
}

func (input *Input) loader() (schemaLoader, error) {
	if input.JSONSchema != nil {
		return input.JSONSchema, nil
	}
	if input.OpenAPI != nil {
		return input.OpenAPI, nil
	}
	if input.KindRegistry != nil {
		return input.KindRegistry, nil
	}
	if input.KindsysCore != nil {
		return &genericCueLoader{CueInput: input.KindsysCore, loader: kindsysCoreLoader}, nil
	}
	if input.KindsysComposable != nil {
		return &genericCueLoader{CueInput: input.KindsysComposable, loader: kindsysComposableLoader}, nil
	}
	if input.Cue != nil {
		return &genericCueLoader{CueInput: input.Cue, loader: cueLoader}, nil
	}

	return nil, fmt.Errorf("empty input")
}

func (input *Input) LoadSchemas(ctx context.Context, logger *slog.Logger) (ir.Schemas, error) {
	var err error

	loader, err := input.loader()
	if err != nil {
		return nil, err
	}

	schemas, err := loader.LoadSchemas(ctx)
	if err != nil {
		return nil, err
	}

	if transformableLoader, ok := loader.(transformable); ok {
		passes, err := transformableLoader.commonPasses()
		if err != nil {
			return nil, err
		}

		schemas, err = passes.Process(logger, schemas)
		if err != nil {
			return nil, err
		}
	}

	return schemas, nil
}
