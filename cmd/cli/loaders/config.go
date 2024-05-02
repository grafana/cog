package loaders

import (
	"fmt"
	"os"
	"path"
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

func ConfigFromFile(configFile string) (Config, error) {
	var err error
	if !filepath.IsAbs(configFile) {
		configFile, err = filepath.Abs(configFile)
		if err != nil {
			return Config{}, err
		}
	}

	file, err := os.Open(configFile)
	if err != nil {
		return Config{}, err
	}
	defer func() { _ = file.Close() }()

	decoder := yaml.NewDecoder(file)
	decoder.KnownFields(true)

	config := Config{
		RootDir: filepath.Dir(configFile),
	}
	if err := decoder.Decode(&config); err != nil {
		return Config{}, err
	}

	return config, nil
}

type Config struct {
	RootDir string `yaml:"-"`
	Debug   bool   `yaml:"debug"`

	Inputs     []Input    `yaml:"inputs"`
	Transforms Transforms `yaml:"transforms"`
	Output     Output     `yaml:"output"`
}

func (config Config) Path(inputPath string) string {
	if path.IsAbs(inputPath) {
		return inputPath
	}

	return filepath.Clean(path.Join(config.RootDir, inputPath))
}

func (config Config) JenniesConfig() common.Config {
	return common.Config{
		Debug:    config.Debug,
		Types:    config.Output.Types,
		Builders: config.Output.Builders,
	}
}

func (config Config) CommonCompilerPasses() (compiler.Passes, error) {
	return cogyaml.NewCompilerLoader().PassesFrom(
		tools.Map(config.Transforms.CompilerPassesFiles, config.Path),
	)
}

func (config Config) veneerFiles() ([]string, error) {
	var veneers []string

	for _, dir := range config.Transforms.VeneersDirectories {
		globPattern := filepath.Join(config.Path(dir), "*.yaml")
		matches, err := filepath.Glob(globPattern)
		if err != nil {
			return nil, err
		}

		veneers = append(veneers, matches...)
	}

	return veneers, nil
}

func (config Config) Veneers() (*rewrite.Rewriter, error) {
	veneerFiles, err := config.veneerFiles()
	if err != nil {
		return nil, err
	}

	return cogyaml.NewVeneersLoader().RewriterFrom(veneerFiles, rewrite.Config{
		Debug: config.Debug,
	})
}

func (config Config) OutputDir(relativeToDir string) (string, error) {
	return filepath.Rel(relativeToDir, config.Path(config.Output.Directory))
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
		schemas, err := input.LoadSchemas(config)
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
		if output.Go != nil {
			outputs[golang.LanguageRef] = golang.New(golang.Config{
				GenerateGoMod: output.Go.GenerateGoMod,
				PackageRoot:   output.Go.PackageRoot,
			})
		} else if output.Java != nil {
			outputs[java.LanguageRef] = java.New(java.Config{
				GenGettersAndSetters: output.Java.GenGettersAndSetters,
			})
		} else if output.JSONSchema != nil {
			outputs[jsonschema.LanguageRef] = jsonschema.New()
		} else if output.OpenAPI != nil {
			outputs[openapi.LanguageRef] = openapi.New()
		} else if output.Python != nil {
			outputs[python.LanguageRef] = python.New(python.Config{
				PathPrefix: output.Python.PathPrefix,
			})
		} else if output.Typescript != nil {
			outputs[typescript.LanguageRef] = typescript.New()
		} else {
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

func (input Input) LoadSchemas(config Config) (ast.Schemas, error) {
	if input.JSONSchema != nil {
		return input.JSONSchema.LoadSchemas(config)
	}
	if input.OpenAPI != nil {
		return input.OpenAPI.LoadSchemas(config)
	}
	if input.KindRegistry != nil {
		return input.KindRegistry.LoadSchemas(config)
	}
	if input.KindsysCore != nil {
		return kindsysCoreLoader(config, *input.KindsysCore)
	}
	if input.KindsysComposable != nil {
		return kindsysComposableLoader(config, *input.KindsysCore)
	}
	if input.Cue != nil {
		return cueLoader(config, *input.Cue)
	}

	return nil, fmt.Errorf("empty input")
}

type Transforms struct {
	CompilerPassesFiles []string `yaml:"compiler_passes"`
	VeneersDirectories  []string `yaml:"veneers"`
}

type Output struct {
	Directory string `yaml:"directory"`

	Types    bool `yaml:"types"`
	Builders bool `yaml:"builders"`

	Languages []OutputLanguage `yaml:"languages"`

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

type OutputLanguage struct {
	Go         *OutputGo         `yaml:"go"`
	Java       *OutputJava       `yaml:"java"`
	JSONSchema *OutputJSONSchema `yaml:"jsonschema"`
	OpenAPI    *OutputOpenAPI    `yaml:"openapi"`
	Python     *OutputPython     `yaml:"python"`
	Typescript *OutputTypescript `yaml:"typescript"`
}

type OutputGo struct {
	// GenerateGoMod indicates whether a go.mod file should be generated.
	// If enabled, PackageRoot is used as module path.
	GenerateGoMod bool `yaml:"go_mod"`

	// Root path for imports.
	// Ex: github.com/grafana/cog/generated
	PackageRoot string `yaml:"package_root"`
}

type OutputJava struct {
	GenGettersAndSetters bool `yaml:"getters_and_setters"`
}

type OutputJSONSchema struct {
}

type OutputOpenAPI struct {
}

type OutputPython struct {
	PathPrefix string `yaml:"path_prefix"`
}

type OutputTypescript struct {
}
