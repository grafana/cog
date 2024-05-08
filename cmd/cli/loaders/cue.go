package loaders

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"testing/fstest"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/simplecue"
	"github.com/grafana/cog/internal/tools"
	"github.com/yalue/merged_fs"
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

	// CueImports allows importing additional libraries.
	// Format: [path]:[import]. Example: '../grafana/common-library:github.com/grafana/grafana/packages/grafana-schema/src/common
	CueImports []string `yaml:"cue_imports"`
}

func (input *CueInput) InterpolateParameters(interpolator ParametersInterpolator) {
	input.InputBase.InterpolateParameters(interpolator)

	input.Entrypoint = interpolator(input.Entrypoint)
	input.CueImports = tools.Map(input.CueImports, interpolator)
}

func cueLoader(input CueInput) (ast.Schemas, error) {
	libraries, err := simplecue.ParseImports(input.CueImports)
	if err != nil {
		return nil, err
	}

	pkg := filepath.Base(input.Entrypoint)

	schemaRootValue, err := parseCueEntrypoint(input.Entrypoint, libraries, pkg)
	if err != nil {
		return nil, err
	}

	schema, err := simplecue.GenerateAST(schemaRootValue, simplecue.Config{
		Package:        pkg,              // TODO: extract from input schema/?
		SchemaMetadata: ast.SchemaMeta{}, // TODO: extract these from somewhere
		Libraries:      libraries,
	})
	if err != nil {
		return nil, err
	}

	return input.filterSchema(schema)
}

func parseCueEntrypoint(entrypoint string, imports []simplecue.LibraryInclude, expectedCuePkgName string) (cue.Value, error) {
	cueFsOverlay, err := buildCueOverlay(imports, entrypoint, expectedCuePkgName)
	if err != nil {
		return cue.Value{}, err
	}

	// Load Cue files into Cue build.Instances slice
	bis := load.Instances([]string{"/" + expectedCuePkgName}, &load.Config{
		Overlay:    cueFsOverlay,
		ModuleRoot: "/",
	})

	value := cuecontext.New().BuildInstance(bis[0])
	if err := value.Err(); err != nil {
		return cue.Value{}, err
	}

	return value, nil
}

func buildCueOverlay(imports []simplecue.LibraryInclude, entrypoint string, expectedCuePkgName string) (map[string]load.Source, error) {
	mockKindsysFS := buildMockKindsysFS()
	libFs, err := buildBaseFSWithLibraries(imports)
	if err != nil {
		return nil, err
	}

	entrypointFS, err := dirToPrefixedFS(entrypoint, expectedCuePkgName)
	if err != nil {
		return nil, err
	}

	mergedFS := merged_fs.MergeMultiple(append(libFs, mockKindsysFS, entrypointFS)...)

	overlay := make(map[string]load.Source)
	if err := toCueOverlay("/", mergedFS, overlay); err != nil {
		return nil, err
	}

	return overlay, nil
}

func buildMockKindsysFS() fs.FS {
	return fstest.MapFS{
		"cue.mod/pkg/github.com/grafana/kindsys/composable.cue": &fstest.MapFile{
			Data: []byte(`package kindsys
Composable: {
	...
}`),
		},
		"cue.mod/pkg/github.com/grafana/kindsys/core.cue": &fstest.MapFile{
			Data: []byte(`package kindsys
Core: {
	...
}`),
		},
		"cue.mod/pkg/github.com/grafana/kindsys/custom.cue": &fstest.MapFile{
			Data: []byte(`package kindsys
Custom: {
	...
}`),
		},
	}
}

func buildBaseFSWithLibraries(imports []simplecue.LibraryInclude) ([]fs.FS, error) {
	var librariesFS []fs.FS
	for _, importDefinition := range imports {
		absPath, err := filepath.Abs(importDefinition.FSPath)
		if err != nil {
			return nil, err
		}

		libraryFS, err := dirToPrefixedFS(absPath, "cue.mod/pkg/"+importDefinition.ImportPath)
		if err != nil {
			return nil, err
		}

		librariesFS = append(librariesFS, libraryFS)
	}

	return librariesFS, nil
}

func dirToPrefixedFS(directory string, prefix string) (fs.FS, error) {
	dirHandle, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	commonFS := fstest.MapFS{}
	for _, file := range dirHandle {
		if file.IsDir() {
			continue
		}

		content, err := os.ReadFile(filepath.Join(directory, file.Name()))
		if err != nil {
			return nil, err
		}

		commonFS[filepath.Join(prefix, file.Name())] = &fstest.MapFile{Data: content}
	}

	return commonFS, nil
}

// ToOverlay converts a fs.FS into a CUE loader overlay.
func toCueOverlay(prefix string, vfs fs.FS, overlay map[string]load.Source) error {
	// TODO why not just stick the prefix on automatically...?
	if !filepath.IsAbs(prefix) {
		return fmt.Errorf("must provide absolute path prefix when generating cue overlay, got %q", prefix)
	}
	err := fs.WalkDir(vfs, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		f, err := vfs.Open(path)
		if err != nil {
			return err
		}
		defer func() { _ = f.Close() }()

		b, err := io.ReadAll(f)
		if err != nil {
			return err
		}

		overlay[filepath.Join(prefix, path)] = load.FromBytes(b)

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
