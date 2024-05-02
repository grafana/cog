package inspect

import (
	"encoding/json"
	"fmt"

	"github.com/grafana/cog/cmd/cli/loaders"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/spf13/cobra"
)

type options struct {
	BuilderIR  bool
	ConfigPath string
}

func Command() *cobra.Command {
	opts := options{}

	// TODO:
	// 	- support inspecting our different IRs: types, builders
	//  - support inspecting "transformed" IRs: language-specific compiler passes applied
	// 	  on types IR, veneers applied on the builders IR
	cmd := &cobra.Command{
		Use:   "inspect",
		Short: "Inspects the intermediate representation.", // TODO: better descriptions
		Long:  `Inspects the intermediate representation.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return doInspect(opts)
		},
	}

	cmd.Flags().BoolVar(&opts.BuilderIR, "builder-ir", false, "Inspect the \"builder IR\" instead of the \"types\" one.") // TODO: better usage text

	cmd.Flags().StringVar(&opts.ConfigPath, "config", "", "Configuration file.")
	_ = cmd.MarkFlagFilename("config")
	_ = cmd.MarkFlagRequired("config")

	return cmd
}

func doInspect(opts options) error {
	config, err := loaders.ConfigFromFile(opts.ConfigPath)
	if err != nil {
		return err
	}

	schemas, err := config.LoadSchemas()
	if err != nil {
		return err
	}

	if opts.BuilderIR {
		return inspectBuilderIR(schemas)
	}

	return prettyPrintJSON(schemas)
}

func inspectBuilderIR(schemas []*ast.Schema) error {
	generator := &ast.BuilderGenerator{}

	return prettyPrintJSON(common.Context{
		Schemas:  schemas,
		Builders: generator.FromAST(schemas),
	})
}

func prettyPrintJSON(input any) error {
	marshaled, err := json.MarshalIndent(input, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(marshaled))

	return nil
}
