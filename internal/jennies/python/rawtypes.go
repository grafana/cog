package python

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/orderedmap"
	"github.com/grafana/cog/internal/tools"
)

type RawTypes struct {
	typeFormatter *typeFormatter
	importModule  moduleImporter
	importPkg     pkgImporter
}

func (jenny RawTypes) JennyName() string {
	return "PythonRawTypes"
}

func (jenny RawTypes) Generate(context common.Context) (codejen.Files, error) {
	files := make(codejen.Files, 0, len(context.Schemas))

	for _, schema := range context.Schemas {
		output, err := jenny.generateSchema(context, schema)
		if err != nil {
			return nil, err
		}

		filename := filepath.Join("models", schema.Package+".py")

		files = append(files, *codejen.NewFile(filename, output, jenny))
	}

	return files, nil
}

func (jenny RawTypes) generateSchema(context common.Context, schema *ast.Schema) ([]byte, error) {
	var buffer strings.Builder
	var err error

	imports := NewImportMap()
	jenny.importModule = func(alias string, pkg string, module string) string {
		if module == schema.Package {
			return ""
		}

		return imports.AddModule(alias, pkg, module)
	}
	jenny.importPkg = func(alias string, pkg string) string {
		if strings.TrimPrefix(pkg, ".") == schema.Package {
			return ""
		}

		return imports.AddPackage(alias, pkg)
	}
	jenny.typeFormatter = defaultTypeFormatter(context, jenny.importPkg, jenny.importModule)

	i := 0
	schema.Objects.Iterate(func(_ string, object ast.Object) {
		objectOutput, innerErr := jenny.typeFormatter.formatObject(object)
		if innerErr != nil {
			err = innerErr
			return
		}
		buffer.WriteString(objectOutput)

		if object.Type.Kind == ast.KindStruct {
			buffer.WriteString("\n\n")
			buffer.WriteString(jenny.generateInitMethod(context.Schemas, object))

			buffer.WriteString("\n\n")
			buffer.WriteString(jenny.generateToJSONMethod(object))

			buffer.WriteString("\n\n")
			buffer.WriteString(jenny.generateFromJSONMethod(context, object))
		}

		if object.Type.ImplementedVariant() == string(ast.SchemaVariantDataQuery) && !object.Type.HasHint(ast.HintSkipVariantPluginRegistration) {
			buffer.WriteString("\n\n\n")
			buffer.WriteString(jenny.generateDataqueryVariantConfigFunc(schema, object))
		}

		// we want two blank lines between objects, except at the end of the file
		if i != schema.Objects.Len()-1 {
			buffer.WriteString("\n\n\n")
		}
	})
	if err != nil {
		return nil, err
	}

	if schema.Metadata.Kind == ast.SchemaKindComposable && schema.Metadata.Variant == ast.SchemaVariantPanel {
		buffer.WriteString("\n\n\n")
		buffer.WriteString(jenny.generatePanelCfgVariantConfigFunc(schema))
	}

	buffer.WriteString("\n")

	importStatements := imports.String()
	if importStatements != "" {
		importStatements += "\n\n\n"
	}

	return []byte(importStatements + buffer.String()), nil
}

func (jenny RawTypes) generateInitMethod(schemas ast.Schemas, object ast.Object) string {
	var buffer strings.Builder

	var args []string
	var assignments []string

	for _, field := range object.Type.AsStruct().Fields {
		fieldName := formatIdentifier(field.Name)
		fieldType := jenny.typeFormatter.formatType(field.Type)
		defaultValue := (any)(nil)

		if !field.Type.Nullable || field.Type.Default != nil {
			var defaultsOverrides map[string]any
			if overrides, ok := field.Type.Default.(map[string]interface{}); ok {
				defaultsOverrides = overrides
			}

			defaultValue = defaultValueForType(schemas, field.Type, jenny.importModule, orderedmap.FromMap(defaultsOverrides))
		}

		if field.Type.IsScalar() && field.Type.AsScalar().IsConcrete() {
			assignments = append(assignments, fmt.Sprintf("        self.%s = %s", fieldName, formatValue(field.Type.AsScalar().Value)))
			continue
		} else if field.Type.IsAnyOf(ast.KindStruct, ast.KindRef, ast.KindEnum, ast.KindMap, ast.KindArray) {
			if !field.Type.Nullable {
				typingPkg := jenny.importPkg("typing", "typing")
				fieldType = fmt.Sprintf("%s.Optional[%s]", typingPkg, fieldType)
			}

			args = append(args, fmt.Sprintf("%s: %s = None", fieldName, fieldType))

			if defaultValue == nil {
				assignments = append(assignments, fmt.Sprintf("        self.%[1]s = %[1]s", fieldName))
			} else {
				assignments = append(assignments, fmt.Sprintf("        self.%[1]s = %[1]s if %[1]s is not None else %[2]s", fieldName, formatValue(defaultValue)))
			}
			continue
		}

		args = append(args, fmt.Sprintf("%s: %s = %s", fieldName, fieldType, formatValue(defaultValue)))
		assignments = append(assignments, fmt.Sprintf("        self.%[1]s = %[1]s", fieldName))
	}

	buffer.WriteString(fmt.Sprintf("    def __init__(self, %s):\n", strings.Join(args, ", ")))
	buffer.WriteString(strings.Join(assignments, "\n"))

	return strings.TrimSuffix(buffer.String(), "\n")
}

