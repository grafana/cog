package cog

import (
	"context"
	"fmt"

	"cuelang.org/go/cue"
	"github.com/grafana/cog/internal/ast/compiler"
	"github.com/grafana/cog/internal/codegen"
	"github.com/grafana/cog/internal/jennies/golang"
	"github.com/grafana/cog/internal/jennies/typescript"
)

// singleSchemaPipeline represents a simplified codegen.Pipeline, meant to
// take a single input schema and generates types for it in a single output
// language.
type singleSchemaPipeline struct {
	debug       bool
	input       *codegen.Input
	finalPasses compiler.Passes
	output      *codegen.OutputLanguage
}

func AppendCommentToObjects(comment string) compiler.Pass {
	return &compiler.AppendCommentObjects{
		Comment: comment,
	}
}

func PrefixObjectsNames(prefix string) compiler.Pass {
	return &compiler.PrefixObjectNames{
		Prefix: prefix,
	}
}

type SingleSchemaOption func(*singleSchemaPipeline)

type CUEOption func(*codegen.CueInput)

func ForceEnvelope(envelopeName string) CUEOption {
	return func(input *codegen.CueInput) {
		input.ForcedEnvelope = envelopeName
	}
}

func Debug(enabled bool) SingleSchemaOption {
	return func(pipeline *singleSchemaPipeline) {
		pipeline.debug = enabled
	}
}

func CUEValue(pkgName string, value cue.Value, opts ...CUEOption) SingleSchemaOption {
	return func(pipeline *singleSchemaPipeline) {
		cueInput := &codegen.CueInput{
			Package: pkgName,
			Value:   &value,
		}

		for _, opt := range opts {
			opt(cueInput)
		}

		pipeline.input = &codegen.Input{
			Cue: cueInput,
		}
	}
}

func GoTypes() SingleSchemaOption {
	return func(pipeline *singleSchemaPipeline) {
		pipeline.output = &codegen.OutputLanguage{
			Go: &golang.Config{
				GenerateGoMod: false,
				SkipRuntime:   true,
			},
		}
	}
}

func TypescriptTypes() SingleSchemaOption {
	return func(pipeline *singleSchemaPipeline) {
		pipeline.output = &codegen.OutputLanguage{
			Typescript: &typescript.Config{
				SkipRuntime: true,
				SkipIndex:   true,
			},
		}
	}
}

func SchemaTransformations(passes ...compiler.Pass) SingleSchemaOption {
	return func(pipeline *singleSchemaPipeline) {
		pipeline.finalPasses = passes
	}
}

// TypesFromSchema generates types from a single input schema and a single
// output language.
func TypesFromSchema(ctx context.Context, options ...SingleSchemaOption) ([]byte, error) {
	simplePipeline := &singleSchemaPipeline{}

	for _, option := range options {
		option(simplePipeline)
	}

	pipeline, err := codegen.NewPipeline()
	if err != nil {
		return nil, err
	}
	pipeline.Debug = simplePipeline.debug

	// Inputs
	if simplePipeline.input == nil {
		return nil, fmt.Errorf("no input configured")
	}
	pipeline.Inputs = []*codegen.Input{simplePipeline.input}

	// Transformations
	pipeline.Transforms.FinalPasses = simplePipeline.finalPasses

	// Outputs
	if simplePipeline.output == nil {
		return nil, fmt.Errorf("no output configured")
	}
	pipeline.Output = codegen.Output{
		Types:     true,
		Languages: []*codegen.OutputLanguage{simplePipeline.output},
	}

	// Run the codegen pipeline and return the generated file's content.
	// Note: since this pipeline is about types with no runtime, we expect a
	// single file to be generated.
	generatedFS, err := pipeline.Run(ctx)
	if err != nil {
		return nil, err
	}

	generatedFiles := generatedFS.AsFiles()
	if len(generatedFiles) != 1 {
		return nil, fmt.Errorf("expected a single generated file, got %d", len(generatedFiles))
	}

	return generatedFiles[0].Data, nil
}
