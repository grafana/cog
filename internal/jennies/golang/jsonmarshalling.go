package golang

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/context"
	"github.com/grafana/cog/internal/tools"
)

type JSONMarshalling struct {
	imports importMap
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
	jenny.imports = newImportMap()

	packageMapper := func(pkg string) string {
		if pkg == schema.Package {
			return ""
		}

		jenny.imports.Add(pkg, "github.com/grafana/cog/generated/"+pkg)

		return pkg
	}

	for _, object := range schema.Objects {
		if jenny.objectNeedsCustomMarshal(object) {
			jsonMarshal, err := jenny.renderCustomMarshal(object)
			if err != nil {
				return nil, err
			}
			buffer.WriteString(jsonMarshal)
		}

		if jenny.objectNeedsCustomUnmarshal(context, object) {
			jsonUnmarshal, err := jenny.renderCustomUnmarshal(context, object, packageMapper)
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

	importStatements := jenny.imports.Format()
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
		return renderTemplate("disjunction_of_scalars.types.json_marshal.tmpl", map[string]any{
			"def": obj,
		})
	}

	if isStruct && obj.Type.HasHint(ast.HintDiscriminatedDisjunctionOfRefs) {
		return renderTemplate("disjunction_of_refs.types.json_marshal.tmpl", map[string]any{
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

func (jenny JSONMarshalling) renderCustomUnmarshal(context context.Builders, obj ast.Object, packageMapper pkgMapper) (string, error) {
	isStruct := obj.Type.Kind == ast.KindStruct

	if isStruct && obj.Type.HasHint(ast.HintDisjunctionOfScalars) {
		return renderTemplate("disjunction_of_scalars.types.json_unmarshal.tmpl", map[string]any{
			"def":       obj,
			"pkgMapper": packageMapper,
		})
	}

	if isStruct && obj.Type.HasHint(ast.HintDiscriminatedDisjunctionOfRefs) {
		return renderTemplate("disjunction_of_refs.types.json_unmarshal.tmpl", map[string]any{
			"def":  obj,
			"hint": obj.Type.Hints[ast.HintDiscriminatedDisjunctionOfRefs],
		})
	}

	return jenny.renderCustomComposableSlotUnmarshal(context, obj)
}

func (jenny JSONMarshalling) renderCustomComposableSlotUnmarshal(context context.Builders, obj ast.Object) (string, error) {
	var buffer strings.Builder
	fields := obj.Type.AsStruct().Fields

	jenny.imports.Add("cog", "github.com/grafana/cog/generated/cog")

	// unmarshal "normal" fields (ie: with no composable slot)
	for _, field := range fields {
		if _, ok := context.ResolveToComposableSlot(field.Type); ok {
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
			source := jenny.renderUnmarshalDataqueryField(field)
			buffer.WriteString(source)
		case ast.SchemaVariantPanel:
			// TODO
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

func (jenny JSONMarshalling) renderUnmarshalDataqueryField(field ast.StructField) string {
	if field.Type.Kind == ast.KindArray {
		return fmt.Sprintf(`
	%[2]s, err := cog.UnmarshalDataqueryArray(fields["%[2]s"])
	if err != nil {
		return err
	}
	resource.%[1]s = %[2]s
`, tools.UpperCamelCase(field.Name), field.Name)
	}

	return fmt.Sprintf(`
	%[2]s, err := cog.UnmarshalDataquery(fields["%[2]s"])
	if err != nil {
		return err
	}
	resource.%[1]s = %[2]s
`, tools.UpperCamelCase(field.Name), field.Name)
}

func (jenny JSONMarshalling) renderPanelcfgVariantUnmarshal(schema *ast.Schema) (string, error) {
	jenny.imports.Add("cogvariants", "github.com/grafana/cog/generated/cog/variants")

	return renderTemplate("variant_panelcfg.types.json_unmarshal.tmpl", map[string]any{
		"schema":         schema,
		"hasOptions":     schema.LocateObject("Options").Name != "",
		"hasFieldConfig": schema.LocateObject("FieldConfig").Name != "",
	})
}

func (jenny JSONMarshalling) renderDataqueryVariantUnmarshal(schema *ast.Schema, obj ast.Object) (string, error) {
	jenny.imports.Add("cogvariants", "github.com/grafana/cog/generated/cog/variants")

	return renderTemplate("variant_dataquery.types.json_unmarshal.tmpl", map[string]any{
		"schema": schema,
		"object": obj,
	})
}
