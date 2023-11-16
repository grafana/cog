package golang

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/context"
	"github.com/grafana/cog/internal/jennies/template"
	"github.com/grafana/cog/internal/tools"
)

type JSONMarshalling struct {
	Config Config

	packageMapper func(string) string
	typeFormatter *typeFormatter
}

func (jenny JSONMarshalling) JennyName() string {
	return "GoJSONMarshalling"
}

func (jenny JSONMarshalling) Generate(context context.Builders) (codejen.Files, error) {
	files := make(codejen.Files, 0, len(context.Schemas))

	for _, schema := range context.Schemas {
		output, err := jenny.generateSchema(context, schema)
		if err != nil {
			return nil, err
		}
		if output == nil {
			continue
		}

		filename := filepath.Join(
			strings.ToLower(schema.Package),
			"types_json_marshalling_gen.go",
		)

		files = append(files, *codejen.NewFile(filename, output, jenny))
	}

	return files, nil
}

func (jenny JSONMarshalling) generateSchema(context context.Builders, schema *ast.Schema) ([]byte, error) {
	var buffer strings.Builder

	imports := template.NewImportMap()
	jenny.packageMapper = func(pkg string) string {
		if pkg == schema.Package {
			return ""
		}

		return imports.Add(pkg, jenny.Config.importPath(pkg))
	}
	jenny.typeFormatter = defaultTypeFormatter(jenny.packageMapper)

	for _, object := range schema.Objects {
		if jenny.objectNeedsCustomMarshal(object) {
			jsonMarshal, err := jenny.renderCustomMarshal(object)
			if err != nil {
				return nil, err
			}
			buffer.WriteString(jsonMarshal)
		}

		if jenny.objectNeedsCustomUnmarshal(context, object) {
			jsonUnmarshal, err := jenny.renderCustomUnmarshal(context, object)
			if err != nil {
				return nil, err
			}
			buffer.WriteString(jsonUnmarshal)
		}

		if object.Type.ImplementsVariant() && object.Type.ImplementedVariant() == string(ast.SchemaVariantDataQuery) {
			variantUnmarshal, err := jenny.renderDataqueryVariantUnmarshal(schema, object)
			if err != nil {
				return nil, err
			}
			buffer.WriteString(variantUnmarshal)
		}
	}

	if schema.Metadata.Kind == ast.SchemaKindComposable && schema.Metadata.Variant == ast.SchemaVariantPanel {
		variantUnmarshal, err := jenny.renderPanelcfgVariantUnmarshal(schema)
		if err != nil {
			return nil, err
		}
		buffer.WriteString(variantUnmarshal)
	}

	if buffer.Len() == 0 {
		return nil, nil
	}

	importStatements := formatImports(imports)
	if importStatements != "" {
		importStatements += "\n\n"
	}

	return []byte(fmt.Sprintf(`package %[1]s
%[2]s%[3]s`, strings.ToLower(schema.Package), importStatements, buffer.String())), nil
}

func (jenny JSONMarshalling) objectNeedsCustomMarshal(obj ast.Object) bool {
	// the only case for which we need a custom marshaller is for structs
	// that are generated from a disjunction by the `DisjunctionToType` compiler pass.

	return obj.Type.IsStructGeneratedFromDisjunction()
}

func (jenny JSONMarshalling) renderCustomMarshal(obj ast.Object) (string, error) {
	// There are only two types of disjunctions we support:
	//  * undiscriminated: string | bool | ..., where all the disjunction branches are scalars (or an array)
	//  * discriminated: SomeStruct | SomeOtherStruct, where all the disjunction branches are references to
	// 	  structs and these structs have a common "discriminator" field.
	isStruct := obj.Type.Kind == ast.KindStruct

	if isStruct && obj.Type.HasHint(ast.HintDisjunctionOfScalars) {
		return jenny.renderTemplate("types/disjunction_of_scalars.json_marshal.tmpl", map[string]any{
			"def": obj,
		})
	}

	if isStruct && obj.Type.HasHint(ast.HintDiscriminatedDisjunctionOfRefs) {
		return jenny.renderTemplate("types/disjunction_of_refs.json_marshal.tmpl", map[string]any{
			"def": obj,
		})
	}

	return "", fmt.Errorf("could not determine how to render custom marshal")
}

