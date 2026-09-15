package config

import (
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/grafana/cog/internal/codegen"
	"github.com/spf13/cobra"
)

func view(configOpts *options) *cobra.Command {
	return &cobra.Command{
		Use:   "view",
		Short: "Inspects the configuration.",
		Long:  `Inspects the configuration.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return doView(configOpts)
		},
	}
}

func doView(opts *options) error {
	pipelineOpts := []codegen.PipelineOption{
		codegen.Parameters(opts.ExtraParameters),
	}

	pipeline, err := codegen.PipelineFromFile(opts.ConfigPath, pipelineOpts...)
	if err != nil {
		return err
	}

	if err := pipeline.LoadUnits(); err != nil {
		return err
	}

	payload, err := yaml.MarshalWithOptions(pipeline, yaml.Indent(2), yaml.IndentSequence(true))
	if err != nil {
		return err
	}

	fmt.Println(string(payload))

	return nil
}
