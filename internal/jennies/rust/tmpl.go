package rust

import (
	"embed"
	"fmt"

	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/jennies/template"
	"github.com/grafana/cog/internal/languages"
)

// templatesFS embeds the templates directory. The directory currently holds
// only a .gitkeep placeholder. Embedding the directory itself (rather than a
// glob over .tmpl files) keeps this compiling until real templates land.
//
//go:embed all:templates
//nolint:gochecknoglobals
var templatesFS embed.FS

func initTemplates(config Config, apiRefCollector *common.APIReferenceCollector) *template.Template {
	tmpl, err := template.New(
		"rust",

		template.Funcs(common.TypeResolvingTemplateHelpers(languages.Context{})),
		template.Funcs(common.TypesTemplateHelpers(languages.Context{})),
		template.Funcs(common.APIRefTemplateHelpers(apiRefCollector)),

		template.ParseFS(templatesFS, "templates"),
	)
	if err != nil {
		panic(fmt.Errorf("could not initialize templates: %w", err))
	}

	return tmpl
}
