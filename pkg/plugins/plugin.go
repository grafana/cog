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

// Option represents a configuration setting for a language plugin runner.
type Option func(runner *pluginRunner)

// LanguagePlugin sets the implementation to run.
func LanguagePlugin(language Language) Option {
	return func(runner *pluginRunner) {
		runner.languagePlugin = &LanguagePluginRunner{Impl: language}
	}
}

// Logger sets the logger used by the plugin.
func Logger(logger *slog.Logger) Option {
	return func(runner *pluginRunner) {
		runner.logger = logger
	}
}

// Run starts and runs a language plugin.
func Run(opts ...Option) {
	runner := &pluginRunner{
		logger:          DefaultLogger,
		handshakeConfig: LanguagePluginHandshakeConfig,
	}

	for _, opt := range opts {
		opt(runner)
	}

	runner.logger.Debug("starting language plugin runner")
	runner.run()
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
