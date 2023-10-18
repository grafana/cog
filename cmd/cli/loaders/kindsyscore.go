package loaders

import (
	"path/filepath"

	"cuelang.org/go/cue"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/simplecue"
)

func kindsysCoreLoader(opts Options) ([]*ast.Schema, error) {
	libraries, err := opts.cueIncludeImports()
	if err != nil {
		return nil, err
	}

	allSchemas := make([]*ast.Schema, 0, len(opts.KindsysCoreEntrypoints))
	for _, entrypoint := range opts.KindsysCoreEntrypoints {
		pkg := filepath.Base(entrypoint)

		schemaRootValue, err := parseCueEntrypoint(opts, entrypoint)
		if err != nil {
			return nil, err
		}

		kindIdentifier, err := inferCoreKindIdentifier(schemaRootValue)
		if err != nil {
			return nil, err
		}

		schemaAst, err := simplecue.GenerateAST(schemaFromThemaLineage(schemaRootValue), simplecue.Config{
			Package: pkg, // TODO: extract from somewhere else?
			SchemaMetadata: ast.SchemaMeta{
				Kind:       ast.SchemaKindCore,
				Identifier: kindIdentifier,
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

func inferCoreKindIdentifier(kindRoot cue.Value) (string, error) {
	return kindRoot.LookupPath(cue.ParsePath("name")).String()
}

func schemaFromThemaLineage(kindRoot cue.Value) cue.Value {
	return kindRoot.LookupPath(cue.ParsePath("lineage.schemas[0].schema"))
}
