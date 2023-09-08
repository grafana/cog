package main

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"testing/fstest"

	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies"
	"github.com/grafana/cog/internal/simplecue"
	"github.com/yalue/merged_fs"
)

func main() {
	entrypoints := []string{
		"./schemas/cue/core/dashboard/",
		// "./schemas/cue/core/playlist/",

		"./schemas/cue/composable/timeseries/",

		"github.com/grafana/grafana/packages/grafana-schema/src/common",
	}

	cueFsOverlay, err := buildCueOverlay()
	if err != nil {
		panic(err)
	}

	allSchemas := make([]*ast.File, 0, len(entrypoints))
	for _, entrypoint := range entrypoints {
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
			panic(err)
		}

		schemaAst, err := simplecue.GenerateAST(values[0], simplecue.Config{
			Package: pkg, // TODO: extract from input schema/?
		})
		if err != nil {
			panic(err)
		}

		allSchemas = append(allSchemas, schemaAst)
	}

	// Here begins the code generation setup
	targetsByLanguage := jennies.All()
	rootCodeJenFS := codejen.NewFS()

	for language, target := range targetsByLanguage {
		fmt.Printf("Running '%s' jennies...\n", language)

		var err error
		processedAsts := allSchemas

		for _, compilerPass := range target.CompilerPasses {
			processedAsts, err = compilerPass.Process(processedAsts)
			if err != nil {
				panic(err)
			}
		}

		fs, err := target.Jennies.GenerateFS(processedAsts)
		if err != nil {
			panic(err)
		}

		err = rootCodeJenFS.Merge(fs)
		if err != nil {
			panic(err)
		}
	}

	err = rootCodeJenFS.Write(context.Background(), "generated")
	if err != nil {
		panic(err)
	}
}

func buildCueOverlay() (map[string]load.Source, error) {
	libFs, err := buildBaseFSWithLibraries()
	if err != nil {
		return nil, err
	}

	overlay := make(map[string]load.Source)
	if err := ToCueOverlay("/", libFs, overlay); err != nil {
		return nil, err
	}

	return overlay, nil
}

func buildBaseFSWithLibraries() (fs.FS, error) {
	// TODO: these should be received as inputs/arguments/parameters
	importDefinitions := [][2]string{
		{
			"github.com/grafana/grafana/packages/grafana-schema/src/common",
			"../kind-registry/grafana/next/common",
		},
		{
			"github.com/grafana/cog",
			".",
		},
	}

	var librariesFS []fs.FS
	for _, importDefinition := range importDefinitions {
		absPath, err := filepath.Abs(importDefinition[1])
		if err != nil {
			return nil, err
		}

		fmt.Printf("Loading '%s' module from '%s'\n", importDefinition[0], absPath)

		libraryFS, err := dirToPrefixedFS(absPath, "cue.mod/pkg/"+importDefinition[0])
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

// ToOverlay converts an fs.FS into a CUE loader overlay.
func ToCueOverlay(prefix string, vfs fs.FS, overlay map[string]load.Source) error {
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
		defer f.Close() // nolint: errcheck

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
