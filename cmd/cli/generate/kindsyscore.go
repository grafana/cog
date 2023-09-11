package generate

import (
	"fmt"
	"path/filepath"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/simplecue"
	"github.com/grafana/kindsys"
	"github.com/grafana/thema"
)

func kindsysCoreLoader(opts options) ([]*ast.File, error) {
	themaRuntime := thema.NewRuntime(cuecontext.New())

	allSchemas := make([]*ast.File, 0, len(opts.entrypoints))
	for _, entrypoint := range opts.entrypoints {
		pkg := filepath.Base(entrypoint)

		overlayFS, err := dirToPrefixedFS(entrypoint, "")
		if err != nil {
			return nil, err
		}

		cueInstance, err := kindsys.BuildInstance(themaRuntime.Context(), ".", "kind", overlayFS)
		if err != nil {
			return nil, fmt.Errorf("could not load kindsys instance: %w", err)
		}

		props, err := kindsys.ToKindProps[kindsys.CoreProperties](cueInstance)
		if err != nil {
			return nil, fmt.Errorf("could not convert cue value to kindsys props: %w", err)
		}

		kindDefinition := kindsys.Def[kindsys.CoreProperties]{
			V:          cueInstance,
			Properties: props,
		}

		boundKind, err := kindsys.BindCore(themaRuntime, kindDefinition)
		if err != nil {
			return nil, fmt.Errorf("could not bind kind definition to kind: %w", err)
		}

		rawLatestSchemaAsCue := boundKind.Lineage().Latest().Underlying()
		latestSchemaAsCue := rawLatestSchemaAsCue.LookupPath(cue.MakePath(cue.Hid("_#schema", "github.com/grafana/thema")))

		schemaAst, err := simplecue.GenerateAST(latestSchemaAsCue, simplecue.Config{
			Package: pkg, // TODO: extract from input schema/folder?
		})
		if err != nil {
			return nil, err
		}

		allSchemas = append(allSchemas, schemaAst)
	}

	return allSchemas, nil
}
