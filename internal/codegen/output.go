package codegen

import (
	"github.com/grafana/cog/internal/jennies/golang"
	"github.com/grafana/cog/internal/jennies/java"
	"github.com/grafana/cog/internal/jennies/jsonschema"
	"github.com/grafana/cog/internal/jennies/openapi"
	"github.com/grafana/cog/internal/jennies/php"
	"github.com/grafana/cog/internal/jennies/python"
	"github.com/grafana/cog/internal/jennies/terraform"
	"github.com/grafana/cog/internal/jennies/typescript"
)

type Output struct {
	Directory string `yaml:"directory"`

	Types        bool `yaml:"types"`
	Builders     bool `yaml:"builders"`
	Converters   bool `yaml:"converters"`
	APIReference bool `yaml:"api_reference"`

	Languages []*OutputLanguage `yaml:"languages"`

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
	RepositoryTemplates string `yaml:"repository_templates,omitempty"`

	// TemplatesData holds data that will be injected into package and
	// repository templates when rendering them.
	TemplatesData map[string]string `yaml:"templates_data,omitempty"`

	// OutputOptions configures the output of the file
	OutputOptions OutputOptions `yaml:"output_options,omitempty"`

	LanguagePlugins map[string]OutputLanguagePlugin `yaml:"language_plugins,omitempty"`
}

func (output *Output) interpolateParameters(interpolator ParametersInterpolator) {
	output.Directory = interpolator(output.Directory)
	output.RepositoryTemplates = interpolator(output.RepositoryTemplates)

	for key, value := range output.TemplatesData {
		output.TemplatesData[key] = interpolator(value)
	}

	for _, outputLanguage := range output.Languages {
		outputLanguage.interpolateParameters(output, interpolator)
	}

	for _, config := range output.LanguagePlugins {
		config.interpolateParameters(interpolator)
	}
}

type OutputLanguage struct {
	Go         *golang.Config     `yaml:"go,omitempty"`
	Java       *java.Config       `yaml:"java,omitempty"`
	JSONSchema *jsonschema.Config `yaml:"jsonschema,omitempty"`
	OpenAPI    *openapi.Config    `yaml:"openapi,omitempty"`
	PHP        *php.Config        `yaml:"php,omitempty"`
	Python     *python.Config     `yaml:"python,omitempty"`
	Typescript *typescript.Config `yaml:"typescript,omitempty"`
	Terraform  *terraform.Config  `yaml:"terraform,omitempty"`
}

func (outputLanguage *OutputLanguage) interpolateParameters(output *Output, interpolator ParametersInterpolator) {
	if outputLanguage.Go != nil {
		outputLanguage.Go.InterpolateParameters(interpolator)
		outputLanguage.Go.ExtraFilesTemplatesData = output.TemplatesData
	}
	if outputLanguage.PHP != nil {
		outputLanguage.PHP.InterpolateParameters(interpolator)
		outputLanguage.PHP.ExtraFilesTemplatesData = output.TemplatesData
	}
	if outputLanguage.Python != nil {
		outputLanguage.Python.InterpolateParameters(interpolator)
		outputLanguage.Python.ExtraFilesTemplatesData = output.TemplatesData
	}
	if outputLanguage.Java != nil {
		outputLanguage.Java.InterpolateParameters(interpolator)
		outputLanguage.Java.ExtraFilesTemplatesData = output.TemplatesData
	}
	if outputLanguage.Typescript != nil {
		outputLanguage.Typescript.InterpolateParameters(interpolator)
		outputLanguage.Typescript.ExtraFilesTemplatesData = output.TemplatesData
	}

	if outputLanguage.Terraform != nil {
		outputLanguage.Terraform.InterpolateParameters(interpolator)
	}
}

type OutputOptions struct {
	// ReplaceExtension updates file extensions to the new one
	ReplaceExtension map[string]string `yaml:"replace_extension"`
}

type OutputLanguagePlugin map[string]any

func (outputLanguagePlugin OutputLanguagePlugin) interpolateParameters(interpolator ParametersInterpolator) {
	interpolateAny := func(input any) any {
		stringVal, ok := input.(string)
		if ok {
			return interpolator(stringVal)
		}

		return input
	}

	for key, val := range outputLanguagePlugin {
		outputLanguagePlugin[key] = interpolateAny(val)
	}
}
