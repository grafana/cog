package golang

import (
	"fmt"
	"strings"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/jennies/template"
	"github.com/grafana/cog/internal/languages"
	"github.com/grafana/cog/internal/tools"
)

type JSONMarshalling struct {
	tmpl            *template.Template
	config          Config
	imports         *common.DirectImportMap
	packageMapper   func(string) string
	typeFormatter   *typeFormatter
	apiRefCollector *common.APIReferenceCollector
}

func NewJSONMarshalling(config Config, tmpl *template.Template, imports *common.DirectImportMap, packageMapper func(string) string, typeFormatter *typeFormatter, apiRefCollector *common.APIReferenceCollector) JSONMarshalling {
	return JSONMarshalling{
		config: config,
		tmpl: tmpl.Funcs(template.FuncMap{
			"formatType": typeFormatter.formatType,
			"importStdPkg": func(pkg string) string {
				return imports.Add(pkg, pkg)
			},
		}),
		imports:         imports,
		packageMapper:   packageMapper,
		typeFormatter:   typeFormatter,
		apiRefCollector: apiRefCollector,
	}
}

func (jenny JSONMarshalling) generateForSchema(buffer *strings.Builder, schema *ast.Schema) error {
	if schema.Metadata.Kind != ast.SchemaKindComposable || schema.Metadata.Variant != ast.SchemaVariantPanel {
		return nil
	}

	variantUnmarshal, err := jenny.renderPanelcfgVariantUnmarshal(schema)
	if err != nil {
		return err
	}
	buffer.WriteString(variantUnmarshal)

	return nil
}

func (jenny JSONMarshalling) generateForObject(buffer *strings.Builder, context languages.Context, schema *ast.Schema, object ast.Object) error {
	if jenny.objectNeedsCustomMarshal(object) {
		jsonMarshal, err := jenny.renderCustomMarshal(object)
		if err != nil {
			return err
		}
		buffer.WriteString(jsonMarshal)
		buffer.WriteString("\n")
	}

	if jenny.objectNeedsCustomUnmarshal(context, object) {
		jsonUnmarshal, err := jenny.renderCustomUnmarshal(context, object)
		if err != nil {
			return err
		}
		buffer.WriteString(jsonUnmarshal)
		buffer.WriteString("\n")
	}

	if object.Type.ImplementedVariant() == string(ast.SchemaVariantDataQuery) && !object.Type.HasHint(ast.HintSkipVariantPluginRegistration) && !object.Type.IsStructGeneratedFromDisjunction() {
		variantUnmarshal, err := jenny.renderDataqueryVariantUnmarshal(schema, object)
		if err != nil {
			return err
		}
		buffer.WriteString(variantUnmarshal)
		buffer.WriteString("\n")
	}

	return nil
}

func (jenny JSONMarshalling) objectNeedsCustomMarshal(obj ast.Object) bool {
	// the only case for which we need a custom marshaller is for structs
	// that are generated from a disjunction by the `DisjunctionToType` compiler pass.

	return obj.Type.IsStructGeneratedFromDisjunction()
}

func (jenny JSONMarshalling) renderCustomMarshal(obj ast.Object) (string, error) {
	jenny.apiRefCollector.RegisterMethod(obj, common.MethodReference{
		Name: "MarshalJSON",
		Comments: []string{
			fmt.Sprintf("MarshalJSON implements a custom JSON marshalling logic to encode `%s` as JSON.", tools.UpperCamelCase(obj.Name)),
		},
		Return: "([]byte, error)",
	})

	// There are only two types of disjunctions we support:
	//  * undiscriminated: string | bool | ..., where all the disjunction branches are scalars (or an array)
	//  * discriminated: SomeStruct | SomeOtherStruct, where all the disjunction branches are references to
	// 	  structs and these structs have a common "discriminator" field.
	if obj.Type.IsStruct() && obj.Type.HasHint(ast.HintDisjunctionOfScalars) {
		return jenny.tmpl.Render("types/disjunction_of_scalars.json_marshal.tmpl", map[string]any{
			"def": obj,
		})
	}

	if obj.Type.IsStruct() && obj.Type.HasHint(ast.HintDiscriminatedDisjunctionOfRefs) {
		return jenny.tmpl.Render("types/disjunction_of_refs.json_marshal.tmpl", map[string]any{
			"def": obj,
		})
	}

	return "", fmt.Errorf("could not determine how to render custom marshal")
}

