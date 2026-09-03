package remote

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"

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

func New(logger *slog.Logger, name string, pluginDirectories []string, config map[string]any) (*Language, error) {
	pluginCmd, err := pluginCommand(name, pluginDirectories)
	if err != nil {
		return nil, err
	}

	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: plugins.LanguagePluginHandshakeConfig,
		Plugins: map[string]plugin.Plugin{
			"language": &plugins.LanguagePluginRunner{},
		},
		Cmd:         pluginCmd,
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

func pluginCommand(name string, pluginDirectories []string) (*exec.Cmd, error) {
	binaryName := "cog-" + name

	if len(pluginDirectories) == 0 {
		// this will look for the binary in the PATH and run it
		return exec.Command(binaryName), nil //nolint
	}

	for _, dir := range pluginDirectories {
		file := filepath.Join(dir, binaryName)

		_, err := os.Stat(file)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return nil, err
		}

		return exec.Command(file), nil //nolint
	}

	return nil, fmt.Errorf("could not locate '%s' binary", binaryName)
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
