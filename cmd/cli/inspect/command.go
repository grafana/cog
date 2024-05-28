package inspect

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/grafana/cog/cmd/cli/codegen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/languages"
	"github.com/spf13/cobra"
)

type options struct {
	BuilderIR       bool
	ConfigPath      string
	ExtraParameters map[string]string
}

func Command() *cobra.Command {
	opts := options{}

	// TODO:
	//  - support inspecting "transformed" IRs: common & language-specific compiler passes, veneers, ...
	cmd := &cobra.Command{
		Use:   "inspect",
		Short: "Inspects the intermediate representation.", // TODO: better descriptions
		Long:  `Inspects the intermediate representation.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return doInspect(opts)
		},
	}

	cmd.Flags().BoolVar(&opts.BuilderIR, "builder-ir", false, "Inspect the \"builder IR\" instead of the \"types\" one.") // TODO: better usage text
	cmd.Flags().StringToStringVar(&opts.ExtraParameters, "parameters", nil, "Sets or overrides parameters used in the config file.")

	cmd.Flags().StringVar(&opts.ConfigPath, "config", "", "Codegen pipeline configuration file.")
	_ = cmd.MarkFlagFilename("config")
	_ = cmd.MarkFlagRequired("config")

	return cmd
}

func doInspect(opts options) error {
	ctx := context.Background()

	pipeline, err := codegen.PipelineFromFile(opts.ConfigPath, codegen.Parameters(opts.ExtraParameters))
	if err != nil {
		return err
	}

	schemas, err := pipeline.LoadSchemas(ctx)
	if err != nil {
		return err
	}

	if opts.BuilderIR {
		return inspectBuilderIR(schemas)
	}

	return prettyPrintJSON(schemas)
}

func inspectBuilderIR(schemas []*ast.Schema) error {
	var err error
	generator := &ast.BuilderGenerator{}

	codegenCtx := common.Context{
		Schemas:  schemas,
		Builders: generator.FromAST(schemas),
	}

	codegenCtx, err = languages.GenerateBuilderNilChecks(nil, codegenCtx)
	if err != nil {
		return err
	}

	return prettyPrintJSON(codegenCtx)
}

func prettyPrintJSON(input any) error {
	marshaled, err := json.MarshalIndent(input, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(marshaled))

	return nil
}
