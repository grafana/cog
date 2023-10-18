package loaders

import (
	"fmt"
	"path/filepath"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/simplecue"
)

func kindsysComposableLoader(opts Options) ([]*ast.Schema, error) {
	cueFsOverlay, err := buildCueOverlay(opts)
	if err != nil {
		return nil, err
	}

	libraries, err := opts.cueIncludeImports()
	if err != nil {
		return nil, err
	}

	allSchemas := make([]*ast.Schema, 0, len(opts.KindsysComposableEntrypoints))
	for _, entrypoint := range opts.KindsysComposableEntrypoints {
		pkg := filepath.Base(entrypoint)

		// Load Cue files into Cue build.Instances slice
		// the second arg is a configuration object, we'll see this later
		bis := load.Instances([]string{entrypoint}, &load.Config{
			Overlay:    cueFsOverlay,
			ModuleRoot: "/",
		})

		values, err := cuecontext.New().BuildInstances(bis)
		if err != nil {
			return nil, err
		}

		schemaRoot := values[0]

		variant, err := schemaVariant(schemaRoot)
		if err != nil {
			return nil, err
		}

		kindIdentifier, err := inferComposableKindIdentifier(schemaRoot)
		if err != nil {
			return nil, err
		}

		schemaAsCueValue := schemaRoot.LookupPath(cue.ParsePath("lineage.schemas[0].schema"))

		schemaAst, err := simplecue.GenerateAST(schemaAsCueValue, simplecue.Config{
			Package:              pkg, // TODO: extract from somewhere else?
			ForceVariantEnvelope: variant == ast.SchemaVariantDataQuery,
			SchemaMetadata: ast.SchemaMeta{
				Kind:       ast.SchemaKindComposable,
				Variant:    variant,
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

func schemaVariant(kindRoot cue.Value) (ast.SchemaVariant, error) {
	schemaInterface, err := kindRoot.LookupPath(cue.ParsePath("schemaInterface")).String()
	if err != nil {
		return "", err
	}

	switch schemaInterface {
	case "PanelCfg":
		return ast.SchemaVariantPanel, nil
	case "DataQuery":
		return ast.SchemaVariantDataQuery, nil
	default:
		return "", fmt.Errorf("unknown schema variant '%s'", schemaInterface)
	}
}

// TODO: the schema should explicitly tell us the "plugin ID"/panel ID/dataquery type/...
func inferComposableKindIdentifier(kindRoot cue.Value) (string, error) {
	schemaInterface, err := kindRoot.LookupPath(cue.ParsePath("schemaInterface")).String()
	if err != nil {
		return "", err
	}

	kindName, err := kindRoot.LookupPath(cue.ParsePath("name")).String()
	if err != nil {
		return "", err
	}

	fmt.Printf("kindName: '%s'\n", kindName)
	fmt.Printf("identifier: '%s'\n", strings.TrimSuffix(kindName, schemaInterface))

	return strings.TrimSuffix(kindName, schemaInterface), nil
}
