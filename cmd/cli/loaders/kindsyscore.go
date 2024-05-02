package loaders

import (
	"path/filepath"

	"cuelang.org/go/cue"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/simplecue"
	"github.com/grafana/cog/internal/tools"
)

func kindsysCoreLoader(config Config, input CueInput) (ast.Schemas, error) {
	libraries, err := simplecue.ParseImports(input.CueImports)
	if err != nil {
		return nil, err
	}

	libraries = tools.Map(libraries, func(library simplecue.LibraryInclude) simplecue.LibraryInclude {
		library.FSPath = config.Path(library.FSPath)
		return library
	})

	schemaRootValue, err := parseCueEntrypoint(config.Path(input.Entrypoint), libraries, "kind")
	if err != nil {
		return nil, err
	}

	kindIdentifier, err := inferCoreKindIdentifier(schemaRootValue)
	if err != nil {
		return nil, err
	}

	schema, err := simplecue.GenerateAST(schemaFromThemaLineage(schemaRootValue), simplecue.Config{
		Package: filepath.Base(input.Entrypoint), // TODO: extract from somewhere else?
		SchemaMetadata: ast.SchemaMeta{
			Kind:       ast.SchemaKindCore,
			Identifier: kindIdentifier,
		},
		Libraries: libraries,
	})
	if err != nil {
		return nil, err
	}

	return ast.Schemas{schema}, nil
}

func inferCoreKindIdentifier(kindRoot cue.Value) (string, error) {
	return kindRoot.LookupPath(cue.ParsePath("name")).String()
}

func schemaFromThemaLineage(kindRoot cue.Value) cue.Value {
	return kindRoot.LookupPath(cue.ParsePath("lineage.schemas[0].schema"))
}
