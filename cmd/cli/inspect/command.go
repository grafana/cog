package inspect

import (
	"encoding/json"
	"fmt"

	"github.com/grafana/cog/cmd/cli/loaders"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/spf13/cobra"
)

type inspectOptions struct {
	LoaderOptions loaders.Options
	BuilderIR     bool
}

func Command() *cobra.Command {
	opts := inspectOptions{}

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

	cmd.Flags().StringArrayVar(&opts.LoaderOptions.CueEntrypoints, "cue", nil, "CUE input schema.")                                                  // TODO: better usage text
	cmd.Flags().StringArrayVar(&opts.LoaderOptions.KindsysCoreEntrypoints, "kindsys-core", nil, "Kindys core kinds input schema.")                   // TODO: better usage text
	cmd.Flags().StringArrayVar(&opts.LoaderOptions.KindsysComposableEntrypoints, "kindsys-composable", nil, "Kindys composable kinds input schema.") // TODO: better usage text
	cmd.Flags().StringArrayVar(&opts.LoaderOptions.KindsysCustomEntrypoints, "kindsys-custom", nil, "Kindys custom kinds input schema.")             // TODO: better usage text
	cmd.Flags().StringArrayVar(&opts.LoaderOptions.JSONSchemaEntrypoints, "jsonschema", nil, "Jsonschema input schema.")                             // TODO: better usage text
	cmd.Flags().StringArrayVar(&opts.LoaderOptions.OpenAPIEntrypoints, "openapi", nil, "Openapi input schema.")                                      // TODO: better usage text
	cmd.Flags().StringVar(&opts.LoaderOptions.KindRegistryPath, "kind-registry", "", "Kind registry input.")                                         // TODO: better usage text
	cmd.Flags().StringVar(&opts.LoaderOptions.JSONSchemaRegistryPath, "jsonschema-registry", "", "JSONschema registry input.")                       // TODO: better usage text

	cmd.Flags().StringArrayVarP(&opts.LoaderOptions.CueImports, "include-cue-import", "I", nil, "Specify an additional library import directory. Format: [path]:[import]. Example: '../grafana/common-library:github.com/grafana/grafana/packages/grafana-schema/src/common")
	cmd.Flags().StringVar(&opts.LoaderOptions.KindRegistryVersion, "kind-registry-version", "next", "Schemas version")

	_ = cmd.MarkFlagDirname("cue")
	_ = cmd.MarkFlagDirname("kindsys-core")
	_ = cmd.MarkFlagDirname("kindsys-custom")
	_ = cmd.MarkFlagDirname("kind-registry")
	_ = cmd.MarkFlagDirname("jsonschema-registry")
	_ = cmd.MarkFlagFilename("jsonschema")
	_ = cmd.MarkFlagDirname("openapi")

	return cmd
}

func doInspect(opts inspectOptions) error {
	schemas, err := loaders.LoadAll(opts.LoaderOptions)
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
