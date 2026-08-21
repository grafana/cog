// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"os"
	"path/filepath"

	"github.com/davecgh/go-spew/spew"
	"github.com/grafana/codejen"
	"github.com/grafana/cog/pkg/languages"
	cogplugin "github.com/grafana/cog/pkg/plugin"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
)

var _ cogplugin.Language = (*DummyLanguage)(nil)

// Here is a real implementation of Greeter
type DummyLanguage struct {
	logger hclog.Logger
}

func (g *DummyLanguage) ValidateConfig(config map[string]any) error {
	g.logger.Debug("message from DummyLanguage.ValidateConfig", spew.Sprint(config))

	return nil
}

func (g *DummyLanguage) Transform(codegenConfig languages.Config, config map[string]any, context languages.Context) (languages.Context, error) {
	g.logger.Debug("message from DummyLanguage.Transform")

	return context, nil
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
	logger := hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Trace,
		Output:     os.Stderr,
		JSONFormat: true,
	})

	lang := &DummyLanguage{
		logger: logger,
	}

	// pluginMap is the map of plugins we can dispense.
	logger.Debug("message from plugin", "foo", "bar")

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: cogplugin.LanguagePluginHandshakeConfig,
		Plugins: map[string]plugin.Plugin{
			"language": &cogplugin.LanguagePlugin{Impl: lang},
		},
	})
}
