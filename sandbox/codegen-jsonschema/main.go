package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies"
	"github.com/grafana/cog/internal/jsonschema"
)

func main() {
	entrypoints := []string{
		"./schemas/jsonschema/core/playlist/playlist.json",
		"./schemas/jsonschema/core/dockerd/dockerd.json",
	}

	allSchemas := make([]*ast.File, 0, len(entrypoints))
	for _, entrypoint := range entrypoints {
		pkg := filepath.Base(filepath.Dir(entrypoint))

		reader, err := os.Open(entrypoint)
		if err != nil {
			panic(err)
		}

		schemaAst, err := jsonschema.GenerateAST(reader, jsonschema.Config{
			Package: pkg, // TODO: extract from input schema/folder?
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

	err := rootCodeJenFS.Write(context.Background(), "generated")
	if err != nil {
		panic(err)
	}
}
