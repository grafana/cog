package inspect

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/ast/compiler"
	"github.com/grafana/cog/internal/codegen"
	"github.com/grafana/cog/internal/languages"
	"github.com/grafana/cog/internal/tools"
	"github.com/spf13/cobra"
)

type options struct {
	IRType          string
	ConfigPath      string
	ExtraParameters map[string]string
	Selector        string
	Language        string
}

func Command() *cobra.Command {
	opts := options{}

	cmd := &cobra.Command{
		Use:   "inspect",
		Short: "Inspects the intermediate representation.",
		Long: `Inspects the intermediate representation of types and builders.

Common and schema-specific transformations are applied.
Language-specific transformations are NOT applied until a language is specified with the --language flag.

Common builder transformations are applied.
Language-specific builder transformations are NOT applied until a language is specified with the --language flag.
`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if opts.IRType != "types" && opts.IRType != "builders" && opts.IRType != "converters" {
				return fmt.Errorf("invalid IR type '%s'. Valid values: types, builders, converters", opts.IRType)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return doInspect(opts)
		},
	}

	cmd.Flags().StringVar(&opts.IRType, "ir", "types", "Type of intermediate representation to Inspect. Valid values: types, builders, converters.")
	cmd.Flags().StringToStringVar(&opts.ExtraParameters, "parameters", nil, "Sets or overrides parameters used in the config file.")
	cmd.Flags().StringVar(&opts.Selector, "selector", "", "Selector allowing to narrow down the result of the inspection to selected objects. Format: package.[object] for types, package.[builder].[option] for builders.")
	cmd.Flags().StringVar(&opts.Language, "language", "", "Language to use when applying language-specific schema and builder transformations. If left empty, only common transformations are applied.")

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

	// Create a codegen context
	schemas, err := pipeline.LoadSchemas(context.Background())
	if err != nil {
		return err
	}

	language, err := inspectedLanguage(pipeline, opts)
	if err != nil {
		return err
	}

	codegenCtx, err := pipeline.ContextForLanguage(language, schemas)
	if err != nil {
		return err
	}

	// Inspect the desired IR
	if opts.IRType == "types" {
		return inspectTypesIR(codegenCtx, opts.Selector)
	}

	if opts.IRType == "builders" {
		return inspectBuildersIR(codegenCtx, opts.Selector)
	}

	return inspectConvertersIR(codegenCtx, language, opts.Selector)
}

func inspectedLanguage(pipeline *codegen.Pipeline, opts options) (languages.Language, error) {
	if opts.Language == "" {
		return dummyLanguage{}, nil
	}

	languagesMap, err := pipeline.OutputLanguages()
	if err != nil {
		return nil, err
	}

	language := languagesMap[opts.Language]
	if language == nil {
		return nil, fmt.Errorf("language \"%s\" is not supported. Valid languages are: %s", opts.Language, strings.Join(tools.Keys(languagesMap), ", "))
	}

	return language, nil
}

func inspectTypesIR(codegenCtx languages.Context, selector string) error {
	selectedResult, err := applyTypesIRSelector(codegenCtx, selector)
	if err != nil {
		return err
	}

	return prettyPrintJSON(selectedResult)
}

func inspectBuildersIR(codegenCtx languages.Context, selector string) error {
	selectedResult, err := applyBuilderIRSelector(codegenCtx, selector)
	if err != nil {
		return err
	}

	return prettyPrintJSON(selectedResult)
}

func inspectConvertersIR(codegenCtx languages.Context, language languages.Language, selector string) error {
	selectorParts := strings.Split(selector, ".")
	if len(selectorParts) != 2 {
		return fmt.Errorf("inspecting the converters IR requires a builder selector: [package].[builder]")
	}

	// select builders within a package
	builders := tools.Filter(codegenCtx.Builders, func(b ast.Builder) bool {
		return strings.EqualFold(b.Package, selectorParts[0]) && strings.EqualFold(b.Name, selectorParts[1])
	})
	if len(builders) == 0 {
		return fmt.Errorf("builder '%s' not found", selector)
	}

	nullableConfig, ok := language.(languages.NullableKindsProvider)
	if !ok {
		return fmt.Errorf("language '%s' does not appear to support converters", language.Name())
	}

	converter := languages.NewConverterGenerator(nullableConfig.NullableKinds(), codegenCtx.ConverterConfig).FromBuilder(codegenCtx, builders[0])

	return prettyPrintJSON(converter)
}

func applyTypesIRSelector(codegenCtx languages.Context, selector string) (any, error) {
	schemas := codegenCtx.Schemas

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

func prettyPrintJSON(input any) error {
	marshaled, err := json.MarshalIndent(input, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(marshaled))

	return nil
}

type dummyLanguage struct {
}

func (language dummyLanguage) Name() string {
	return "dummy"
}

func (language dummyLanguage) Jennies(_ languages.Config) *codejen.JennyList[languages.Context] {
	return nil
}

func (language dummyLanguage) CompilerPasses() compiler.Passes {
	return nil
}
