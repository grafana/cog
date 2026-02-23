package generate

import (
	"context"
	"log/slog"
	"os"

	"github.com/grafana/cog/internal/codegen"
	"github.com/grafana/cog/internal/logs"
	"github.com/spf13/cobra"
)

type options struct {
	Debug           bool
	ConfigPath      string
	ExtraParameters map[string]string
}

func Command() *cobra.Command {
	opts := options{}
	verbosity := 0
	var logger *slog.Logger

	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generates code from schemas.", // TODO: better descriptions
		Long:  `Generates code from schemas.`,
		PreRun: func(cmd *cobra.Command, _ []string) {
			logLevel := new(slog.LevelVar)
			logLevel.Set(slog.LevelInfo)
			// Multiplying the number of occurrences of the `-v` flag by 4 (gap between log levels in slog)
			// allows us to increase the logger's verbosity.
			logLevel.Set(logLevel.Level() - slog.Level(min(verbosity, 3)*4))

			logHandler := logs.NewHandler(os.Stderr, &logs.Options{
				Level: logLevel,
			})
			logger = slog.New(logHandler)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return doGenerate(cmd.Context(), logger, opts)
		},
	}

	cmd.Flags().BoolVar(&opts.Debug, "debug", false, "Debug mode.")
	cmd.Flags().StringVar(&opts.ConfigPath, "config", "", "Codegen pipeline configuration file.")
	_ = cmd.MarkFlagFilename("config")
	_ = cmd.MarkFlagRequired("config")

	cmd.Flags().StringToStringVar(&opts.ExtraParameters, "parameters", nil, "Sets or overrides parameters used in the config file.")
	cmd.Flags().CountVarP(&verbosity, "verbose", "v", "Verbose mode. Multiple -v options increase the verbosity (maximum: 3).")

	return cmd
}

func doGenerate(ctx context.Context, logger *slog.Logger, opts options) error {
	pipelineOpts := []codegen.PipelineOption{
		codegen.Parameters(opts.ExtraParameters),
		codegen.Logger(logger),
		codegen.Debug(opts.Debug),
	}

	pipeline, err := codegen.PipelineFromFile(opts.ConfigPath, pipelineOpts...)
	if err != nil {
		return err
	}

	generatedFS, err := pipeline.Run(ctx)
	if err != nil {
		return err
	}

	return generatedFS.Write(ctx, "")
}
