package loaders

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jsonschema"
)

func jsonschemaRegistryLoader(opts Options) ([]*ast.Schema, error) {
	var allSchemas []*ast.Schema

	if opts.JSONSchemaRegistryPath == "" {
		return nil, nil
	}

	coreEntrypoints, err := locateJSONEntrypoints(opts, "core")
	if err != nil {
		return nil, fmt.Errorf("could not locate core entrypoints: %w", err)
	}

	panelsEntrypoints, err := locateJSONEntrypoints(opts, "panels")
	if err != nil {
		return nil, fmt.Errorf("could not locate panel entrypoints: %w", err)
	}

	dataqueriesEntrypoints, err := locateJSONEntrypoints(opts, "dataqueries")
	if err != nil {
		return nil, fmt.Errorf("could not locate dataqueries entrypoints: %w", err)
	}

	coreSchemas, err := loadJSONEntryPoints(coreEntrypoints, ast.SchemaMeta{
		Kind: ast.SchemaKindCore,
	})
	if err != nil {
		return nil, err
	}

	panelSchemas, err := loadJSONEntryPoints(panelsEntrypoints, ast.SchemaMeta{
		Kind:    ast.SchemaKindComposable,
		Variant: ast.SchemaVariantPanel,
	})
	if err != nil {
		return nil, err
	}

	dataqueriesSchemas, err := loadJSONEntryPoints(dataqueriesEntrypoints, ast.SchemaMeta{
		Kind:    ast.SchemaKindComposable,
		Variant: ast.SchemaVariantDataQuery,
	})
	if err != nil {
		return nil, err
	}

	allSchemas = append(allSchemas, coreSchemas...)
	allSchemas = append(allSchemas, panelSchemas...)
	allSchemas = append(allSchemas, dataqueriesSchemas...)

	return allSchemas, nil
}

func loadJSONEntryPoints(entrypoints []string, schemaMeta ast.SchemaMeta) ([]*ast.Schema, error) {
	allSchemas := make([]*ast.Schema, 0, len(entrypoints))
	for _, entrypoint := range entrypoints {
		reader, err := os.Open(entrypoint)
		if err != nil {
			return nil, err
		}

		pkg := guessPackageFromFilename(entrypoint)

		meta := schemaMeta
		meta.Identifier = pkg

		schemaAst, err := jsonschema.GenerateAST(reader, jsonschema.Config{
			Package:        pkg,
			SchemaMetadata: meta,
		})
		if err != nil {
			return nil, err
		}

		allSchemas = append(allSchemas, schemaAst)
	}

	return allSchemas, nil
}

func locateJSONEntrypoints(opts Options, kind string) ([]string, error) {
	directory := filepath.Join(opts.JSONSchemaRegistryPath, kind)

	var results []string
	err := filepath.WalkDir(directory, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		results = append(results, path)

		return nil
	})
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return results, nil
}
