package generate

import (
	"context"

	"github.com/grafana/cog/internal/codegen"
	"github.com/spf13/cobra"
)

type options struct {
	ConfigPath      string
	ExtraParameters map[string]string
}

func Command() *cobra.Command {
	opts := options{}

	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generates code from schemas.", // TODO: better descriptions
		Long:  `Generates code from schemas.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return doGenerate(opts)
		},
	}

	cmd.Flags().StringVar(&opts.ConfigPath, "config", "", "Codegen pipeline configuration file.")
	_ = cmd.MarkFlagFilename("config")
	_ = cmd.MarkFlagRequired("config")

	cmd.Flags().StringToStringVar(&opts.ExtraParameters, "parameters", nil, "Sets or overrides parameters used in the config file.")

	return cmd
}

func doGenerate(opts options) error {
	ctx := context.Background()

	pipeline, err := codegen.PipelineFromFile(opts.ConfigPath, codegen.Parameters(opts.ExtraParameters), codegen.Reporter(codegen.StdoutReporter))
	if err != nil {
		return err
	}

	generatedFS, err := pipeline.Run(ctx)
	if err != nil {
		return err
	}

	return generatedFS.Write(context.Background(), "")
}