func (jenny JSONMarshalling) objectNeedsCustomUnmarshal(context languages.Context, obj ast.Object) bool {
	// an object needs a custom unmarshal if:
	// - it is a struct that was generated from a disjunction by the `DisjunctionToType` compiler pass.
	// - it is a struct and one or more of its fields is a KindComposableSlot, or an array of KindComposableSlot

	if !obj.Type.IsStruct() {
		return false
	}

	// is there a custom unmarshal template block?
	if jenny.tmpl.Exists(template.CustomObjectUnmarshalBlock(obj)) {
		return true
	}

	// is it a struct generated from a disjunction?
	if obj.Type.IsStructGeneratedFromDisjunction() {
		return true
	}

	// is there a KindComposableSlot field somewhere?
	for _, field := range obj.Type.AsStruct().Fields {
		if _, ok := context.ResolveToComposableSlot(field.Type); ok {
			return true
		}
	}

	return false
}

func (jenny JSONMarshalling) renderCustomUnmarshal(context languages.Context, obj ast.Object) (string, error) {
	jenny.apiRefCollector.RegisterMethod(obj, common.MethodReference{
		Name: "UnmarshalJSON",
		Arguments: []common.ArgumentReference{
			{Name: "raw", Type: "[]byte"},
		},
		Comments: []string{
			fmt.Sprintf("UnmarshalJSON implements a custom JSON unmarshalling logic to decode `%s` from JSON.", tools.UpperCamelCase(obj.Name)),
		},
		Return: "error",
	})

	customUnmarshalTmpl := template.CustomObjectUnmarshalBlock(obj)
	if jenny.tmpl.Exists(customUnmarshalTmpl) {
		return jenny.tmpl.Render(customUnmarshalTmpl, map[string]any{
			"Object": obj,
		})
	}

	if obj.Type.IsStruct() && obj.Type.HasHint(ast.HintDisjunctionOfScalars) {
		return jenny.tmpl.Render("types/disjunction_of_scalars.json_unmarshal.tmpl", map[string]any{
			"def": obj,
		})
	}

	if obj.Type.IsStruct() && obj.Type.HasHint(ast.HintDiscriminatedDisjunctionOfRefs) {
		return jenny.tmpl.Render("types/disjunction_of_refs.json_unmarshal.tmpl", map[string]any{
			"def":  obj,
			"hint": obj.Type.Hints[ast.HintDiscriminatedDisjunctionOfRefs],
		})
	}

	return jenny.renderCustomComposableSlotUnmarshal(context, obj)
}

func (jenny JSONMarshalling) renderCustomComposableSlotUnmarshal(context languages.Context, obj ast.Object) (string, error) {
	var buffer strings.Builder
	fields := obj.Type.AsStruct().Fields

	jenny.packageMapper("cog")

	// unmarshal "normal" fields (ie: with no composable slot)
	for _, field := range fields {
		if _, ok := context.ResolveToComposableSlot(field.Type); ok {
			continue
		}

		jenny.imports.Add("fmt", "fmt")
		buffer.WriteString(fmt.Sprintf(`
	if fields["%[1]s"] != nil {
		if err := json.Unmarshal(fields["%[1]s"], &resource.%[2]s); err != nil {
			return fmt.Errorf("error decoding field '%[1]s': %%w", err)
		}
	}
`, field.Name, tools.UpperCamelCase(field.Name)))
	}

	// unmarshal "composable slot" fields
	for _, field := range fields {
		composableSlotType, resolved := context.ResolveToComposableSlot(field.Type)
		if !resolved {
			continue
		}

		variant := composableSlotType.AsComposableSlot().Variant

		switch variant {
		case ast.SchemaVariantDataQuery:
			source := jenny.renderUnmarshalDataqueryField(obj, field)
			buffer.WriteString(source)
		default:
			return "", fmt.Errorf("can not generate custom unmarshal function for composable slot with variant '%s'", variant)
		}
	}

	jenny.imports.Add("json", "encoding/json")

	return fmt.Sprintf(`// UnmarshalJSON implements a custom JSON unmarshalling logic to decode %[1]s from JSON.
func (resource *%[1]s) UnmarshalJSON(raw []byte) error {
	if raw == nil {
		return nil
	}
	fields := make(map[string]json.RawMessage)
	if err := json.Unmarshal(raw, &fields); err != nil {
		return err
	}
	%[2]s
	return nil
}
`, tools.UpperCamelCase(obj.Name), buffer.String()), nil
}

