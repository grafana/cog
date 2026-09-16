package config

import (
	"github.com/spf13/cobra"
)

type options struct {
	ConfigPath      string
	ExtraParameters map[string]string
}

func Command() *cobra.Command {
	configOpts := &options{}

	cmd := &cobra.Command{
		Use:   "config",
		Short: "View or manipulate the configuration",
	}

	cmd.PersistentFlags().StringToStringVar(&configOpts.ExtraParameters, "parameters", nil, "Sets or overrides parameters used in the config file.")

	cmd.PersistentFlags().StringVar(&configOpts.ConfigPath, "config", "", "Codegen pipeline configuration file.")
	_ = cmd.MarkFlagFilename("config")
	_ = cmd.MarkFlagRequired("config")

	cmd.AddCommand(updateInputs(configOpts))
	cmd.AddCommand(view(configOpts))

	return cmd
}