func (jenny RawTypes) generateToJSONMethod(object ast.Object) string {
	var buffer strings.Builder

	buffer.WriteString("    def to_json(self) -> dict[str, object]:\n")
	buffer.WriteString("        payload: dict[str, object] = {\n")

	for _, field := range object.Type.AsStruct().Fields {
		if !field.Required {
			continue
		}

		buffer.WriteString(fmt.Sprintf(`            "%s": self.%s,`+"\n", field.Name, formatIdentifier(field.Name)))
	}

	buffer.WriteString("        }\n")

	for _, field := range object.Type.AsStruct().Fields {
		if field.Required {
			continue
		}

		fieldName := formatIdentifier(field.Name)

		buffer.WriteString(fmt.Sprintf("        if self.%s is not None:\n", fieldName))
		buffer.WriteString(fmt.Sprintf(`            payload["%s"] = self.%s`+"\n", field.Name, fieldName))
	}

	buffer.WriteString("        return payload")

	return buffer.String()
}

func (jenny RawTypes) generateFromJSONMethod(context common.Context, object ast.Object) string {
	var buffer strings.Builder

	typingPkg := jenny.importPkg("typing", "typing")

	buffer.WriteString("    @classmethod\n")
	buffer.WriteString(fmt.Sprintf("    def from_json(cls, data: dict[str, %[1]s.Any]) -> %[1]s.Self:\n", typingPkg))

	buffer.WriteString(fmt.Sprintf("        args: dict[str, %s.Any] = {}\n", typingPkg))
	var assignments []string
	for _, field := range object.Type.AsStruct().Fields {
		fieldName := formatIdentifier(field.Name)
		value := fmt.Sprintf(`data["%s"]`, field.Name)

		// Special cases to properly parse dashboard.Panel options
		if object.SelfRef.ReferredPkg == "dashboard" && strings.EqualFold(object.Name, "panel") && field.Name == "options" {
			cogruntime := jenny.importModule("cogruntime", "..cog", "runtime")
			assignment := fmt.Sprintf(`        if "options" in data:
            config = %[1]s.panelcfg_config(data.get("type", ""))
            if config is not None and config.options_from_json_hook is not None:
                args["%[2]s"] = config.options_from_json_hook(data["options"])
            else:
                args["%[2]s"] = data["options"]`, cogruntime, fieldName)

			assignments = append(assignments, assignment)
			continue
		}

		// Special cases to properly parse dashboard.Panel fieldConfig
		if object.SelfRef.ReferredPkg == "dashboard" && strings.EqualFold(object.Name, "panel") && field.Name == "fieldConfig" {
			cogruntime := jenny.importModule("cogruntime", "..cog", "runtime")
			assignment := fmt.Sprintf(`        if "fieldConfig" in data:
            config = %[1]s.panelcfg_config(data.get("type", ""))
            field_config = FieldConfigSource.from_json(data["fieldConfig"])

            if config is not None and config.field_config_from_json_hook is not None:
                custom_field_config = data["fieldConfig"].get("defaults", {}).get("custom", {})
                field_config.defaults.custom = config.field_config_from_json_hook(custom_field_config)

            args["%[2]s"] = field_config`, cogruntime, fieldName)

			assignments = append(assignments, assignment)
			continue
		}

		// No need to unmarshal constant scalar fields since they're set in
		// the object's constructor
		if field.Type.IsConcreteScalar() {
			continue
		}

		if _, ok := context.ResolveToComposableSlot(field.Type); ok {
			value = jenny.composableSlotFromJSON(context, object.Type.AsStruct(), field)
		} else if field.Type.IsRef() { //nolint:gocritic
			ref := field.Type.AsRef()
			referredObject, found := context.LocateObject(ref.ReferredPkg, ref.ReferredType)
			if found && referredObject.Type.IsStruct() {
				formattedRef := jenny.typeFormatter.formatFullyQualifiedRef(ref, false)

				value = fmt.Sprintf(`%s.from_json(data["%s"])`, formattedRef, field.Name)
			}
		} else if field.Type.IsArray() && field.Type.Array.ValueType.IsDisjunction() {
			valueType := field.Type.Array.ValueType
			disjunctionFromJSON := jenny.disjunctionFromJSON(valueType.AsDisjunction(), "item")

			value = fmt.Sprintf(`[%[2]s for item in data["%[1]s"]]`, field.Name, disjunctionFromJSON)
		} else if field.Type.IsDisjunction() {
			value = jenny.disjunctionFromJSON(field.Type.AsDisjunction(), fmt.Sprintf(`data["%s"]`, field.Name))
		}

		assignment := fmt.Sprintf(`        if "%s" in data:
            args["%s"] = %s`, field.Name, fieldName, value)

		assignments = append(assignments, assignment)
	}

	if len(assignments) != 0 {
		buffer.WriteString("        \n")
		buffer.WriteString(strings.Join(assignments, "\n"))
		buffer.WriteString("        \n\n")
	}

	buffer.WriteString("        return cls(**args)")

	return buffer.String()
}

