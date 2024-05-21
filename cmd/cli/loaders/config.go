package loaders

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/expr-lang/expr"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/ast/compiler"
	"github.com/grafana/cog/internal/jennies"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/jennies/golang"
	"github.com/grafana/cog/internal/jennies/java"
	"github.com/grafana/cog/internal/jennies/jsonschema"
	"github.com/grafana/cog/internal/jennies/openapi"
	"github.com/grafana/cog/internal/jennies/python"
	"github.com/grafana/cog/internal/jennies/typescript"
	"github.com/grafana/cog/internal/semver"
	"github.com/grafana/cog/internal/tools"
	"github.com/grafana/cog/internal/veneers/rewrite"
	cogyaml "github.com/grafana/cog/internal/yaml"
	"gopkg.in/yaml.v3"
)

type ParametersInterpolator func(input string) string

type transformable interface {
	CompilerPasses() (compiler.Passes, error)
}

type interpolable interface {
	InterpolateParameters(interpolator ParametersInterpolator)
}

type schemaLoader interface {
	LoadSchemas(ctx context.Context) (ast.Schemas, error)
}

func ConfigFromFile(configFile string) (Config, error) {
	var err error
	if !filepath.IsAbs(configFile) {
		configFile, err = filepath.Abs(configFile)
		if err != nil {
			return Config{}, err
		}
	}

	currentDir, err := os.Getwd()
	if err != nil {
		return Config{}, err
	}

	file, err := os.Open(configFile)
	if err != nil {
		return Config{}, err
	}
	defer func() { _ = file.Close() }()

	decoder := yaml.NewDecoder(file)
	decoder.KnownFields(true)

	config := Config{
		Parameters: map[string]string{
			"__config_dir":  filepath.Dir(configFile),
			"__current_dir": currentDir,
		},
	}
	if err := decoder.Decode(&config); err != nil {
		return Config{}, err
	}

	return config, nil
}

type Config struct {
	Debug bool `yaml:"debug"`

	Inputs     []*Input   `yaml:"inputs"`
	Transforms Transforms `yaml:"transformations"`
	Output     Output     `yaml:"output"`

	Parameters map[string]string `yaml:"parameters"`
}

func (config Config) WithParameters(extraParameters map[string]string) Config {
	for key, value := range extraParameters {
		config.Parameters[key] = value
	}

	return config
}

func (config Config) InterpolateParameters() Config {
	interpolated := config

	for _, input := range config.Inputs {
		// An error can only happen with the input isn't descriptive.
		// This case should have already been handled before
		// InterpolateParameters() is called.
		_ = input.InterpolateParameters(config.interpolate)
	}

	interpolated.Transforms.InterpolateParameters(config.interpolate)
	interpolated.Output.InterpolateParameters(config.interpolate)

	return interpolated
}

func (config Config) interpolate(input string) string {
	interpolated := input

	for key, value := range config.Parameters {
		interpolated = strings.ReplaceAll(interpolated, "%"+key+"%", value)
	}

	return interpolated
}

func (config Config) JenniesConfig() common.Config {
	return common.Config{
		Debug:    config.Debug,
		Types:    config.Output.Types,
		Builders: config.Output.Builders,
	}
}

func (config Config) CompilerPasses() (compiler.Passes, error) {
	return cogyaml.NewCompilerLoader().PassesFrom(config.Transforms.CompilerPassesFiles)
}

func (config Config) Veneers() (*rewrite.Rewriter, error) {
	var veneers []string

	for _, dir := range config.Transforms.VeneersDirectories {
		globPattern := filepath.Join(dir, "*.yaml")
		matches, err := filepath.Glob(globPattern)
		if err != nil {
			return nil, err
		}

		veneers = append(veneers, matches...)
	}

	return cogyaml.NewVeneersLoader().RewriterFrom(veneers, rewrite.Config{
		Debug: config.Debug,
	})
}

func (config Config) OutputDir(relativeToDir string) (string, error) {
	if !filepath.IsAbs(config.Output.Directory) {
		return config.Output.Directory, nil
	}

	return filepath.Rel(relativeToDir, config.Output.Directory)
}

