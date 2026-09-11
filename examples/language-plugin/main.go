package main

import (
	"log/slog"
	"path/filepath"

	"github.com/davecgh/go-spew/spew"
	"github.com/grafana/codejen"
	"github.com/grafana/cog/pkg/ir"
	"github.com/grafana/cog/pkg/languages"
	"github.com/grafana/cog/pkg/plugins"
)

var _ plugins.Language = (*DummyLanguage)(nil)

type DummyLanguage struct {
	logger *slog.Logger
}

func (g *DummyLanguage) ValidateConfig(config map[string]any) error {
	g.logger.Debug("message from DummyLanguage.ValidateConfig", "config", spew.Sprint(config))

	return nil
}

func (g *DummyLanguage) NullableConfig(config map[string]any) (languages.NullableConfig, error) {
	g.logger.Debug("message from DummyLanguage.NullableConfig", "config", spew.Sprint(config))

	return languages.NullableConfig{
		Kinds:         []ir.Kind{ir.KindArray},
		AnyIsNullable: true,
	}, nil
}

func (g *DummyLanguage) TransformSchemas(config map[string]any, schemas ir.Schemas) (ir.Schemas, error) {
	g.logger.Debug("message from DummyLanguage.TransformSchemas", "config", spew.Sprint(config))

	return schemas, nil
}

func (g *DummyLanguage) Generate(codegenConfig languages.Config, config map[string]any, context languages.Context) (codejen.Files, error) {
	g.logger.Debug("message from DummyLanguage.Generate")

	var files codejen.Files

	for _, schema := range context.Schemas {
		filename := filepath.Join(
			schema.Package,
			"types_gen.rs",
		)

		files = append(files, *codejen.NewFile(filename, []byte("lala")))
	}

	return files, nil
}

func main() {
	lang := &DummyLanguage{
		logger: plugins.DefaultLogger,
	}

	plugins.Run(plugins.LanguagePlugin(lang))
}
