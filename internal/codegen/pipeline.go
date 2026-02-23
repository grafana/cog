package codegen

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/ast/compiler"
	"github.com/grafana/cog/internal/jennies/golang"
	"github.com/grafana/cog/internal/jennies/java"
	"github.com/grafana/cog/internal/jennies/jsonschema"
	"github.com/grafana/cog/internal/jennies/openapi"
	"github.com/grafana/cog/internal/jennies/php"
	"github.com/grafana/cog/internal/jennies/python"
	"github.com/grafana/cog/internal/jennies/terraform"
	"github.com/grafana/cog/internal/jennies/typescript"
	"github.com/grafana/cog/internal/languages"
	"github.com/grafana/cog/internal/tools"
	"github.com/grafana/cog/internal/veneers/rewrite"
	cogyaml "github.com/grafana/cog/internal/yaml"
)

type ParametersInterpolator func(input string) string

type ProgressReporter func(msg string)

type PipelineOption func(pipeline *Pipeline)

func PipelineFromFile(file string, opts ...PipelineOption) (*Pipeline, error) {
	var err error
	if !filepath.IsAbs(file) {
		file, err = filepath.Abs(file)
		if err != nil {
			return nil, err
		}
	}

	fileHandle, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("could not open pipeline config: %w", err)
	}
	defer func() { _ = fileHandle.Close() }()

	pipeline, err := NewPipeline()
	if err != nil {
		return nil, err
	}

	decoder := yaml.NewDecoder(fileHandle, yaml.DisallowUnknownField())
	if err := decoder.Decode(pipeline); err != nil {
		return nil, fmt.Errorf("could not parse pipeline config:\n%s", yaml.FormatError(err, true, true))
	}

	pipeline.Parameters["__config_dir"] = filepath.Dir(file)

	pipeline.interpolator = createInterpolator(pipeline.Parameters)

	for _, opt := range opts {
		opt(pipeline)
	}

	return pipeline, nil
}

type Pipeline struct {
	Debug bool `yaml:"debug"`

	UnitsFrom  []string   `yaml:"units_from"`
	Inputs     []*Input   `yaml:"inputs"`
	Transforms Transforms `yaml:"transformations"`
	Output     Output     `yaml:"output"`

	Parameters map[string]string `yaml:"parameters"`

	currentDirectory string
	reporter         ProgressReporter
	veneersRewriter  *rewrite.Rewriter
	converterConfig  *languages.ConverterConfig
	interpolator     ParametersInterpolator
	unitsMerged      bool
}

func NewPipeline() (*Pipeline, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	pipeline := &Pipeline{
		reporter:         func(msg string) {},
		currentDirectory: currentDir,
		Parameters: map[string]string{
			"__config_dir": currentDir,
		},
	}

	pipeline.interpolator = createInterpolator(pipeline.Parameters)

	return pipeline, nil
}

func (pipeline *Pipeline) interpolateParameters() {
	pipeline.UnitsFrom = tools.Map(pipeline.UnitsFrom, pipeline.interpolator)

	for _, input := range pipeline.Inputs {
		// An error can only happen with the input isn't descriptive.
		// This case should have already been handled before
		// interpolateParameters() is called.
		_ = input.InterpolateParameters(pipeline.interpolator)
	}

	pipeline.Transforms.interpolateParameters(pipeline.interpolator)
	pipeline.Output.interpolateParameters(pipeline.interpolator)
}

func (pipeline *Pipeline) jenniesConfig() languages.Config {
	return languages.Config{
		Debug:        pipeline.Debug,
		Types:        pipeline.Output.Types,
		Builders:     pipeline.Output.Builders,
		Converters:   pipeline.Output.Converters,
		APIReference: pipeline.Output.APIReference,
	}
}

func (pipeline *Pipeline) commonPasses() (compiler.Passes, error) {
	if pipeline.Transforms.CommonPasses != nil {
		return pipeline.Transforms.CommonPasses, nil
	}

	return cogyaml.NewCompilerLoader().PassesFrom(pipeline.Transforms.CommonPassesFiles)
}

func (pipeline *Pipeline) finalPasses() compiler.Passes {
	return pipeline.Transforms.FinalPasses
}

func (pipeline *Pipeline) veneers() (*rewrite.Rewriter, error) {
	if pipeline.veneersRewriter != nil {
		return pipeline.veneersRewriter, nil
	}

	var veneerFiles []string
	for _, file := range pipeline.Transforms.VeneersPaths {
		isSingleFile, err := isFile(file)
		if err != nil {
			return nil, err
		}

		if isSingleFile {
			veneerFiles = append(veneerFiles, file)
		} else {
			globPattern := filepath.Join(file, "*.yaml")
			matches, err := filepath.Glob(globPattern)
			if err != nil {
				return nil, err
			}

			veneerFiles = append(veneerFiles, matches...)
		}
	}

	rewriter, err := cogyaml.NewVeneersLoader().RewriterFrom(veneerFiles, rewrite.Config{
		Debug: pipeline.Debug,
	})
	if err != nil {
		return nil, err
	}

	pipeline.veneersRewriter = rewriter

	return pipeline.veneersRewriter, nil
}

