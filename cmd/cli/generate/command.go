package generate

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/cmd/cli/loaders"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/ast/compiler"
	"github.com/grafana/cog/internal/jennies"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/veneers/rewrite"
	"github.com/grafana/cog/internal/yaml"
	"github.com/spf13/cobra"
)

type Options struct {
	loaders.Options

	JenniesConfig           common.Config
	Languages               []string
	VeneerConfigFiles       []string
	VeneerConfigDirectories []string
	CompilerConfigFiles     []string
	OutputDir               string

	// PackageTemplates is the path to a directory containing "package templates".
	// These templates are used to add arbitrary files to the generated code, with
	// the goal of turning it into a fully-fledged package.
	// Templates in that directory are expected to be organized by language:
	// ```
	// package_templates
	// ├── go
	// │   ├── LICENSE.md
	// │   └── README.md
	// └── typescript
	//     ├── babel.config.json
	//     ├── package.json
	//     ├── README.md
	//     └── tsconfig.json
	// ```
	PackageTemplates string

	// RepositoryTemplates is the path to a directory containing
	// "repository-level templates".
	// These templates are used to add arbitrary files to the repository, such as CI pipelines.
	//
	// Templates in that directory are expected to be organized by language:
	// ```
	// repository_templates
	// ├── go
	// │   └── .github
	// │   	   └── workflows
	// │   	       └── go-ci.yaml
	// └── typescript
	//     └── .github
	//     	   └── workflows
	//     	       └── typescript-ci.yaml
	// ```
	RepositoryTemplates string

	// TemplatesData holds data that will be injected into package and
	// repository templates when rendering them.
	TemplatesData map[string]string
}

