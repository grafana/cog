package generate

import (
	"context"
	"fmt"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/cmd/cli/loaders"
	"github.com/grafana/cog/internal/ast"
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

	cmd.Flags().StringVarP(&opts.SchemasType, "loader", "l", "cue", "Schemas type.")         // TODO: better usage text
	cmd.Flags().StringVarP(&opts.OutputDir, "output", "o", "generated", "Output directory.") // TODO: better usage text
	cmd.Flags().StringArrayVarP(&opts.Entrypoints, "input", "i", nil, "Schema.")             // TODO: better usage text
	cmd.Flags().StringArrayVarP(&opts.CueImports, "include-cue-import", "I", nil, "Specify an additional library import directory. Format: [path]:[import]. Example: '../grafana/common-library:github.com/grafana/grafana/packages/grafana-schema/src/common")

	_ = cmd.MarkFlagRequired("input")
	_ = cmd.MarkFlagDirname("input")
	_ = cmd.MarkFlagDirname("output")

	return cmd
}

func doGenerate(opts Options) error {
	loader, err := loaders.ForSchemaType(opts.SchemasType)
	if err != nil {
		return err
	}

	schemas, err := loader(opts.Options)
	if err != nil {
		return err
	}

	// Here begins the code generation setup
	targetsByLanguage := jennies.All()
	rootCodeJenFS := codejen.NewFS()

	for language, target := range targetsByLanguage {
		fmt.Printf("Running '%s' jennies...\n", language)

		var err error
		processedAsts := ast.Files(schemas).DeepCopy()

		for _, compilerPass := range target.CompilerPasses {
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
