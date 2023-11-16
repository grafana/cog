package generate

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/cmd/cli/loaders"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/ast/compiler"
	"github.com/grafana/cog/internal/jennies"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/veneers/yaml"
	"github.com/spf13/cobra"
)

type Options struct {
	loaders.Options

	Targets                 common.Targets
	Languages               []string
	VeneerConfigFiles       []string
	VeneerConfigDirectories []string
	OutputDir               string
}

func (opts Options) VeneerFiles() ([]string, error) {
	veneers := make([]string, 0, len(opts.VeneerConfigFiles))
	veneers = append(veneers, opts.VeneerConfigFiles...)

	for _, dir := range opts.VeneerConfigDirectories {
		globPattern := filepath.Join(filepath.Clean(dir), "*.yaml")
		matches, err := filepath.Glob(globPattern)
		if err != nil {
			return nil, err
		}

		veneers = append(veneers, matches...)
	}

	return veneers, nil
}

func Command() *cobra.Command {
	opts := Options{}
	languageJennies := jennies.All()

	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generates code from schemas.", // TODO: better descriptions
		Long:  `Generates code from schemas.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return doGenerate(languageJennies, opts)
		},
	}

	cmd.Flags().BoolVar(&opts.Targets.Types, "generate-types", true, "Generate types.")          // TODO: better usage text
	cmd.Flags().BoolVar(&opts.Targets.Builders, "generate-builders", true, "Generate builders.") // TODO: better usage text

	cmd.Flags().StringVarP(&opts.OutputDir, "output", "o", "generated", "Output directory.") // TODO: better usage text
	cmd.Flags().StringArrayVarP(&opts.Languages, "language", "l", nil, "Language to generate. If left empty, all supported languages will be generated.")
	cmd.Flags().StringArrayVarP(&opts.VeneerConfigFiles, "veneer", "c", nil, "Veneer configuration file.")
	cmd.Flags().StringArrayVar(&opts.VeneerConfigDirectories, "veneers", nil, "Veneer configuration directories.")

	cmd.Flags().StringArrayVar(&opts.CueEntrypoints, "cue", nil, "CUE input schema.")                                                  // TODO: better usage text
	cmd.Flags().StringArrayVar(&opts.KindsysCoreEntrypoints, "kindsys-core", nil, "Kindys core kinds input schema.")                   // TODO: better usage text
	cmd.Flags().StringArrayVar(&opts.KindsysComposableEntrypoints, "kindsys-composable", nil, "Kindys composable kinds input schema.") // TODO: better usage text
	cmd.Flags().StringArrayVar(&opts.KindsysCustomEntrypoints, "kindsys-custom", nil, "Kindys custom kinds input schema.")             // TODO: better usage text
	cmd.Flags().StringArrayVar(&opts.JSONSchemaEntrypoints, "jsonschema", nil, "Jsonschema input schema.")                             // TODO: better usage text
	cmd.Flags().StringArrayVar(&opts.OpenAPIEntrypoints, "openapi", nil, "Openapi input schema.")                                      // TODO: better usage text
	cmd.Flags().StringVar(&opts.KindRegistryPath, "kind-registry", "", "Kind registry input.")                                         // TODO: better usage text

	cmd.Flags().StringArrayVarP(&opts.CueImports, "include-cue-import", "I", nil, "Specify an additional library import directory. Format: [path]:[import]. Example: '../grafana/common-library:github.com/grafana/grafana/packages/grafana-schema/src/common")
	cmd.Flags().StringVar(&opts.KindRegistryVersion, "kind-registry-version", "next", "Schemas version")

	for _, jenny := range languageJennies {
		jenny.RegisterCliFlags(cmd)
	}

	_ = cmd.MarkFlagDirname("cue")
	_ = cmd.MarkFlagDirname("kindsys-core")
	_ = cmd.MarkFlagDirname("kindsys-custom")
	_ = cmd.MarkFlagDirname("kind-registry")
	_ = cmd.MarkFlagFilename("jsonschema")
	_ = cmd.MarkFlagDirname("openapi")
	_ = cmd.MarkFlagDirname("output")
	_ = cmd.MarkFlagFilename("veneer")
	_ = cmd.MarkFlagDirname("veneers")

	return cmd
}

func doGenerate(allTargets jennies.LanguageJennies, opts Options) error {
	veneerFiles, err := opts.VeneerFiles()
	if err != nil {
		return err
	}

	veneers, err := yaml.NewLoader().RewriterFromFiles(veneerFiles)
	if err != nil {
		return err
	}

	// Here begins the code generation setup
	targetsByLanguage, err := allTargets.ForLanguages(opts.Languages)
	if err != nil {
		return err
	}

	schemas, err := loaders.LoadAll(opts.Options)
	if err != nil {
		return err
	}

	rootCodeJenFS := codejen.NewFS()

	for language, target := range targetsByLanguage {
		fmt.Printf("Running '%s' jennies...\n", language)

		compilerPasses := compiler.CommonPasses().Concat(target.CompilerPasses())
		processedSchemas, err := compilerPasses.Process(schemas)
		if err != nil {
			return err
		}

		// from these types, create builders
		builderGenerator := &ast.BuilderGenerator{}
		builders := builderGenerator.FromAST(processedSchemas)

		// apply the builder veneers
		builders, err = veneers.ApplyTo(builders, language)
		if err != nil {
			return err
		}

		// then delegate the actual codegen to the jennies
		fs, err := target.Jennies(opts.Targets).GenerateFS(common.Context{
			Schemas:  processedSchemas,
			Builders: builders,
		})
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
