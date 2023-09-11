package generate

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"testing/fstest"

	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/simplecue"
	"github.com/yalue/merged_fs"
)

type cueIncludeImport struct {
	fsPath     string // path of the library on the filesystem
	importPath string // path used in CUE files to import that library
}

func cueLoader(opts options) ([]*ast.File, error) {
	cueFsOverlay, err := buildCueOverlay(opts)
	if err != nil {
		return nil, err
	}

	allSchemas := make([]*ast.File, 0, len(opts.entrypoints))
	for _, entrypoint := range opts.entrypoints {
		pkg := filepath.Base(entrypoint)

		// Load Cue files into Cue build.Instances slice
		// the second arg is a configuration object, we'll see this later
		bis := load.Instances([]string{entrypoint}, &load.Config{
			Overlay: cueFsOverlay,
			//Module:     "github.com/grafana/cog", // TODO: is that needed?
			ModuleRoot: "/",
		})

		values, err := cuecontext.New().BuildInstances(bis)
		if err != nil {
			return nil, err
		}

		schemaAst, err := simplecue.GenerateAST(values[0], simplecue.Config{
			Package: pkg, // TODO: extract from input schema/?
		})
		if err != nil {
			return nil, err
		}

		allSchemas = append(allSchemas, schemaAst)
	}

	return allSchemas, nil
}

func buildCueOverlay(opts options) (map[string]load.Source, error) {
	libFs, err := buildBaseFSWithLibraries(opts)
	if err != nil {
		return nil, err
	}

	overlay := make(map[string]load.Source)
	if err := toCueOverlay("/", libFs, overlay); err != nil {
		return nil, err
	}

	return overlay, nil
}

func buildBaseFSWithLibraries(opts options) (fs.FS, error) {
	importDefinitions, err := opts.cueIncludeImports()
	if err != nil {
		return nil, err
	}

	var librariesFS []fs.FS
	for _, importDefinition := range importDefinitions {
		absPath, err := filepath.Abs(importDefinition.fsPath)
		if err != nil {
			return nil, err
		}

		fmt.Printf("Loading '%s' module from '%s'\n", importDefinition.importPath, absPath)

		libraryFS, err := dirToPrefixedFS(absPath, "cue.mod/pkg/"+importDefinition.importPath)
		if err != nil {
			return nil, err
		}

		librariesFS = append(librariesFS, libraryFS)
	}

	return merged_fs.MergeMultiple(librariesFS...), nil
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
