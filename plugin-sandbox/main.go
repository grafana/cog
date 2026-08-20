// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/grafana/codejen"
	cogplugin "github.com/grafana/cog/pkg/plugin"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
)

// Here is a real implementation of Greeter
type DummyLanguage struct {
	logger hclog.Logger
}

func (g *DummyLanguage) ValidateConfig(config map[string]any) error {
	g.logger.Debug("message from DummyLanguage.ValidateConfig", spew.Sprint(config))

	return nil
}

func (g *DummyLanguage) Generate() (codejen.Files, error) {
	g.logger.Debug("message from DummyLanguage.Generate")

	var files codejen.Files

	files = append(files, *codejen.NewFile("foo/bar/baz.rs", []byte("lala")))

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
