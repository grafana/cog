package remote

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/languages"
	"github.com/grafana/cog/internal/tools"
	cogplugin "github.com/grafana/cog/pkg/plugin"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
)

type remote struct {
	config map[string]any
}

func (jenny remote) JennyName() string {
	return "Remote"
}

func (jenny remote) Generate(context languages.Context) (codejen.Files, error) {
	// TODO: get a logger from upstream, and adapt it into an hcl one
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "plugin",
		Output: os.Stdout,
		Level:  hclog.Debug,
	})

	// We're a host! Start by launching the plugin process.
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: cogplugin.LanguagePluginHandshakeConfig,
		Plugins: map[string]plugin.Plugin{
			"language": &cogplugin.LanguagePlugin{},
		},
		Cmd:    exec.Command("./plugin-sandbox/remote-lang"), // TODO: this should be injected into the constructor
		Logger: logger,
	})
	defer client.Kill()

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

	// We should have a Greeter now! This feels like a normal interface
	// implementation but is in fact over an RPC connection.
	remoteLang := raw.(cogplugin.Language)

	if err := remoteLang.ValidateConfig(jenny.config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	files, err := remoteLang.Generate()
	if err != nil {
		return nil, err
	}

	return tools.Map(files, func(file codejen.File) codejen.File {
		file.From = append(file.From, jenny)
		return file
	}), nil
}
