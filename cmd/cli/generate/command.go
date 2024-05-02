package generate

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/cmd/cli/loaders"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/tools"
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

	cmd.Flags().StringVar(&opts.ConfigPath, "config", "", "Configuration file.")
	_ = cmd.MarkFlagFilename("config")
	_ = cmd.MarkFlagRequired("config")

	cmd.Flags().StringToStringVar(&opts.ExtraParameters, "parameters", nil, "Sets or overrides parameters used in the config file.")

	return cmd
}

func doGenerate(opts options) error {
	config, err := loaders.ConfigFromFile(opts.ConfigPath)
	if err != nil {
		return err
	}

	config = config.WithParameters(opts.ExtraParameters)
	config = config.InterpolateParameters()

	workingDir, err := os.Getwd()
	if err != nil {
		return err
	}

	veneers, err := config.Veneers()
	if err != nil {
		return err
	}

	commonCompilerPasses, err := config.CompilerPasses()
	if err != nil {
		return err
	}

	// Here begins the code generation setup
	targetsByLanguage, err := config.OutputLanguages()
	if err != nil {
		return err
	}

	fmt.Printf("Parsing inputs...\n")
	schemas, err := config.LoadSchemas()
	if err != nil {
		return err
	}

	rootCodeJenFS := codejen.NewFS()
	for language, target := range targetsByLanguage {
		fmt.Printf("Running '%s' jennies...\n", language)

		languageOutputDir, err := config.LanguageOutputDir(workingDir, language)
		if err != nil {
			return err
		}

		compilerPasses := commonCompilerPasses.Concat(target.CompilerPasses())
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

		// prepare the jennies
		languageJennies := target.Jennies(config.JenniesConfig())
		languageJennies.AddPostprocessors(common.PathPrefixer(languageOutputDir))

		jenniesInput := common.Context{
			Schemas:  processedSchemas,
			Builders: builders,
		}

		// then delegate the codegen to the jennies
		if err := runJenny(languageJennies, jenniesInput, rootCodeJenFS); err != nil {
			return err
		}

		if config.Output.PackageTemplates != "" {
			packageJennies, err := packageTemplatesJenny(language, config, workingDir)
			if err != nil {
				return err
			}

			if err := runJenny(packageJennies, jenniesInput, rootCodeJenFS); err != nil {
				return err
			}
		}
	}

	if config.Output.RepositoryTemplates != "" {
		repoTemplatesJenny, err := repositoryTemplatesJenny(config, workingDir)
		if err != nil {
			return err
		}
		jennyInput := common.BuildOptions{
			Languages: targetsByLanguage.AsLanguageRefs(),
		}

		if err := runJenny(repoTemplatesJenny, jennyInput, rootCodeJenFS); err != nil {
			return err
		}
	}

	err = rootCodeJenFS.Write(context.Background(), "")
	if err != nil {
		return err
	}

	return nil
}

func repositoryTemplatesJenny(config loaders.Config, workDir string) (*codejen.JennyList[common.BuildOptions], error) {
	outputDir, err := config.OutputDir(workDir)
	if err != nil {
		return nil, err
	}

	repoTemplatesJenny := codejen.JennyListWithNamer[common.BuildOptions](func(_ common.BuildOptions) string {
		return "RepositoryTemplates"
	})
	repoTemplatesJenny.AppendOneToMany(&common.RepositoryTemplate{
		TemplateDir: config.Path(config.Output.RepositoryTemplates),
		ExtraData:   config.Output.TemplatesData,
	})
	repoTemplatesJenny.AddPostprocessors(
		common.GeneratedCommentHeader(config.JenniesConfig()),
		common.PathPrefixer(strings.ReplaceAll(outputDir, "%l", ".")),
	)

	return repoTemplatesJenny, nil
}

func packageTemplatesJenny(language string, config loaders.Config, workDir string) (*codejen.JennyList[common.Context], error) {
	outputDir, err := config.LanguageOutputDir(workDir, language)
	if err != nil {
		return nil, err
	}

	pkgTemplatesJenny := codejen.JennyListWithNamer[common.Context](func(_ common.Context) string {
		return "PackageTemplates" + tools.UpperCamelCase(language)
	})
	pkgTemplatesJenny.AppendOneToMany(&common.PackageTemplate{
		Language:    language,
		TemplateDir: config.Path(config.Output.PackageTemplates),
		ExtraData:   config.Output.TemplatesData,
	})
	pkgTemplatesJenny.AddPostprocessors(common.PathPrefixer(outputDir))

	return pkgTemplatesJenny, nil
}

func runJenny[I any](jenny *codejen.JennyList[I], input I, destinationFS *codejen.FS) error {
	fs, err := jenny.GenerateFS(input)
	if err != nil {
		return err
	}

	return destinationFS.Merge(fs)
}
