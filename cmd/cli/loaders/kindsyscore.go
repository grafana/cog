package loaders

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/simplecue"
	"github.com/grafana/kindsys"
	"github.com/grafana/thema"
	"github.com/yalue/merged_fs"
)

func kindsysCoreLoader(opts Options) ([]*ast.Schema, error) {
	themaRuntime := thema.NewRuntime(cuecontext.New())

	allSchemas := make([]*ast.Schema, 0, len(opts.KindsysCoreEntrypoints))
	for _, entrypoint := range opts.KindsysCoreEntrypoints {
		pkg := filepath.Base(entrypoint)

		overlayFS, err := buildKindsysEntrypointFS(opts, entrypoint)
		if err != nil {
			return nil, err
		}

		cueInstance, err := kindsys.BuildInstance(themaRuntime.Context(), ".", "kind", overlayFS)
		if err != nil {
			return nil, fmt.Errorf("could not load kindsys composable kind %s: %w", pkg, err)
		}

		props, err := kindsys.ToKindProps[kindsys.CoreProperties](cueInstance)
		if err != nil {
			return nil, fmt.Errorf("could not convert cue value to kindsys core props: %w", err)
		}

		kindDefinition := kindsys.Def[kindsys.CoreProperties]{
			V:          cueInstance,
			Properties: props,
		}

		boundKind, err := kindsys.BindCore(themaRuntime, kindDefinition)
		if err != nil {
			return nil, fmt.Errorf("could not bind kind definition to kind: %w", err)
		}

		schemaAst, err := simplecue.GenerateAST(kindToLatestSchema(boundKind), simplecue.Config{
			Package: pkg, // TODO: extract from input schema/folder?
			SchemaMetadata: ast.SchemaMeta{
				Kind:       ast.SchemaKindCore,
				Identifier: pkg, // TODO: maybe even core kinds could have one explicitly set in their schema?
			},
		})
		if err != nil {
			return nil, err
		}

		allSchemas = append(allSchemas, schemaAst)
	}

	return allSchemas, nil
}

func buildKindsysEntrypointFS(opts Options, entrypoint string) (fs.FS, error) {
	libFs, err := buildBaseFSWithLibraries(opts)
	if err != nil {
		return nil, err
	}

	overlayFS, err := dirToPrefixedFS(entrypoint, "")
	if err != nil {
		return nil, err
	}

	return merged_fs.MergeMultiple(libFs, overlayFS), nil
}

func kindToLatestSchema(kind kindsys.Kind) cue.Value {
	rawLatestSchemaAsCue := kind.Lineage().Latest().Underlying()

	return rawLatestSchemaAsCue.LookupPath(cue.MakePath(cue.Hid("_#schema", "github.com/grafana/thema")))
}