func (jenny JSONMarshalling) renderUnmarshalDataqueryField(parentStruct ast.Object, field ast.StructField) string {
	// First: try to locate a field that would contain the type of datasource being used.
	// We're looking for a field defined as a reference to the `DataSourceRef` type.
	var hintField *ast.StructField
	for i, candidate := range parentStruct.Type.AsStruct().Fields {
		if !candidate.Type.IsRef() {
			continue
		}
		if candidate.Type.AsRef().ReferredType != "DataSourceRef" {
			continue
		}

		hintField = &parentStruct.Type.AsStruct().Fields[i]
	}

	// then: unmarshalling boilerplate
	hintValue := `dataqueryTypeHint := ""
`

	if hintField != nil {
		hintValue += fmt.Sprintf(`if resource.%[1]s != nil && resource.%[1]s.Type != nil {
dataqueryTypeHint = *resource.%[1]s.Type
}
`, tools.UpperCamelCase(hintField.Name))
	}

	jenny.packageMapper("cog")

	if field.Type.IsArray() {
		jenny.packageMapper("cog")
		return fmt.Sprintf(`
	%[3]s
	if fields["%[2]s"] != nil {
		%[2]s, err := cog.UnmarshalDataqueryArray(fields["%[2]s"], dataqueryTypeHint)
		if err != nil {
			return err
		}
		resource.%[1]s = %[2]s
	}
`, tools.UpperCamelCase(field.Name), field.Name, hintValue)
	}

	return fmt.Sprintf(`
	%[3]s
	if fields["%[2]s"] != nil {
		%[2]s, err := cog.UnmarshalDataquery(fields["%[2]s"], dataqueryTypeHint)
		if err != nil {
			return err
		}
		resource.%[1]s = %[2]s
	}
`, tools.UpperCamelCase(field.Name), field.Name, hintValue)
}

func (jenny JSONMarshalling) renderPanelcfgVariantUnmarshal(schema *ast.Schema) (string, error) {
	jenny.packageMapper("cog/variants")

	_, hasOptions := schema.LocateObject("Options")
	_, hasFieldConfig := schema.LocateObject("FieldConfig")

	if jenny.config.generateConverters {
		jenny.packageMapper("dashboard")
	}

	jenny.apiRefCollector.RegisterFunction(schema.Package, common.FunctionReference{
		Name: "VariantConfig",
		Comments: []string{
			fmt.Sprintf("VariantConfig returns the configuration related to %s panels.", strings.ToLower(schema.Metadata.Identifier)),
			"This configuration describes how to unmarshal it, convert it to code, …",
		},
		Return: "variants.PanelcfgConfig",
	})

	return jenny.tmpl.Render("types/variant_panelcfg.json_unmarshal.tmpl", map[string]any{
		"schema":         schema,
		"hasOptions":     hasOptions,
		"hasFieldConfig": hasFieldConfig,
		"hasConverter":   jenny.config.generateConverters,
	})
}

func (jenny JSONMarshalling) renderDataqueryVariantUnmarshal(schema *ast.Schema, obj ast.Object) (string, error) {
	jenny.packageMapper("cog/variants")

	jenny.apiRefCollector.RegisterFunction(schema.Package, common.FunctionReference{
		Name: "VariantConfig",
		Comments: []string{
			fmt.Sprintf("VariantConfig returns the configuration related to %s dataqueries.", strings.ToLower(schema.Metadata.Identifier)),
			"This configuration describes how to unmarshal it, convert it to code, …",
		},
		Return: "variants.DataqueryConfig",
	})

	var disjunctionStruct *ast.StructType

	if obj.Type.IsRef() {
		resolved, _ := schema.Resolve(obj.Type)
		if resolved.IsStructGeneratedFromDisjunction() {
			disjunctionStruct = resolved.Struct
		}
	}

	return jenny.tmpl.Render("types/variant_dataquery.json_unmarshal.tmpl", map[string]any{
		"schema":            schema,
		"object":            obj,
		"hasConverter":      jenny.config.generateConverters,
		"disjunctionStruct": disjunctionStruct,
	})
}
