package remote

import (
	"fmt"
	"log/slog"
	"os/exec"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/pkg/ir"
	"github.com/grafana/cog/pkg/languages"
	"github.com/grafana/cog/pkg/logs"
	"github.com/grafana/cog/pkg/plugins"
	"github.com/hashicorp/go-plugin"
)

type Language struct {
	logger         *slog.Logger
	name           string
	config         map[string]any
	pluginClient   *plugin.Client
	plugin         plugins.Language
	nullableConfig languages.NullableConfig
}

func New(logger *slog.Logger, name string, config map[string]any) (*Language, error) {
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: plugins.LanguagePluginHandshakeConfig,
		Plugins: map[string]plugin.Plugin{
			"language": &plugins.LanguagePluginRunner{},
		},
		// this will look for "cog-{name}" in the PATH and run it
		Cmd:         exec.Command("cog-" + name), //nolint
		Logger:      logs.HCLLoggerFromSlog(logger),
		SkipHostEnv: true,
	})

	// Connect via RPC
	rpcClient, err := client.Client()
	if err != nil {
		return nil, fmt.Errorf("could not create plugin client: %w", err)
	}

	// Request the plugin
	raw, err := rpcClient.Dispense("language")
	if err != nil {
		return nil, fmt.Errorf("could not dispense language plugin: %w", err)
	}

	languagePlugin := raw.(plugins.Language)

	if err := languagePlugin.ValidateConfig(config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	nullableConfig, err := languagePlugin.NullableConfig(config)
	if err != nil {
		return nil, fmt.Errorf("could not load 'nullable configuration': %w", err)
	}

	return &Language{
		logger:         logger,
		name:           name,
		config:         config,
		pluginClient:   client,
		plugin:         languagePlugin,
		nullableConfig: nullableConfig,
	}, nil
}

func (language *Language) Terminate() {
	language.pluginClient.Kill()
}

func (language *Language) Name() string {
	return language.name
}

func (language *Language) Jennies(globalConfig languages.Config) *codejen.JennyList[languages.Context] {
	jenny := codejen.JennyListWithNamer(func(_ languages.Context) string {
		return language.Name()
	})

	jenny.AppendOneToMany(remote{
		globalConfig: globalConfig,
		config:       language.config,
		plugin:       language.plugin,
	})

	return jenny
}

func (language *Language) Transform(schemas ir.Schemas) (ir.Schemas, error) {
	return language.plugin.TransformSchemas(language.config, schemas)
}

func (language *Language) NullableKinds() languages.NullableConfig {
	return language.nullableConfig
}
