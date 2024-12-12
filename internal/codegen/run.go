package codegen

import (
	"context"
	"fmt"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/ast/compiler"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/languages"
	"github.com/grafana/cog/internal/veneers/rewrite"
)

func (pipeline *Pipeline) Run(ctx context.Context) (*codejen.FS, error) {
	veneers, err := pipeline.veneers()
	if err != nil {
		return nil, err
	}

	commonPasses, err := pipeline.commonPasses()
	if err != nil {
		return nil, err
	}

	finalPasses := pipeline.finalPasses()

	// Here begins the code generation setup
	targetsByLanguage, err := pipeline.outputLanguages()
	if err != nil {
		return nil, err
	}

	pipeline.reporter("Parsing inputs...")
	schemas, err := pipeline.LoadSchemas(ctx)
	if err != nil {
		return nil, err
	}

	generatedFS := codejen.NewFS()
	for language, target := range targetsByLanguage {
		pipeline.reporter(fmt.Sprintf("Running '%s' jennies...", language))

		languageOutputDir, err := pipeline.languageOutputDir(pipeline.currentDirectory, language)
		if err != nil {
			return nil, err
		}

		jenniesInput, err := pipeline.jenniesInputForLanguage(target, schemas, commonPasses, finalPasses, veneers)
		if err != nil {
			return nil, err
		}

		// prepare the jennies
		languageJennies := target.Jennies(pipeline.jenniesConfig())
		languageJennies.AddPostprocessors(common.PathPrefixer(languageOutputDir))

		// then delegate the codegen to the jennies
		if err := runJenny(languageJennies, jenniesInput, generatedFS); err != nil {
			return nil, err
		}
	}

	if pipeline.Output.RepositoryTemplates != "" {
		repoTemplatesJenny, err := repositoryTemplatesJenny(pipeline)
		if err != nil {
			return nil, err
		}
		jennyInput := common.BuildOptions{
			Languages: targetsByLanguage.AsLanguageRefs(),
		}

		if err := runJenny(repoTemplatesJenny, jennyInput, generatedFS); err != nil {
			return nil, err
		}
	}

	return generatedFS, nil
}

func (pipeline *Pipeline) jenniesInputForLanguage(language languages.Language, schemas ast.Schemas, commonPasses compiler.Passes, finalPasses compiler.Passes, veneers *rewrite.Rewriter) (languages.Context, error) {
	var err error
	jenniesInput := languages.Context{
		Schemas: schemas,
	}

	// apply common and language-specific compiler passes
	compilerPasses := commonPasses.Concat(language.CompilerPasses()).Concat(finalPasses)
	jenniesInput.Schemas, err = compilerPasses.Process(jenniesInput.Schemas)
	if err != nil {
		return languages.Context{}, err
	}

	if !pipeline.Output.Builders {
		return jenniesInput, nil
	}

	// from schemas, derive builders
	builderGenerator := &ast.BuilderGenerator{}
	jenniesInput.Builders = builderGenerator.FromAST(jenniesInput.Schemas)

	// apply veneers to builders
	jenniesInput.Builders, err = veneers.ApplyTo(jenniesInput.Schemas, jenniesInput.Builders, language.Name())
	if err != nil {
		return languages.Context{}, err
	}

	// with the veneers applied, generate "nil-checks" for assignments
	jenniesInput, err = languages.GenerateBuilderNilChecks(language, jenniesInput)
	if err != nil {
		return languages.Context{}, err
	}

	return jenniesInput, nil
}
