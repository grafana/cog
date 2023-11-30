package template

import (
	"github.com/grafana/cog/internal/jennies/common"
	"text/template"
)

type Config struct {
	Debug bool

	// GenerateGoMod indicates whether a go.mod file should be generated.
	// If enabled, PackageRoot is used as module path.
	GenerateGoMod bool

	// Root path for imports.
	// Ex: github.com/grafana/cog/generated
	PackageRoot string

	FileExtension string

	// Configuration for template for a specific language
	TemplateConfig FunctionsConfig

	ImportMapper    *common.DirectImportMap
	BuilderTemplate BuilderFormatter
	Formatter       TypeFormatter
}

type FunctionsConfig struct {
	Name string

	FuncMap template.FuncMap
}