func (config Config) LanguageOutputDir(relativeToDir string, language string) (string, error) {
	outputDir, err := config.OutputDir(relativeToDir)
	if err != nil {
		return "", err
	}

	return strings.ReplaceAll(outputDir, "%l", language), nil
}

func (config Config) LoadSchemas(ctx context.Context) (ast.Schemas, error) {
	var allSchemas ast.Schemas

	for _, input := range config.Inputs {
		schemas, err := input.LoadSchemas(ctx)
		if err != nil {
			return nil, err
		}

		allSchemas = append(allSchemas, schemas...)
	}

	if allSchemas == nil {
		return nil, nil
	}

	return allSchemas.Consolidate()
}

func (config Config) OutputLanguages() (jennies.LanguageJennies, error) {
	outputs := make(jennies.LanguageJennies)

	for _, output := range config.Output.Languages {
		switch {
		case output.Go != nil:
			outputs[golang.LanguageRef] = golang.New(*output.Go)
		case output.Java != nil:
			outputs[java.LanguageRef] = java.New(*output.Java)
		case output.JSONSchema != nil:
			outputs[jsonschema.LanguageRef] = jsonschema.New()
		case output.OpenAPI != nil:
			outputs[openapi.LanguageRef] = openapi.New()
		case output.Python != nil:
			outputs[python.LanguageRef] = python.New(*output.Python)
		case output.Typescript != nil:
			outputs[typescript.LanguageRef] = typescript.New()
		default:
			return nil, fmt.Errorf("empty language configuration")
		}
	}

	return outputs, nil
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
}

func (input *InputBase) CompilerPasses() (compiler.Passes, error) {
	return cogyaml.NewCompilerLoader().PassesFrom(input.Transforms)
}

func (input *InputBase) InterpolateParameters(interpolator ParametersInterpolator) {
	input.AllowedObjects = tools.Map(input.AllowedObjects, interpolator)
	input.Transforms = tools.Map(input.Transforms, interpolator)
}

func (input *InputBase) filterSchema(schema *ast.Schema) (ast.Schemas, error) {
	if len(input.AllowedObjects) == 0 {
		return ast.Schemas{schema}, nil
	}

	filterPass := compiler.FilterSchemas{
		AllowedObjects: tools.Map(input.AllowedObjects, func(objectName string) compiler.ObjectReference {
			return compiler.ObjectReference{Package: schema.Package, Object: objectName}
		}),
	}

	return filterPass.Process(ast.Schemas{schema})
}

type Input struct {
	If string `yaml:"if"`

	JSONSchema *JSONSchemaInput `yaml:"jsonschema"`
	OpenAPI    *OpenAPIInput    `yaml:"openapi"`

	KindRegistry      *KindRegistryInput `yaml:"kind_registry"`
	KindsysCore       *CueInput          `yaml:"kindsys_core"`
	KindsysComposable *CueInput          `yaml:"kindsys_composable"`
	Cue               *CueInput          `yaml:"cue"`
}

