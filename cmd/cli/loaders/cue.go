package loaders

import (
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
	"github.com/yalue/merged_fs"
)

func cueLoader(opts Options) ([]*ast.Schema, error) {
	libraries, err := opts.cueIncludeImports()
	if err != nil {
		return nil, err
	}

	allSchemas := make([]*ast.Schema, 0, len(opts.CueEntrypoints))
	for _, entrypoint := range opts.CueEntrypoints {
		pkg := filepath.Base(entrypoint)

		schemaRootValue, err := parseCueEntrypoint(opts, entrypoint)
		if err != nil {
			return nil, err
		}

		schemaAst, err := simplecue.GenerateAST(schemaRootValue, simplecue.Config{
			Package:        pkg, // TODO: extract from input schema/?
			SchemaMetadata: ast.SchemaMeta{
				// TODO: extract these from somewhere
			},
			Libraries: libraries,
		})
		if err != nil {
			return nil, err
		}

		allSchemas = append(allSchemas, schemaAst)
	}

	return allSchemas, nil
}

func parseCueEntrypoint(opts Options, entrypoint string) (cue.Value, error) {
	// this is wasteful: we rebuild this overlay for every entrypoint we parse
	cueFsOverlay, err := buildCueOverlay(opts)
	if err != nil {
		return cue.Value{}, err
	}

	// Load Cue files into Cue build.Instances slice
	// the second arg is a configuration object, we'll see this later
	bis := load.Instances([]string{entrypoint}, &load.Config{
		Overlay:    cueFsOverlay,
		ModuleRoot: "/",
	})

	values, err := cuecontext.New().BuildInstances(bis)
	if err != nil {
		return cue.Value{}, err
	}

	return values[0], nil
}

func buildCueOverlay(opts Options) (map[string]load.Source, error) {
	mockKindsysFS := buildMockKindsysFS()
	libFs, err := buildBaseFSWithLibraries(opts)
	if err != nil {
		return nil, err
	}

	mergedFS := merged_fs.MergeMultiple(append(libFs, mockKindsysFS)...)

	overlay := make(map[string]load.Source)
	if err := toCueOverlay("/", mergedFS, overlay); err != nil {
		return nil, err
	}

	return overlay, nil
}

func buildMockKindsysFS() fs.FS {
	mockFS := fstest.MapFS{
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

	return mockFS
}

func buildBaseFSWithLibraries(opts Options) ([]fs.FS, error) {
	importDefinitions, err := opts.cueIncludeImports()
	if err != nil {
		return nil, err
	}

	var librariesFS []fs.FS
	for _, importDefinition := range importDefinitions {
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
