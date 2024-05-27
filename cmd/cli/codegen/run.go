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

	commonCompilerPasses, err := pipeline.compilerPasses()
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

		jenniesInput, err := jenniesInputForLanguage(target, schemas, commonCompilerPasses, veneers)
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
		jennyInput := common.BuildOptions{
			Languages: targetsByLanguage.AsLanguageRefs(),
		}

		if err := runJenny(repoTemplatesJenny, jennyInput, generatedFS); err != nil {
			return nil, err
		}
	}

	return generatedFS, nil
}

func jenniesInputForLanguage(language languages.Language, schemas ast.Schemas, commonCompilerPasses compiler.Passes, veneers *rewrite.Rewriter) (common.Context, error) {
	// apply common and language-specific compiler passes
	compilerPasses := commonCompilerPasses.Concat(language.CompilerPasses())
	processedSchemas, err := compilerPasses.Process(schemas)
	if err != nil {
		return common.Context{}, err
	}

	// from these types, create builders
	builderGenerator := &ast.BuilderGenerator{}
	builders := builderGenerator.FromAST(processedSchemas)

	// apply the builder veneers
	builders, err = veneers.ApplyTo(processedSchemas, builders, language.Name())
	if err != nil {
		return common.Context{}, err
	}

	jenniesInput := common.Context{
		Schemas:  processedSchemas,
		Builders: builders,
	}

	// with the veneers applied, generate "nil-checks" for assignments
	jenniesInput, err = languages.GenerateBuilderNilChecks(language, jenniesInput)
	if err != nil {
		return common.Context{}, err
	}

	return jenniesInput, nil
}
