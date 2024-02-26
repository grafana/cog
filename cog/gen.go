package cog

import (
	"fmt"
	"strings"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/cmd/cli/generate"
	"github.com/grafana/cog/cmd/cli/loaders"
	"github.com/grafana/cog/internal/ast/compiler"
	"github.com/grafana/cog/internal/jennies"
	"github.com/grafana/cog/internal/jennies/common"
)

type Gen struct {
	cfg                  Config
	cueEntries           []string
	cueCoreEntries       []string
	cueComposableEntries []string
	cueCustomEntries     []string
	jsonSchemaEntries    []string
	openAPIEntries       []string
}

func NewGen(cfg Config) (*Gen, error) {
	g := &Gen{
		cfg:                  cfg,
		cueEntries:           make([]string, 0),
		cueCoreEntries:       make([]string, 0),
		cueComposableEntries: make([]string, 0),
		cueCustomEntries:     make([]string, 0),
		jsonSchemaEntries:    make([]string, 0),
		openAPIEntries:       make([]string, 0),
	}

	for _, file := range cfg.FileDirs {
		switch cfg.Kind {
		case CoreKind:
			g.cueCoreEntries = append(g.cueCoreEntries, file)
		case ComposableKind:
			g.cueComposableEntries = append(g.cueComposableEntries, file)
		case CustomKind:
			g.cueCustomEntries = append(g.cueCustomEntries, file)
		default:
			g.cueEntries = append(g.cueEntries, file)
		}
	}

	return g, nil
}

func (g *Gen) Generate() (codejen.Files, error) {
	opts := generate.Options{
		Options: loaders.Options{
			CueEntrypoints:               g.cueEntries,
			KindsysCoreEntrypoints:       g.cueCoreEntries,
			KindsysComposableEntrypoints: g.cueComposableEntries,
			KindsysCustomEntrypoints:     g.cueCustomEntries,
			OpenAPIEntrypoints:           g.openAPIEntries,
			JSONSchemaEntrypoints:        g.jsonSchemaEntries,
		},
		JenniesConfig: common.Config{
			Debug:    g.cfg.Debug,
			Builders: false,
			Types:    true,
			GoConfig: common.GoConfig{
				PackageRoot: g.cfg.GoConfig.PackageRoot,
			},
			TSConfig: common.TSConfig{
				GenTSIndex:   g.cfg.TSConfig.GenTSIndex,
				GenRuntime:   g.cfg.TSConfig.GenRuntime,
				RuntimePath:  g.cfg.TSConfig.RuntimePath,
				ImportMapper: g.cfg.TSConfig.ImportMapper,
			},
			RenameOutputFunc: g.cfg.RenameOutputFunc,
		},
		Languages: g.cfg.Languages.languages(),
		OutputDir: g.cfg.OutputDir,
	}

	return doGenerate(jennies.All(), opts)
}

func doGenerate(allTargets jennies.LanguageJennies, opts generate.Options) (codejen.Files, error) {
	// Here begins the code generation setup
	targetsByLanguage, err := allTargets.ForLanguages(opts.Languages)
	if err != nil {
		return nil, err
	}

	schemas, err := loaders.LoadAll(opts.Options)
	if err != nil {
		return nil, err
	}

	var files codejen.Files
	for language, target := range targetsByLanguage {
		fmt.Printf("Running '%s' jennies...\n", language)

		compilerPasses := compiler.CommonPasses().Concat(target.CompilerPasses())
		processedSchemas, err := compilerPasses.Process(schemas)
		if err != nil {
			return nil, err
		}

		// prepare the jennies
		outputDir := strings.ReplaceAll(opts.OutputDir, "%l", language)
		languageJennies := target.Jennies(opts.JenniesConfig)
		languageJennies.AddPostprocessors(common.PathPrefixer(outputDir))

		// then delegate the codegen to the jennies
		files, err = languageJennies.Generate(common.Context{
			Schemas: processedSchemas,
		})
		if err != nil {
			return nil, err
		}
	}

	return files, nil
}
