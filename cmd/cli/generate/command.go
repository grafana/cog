package generate

import (
	"context"
	"fmt"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/cmd/cli/loaders"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/ast/compiler"
	"github.com/grafana/cog/internal/jennies"
	"github.com/spf13/cobra"
)

type Options struct {
	loaders.Options

	OutputDir string
}

func Command() *cobra.Command {
	opts := Options{}

	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generates code from schemas.", // TODO: better descriptions
		Long:  `Generates code from schemas.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return doGenerate(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.OutputDir, "output", "o", "generated", "Output directory.") // TODO: better usage text

	cmd.Flags().StringArrayVar(&opts.CueEntrypoints, "cue", nil, "CUE input schema.")                                                  // TODO: better usage text
	cmd.Flags().StringArrayVar(&opts.KindsysCoreEntrypoints, "kindsys-core", nil, "Kindys core kinds input schema.")                   // TODO: better usage text
	cmd.Flags().StringArrayVar(&opts.KindsysComposableEntrypoints, "kindsys-composable", nil, "Kindys composable kinds input schema.") // TODO: better usage text
	cmd.Flags().StringArrayVar(&opts.KindsysCustomEntrypoints, "kindsys-custom", nil, "Kindys custom kinds input schema.")             // TODO: better usage text
	cmd.Flags().StringArrayVar(&opts.JSONSchemaEntrypoints, "jsonschema", nil, "Jsonschema input schema.")                             // TODO: better usage text
	cmd.Flags().StringArrayVar(&opts.OpenAPIEntrypoints, "openapi", nil, "Openapi input schema.")                                      // TODO: better usage text
	cmd.Flags().StringVar(&opts.KindRegistryPath, "kind-registry", "", "Kind registry input.")                                         // TODO: better usage text

	cmd.Flags().StringArrayVarP(&opts.CueImports, "include-cue-import", "I", nil, "Specify an additional library import directory. Format: [path]:[import]. Example: '../grafana/common-library:github.com/grafana/grafana/packages/grafana-schema/src/common")
	cmd.Flags().StringVar(&opts.KindRegistryVersion, "kind-registry-version", "next", "Schemas version")

	_ = cmd.MarkFlagDirname("cue")
	_ = cmd.MarkFlagDirname("kindsys-core")
	_ = cmd.MarkFlagDirname("kindsys-custom")
	_ = cmd.MarkFlagDirname("kind-registry")
	_ = cmd.MarkFlagFilename("jsonschema")
	_ = cmd.MarkFlagDirname("openapi")
	_ = cmd.MarkFlagDirname("output")

	return cmd
}

func doGenerate(opts Options) error {
	schemas, err := loaders.LoadAll(opts.Options)
	if err != nil {
		return err
	}

	// Here begins the code generation setup
	targetsByLanguage := jennies.All()
	rootCodeJenFS := codejen.NewFS()

	for language, target := range targetsByLanguage {
		fmt.Printf("Running '%s' jennies...\n", language)

		var err error
		processedAsts := ast.Schemas(schemas).DeepCopy()

		var compilerPasses []compiler.Pass
		compilerPasses = append(compilerPasses, compiler.CommonPasses()...)
		compilerPasses = append(compilerPasses, target.CompilerPasses...)

		for _, compilerPass := range compilerPasses {
			processedAsts, err = compilerPass.Process(processedAsts)
			if err != nil {
				return err
			}
		}

		fs, err := target.Jennies.GenerateFS(processedAsts)
		if err != nil {
			return err
		}

		err = rootCodeJenFS.Merge(fs)
		if err != nil {
			return err
		}
	}

	err = rootCodeJenFS.Write(context.Background(), opts.OutputDir)
	if err != nil {
		return err
	}

	return nil
}