func (pipeline *Pipeline) readConverterConfig() (*languages.ConverterConfig, error) {
	if pipeline.converterConfig != nil {
		return pipeline.converterConfig, nil
	}

	cfg, err := cogyaml.NewConverterConfigReader().ReadConverterConfig(pipeline.Transforms.ConverterConfig)
	if err != nil {
		return nil, err
	}

	runtimeConfig := make([]languages.RuntimeConfig, len(cfg.Runtime))
	for i, rm := range cfg.Runtime {
		runtimeConfig[i] = languages.RuntimeConfig{
			Package:            rm.Package,
			Name:               rm.Name,
			NameFunc:           rm.NameFunc,
			DiscriminatorField: rm.DiscriminatorField,
		}
	}

	pipeline.converterConfig = &languages.ConverterConfig{
		RuntimeConfig: runtimeConfig,
	}
	return pipeline.converterConfig, err
}

func (pipeline *Pipeline) outputDir(relativeToDir string) (string, error) {
	if !filepath.IsAbs(pipeline.Output.Directory) {
		return pipeline.Output.Directory, nil
	}

	return filepath.Rel(relativeToDir, pipeline.Output.Directory)
}

func (pipeline *Pipeline) languageOutputDir(relativeToDir string, language string) (string, error) {
	outputDir, err := pipeline.outputDir(relativeToDir)
	if err != nil {
		return "", err
	}

	return strings.ReplaceAll(outputDir, "%l", language), nil
}

// mergeUnit merges the given unit into the pipeline.
func (pipeline *Pipeline) mergeUnit(unit *Unit) {
	for _, input := range unit.Inputs {
		pipeline.Inputs = append(pipeline.Inputs, input)
	}

	for _, builderTransform := range unit.BuilderTransforms {
		pipeline.Transforms.VeneersPaths = append(pipeline.Transforms.VeneersPaths, builderTransform)
	}
}

// loadUnits loads and merges the codegen units into the current pipeline.
func (pipeline *Pipeline) loadUnits() error {
	// Just making sure that this function is idempotent.
	if pipeline.unitsMerged {
		return nil
	}

	for _, unit := range pipeline.UnitsFrom {
		matches, err := filepath.Glob(unit)
		if err != nil {
			return err
		}

		for _, match := range matches {
			unit, err := unitFromFile(match, pipeline.Parameters)
			if err != nil {
				return err
			}

			pipeline.mergeUnit(unit)
		}
	}

	pipeline.unitsMerged = true

	return nil
}

// LoadSchemas parses the schemas described by the pipeline and applies common
// transformations.
// Note: input-specific transformations are applied.
func (pipeline *Pipeline) LoadSchemas(ctx context.Context) (ast.Schemas, error) {
	var allSchemas ast.Schemas
	var err error

	// Merge additional units into the main pipeline
	if err := pipeline.loadUnits(); err != nil {
		return nil, err
	}

	// Parse inputs
	for _, input := range pipeline.Inputs {
		schemas, err := input.LoadSchemas(ctx)
		if err != nil {
			return nil, err
		}

		allSchemas = append(allSchemas, schemas...)
	}

	if allSchemas == nil {
		return nil, nil
	}

	allSchemas, err = allSchemas.Consolidate()
	if err != nil {
		return nil, err
	}

	// Apply common and final compiler passes
	commonPasses, err := pipeline.commonPasses()
	if err != nil {
		return nil, err
	}

	return commonPasses.Process(allSchemas)
}

func (pipeline *Pipeline) OutputLanguages() (languages.Languages, error) {
	outputs := make(languages.Languages)

	for _, output := range pipeline.Output.Languages {
		switch {
		case output.Go != nil:
			outputs[golang.LanguageRef] = golang.New(*output.Go)
		case output.Java != nil:
			outputs[java.LanguageRef] = java.New(*output.Java)
		case output.JSONSchema != nil:
			outputs[jsonschema.LanguageRef] = jsonschema.New(*output.JSONSchema)
		case output.OpenAPI != nil:
			outputs[openapi.LanguageRef] = openapi.New(*output.OpenAPI)
		case output.PHP != nil:
			outputs[php.LanguageRef] = php.New(*output.PHP)
		case output.Python != nil:
			outputs[python.LanguageRef] = python.New(*output.Python)
		case output.Typescript != nil:
			outputs[typescript.LanguageRef] = typescript.New(*output.Typescript)
		case output.Terraform != nil:
			outputs[terraform.LanguageRef] = terraform.New(*output.Terraform)
		default:
			return nil, fmt.Errorf("empty language configuration")
		}
	}

	return outputs, nil
}
