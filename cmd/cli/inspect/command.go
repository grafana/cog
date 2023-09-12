package inspect

import (
	"encoding/json"
	"fmt"

	"github.com/grafana/cog/cmd/cli/loaders"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	opts := loaders.Options{}

	// TODO:
	// 	- support inspecting our different IRs: types, builders
	//  - support inspecting "transformed" IRs: language-specific compiler passes applied
	// 	  on types IR, veneers applied on the builders IR
	cmd := &cobra.Command{
		Use:   "inspect",
		Short: "Inspects the intermediate representation.", // TODO: better descriptions
		Long:  `Inspects the intermediate representation.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return doGenerate(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.SchemasType, "loader", "l", "cue", "Schemas type.") // TODO: better usage text
	cmd.Flags().StringArrayVarP(&opts.Entrypoints, "input", "i", nil, "Schema.")     // TODO: better usage text
	cmd.Flags().StringArrayVarP(&opts.CueImports, "include-cue-import", "I", nil, "Specify an additional library import directory. Format: [path]:[import]. Example: '../grafana/common-library:github.com/grafana/grafana/packages/grafana-schema/src/common")

	_ = cmd.MarkFlagRequired("input")
	_ = cmd.MarkFlagDirname("input")

	return cmd
}

func doGenerate(opts loaders.Options) error {
	loader, err := loaders.ForSchemaType(opts.SchemasType)
	if err != nil {
		return err
	}

	schemas, err := loader(opts)
	if err != nil {
		return err
	}

	marshaled, err := json.MarshalIndent(schemas, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(marshaled))

	return nil
}
