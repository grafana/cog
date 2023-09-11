package generate

import (
	"context"
	"fmt"
	"strings"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies"
	"github.com/spf13/cobra"
)

type options struct {
	outputDir   string
	entrypoints []string
	schemasType string

	// Cue-specific options
	cueImports []string
}

func (opts options) cueIncludeImports() ([]cueIncludeImport, error) {
	if len(opts.cueImports) == 0 {
		return nil, nil
	}

	imports := make([]cueIncludeImport, len(opts.cueImports))
	for i, importDefinition := range opts.cueImports {
		parts := strings.Split(importDefinition, ":")
		if len(parts) != 2 {
			return nil, fmt.Errorf("'%s' is not a valid import definition", importDefinition)
		}

		imports[i].fsPath = parts[0]
		imports[i].importPath = parts[1]
	}

	return imports, nil
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

	cmd.Flags().StringVarP(&opts.schemasType, "loader", "l", "cue", "Schemas type.")         // TODO: better usage text
	cmd.Flags().StringVarP(&opts.outputDir, "output", "o", "generated", "Output directory.") // TODO: better usage text
	cmd.Flags().StringArrayVarP(&opts.entrypoints, "input", "i", nil, "Schema.")             // TODO: better usage text
	cmd.Flags().StringArrayVarP(&opts.cueImports, "include-cue-import", "I", nil, "Specify an additional library import directory. Format: [path]:[import]. Example: '../grafana/common-library:github.com/grafana/grafana/packages/grafana-schema/src/common")

	_ = cmd.MarkFlagRequired("input")
	_ = cmd.MarkFlagDirname("input")
	_ = cmd.MarkFlagDirname("output")

	return cmd
}

func doGenerate(opts options) error {
	loaders := map[string]func(opts options) ([]*ast.File, error){
		"cue":            cueLoader,
		"kindsys-core":   kindsysCoreLoader,
		"kindsys-custom": kindsysCustomLoader,
		"jsonschema":     jsonschemaLoader,
	}
	loader, ok := loaders[opts.schemasType]
	if !ok {
		return fmt.Errorf("no loader found for '%s'", opts.schemasType)
	}

	schemas, err := loader(opts)
	if err != nil {
		return err
	}

	// Here begins the code generation setup
	targetsByLanguage := jennies.All()
	rootCodeJenFS := codejen.NewFS()

	for language, target := range targetsByLanguage {
		fmt.Printf("Running '%s' jennies...\n", language)

		var err error
		processedAsts := schemas

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

	err = rootCodeJenFS.Write(context.Background(), opts.outputDir)
	if err != nil {
		return err
	}

	return nil
}
