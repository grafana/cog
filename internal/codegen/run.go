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

	commonPasses, err := pipeline.compilerPasses()
	if err != nil {
		return nil, err
	}

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

		jenniesInput, err := pipeline.jenniesInputForLanguage(target, schemas, commonPasses, veneers)
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

		if pipeline.Output.PackageTemplates != "" {
			packageJennies, err := packageTemplatesJenny(pipeline, language)
			if err != nil {
				return nil, err
			}

			if err := runJenny(packageJennies, jenniesInput, generatedFS); err != nil {
				return nil, err
			}
		}
	}

	if pipeline.Output.RepositoryTemplates != "" {
		repoTemplatesJenny, err := repositoryTemplatesJenny(pipeline)
		if err != nil {
			return nil, err
		}
		jennyInput := languages.BuildOptions{
			Languages: targetsByLanguage.AsLanguageRefs(),
		}

		if err := runJenny(repoTemplatesJenny, jennyInput, generatedFS); err != nil {
			return nil, err
		}
	}

	return generatedFS, nil
}

func (pipeline *Pipeline) jenniesInputForLanguage(language languages.Language, schemas ast.Schemas, commonCompilerPasses compiler.Passes, veneers *rewrite.Rewriter) (languages.Context, error) {
	var err error
	jenniesInput := languages.Context{
		Schemas: schemas,
	}

	// apply common and language-specific compiler passes
	compilerPasses := commonCompilerPasses.Concat(language.CompilerPasses())
	jenniesInput.Schemas, err = compilerPasses.Process(jenniesInput.Schemas)
	if err != nil {
		return languages.Context{}, err
	}

	if !pipeline.Output.Builders {
		return pipeline.formatIdentifiers(language, jenniesInput)
	}

	// from these types, create builders
	builderGenerator := &ast.BuilderGenerator{}
	builders := builderGenerator.FromAST(jenniesInput.Schemas)

	// apply the builder veneers
	builders, err = veneers.ApplyTo(jenniesInput.Schemas, builders, language.Name())
	if err != nil {
		return languages.Context{}, err
	}

	// ensure identifiers are properly formatted
	jenniesInput, err = pipeline.formatIdentifiers(language, jenniesInput)
	if err != nil {
		return languages.Context{}, err
	}

	// if the language defines an identifier formatter, let's apply it.
	if formatterProvider, ok := language.(languages.IdentifiersFormatterProvider); ok {
		jenniesInput, err = languages.FormatIdentifiers(formatterProvider, jenniesInput)
		if err != nil {
			return languages.Context{}, err
		}
	}

	// with the veneers applied, generate "nil-checks" for assignments
	jenniesInput, err = languages.GenerateBuilderNilChecks(language, jenniesInput)
	if err != nil {
		return languages.Context{}, err
	}

	return jenniesInput, nil
}

func (pipeline *Pipeline) formatIdentifiers(language languages.Language, jenniesInput languages.Context) (languages.Context, error) {
	var err error

	// if the language defines an identifier formatter, let's apply it.
	formatterProvider, ok := language.(languages.IdentifiersFormatterProvider)

	if !ok {
		return jenniesInput, nil
	}

	formatter := formatterProvider.IdentifiersFormatter()

	formatterPass := languages.NewIdentifierFormatterPass(jenniesInput.Schemas, formatter)
	jenniesInput.Schemas, err = formatterPass.Process(jenniesInput.Schemas)
	if err != nil {
		return jenniesInput, err
	}

	buildersRewriter := languages.NewIdentifierFormatterBuilderRewriter(jenniesInput.Schemas, formatter)
	jenniesInput.Builders, err = buildersRewriter.Rewrite(jenniesInput.Schemas, jenniesInput.Builders)
	if err != nil {
		return jenniesInput, err
	}

	return jenniesInput, nil
}
