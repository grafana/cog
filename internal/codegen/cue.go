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
	cueast "cuelang.org/go/cue/ast"
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

type cueOverlay struct {
	files     map[string]load.Source
	instances map[string]bool
}

func (input *CueInput) parseCueEntrypoint(ctx context.Context, imports []simplecue.LibraryInclude) (cue.Value, error) {
	var err error
	var instanceName string

	prefix := "/cog/vfs/cue.mod/pkg"
	overlay := cueOverlay{
		files: map[string]load.Source{
			"/cog/vfs/cue.mod/module.cue": load.FromBytes([]byte(
				"language: { version: \"v0.10.1\" }\nmodule: \"cog.vfs\"\n",
			)),
		},
		instances: map[string]bool{},
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
		instanceName, err = addDirToOverlay(ctx, overlay, prefix, input.Entrypoint)
		if err != nil {
			return cue.Value{}, fmt.Errorf("could not add entrypoint '%s' to overlay: %w", input.Entrypoint, err)
		}
	}

	// Load Cue files into Cue build.Instances slice
	bis := load.Instances([]string{instanceName}, &load.Config{
		Overlay: overlay.files,
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

func addDirToOverlay(ctx context.Context, overlay cueOverlay, prefix string, directory string) (string, error) {
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

		file, err := parser.ParseFile(path, content, parser.ParseComments)
		if err != nil {
			return fmt.Errorf("parse error in '%s': %w", path, err)
		}

		if instanceName == "" {
			instanceName = "github.com/cog-vfs/" + file.PackageName()
		}

		err = addFileImportsToOverlay(ctx, overlay, prefix, file)
		if err != nil {
			return fmt.Errorf("could add file imports to overlay: %w", err)
		}

		instance := prefix + "/" + instanceName

		overlay.files[instance+"/"+info.Name()] = load.FromBytes(content)
		overlay.instances[instance] = true

		return nil
	})
	if err != nil {
		return "", err
	}

	return instanceName, nil
}

func addLibrariesToOverlay(overlay cueOverlay, prefix string, imports []simplecue.LibraryInclude) error {
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

			instance := prefix + "/" + importDefinition.ImportPath

			overlay.files[instance+"/"+info.Name()] = load.FromBytes(content)
			overlay.instances[instance] = true

			return nil
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func addURLToOverlay(ctx context.Context, overlay cueOverlay, overlayPrefix string, entrypoint string) (string, error) {
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

	var instanceName string
	for name, content := range dirFiles {
		file, err := parser.ParseFile(name, content, parser.ParseComments)
		if err != nil {
			return "", fmt.Errorf("parse error in '%s': %w", name, err)
		}

		if instanceName == "" {
			instanceName = "github.com/cog-vfs/" + file.PackageName()
		}

		instance := overlayPrefix + "/" + instanceName

		overlay.files[instance+"/"+name] = load.FromBytes(content)
		overlay.instances[instance] = true

		err = addFileImportsToOverlay(ctx, overlay, overlayPrefix, file)
		if err != nil {
			return "", fmt.Errorf("could add file imports to overlay: %w", err)
		}
	}

	return instanceName, nil
}

func addFileImportsToOverlay(ctx context.Context, overlay cueOverlay, overlayPrefix string, file *cueast.File) error {
	for _, imp := range file.Imports {
		importPath := strings.Trim(imp.Path.Value, `"`)
		instance := overlayPrefix + "/" + importPath

		if !strings.HasPrefix(importPath, "github.com") || overlay.instances[instance] {
			// TODO: log? (could also just be an import if a standard pkg)
			continue
		}

		// TODO: make the ref configurable?
		ghInfo := github.URLToRepoDescriptor("https://"+importPath, "main")
		if ghInfo == nil {
			return fmt.Errorf("could not parse import path as repo descriptor: %s", importPath)
		}

		repoPath := strings.TrimPrefix(importPath, fmt.Sprintf("github.com/%s/%s/", ghInfo.Owner, ghInfo.Name))
		pkgFiles, err := github.FetchDirectory(ctx, *ghInfo, repoPath, ".cue")
		if err != nil {
			return fmt.Errorf("could not fetch import path %s: %w", importPath, err)
		}

		for fileName, fileContent := range pkgFiles {
			overlay.files[instance+"/"+fileName] = load.FromBytes(fileContent)
			overlay.instances[instance] = true
		}
	}

	return nil
}
