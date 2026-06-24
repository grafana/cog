package codegen

import (
	"context"

	"cuelang.org/go/cue"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/simplecue"
	"github.com/grafana/cog/internal/tools"
	cogcue "github.com/grafana/cog/pkg/cue"
)

type genericCueLoader struct {
	*CueInput

	loader func(input CueInput) (ast.Schemas, error)
}

func (loader *genericCueLoader) LoadSchemas(_ context.Context) (ast.Schemas, error) {
	return loader.loader(*loader.CueInput)
}

type CueInput struct {
	InputBase `yaml:",inline"`

	// Entrypoint refers to a directory containing CUE files.
	Entrypoint string `yaml:"entrypoint"`

	// URL to a cue file
	URL string `yaml:"url"`

	// Subpath is used when the schema is not in the root of the schema
	Subpath string `yaml:"subpath"`

	// Value represents the CUE value to use as an input. If specified, it
	// supersedes the Entrypoint and URL options.
	Value *cue.Value `yaml:"-"`

	// ForcedEnvelope decorates the parsed cue Value with an envelope whose
	// name is given. This is useful for dataqueries for example, where the
	// schema doesn't define any suitable top-level object.
	ForcedEnvelope string `yaml:"forced_envelope"`

	// Package name to use for the input schema. If empty, it will be guessed.
	Package string `yaml:"package"`

	// CueImports allows importing additional libraries.
	// Format: [path]:[import]. Example: '../grafana/common-library:github.com/grafana/grafana/packages/grafana-schema/src/common
	CueImports []string `yaml:"cue_imports"`

	// NameFunc allows users to specify an alternative naming strategy for
	// objects and references. It is called with the value passed to the top
	// level method or function and the path to the entity being parsed.
	NameFunc simplecue.NameFunc `yaml:"-"`

	// InlineExternalReference instructs the parser to follow external
	// references (ie: references to objects outside the current schema)
	// and inline them.
	// By default, external references are parsed as actual `ast.Ref` to the
	// external objects.
	InlineExternalReference bool `yaml:"-"`
}

func (input *CueInput) schemaRootValue(ctx context.Context) (cue.Value, []simplecue.LibraryInclude, error) {
	cueInput := cogcue.Input{
		Entrypoint: input.Entrypoint,
		URL:        input.URL,
		Subpath:    input.Subpath,
		Value:      input.Value,
		CueImports: input.CueImports,
	}

	result, err := cogcue.Parse(ctx, cueInput)
	if err != nil {
		return cue.Value{}, nil, err
	}

	libraries := tools.Map(result.Libraries, func(l cogcue.LibraryInclude) simplecue.LibraryInclude {
		return simplecue.LibraryInclude{FSPath: l.FSPath, ImportPath: l.ImportPath}
	})

	if input.Package == "" {
		input.Package = result.Package
	}

	return result.Value, libraries, nil
}

func (input *CueInput) interpolateParameters(interpolator ParametersInterpolator) {
	input.InputBase.interpolateParameters(interpolator)

	input.Entrypoint = interpolator(input.Entrypoint)
	input.URL = interpolator(input.URL)
	input.CueImports = tools.Map(input.CueImports, interpolator)
}

func cueLoader(input CueInput) (ast.Schemas, error) {
	schemaRootValue, libraries, err := input.schemaRootValue(context.Background())
	if err != nil {
		return nil, err
	}

	schema, err := simplecue.GenerateAST(schemaRootValue, simplecue.Config{
		Package:                 input.Package,
		ForceNamedEnvelope:      input.ForcedEnvelope,
		SchemaMetadata:          input.schemaMetadata(),
		Libraries:               libraries,
		NameFunc:                input.NameFunc,
		InlineExternalReference: input.InlineExternalReference,
	})
	if err != nil {
		return nil, err
	}

	return input.filterSchema(schema)
}
