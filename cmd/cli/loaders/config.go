package loaders

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
	"github.com/grafana/cog/internal/tools"
	"github.com/grafana/cog/internal/veneers/rewrite"
	cogyaml "github.com/grafana/cog/internal/yaml"
	"gopkg.in/yaml.v3"
)

type ParametersInterpolator func(input string) string

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

	Inputs     []Input    `yaml:"inputs"`
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
		input.InterpolateParameters(config.interpolate)
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

func (config Config) LoadSchemas() (ast.Schemas, error) {
	var allSchemas ast.Schemas

	for _, input := range config.Inputs {
		schemas, err := input.LoadSchemas()
		if err != nil {
			return nil, err
		}

		allSchemas = append(allSchemas, schemas...)
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

type Input struct {
	JSONSchema *JSONSchemaInput `yaml:"jsonschema"`
	OpenAPI    *OpenAPIInput    `yaml:"openapi"`

	KindRegistry      *KindRegistryInput `yaml:"kind_registry"`
	KindsysCore       *CueInput          `yaml:"kindsys_core"`
	KindsysComposable *CueInput          `yaml:"kindsys_composable"`
	Cue               *CueInput          `yaml:"cue"`
}

func (input *Input) InterpolateParameters(interpolator ParametersInterpolator) {
	if input.JSONSchema != nil {
		input.JSONSchema.InterpolateParameters(interpolator)
	}
	if input.OpenAPI != nil {
		input.OpenAPI.InterpolateParameters(interpolator)
	}
	if input.KindRegistry != nil {
		input.KindRegistry.InterpolateParameters(interpolator)
	}
	if input.KindsysCore != nil {
		input.KindsysCore.InterpolateParameters(interpolator)
	}
	if input.KindsysComposable != nil {
		input.KindsysComposable.InterpolateParameters(interpolator)
	}
	if input.Cue != nil {
		input.Cue.InterpolateParameters(interpolator)
	}
}

func (input *Input) LoadSchemas() (ast.Schemas, error) {
	if input.JSONSchema != nil {
		return input.JSONSchema.LoadSchemas()
	}
	if input.OpenAPI != nil {
		return input.OpenAPI.LoadSchemas()
	}
	if input.KindRegistry != nil {
		return input.KindRegistry.LoadSchemas()
	}
	if input.KindsysCore != nil {
		return kindsysCoreLoader(*input.KindsysCore)
	}
	if input.KindsysComposable != nil {
		return kindsysComposableLoader(*input.KindsysCore)
	}
	if input.Cue != nil {
		return cueLoader(*input.Cue)
	}

	return nil, fmt.Errorf("empty input")
}

type Transforms struct {
	CompilerPassesFiles []string `yaml:"schemas"`
	VeneersDirectories  []string `yaml:"builders"`
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
