package main

import (
	"os"

	"github.com/grafana/cog/cmd/cli/generate"
	"github.com/grafana/cog/cmd/cli/inspect"
	"github.com/spf13/cobra"
)

var version = "SNAPSHOT"

func main() {
	rootCmd := &cobra.Command{
		Use:          "cog <command>",
		Short:        "A tool for working with Grafana objects from code",
		SilenceUsage: true,
		Version:      version,
	}

	rootCmd.AddCommand(generate.Command())
	rootCmd.AddCommand(inspect.Command())

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