func (jenny JSONMarshalling) objectNeedsCustomUnmarshal(context context.Builders, obj ast.Object) bool {
	// an object needs a custom unmarshal if:
	// - it is a struct that was generated from a disjunction by the `DisjunctionToType` compiler pass.
	// - it is a struct and one or more of its fields is a KindComposableSlot, or an array of KindComposableSlot

	if obj.Type.Kind != ast.KindStruct {
		return false
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

func (jenny JSONMarshalling) renderCustomUnmarshal(context context.Builders, obj ast.Object) (string, error) {
	isStruct := obj.Type.Kind == ast.KindStruct

	if isStruct && obj.Type.HasHint(ast.HintDisjunctionOfScalars) {
		return jenny.renderTemplate("types/disjunction_of_scalars.json_unmarshal.tmpl", map[string]any{
			"def": obj,
		})
	}

	if isStruct && obj.Type.HasHint(ast.HintDiscriminatedDisjunctionOfRefs) {
		return jenny.renderTemplate("types/disjunction_of_refs.json_unmarshal.tmpl", map[string]any{
			"def":  obj,
			"hint": obj.Type.Hints[ast.HintDiscriminatedDisjunctionOfRefs],
		})
	}

	return jenny.renderCustomComposableSlotUnmarshal(context, obj)
}

func (jenny JSONMarshalling) renderCustomComposableSlotUnmarshal(context context.Builders, obj ast.Object) (string, error) {
	var buffer strings.Builder
	fields := obj.Type.AsStruct().Fields

	jenny.packageMapper("cog")

	// unmarshal "normal" fields (ie: with no composable slot)
	for _, field := range fields {
		if _, ok := context.ResolveToComposableSlot(field.Type); ok {
			continue
		}

		if obj.SelfRef.ReferredPkg == "dashboard" && obj.Name == "Panel" && field.Name == "options" {
			buffer.WriteString(fmt.Sprintf(`
	if fields["%[1]s"] != nil {
		variantCfg, found := cog.ConfigForPanelcfgVariant(resource.Type)
		if found && variantCfg.OptionsUnmarshaler != nil {
			options, err := variantCfg.OptionsUnmarshaler(fields["%[1]s"])
			if err != nil {
				return err
			}
			resource.%[2]s = options
		} else {
			if err := json.Unmarshal(fields["%[1]s"], &resource.%[2]s); err != nil {
				return err
			}
		}
	}
`, field.Name, tools.UpperCamelCase(field.Name)))
			continue
		}

		if obj.SelfRef.ReferredPkg == "dashboard" && obj.Name == "Panel" && field.Name == "fieldConfig" {
			buffer.WriteString(fmt.Sprintf(`
	if fields["%[1]s"] != nil {
		variantCfg, found := cog.ConfigForPanelcfgVariant(resource.Type)
		if found && variantCfg.FieldConfigUnmarshaler != nil {
			fakeFieldConfigSource := struct{
				Defaults struct {
					Custom json.RawMessage `+"`json:\"custom\"`"+` 
				} `+"`json:\"defaults\"`"+`
			}{}
			if err := json.Unmarshal(fields["%[1]s"], &fakeFieldConfigSource); err != nil {
				return err
			}
			customFieldConfig, err := variantCfg.FieldConfigUnmarshaler(fakeFieldConfigSource.Defaults.Custom)
			if err != nil {
				return err
			}
			if err := json.Unmarshal(fields["%[1]s"], &resource.%[2]s); err != nil {
				return err
			}
	
			resource.%[2]s.Defaults.Custom = customFieldConfig
		} else {
			if err := json.Unmarshal(fields["%[1]s"], &resource.%[2]s); err != nil {
				return err
			}
		}
	}
`, field.Name, tools.UpperCamelCase(field.Name)))
			continue
		}

		buffer.WriteString(fmt.Sprintf(`
	if fields["%[1]s"] != nil {
		if err := json.Unmarshal(fields["%[1]s"], &resource.%[2]s); err != nil {
			return err
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

	return fmt.Sprintf(`func (resource *%[1]s) UnmarshalJSON(raw []byte) error {
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
		if candidate.Type.Kind != ast.KindRef {
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

	if field.Type.Kind == ast.KindArray {
		return fmt.Sprintf(`
	%[3]s
	%[2]s, err := cog.UnmarshalDataqueryArray(fields["%[2]s"], dataqueryTypeHint)
	if err != nil {
		return err
	}
	resource.%[1]s = %[2]s
`, tools.UpperCamelCase(field.Name), field.Name, hintValue)
	}

	return fmt.Sprintf(`
	%[3]s
	%[2]s, err := cog.UnmarshalDataquery(fields["%[2]s"], dataqueryTypeHint)
	if err != nil {
		return err
	}
	resource.%[1]s = %[2]s
`, tools.UpperCamelCase(field.Name), field.Name, hintValue)
}

func (jenny JSONMarshalling) renderPanelcfgVariantUnmarshal(schema *ast.Schema) (string, error) {
	jenny.packageMapper("cog/variants")

	return jenny.renderTemplate("types/variant_panelcfg.json_unmarshal.tmpl", map[string]any{
		"schema":         schema,
		"hasOptions":     schema.LocateObject("Options").Name != "",
		"hasFieldConfig": schema.LocateObject("FieldConfig").Name != "",
	})
}

func (jenny JSONMarshalling) renderDataqueryVariantUnmarshal(schema *ast.Schema, obj ast.Object) (string, error) {
	jenny.packageMapper("cog/variants")

	return jenny.renderTemplate("types/variant_dataquery.json_unmarshal.tmpl", map[string]any{
		"schema": schema,
		"object": obj,
	})
}

func (jenny JSONMarshalling) renderTemplate(templateFile string, data map[string]any) (string, error) {
	buf := bytes.Buffer{}

	tmpls := templates.Funcs(map[string]any{
		"formatType": jenny.typeFormatter.formatType,
	})

	if err := tmpls.ExecuteTemplate(&buf, templateFile, data); err != nil {
		return "", fmt.Errorf("failed executing template: %w", err)
	}

	return buf.String(), nil
}
