package main

import (
	"fmt"
	"os"

	"github.com/grafana/grafana-foundation-sdk/go/cmd/cli/convert"
	"github.com/grafana/grafana-foundation-sdk/go/cmd/cli/validate"
	"github.com/grafana/grafana-foundation-sdk/go/cog/plugins"
	"github.com/urfave/cli/v2"
)

func main() {
	plugins.RegisterDefaultPlugins()

	app := &cli.App{
		Name:    "grafana-foundation-sdk",
		Version: "{{ .Extra.ReleaseBranch }}",
		Commands: cli.Commands{
			convert.Command(),
			validate.Command(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
	}
}
