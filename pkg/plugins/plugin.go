package plugins

import (
	"log/slog"
	"os"

	"github.com/grafana/cog/pkg/logs"
	"github.com/hashicorp/go-plugin"
)

//nolint:gochecknoglobals
var DefaultLogger = slog.New(logs.NewHandler(os.Stderr, &logs.Options{
	Level: slog.LevelDebug,
}))

type Option func(runner *pluginRunner)

func LanguagePlugin(implementation Language) Option {
	return func(runner *pluginRunner) {
		runner.languagePlugin = &LanguagePluginRunner{Impl: implementation}
	}
}

func Logger(logger *slog.Logger) Option {
	return func(runner *pluginRunner) {
		runner.logger = logger
	}
}

type pluginRunner struct {
	logger          *slog.Logger
	handshakeConfig plugin.HandshakeConfig
	languagePlugin  plugin.Plugin
}

func (runner *pluginRunner) run() {
	pluginsMap := map[string]plugin.Plugin{}
	if runner.languagePlugin != nil {
		pluginsMap["language"] = runner.languagePlugin
	}

	plugin.Serve(&plugin.ServeConfig{
		Logger:          logs.HCLLoggerFromSlog(runner.logger),
		HandshakeConfig: runner.handshakeConfig,
		Plugins:         pluginsMap,
	})
}

func Run(opts ...Option) {
	runner := &pluginRunner{
		logger:          DefaultLogger,
		handshakeConfig: LanguagePluginHandshakeConfig,
	}

	for _, opt := range opts {
		opt(runner)
	}

	runner.logger.Debug("starting plugin runner")
	runner.run()
}
