package loaders

import (
	"fmt"
	"path/filepath"
	"strings"

	"cuelang.org/go/cue/cuecontext"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/simplecue"
	"github.com/grafana/kindsys"
	"github.com/grafana/thema"
)

func kindsysComposableLoader(opts Options) ([]*ast.Schema, error) {
	themaRuntime := thema.NewRuntime(cuecontext.New())

	allSchemas := make([]*ast.Schema, 0, len(opts.KindsysComposableEntrypoints))
	for _, entrypoint := range opts.KindsysComposableEntrypoints {
		pkg := filepath.Base(entrypoint)

		overlayFS, err := buildKindsysEntrypointFS(opts, entrypoint)
		if err != nil {
			return nil, err
		}

		cueInstance, err := kindsys.BuildInstance(themaRuntime.Context(), ".", "grafanaplugin", overlayFS)
		if err != nil {
			return nil, fmt.Errorf("could not load kindsys composable kind %s: %w", pkg, err)
		}

		props, err := kindsys.ToKindProps[kindsys.ComposableProperties](cueInstance)
		if err != nil {
			return nil, fmt.Errorf("could not convert cue value to kindsys composable props: %w", err)
		}

		kindDefinition := kindsys.Def[kindsys.ComposableProperties]{
			V:          cueInstance,
			Properties: props,
		}

		boundKind, err := kindsys.BindComposable(themaRuntime, kindDefinition)
		if err != nil {
			return nil, fmt.Errorf("could not bind kind definition to kind: %w", err)
		}

		variant, err := schemaVariant(props)
		if err != nil {
			return nil, err
		}

		schemaAst, err := simplecue.GenerateAST(kindToLatestSchema(boundKind), simplecue.Config{
			Package:              pkg, // TODO: extract from input schema/folder?
			ForceVariantEnvelope: variant == ast.SchemaVariantDataQuery,
			SchemaMetadata: ast.SchemaMeta{
				Kind:       ast.SchemaKindComposable,
				Variant:    variant,
				Identifier: inferComposableKindIdentifier(props),
			},
		})
		if err != nil {
			return nil, err
		}

		allSchemas = append(allSchemas, schemaAst)
	}

	return allSchemas, nil
}

func schemaVariant(kindProps kindsys.ComposableProperties) (ast.SchemaVariant, error) {
	switch kindProps.SchemaInterface {
	case "PanelCfg":
		return ast.SchemaVariantPanel, nil
	case "DataQuery":
		return ast.SchemaVariantDataQuery, nil
	default:
		return "", fmt.Errorf("unknown schema variant '%s'", kindProps.SchemaInterface)
	}
}

// TODO: the schema should explicitly tell us the "plugin ID"/panel ID/dataquery type/...
func inferComposableKindIdentifier(kindProps kindsys.ComposableProperties) string {
	lowerVariant := strings.ToLower(kindProps.SchemaInterface)

	return strings.TrimSuffix(kindProps.MachineName, lowerVariant)
}
