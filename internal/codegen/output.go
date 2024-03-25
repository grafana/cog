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

	Types      bool `yaml:"types"`
	Builders   bool `yaml:"builders"`
	Converters bool `yaml:"converters"`

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

func (output *Output) interpolateParameters(interpolator ParametersInterpolator) {
	output.Directory = interpolator(output.Directory)

	for _, outputLanguage := range output.Languages {
		outputLanguage.interpolateParameters(interpolator)
	}

	output.PackageTemplates = interpolator(output.PackageTemplates)
	output.RepositoryTemplates = interpolator(output.RepositoryTemplates)

	for key, value := range output.TemplatesData {
		output.TemplatesData[key] = interpolator(value)
	}
}

type OutputLanguage struct {
	Go         *golang.Config     `yaml:"go"`
	Java       *java.Config       `yaml:"java"`
	JSONSchema *jsonschema.Config `yaml:"jsonschema"`
	OpenAPI    *openapi.Config    `yaml:"openapi"`
	PHP        *php.Config        `yaml:"php"`
	Python     *python.Config     `yaml:"python"`
	Terraform  *terraform.Config  `yaml:"terraform"`
	Typescript *typescript.Config `yaml:"typescript"`
}

func (outputLanguage *OutputLanguage) interpolateParameters(interpolator ParametersInterpolator) {
	if outputLanguage.Go != nil {
		outputLanguage.Go.InterpolateParameters(interpolator)
	}
	if outputLanguage.PHP != nil {
		outputLanguage.PHP.InterpolateParameters(interpolator)
	}
	if outputLanguage.Python != nil {
		outputLanguage.Python.InterpolateParameters(interpolator)
	}
	if outputLanguage.Java != nil {
		outputLanguage.Java.InterpolateParameters(interpolator)
	}
	if outputLanguage.Terraform != nil {
		outputLanguage.Terraform.InterpolateParameters(interpolator)
	}
}
