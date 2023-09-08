package main

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"testing/fstest"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies"
	"github.com/grafana/cog/internal/simplecue"
	"github.com/grafana/kindsys"
	"github.com/grafana/thema"
)

func main() {
	themaRuntime := thema.NewRuntime(cuecontext.New())

	entrypoints := []string{"./schemas/kindsys/core/dashboard"}
	pkg := "dashboard"

	overlayFS, err := dirToPrefixedFS(entrypoints[0], "")
	if err != nil {
		panic(err)
	}

	cueInstance, err := kindsys.BuildInstance(themaRuntime.Context(), ".", "kind", overlayFS)
	if err != nil {
		panic(fmt.Errorf("could not load kindsys instance: %w", err))
	}

	props, err := kindsys.ToKindProps[kindsys.CoreProperties](cueInstance)
	if err != nil {
		panic(fmt.Errorf("could not convert cue value to kindsys props: %w", err))
	}

	kindDefinition := kindsys.Def[kindsys.CoreProperties]{
		V:          cueInstance,
		Properties: props,
	}

	boundKind, err := kindsys.BindCore(themaRuntime, kindDefinition)
	if err != nil {
		panic(fmt.Errorf("could not bind kind definition to kind: %w", err))
	}

	rawLatestSchemaAsCue := boundKind.Lineage().Latest().Underlying()
	latestSchemaAsCue := rawLatestSchemaAsCue.LookupPath(cue.MakePath(cue.Hid("_#schema", "github.com/grafana/thema")))

	schemaAst, err := simplecue.GenerateAST(latestSchemaAsCue, simplecue.Config{
		Package: pkg, // TODO: extract from input schema/folder?
	})
	if err != nil {
		panic(err)
	}

	// Here begins the code generation setup
	targetsByLanguage := jennies.All()
	rootCodeJenFS := codejen.NewFS()

	for language, target := range targetsByLanguage {
		fmt.Printf("Running '%s' jennies...\n", language)

		var err error
		processedAst := []*ast.File{schemaAst}

		for _, compilerPass := range target.CompilerPasses {
			processedAst, err = compilerPass.Process(processedAst)
			if err != nil {
				panic(err)
			}
		}

		targetFs, err := target.Jennies.GenerateFS(processedAst)
		if err != nil {
			panic(err)
		}

		err = rootCodeJenFS.Merge(targetFs)
		if err != nil {
			panic(err)
		}
	}

	err = rootCodeJenFS.Write(context.Background(), "generated")
	if err != nil {
		panic(err)
	}
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