func (jenny RawTypes) generatePanelCfgVariantConfigFunc(schema *ast.Schema) string {
	cogruntime := jenny.importModule("cogruntime", "..cog", "runtime")
	identifier := schema.Metadata.Identifier

	options := "Options.from_json"
	if _, hasOptions := schema.LocateObject("Options"); !hasOptions {
		options = "None"
	}

	fieldConfig := "FieldConfig.from_json"
	if _, hasFieldConfig := schema.LocateObject("FieldConfig"); !hasFieldConfig {
		fieldConfig = "None"
	}

	return fmt.Sprintf(`def variant_config():
    return %[1]s.PanelCfgConfig(
        identifier="%[2]s",
        options_from_json_hook=%[3]s,
        field_config_from_json_hook=%[4]s,
    )`, cogruntime, identifier, options, fieldConfig)
}

func (jenny RawTypes) generateDataqueryVariantConfigFunc(schema *ast.Schema, object ast.Object) string {
	cogruntime := jenny.importModule("cogruntime", "..cog", "runtime")
	objectName := tools.UpperCamelCase(object.Name)
	identifier := schema.Metadata.Identifier

	// The `from_json_hook` needs to be generated differently if `object.Type` is a disjunction
	// since there is no "from_json" method to call
	fromJSONHook := fmt.Sprintf("%s.from_json", objectName)
	if object.Type.IsDisjunction() {
		fromJSONHook = "lambda data: " + jenny.disjunctionFromJSON(object.Type.AsDisjunction(), "data")
	}

	return fmt.Sprintf(`def variant_config():
    return %[2]s.DataqueryConfig(
        identifier="%[3]s",
        from_json_hook=%[1]s,
    )`, fromJSONHook, cogruntime, identifier)
}

func (jenny RawTypes) disjunctionFromJSON(disjunction ast.DisjunctionType, inputVar string) string {
	// this potentially generates incorrect code, but there isn't much we can do without more information.
	if disjunction.Discriminator == "" || disjunction.DiscriminatorMapping == nil {
		return inputVar
	}

	decodingMap := "({"
	defaultBranch := ""
	for discriminator, objectRef := range disjunction.DiscriminatorMapping {
		if discriminator == ast.DiscriminatorCatchAll {
			continue
		}

		decodingMap += fmt.Sprintf(`"%s": %s, `, discriminator, objectRef)
	}

	if defaultBranchType, ok := disjunction.DiscriminatorMapping[ast.DiscriminatorCatchAll]; ok {
		defaultBranch = fmt.Sprintf(`, %s`, defaultBranchType)
	}

	decodingMap = strings.TrimSuffix(decodingMap, ", ")
	decodingMap += fmt.Sprintf(`}.get(%[3]s["%[1]s"]%[2]s)).from_json(%[3]s)`, disjunction.Discriminator, defaultBranch, inputVar)

	return decodingMap
}

func (jenny RawTypes) composableSlotFromJSON(context common.Context, parentStruct ast.StructType, field ast.StructField) string {
	slot, _ := context.ResolveToComposableSlot(field.Type)
	if slot.AsComposableSlot().Variant != ast.SchemaVariantDataQuery {
		return "unknown composable slot variant"
	}

	cogruntime := jenny.importModule("cogruntime", "..cog", "runtime")

	// First: try to locate a field that would contain the type of datasource being used.
	// We're looking for a field defined as a reference to the `DataSourceRef` type.
	var hintField *ast.StructField
	for i, candidate := range parentStruct.Fields {
		if candidate.Type.Kind != ast.KindRef {
			continue
		}
		if candidate.Type.AsRef().ReferredType != "DataSourceRef" {
			continue
		}

		hintField = &parentStruct.Fields[i]
	}

	// then: unmarshalling boilerplate
	hintValue := `""`
	if hintField != nil {
		hintValue = fmt.Sprintf(`data["%[1]s"]["type"] if data.get("%[1]s") is not None and data["%[1]s"].get("type", "") != "" else ""`, hintField.Name)
	}

	if field.Type.Kind == ast.KindArray {
		return fmt.Sprintf(`[%[3]s.dataquery_from_json(dataquery_json, %[2]s) for dataquery_json in data["%[1]s"]]`, field.Name, hintValue, cogruntime)
	}

	return fmt.Sprintf(`%[3]s.dataquery_from_json(data["%[1]s"], %[2]s)`, field.Name, hintValue, cogruntime)
}
