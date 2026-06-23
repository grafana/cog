package codegen

import (
	"context"
	"fmt"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
	"cuelang.org/go/cue/parser"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/github"
	"github.com/grafana/cog/internal/simplecue"
	"github.com/grafana/cog/internal/tools"
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

func (input *CueInput) schemaRootValue() (cue.Value, []simplecue.LibraryInclude, error) {
	if input.Value != nil {
		libraries, err := simplecue.ParseImports(input.CueImports)
		if err != nil {
			return cue.Value{}, nil, err
		}
		return *input.Value, libraries, nil
	}

	if input.Entrypoint == "" && input.URL == "" {
		return cue.Value{}, nil, fmt.Errorf("no entrypoint or url defined in cue input")
	}

	if input.Entrypoint != "" && input.URL != "" {
		return cue.Value{}, nil, fmt.Errorf("only one entrypoint or url defined in cue input")
	}

	libraries, err := simplecue.ParseImports(input.CueImports)
	if err != nil {
		return cue.Value{}, nil, err
	}

	value, err := input.parseCueEntrypoint(context.Background(), libraries)
	if err != nil {
		return cue.Value{}, nil, fmt.Errorf("could not parse cue entrypoint: %w", err)
	}

	return value, libraries, nil
}

func (input *CueInput) interpolateParameters(interpolator ParametersInterpolator) {
	input.InputBase.interpolateParameters(interpolator)

	input.Entrypoint = interpolator(input.Entrypoint)
	input.URL = interpolator(input.URL)
	input.CueImports = tools.Map(input.CueImports, interpolator)
}

func cueLoader(input CueInput) (ast.Schemas, error) {
	schemaRootValue, libraries, err := input.schemaRootValue()
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

func (input *CueInput) parseCueEntrypoint(ctx context.Context, imports []simplecue.LibraryInclude) (cue.Value, error) {
	var err error
	var instanceName string

	prefix := "/cog/vfs/cue.mod/pkg"
	overlay := map[string]load.Source{
		"/cog/vfs/cue.mod/module.cue": load.FromBytes([]byte(
			"language: { version: \"v0.10.1\" }\nmodule: \"cog.vfs\"\n",
		)),
	}

	if err = addLibrariesToOverlay(overlay, prefix, imports); err != nil {
		return cue.Value{}, fmt.Errorf("could not add libraries to overlay: %w", err)
	}

	if input.URL != "" {
		instanceName, err = addURLToOverlay(ctx, overlay, prefix, input.URL)
		if err != nil {
			return cue.Value{}, fmt.Errorf("could not add URL '%s' to overlay: %w", input.URL, err)
		}
	}

	if input.Entrypoint != "" {
		instanceName, err = addDirToOverlay(overlay, prefix, input.Entrypoint)
		if err != nil {
			return cue.Value{}, fmt.Errorf("could not add entrypoint '%s' to overlay: %w", input.Entrypoint, err)
		}
	}

	// Load Cue files into Cue build.Instances slice
	bis := load.Instances([]string{instanceName}, &load.Config{
		Overlay: overlay,
		// Point cue to a directory defined by the cueFsOverlay as base directory
		// for import path resolution instead of using the current working directory.
		// This ensures that only files/schemas defined in the vfs will be parsed.
		Dir: "/cog/vfs",
	})

	value := cuecontext.New().BuildInstance(bis[0])
	if err := value.Err(); err != nil {
		return cue.Value{}, fmt.Errorf("could not build cue instance: %w", err)
	}

	if input.Subpath != "" {
		value = value.LookupPath(cue.ParsePath(input.Subpath))
		if value.Err() != nil {
			return cue.Value{}, fmt.Errorf("could not read subpath '%s': %w", input.Subpath, value.Err())
		}
	}

	if input.Package == "" {
		input.Package = bis[0].PkgName
	}

	return value, nil
}

func addDirToOverlay(overlay map[string]load.Source, prefix string, directory string) (string, error) {
	absolutePath, err := filepath.Abs(directory)
	if err != nil {
		return "", err
	}

	var instanceName string
	err = filepath.Walk(absolutePath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if !strings.HasSuffix(info.Name(), ".cue") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		if instanceName == "" {
			file, err := parser.ParseFile(path, content, parser.ParseComments)
			if err != nil {
				return fmt.Errorf("parse error in '%s': %w", path, err)
			}
			instanceName = "github.com/cog-vfs/" + file.PackageName()
		}

		overlay[prefix+"/"+instanceName+"/"+info.Name()] = load.FromBytes(content)

		return nil
	})
	if err != nil {
		return "", err
	}

	return instanceName, nil
}

func addLibrariesToOverlay(overlay map[string]load.Source, prefix string, imports []simplecue.LibraryInclude) error {
	for _, importDefinition := range imports {
		err := filepath.Walk(importDefinition.FSPath, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			if !strings.HasSuffix(info.Name(), ".cue") {
				return nil
			}

			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			overlay[prefix+"/"+importDefinition.ImportPath+"/"+info.Name()] = load.FromBytes(content)

			return nil
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func addURLToOverlay(ctx context.Context, overlay map[string]load.Source, overlayPrefix string, entrypoint string) (string, error) {
	u, err := url.Parse(entrypoint)
	if err != nil {
		return "", err
	}

	if !strings.HasSuffix(u.Path, ".cue") {
		return "", fmt.Errorf("entrypoint %s must be a cue url", entrypoint)
	}

	ghInfo := github.RawURLToRepoDescriptor(entrypoint)
	if ghInfo == nil {
		return "", fmt.Errorf("could not parse raw url as repo descriptor: %s", entrypoint)
	}

	dirFiles, err := github.FetchDirectory(ctx, *ghInfo, github.DirPath(entrypoint), ".cue")
	if err != nil {
		return "", fmt.Errorf("fetch directory: %w", err)
	}

	allFiles := make(map[string][]byte, len(dirFiles))
	fetched := make(map[string]bool)

	var instanceName string
	for name, content := range dirFiles {
		file, err := parser.ParseFile(name, content, parser.ParseComments)
		if err != nil {
			continue
		}

		if instanceName == "" {
			instanceName = "github.com/cog-vfs/" + file.PackageName()
		}

		overlay[overlayPrefix+"/"+instanceName+"/"+name] = load.FromBytes(content)
		allFiles[name] = content

		for _, imp := range file.Imports {
			importPath := strings.Trim(imp.Path.Value, `"`)
			if !strings.HasPrefix(importPath, "github.com") || fetched[importPath] {
				continue
			}

			// TODO: check if the import is already included via the "libraries includes"

			fetched[importPath] = true
			// TODO: this isn't a raw url :o
			ghInfo = github.RawURLToRepoDescriptor(importPath)
			if ghInfo == nil {
				return "", fmt.Errorf("could not parse import path as repo descriptor: %s", importPath)
			}

			repoPath := strings.TrimPrefix(importPath, "github.com")
			pkgFiles, err := github.FetchDirectory(ctx, *ghInfo, repoPath, ".cue")
			if err != nil {
				return "", fmt.Errorf("could not fetch import path: %s", importPath)
			}

			for fileName, fileContent := range pkgFiles {
				overlay[overlayPrefix+"/"+importPath+"/"+fileName] = load.FromBytes(fileContent)
			}
		}
	}

	return instanceName, nil
}