func (opts Options) veneerFiles() ([]string, error) {
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

func (opts Options) veneers() (*rewrite.Rewriter, error) {
	veneerFiles, err := opts.veneerFiles()
	if err != nil {
		return nil, err
	}

	config := rewrite.Config{Debug: opts.JenniesConfig.Debug}
	return yaml.NewVeneersLoader().RewriterFrom(veneerFiles, config)
}

func (opts Options) commonCompilerPasses() (compiler.Passes, error) {
	return yaml.NewCompilerLoader().PassesFrom(opts.CompilerConfigFiles)
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

	cmd.Flags().BoolVar(&opts.JenniesConfig.Debug, "debug", false, "Debugging mode.") // TODO: better usage text

	cmd.Flags().BoolVar(&opts.JenniesConfig.Types, "generate-types", true, "Generate types.")          // TODO: better usage text
	cmd.Flags().BoolVar(&opts.JenniesConfig.Builders, "generate-builders", true, "Generate builders.") // TODO: better usage text
	cmd.Flags().StringVar(&opts.PackageTemplates, "package-templates", "", "Directory used as a template used to generate fully fledged packages.")
	cmd.Flags().StringVar(&opts.RepositoryTemplates, "repository-templates", "", "Directory used as a template used to generate additional content at the repository root.")
	cmd.Flags().StringToStringVar(&opts.TemplatesData, "templates-data", nil, "Data to hand over to package and repository templates.")

	cmd.Flags().StringVarP(&opts.OutputDir, "output", "o", "generated", "Output directory.") // TODO: better usage text
	cmd.Flags().StringArrayVarP(&opts.Languages, "language", "l", nil, "Language to generate. If left empty, all supported languages will be generated.")
	cmd.Flags().StringArrayVarP(&opts.VeneerConfigFiles, "veneer", "c", nil, "Veneer configuration file.")
	cmd.Flags().StringArrayVar(&opts.VeneerConfigDirectories, "veneers", nil, "Veneer configuration directories.")
	cmd.Flags().StringArrayVar(&opts.CompilerConfigFiles, "compiler-config", nil, "Compiler configuration file.")

	cmd.Flags().StringArrayVar(&opts.CueEntrypoints, "cue", nil, "CUE input schema.")                                                                                                           // TODO: better usage text
	cmd.Flags().StringArrayVar(&opts.KindsysCoreEntrypoints, "kindsys-core", nil, "Kindys core kinds input schema.")                                                                            // TODO: better usage text
	cmd.Flags().StringArrayVar(&opts.KindsysComposableEntrypoints, "kindsys-composable", nil, "Kindys composable kinds input schema.")                                                          // TODO: better usage text
	cmd.Flags().StringArrayVar(&opts.KindsysCustomEntrypoints, "kindsys-custom", nil, "Kindys custom kinds input schema.")                                                                      // TODO: better usage text
	cmd.Flags().StringArrayVar(&opts.JSONSchemaEntrypoints, "jsonschema", nil, "Jsonschema input schema.")                                                                                      // TODO: better usage text
	cmd.Flags().StringArrayVar(&opts.OpenAPIEntrypoints, "openapi", nil, "Openapi input schema.")                                                                                               // TODO: better usage text
	cmd.Flags().StringVar(&opts.KindRegistryPath, "kind-registry", "", "Kind registry input.")                                                                                                  // TODO: better usage text
	cmd.Flags().StringVar(&opts.JSONSchemaRegistryPath, "jsonschema-registry", "", "JSONschema registry input. This flag is totally experimental and it could be deleted in forward versions.") // TODO: better usage text

	cmd.Flags().StringArrayVarP(&opts.CueImports, "include-cue-import", "I", nil, "Specify an additional library import directory. Format: [path]:[import]. Example: '../grafana/common-library:github.com/grafana/grafana/packages/grafana-schema/src/common")
	cmd.Flags().StringVar(&opts.KindRegistryVersion, "kind-registry-version", "next", "Schemas version")

	for _, jenny := range languageJennies {
		jenny.RegisterCliFlags(cmd)
	}

	_ = cmd.MarkFlagDirname("package-templates")
	_ = cmd.MarkFlagDirname("cue")
	_ = cmd.MarkFlagDirname("kindsys-core")
	_ = cmd.MarkFlagDirname("kindsys-custom")
	_ = cmd.MarkFlagDirname("kind-registry")
	_ = cmd.MarkFlagDirname("jsonschema-registry")
	_ = cmd.MarkFlagFilename("jsonschema")
	_ = cmd.MarkFlagDirname("openapi")
	_ = cmd.MarkFlagDirname("output")
	_ = cmd.MarkFlagFilename("veneer")
	_ = cmd.MarkFlagDirname("veneers")

	return cmd
}

func doGenerate(allTargets jennies.LanguageJennies, opts Options) error {
	veneers, err := opts.veneers()
	if err != nil {
		return err
	}

	commonCompilerPasses, err := opts.commonCompilerPasses()
	if err != nil {
		return err
	}

	// Here begins the code generation setup
	targetsByLanguage, err := allTargets.ForLanguages(opts.Languages)
	if err != nil {
		return err
	}

	fmt.Printf("Parsing inputs...\n")
	schemas, err := loaders.LoadAll(opts.Options)
	if err != nil {
		return err
	}

	rootCodeJenFS := codejen.NewFS()
	for language, target := range targetsByLanguage {
		fmt.Printf("Running '%s' jennies...\n", language)

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
		outputDir := strings.ReplaceAll(opts.OutputDir, "%l", language)
		languageJennies := target.Jennies(opts.JenniesConfig)
		languageJennies.AddPostprocessors(common.PathPrefixer(outputDir))

		if opts.PackageTemplates != "" {
			languageJennies.AppendOneToMany(&common.PackageTemplate{
				Language:    language,
				TemplateDir: opts.PackageTemplates,
				ExtraData:   opts.TemplatesData,
			})
		}

		// then delegate the codegen to the jennies
		fs, err := languageJennies.GenerateFS(common.Context{
			Schemas:  processedSchemas,
			Builders: builders,
		})
		if err != nil {
			return err
		}

		if err = rootCodeJenFS.Merge(fs); err != nil {
			return err
		}
	}

	if opts.RepositoryTemplates != "" {
		globalJenny := codejen.JennyListWithNamer[common.BuildOptions](func(_ common.BuildOptions) string {
			return "Global"
		})
		globalJenny.AppendOneToMany(&common.RepositoryTemplate{
			TemplateDir: opts.RepositoryTemplates,
			ExtraData:   opts.TemplatesData,
		})
		globalJenny.AddPostprocessors(
			common.GeneratedCommentHeader(opts.JenniesConfig),
			common.PathPrefixer(filepath.Clean(strings.ReplaceAll(opts.OutputDir, "%l", "."))),
		)

		fs, err := globalJenny.GenerateFS(common.BuildOptions{
			Languages: targetsByLanguage.AsLanguageRefs(),
		})
		if err != nil {
			return err
		}

		if err = rootCodeJenFS.Merge(fs); err != nil {
			return err
		}
	}

	err = rootCodeJenFS.Write(context.Background(), "")
	if err != nil {
		return err
	}

	return nil
}
