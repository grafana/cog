package terraform

import (
	"fmt"

	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/jennies/template"
	"github.com/grafana/cog/internal/languages"
)

func initTemplates(config Config) *template.Template {
	tmpl, err := template.New("terraform",
		template.Funcs(common.TypeResolvingTemplateHelpers(languages.Context{})),
		template.Funcs(common.TypesTemplateHelpers(languages.Context{})),
		template.ParseDirectories(config.OverridesTemplatesDirectories...),
	)

	if err != nil {
		panic(fmt.Errorf("could not initialize templates: %w", err))
	}
	return tmpl
}
