package main

import (
	"os"

	"github.com/grafana/cog/cmd/cli/generate"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:          "cog <command>",
		Short:        "A tool for working with Grafana objects from code",
		SilenceUsage: true,
	}

	rootCmd.AddCommand(generate.Command())

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
