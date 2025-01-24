package inspect

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/codegen"
	"github.com/grafana/cog/internal/languages"
	"github.com/grafana/cog/internal/tools"
	"github.com/spf13/cobra"
)

type options struct {
	BuilderIR       bool
	ConfigPath      string
	ExtraParameters map[string]string
	Selector        string
}

func Command() *cobra.Command {
	opts := options{}

	// TODO:
	//  - better support for inspecting "transformed" IRs: language-specific compiler passes, veneers, ...
	cmd := &cobra.Command{
		Use:   "inspect",
		Short: "Inspects the intermediate representation.",
		Long: `Inspects the intermediate representation of types and builders.

Common and schema-specific transformations are applied.
Language-specific transformations are NOT applied.

Builder transformations are currently NOT applied.
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return doInspect(opts)
		},
	}

	cmd.Flags().BoolVar(&opts.BuilderIR, "builder-ir", false, "Inspect the \"builder IR\" instead of the \"types\" one.") // TODO: better usage text
	cmd.Flags().StringToStringVar(&opts.ExtraParameters, "parameters", nil, "Sets or overrides parameters used in the config file.")
	cmd.Flags().StringVar(&opts.Selector, "selector", "", "Selector allowing to narrow down the result of the inspection to selected objects. Format: package.[object] for types, package.[builder].[option] for builders.")

	cmd.Flags().StringVar(&opts.ConfigPath, "config", "", "Codegen pipeline configuration file.")
	_ = cmd.MarkFlagFilename("config")
	_ = cmd.MarkFlagRequired("config")

	return cmd
}

func doInspect(opts options) error {
	pipeline, err := codegen.PipelineFromFile(opts.ConfigPath, codegen.Parameters(opts.ExtraParameters))
	if err != nil {
		return err
	}

	// Load schemas
	schemas, err := pipeline.LoadSchemas(context.Background())
	if err != nil {
		return err
	}

	// apply compiler passes

	if opts.BuilderIR {
		return inspectBuilderIR(opts, schemas)
	}

	selectedResult, err := applyTypesIRSelector(schemas, opts.Selector)
	if err != nil {
		return err
	}

	return prettyPrintJSON(selectedResult)
}

func inspectBuilderIR(opts options, schemas []*ast.Schema) error {
	var err error
	codegenCtx := languages.Context{
		Schemas:  schemas,
		Builders: (&ast.BuilderGenerator{}).FromAST(schemas),
	}

	codegenCtx, err = languages.GenerateBuilderNilChecks(nil, codegenCtx)
	if err != nil {
		return err
	}

	selectedResult, err := applyBuilderIRSelector(codegenCtx, opts.Selector)
	if err != nil {
		return err
	}

	return prettyPrintJSON(selectedResult)
}

func prettyPrintJSON(input any) error {
	marshaled, err := json.MarshalIndent(input, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(marshaled))

	return nil
}

func applyTypesIRSelector(schemas ast.Schemas, selector string) (any, error) {
	if selector == "" {
		return schemas, nil
	}

	selectorParts := strings.Split(selector, ".")
	if len(selectorParts) > 2 {
		return nil, fmt.Errorf("invalid selector '%s'", selector)
	}

	// select a package
	schema, found := schemas.Locate(selectorParts[0])
	if !found {
		return nil, fmt.Errorf("package '%s' not found", selectorParts[0])
	}
	if len(selectorParts) == 1 {
		return schema, nil
	}

	// select a specific object
	objects := schema.Objects.Filter(func(_ string, object ast.Object) bool {
		return strings.EqualFold(object.Name, selectorParts[1])
	})
	if objects.Len() == 0 {
		return nil, fmt.Errorf("object '%s.%s' not found", selectorParts[0], selectorParts[1])
	}

	return objects.At(0), nil
}

func applyBuilderIRSelector(context languages.Context, selector string) (any, error) {
	if selector == "" {
		return context, nil
	}

	selectorParts := strings.Split(selector, ".")
	if len(selectorParts) > 3 {
		return nil, fmt.Errorf("invalid selector '%s'", selector)
	}

	// select builders within a package
	builders := tools.Filter(context.Builders, func(b ast.Builder) bool {
		return strings.EqualFold(b.Package, selectorParts[0])
	})
	if len(builders) == 0 {
		return nil, fmt.Errorf("package '%s' not found", selectorParts[0])
	}
	if len(selectorParts) == 1 {
		return builders, nil
	}

	// target a specific builder
	builders = tools.Filter(builders, func(builder ast.Builder) bool {
		return strings.EqualFold(builder.Name, selectorParts[1])
	})
	if len(builders) == 0 {
		return nil, fmt.Errorf("builder '%s.%s' not found", selectorParts[0], selectorParts[1])
	}
	if len(selectorParts) == 2 {
		return builders[0], nil
	}

	// select a specific option within a builder
	opts := tools.Filter(builders[0].Options, func(opt ast.Option) bool {
		return strings.EqualFold(opt.Name, selectorParts[2])
	})
	if len(opts) == 0 {
		return nil, fmt.Errorf("option '%s' not found", selector)
	}

	return opts[0], nil
}