func (input *Input) InterpolateParameters(interpolator ParametersInterpolator) error {
	input.If = interpolator(input.If)

	loader, err := input.loader()
	if err != nil {
		return err
	}

	if interpolableLoader, ok := loader.(interpolable); ok {
		interpolableLoader.InterpolateParameters(interpolator)
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

func (input *Input) shouldLoadSchemas() (bool, error) {
	if input.If == "" {
		return true, nil
	}

	env := map[string]any{
		"sprintf": fmt.Sprintf,
		"semver":  semver.ParseTolerant,
	}

	program, err := expr.Compile(input.If, expr.Env(env))
	if err != nil {
		return false, err
	}

	output, err := expr.Run(program, env)
	if err != nil {
		return false, err
	}

	if _, ok := output.(bool); !ok {
		return false, fmt.Errorf("expected expression to evaluate to a boolean, got %T", output)
	}

	return output.(bool), nil
}

func (input *Input) LoadSchemas(ctx context.Context) (ast.Schemas, error) {
	var err error

	shouldLoad, err := input.shouldLoadSchemas()
	if err != nil {
		return nil, err
	}
	if !shouldLoad {
		return nil, nil
	}

	loader, err := input.loader()
	if err != nil {
		return nil, err
	}

	schemas, err := loader.LoadSchemas(ctx)
	if err != nil {
		return nil, err
	}

	if transformableLoader, ok := loader.(transformable); ok {
		passes, err := transformableLoader.CompilerPasses()
		if err != nil {
			return nil, err
		}

		return passes.Process(schemas)
	}

	return schemas, nil
}

type Transforms struct {
	// CompilerPassesFiles holds a list of paths to files containing compiler
	// passes to apply to all the schemas.
	CompilerPassesFiles []string `yaml:"schemas"`

	// VeneersDirectories holds a list of paths to directories containing
	// veneers to apply to all the builders.
	VeneersDirectories []string `yaml:"builders"`
}

func (transforms *Transforms) InterpolateParameters(interpolator ParametersInterpolator) {
	transforms.CompilerPassesFiles = tools.Map(transforms.CompilerPassesFiles, interpolator)
	transforms.VeneersDirectories = tools.Map(transforms.VeneersDirectories, interpolator)
}

type Output struct {
	Directory string `yaml:"directory"`

	Types    bool `yaml:"types"`
	Builders bool `yaml:"builders"`

	Languages []*OutputLanguage `yaml:"languages"`

	// PackageTemplates is the path to a directory containing "package templates".
	// These templates are used to add arbitrary files to the generated code, with
	// the goal of turning it into a fully-fledged package.
	// Templates in that directory are expected to be organized by language:
	// ```
	// package_templates
	// ├── go
	// │   ├── LICENSE.md
	// │   └── README.md
	// └── typescript
	//     ├── babel.config.json
	//     ├── package.json
	//     ├── README.md
	//     └── tsconfig.json
	// ```
	PackageTemplates string `yaml:"package_templates"`

	// RepositoryTemplates is the path to a directory containing
	// "repository-level templates".
	// These templates are used to add arbitrary files to the repository, such as CI pipelines.
	//
	// Templates in that directory are expected to be organized by language:
	// ```
	// repository_templates
	// ├── go
	// │   └── .github
	// │   	   └── workflows
	// │   	       └── go-ci.yaml
	// └── typescript
	//     └── .github
	//     	   └── workflows
	//     	       └── typescript-ci.yaml
	// ```
	RepositoryTemplates string `yaml:"repository_templates"`

	// TemplatesData holds data that will be injected into package and
	// repository templates when rendering them.
	TemplatesData map[string]string `yaml:"templates_data"`
}

func (output *Output) InterpolateParameters(interpolator ParametersInterpolator) {
	output.Directory = interpolator(output.Directory)

	for _, outputLanguage := range output.Languages {
		outputLanguage.InterpolateParameters(interpolator)
	}

	output.PackageTemplates = interpolator(output.PackageTemplates)
	output.RepositoryTemplates = interpolator(output.RepositoryTemplates)

	for key, value := range output.TemplatesData {
		output.TemplatesData[key] = interpolator(value)
	}
}

type OutputLanguage struct {
	Go         *golang.Config `yaml:"go"`
	Java       *java.Config   `yaml:"java"`
	JSONSchema *NoConfig      `yaml:"jsonschema"`
	OpenAPI    *NoConfig      `yaml:"openapi"`
	Python     *python.Config `yaml:"python"`
	Typescript *NoConfig      `yaml:"typescript"`
}

func (outputLanguage *OutputLanguage) InterpolateParameters(interpolator ParametersInterpolator) {
	if outputLanguage.Go != nil {
		outputLanguage.Go.InterpolateParameters(interpolator)
	}
	if outputLanguage.Python != nil {
		outputLanguage.Python.InterpolateParameters(interpolator)
	}
}

type NoConfig struct {
}
