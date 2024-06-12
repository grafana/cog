package cog

import (
	"context"
	"fmt"

	"cuelang.org/go/cue"
	"github.com/grafana/cog/internal/ast/compiler"
	"github.com/grafana/cog/internal/codegen"
	"github.com/grafana/cog/internal/jennies/golang"
	"github.com/grafana/cog/internal/jennies/typescript"
	"github.com/grafana/cog/internal/simplecue"
)

// SchemaToTypesPipeline represents a simplified codegen.Pipeline, meant to
// take a single input schema and generates types for it in a single output
// language.
type SchemaToTypesPipeline struct {
	debug       bool
	input       *codegen.Input
	finalPasses compiler.Passes
	output      *codegen.OutputLanguage
}

// AppendCommentToObjects adds the given comment to every object definition.
func AppendCommentToObjects(comment string) compiler.Pass {
	return &compiler.AppendCommentObjects{
		Comment: comment,
	}
}

// PrefixObjectsNames adds the given prefix to every object's name.
func PrefixObjectsNames(prefix string) compiler.Pass {
	return &compiler.PrefixObjectNames{
		Prefix: prefix,
	}
}

type CUEOption func(*codegen.CueInput)

// ForceEnvelope decorates the parsed cue Value with an envelope whose
// name is given. This is useful for dataqueries for example, where the
// schema doesn't define any suitable top-level object.
func ForceEnvelope(envelopeName string) CUEOption {
	return func(input *codegen.CueInput) {
		input.ForcedEnvelope = envelopeName
	}
}

// NameFunc specifies the naming strategy used for objects and references.
// It is called with the value passed to the top level method or function and
// the path to the entity being parsed.
func NameFunc(nameFunc simplecue.NameFunc) CUEOption {
	return func(input *codegen.CueInput) {
		input.NameFunc = nameFunc
	}
}

// TypesFromSchema generates types from a single input schema and a single
// output language.
func TypesFromSchema() *SchemaToTypesPipeline {
	return &SchemaToTypesPipeline{}
}

// Debug controls whether debug mode is enabled or not.
// When enabled, more information is included in the generated output,
// such as an audit trail of applied transformations.
func (pipeline *SchemaToTypesPipeline) Debug(enabled bool) *SchemaToTypesPipeline {
	pipeline.debug = enabled
	return pipeline
}

// Run executes the codegen pipeline and returns its output.
func (pipeline *SchemaToTypesPipeline) Run(ctx context.Context) ([]byte, error) {
	// Validation
	if pipeline.input == nil {
		return nil, fmt.Errorf("no input configured")
	}
	if pipeline.output == nil {
		return nil, fmt.Errorf("no output configured")
	}

	codegenPipeline, err := codegen.NewPipeline()
	if err != nil {
		return nil, err
	}
	codegenPipeline.Inputs = []*codegen.Input{pipeline.input}
	codegenPipeline.Transforms = codegen.Transforms{
		FinalPasses: pipeline.finalPasses,
	}
	codegenPipeline.Output = codegen.Output{
		Types:     true,
		Languages: []*codegen.OutputLanguage{pipeline.output},
	}

	// Run the codegen pipeline and return the generated file's content.
	// Note: since this pipeline is about types with no runtime, we expect a
	// single file to be generated.
	generatedFS, err := codegenPipeline.Run(ctx)
	if err != nil {
		return nil, err
	}

	generatedFiles := generatedFS.AsFiles()
	if len(generatedFiles) != 1 {
		return nil, fmt.Errorf("expected a single generated file, got %d", len(generatedFiles))
	}

	return generatedFiles[0].Data, nil
}

/**********
 * Inputs *
 **********/

// CUEValue sets the pipeline's input to the given cue value.
func (pipeline *SchemaToTypesPipeline) CUEValue(pkgName string, value cue.Value, opts ...CUEOption) *SchemaToTypesPipeline {
	cueInput := &codegen.CueInput{
		Package: pkgName,
		Value:   &value,
	}

	for _, opt := range opts {
		opt(cueInput)
	}

	pipeline.input = &codegen.Input{Cue: cueInput}

	return pipeline
}

/*******************
 * Transformations *
 *******************/

// SchemaTransformations adds the given transformations to the set of
// transformations that will be applied to the input schema.
func (pipeline *SchemaToTypesPipeline) SchemaTransformations(passes ...compiler.Pass) *SchemaToTypesPipeline {
	pipeline.finalPasses = append(pipeline.finalPasses, passes...)

	return pipeline
}

/***********
 * Outputs *
 ***********/

// Golang sets the output to Golang types.
func (pipeline *SchemaToTypesPipeline) Golang() *SchemaToTypesPipeline {
	pipeline.output = &codegen.OutputLanguage{
		Go: &golang.Config{
			GenerateGoMod: false,
			SkipRuntime:   true,
		},
	}
	return pipeline
}

// Typescript sets the output to Typescript types.
func (pipeline *SchemaToTypesPipeline) Typescript() *SchemaToTypesPipeline {
	pipeline.output = &codegen.OutputLanguage{
		Typescript: &typescript.Config{
			SkipRuntime: true,
			SkipIndex:   true,
		},
	}
	return pipeline
}
