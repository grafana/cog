package python

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/jennies/template"
	"github.com/grafana/cog/internal/languages"
	"github.com/grafana/cog/internal/orderedmap"
	"github.com/grafana/cog/internal/tools"
)

type RawTypes struct {
	tmpl            *template.Template
	typeFormatter   *typeFormatter
	importModule    moduleImporter
	importPkg       pkgImporter
	apiRefCollector *common.APIReferenceCollector
}

func (jenny RawTypes) JennyName() string {
	return "PythonRawTypes"
}

func (jenny RawTypes) Generate(context languages.Context) (codejen.Files, error) {
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

func (jenny RawTypes) generateSchema(context languages.Context, schema *ast.Schema) ([]byte, error) {
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

	jenny.tmpl = jenny.tmpl.
		Funcs(template.FuncMap{
			"formatFullyQualifiedRef": func(typeDef ast.RefType) string {
				return jenny.typeFormatter.formatFullyQualifiedRef(typeDef, false)
			},
			"importModule": jenny.importModule,
			"importPkg":    jenny.importPkg,
		})

	i := 0
	schema.Objects.Iterate(func(_ string, object ast.Object) {
		objectOutput, innerErr := jenny.typeFormatter.formatObject(object)
		if innerErr != nil {
			err = innerErr
			return
		}
		buffer.WriteString(objectOutput)

		if object.Type.IsStruct() {
			buffer.WriteString("\n\n")
			buffer.WriteString(jenny.generateInitMethod(context.Schemas, object))

			buffer.WriteString("\n\n")
			buffer.WriteString(jenny.generateToJSONMethod(object))

			buffer.WriteString("\n\n")
			fromJSON, innerErr := jenny.generateFromJSONMethod(context, object)
			if innerErr != nil {
				err = innerErr
				return
			}
			buffer.WriteString(fromJSON)
		}

		// TODO(kgz): this shouldn't be done by cog
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

	// TODO(kgz): this shouldn't be done by cog
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

		if field.Type.IsConcreteScalar() {
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

	jenny.apiRefCollector.ObjectMethod(object, common.MethodReference{
		Name: "to_json",
		Comments: []string{
			"Converts this object into a representation that can easily be encoded to JSON.",
		},
		Return: "dict[str, object]",
	})

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

func (jenny RawTypes) generateFromJSONMethod(context languages.Context, object ast.Object) (string, error) {
	jenny.apiRefCollector.ObjectMethod(object, common.MethodReference{
		Name: "from_json",
		Comments: []string{
			"Builds this object from a JSON-decoded dict.",
		},
		Arguments: []common.ArgumentReference{
			{Name: "data", Type: "dict[str, typing.Any]"},
		},
		Return: "typing.Self",
		Static: true,
	})

	customUnmarshalTmpl := template.CustomObjectUnmarshalBlock(object)
	if jenny.tmpl.Exists(customUnmarshalTmpl) {
		return jenny.tmpl.Render(customUnmarshalTmpl, map[string]any{
			"Object": object,
		})
	}

	var buffer strings.Builder

	typingPkg := jenny.importPkg("typing", "typing")

	buffer.WriteString("    @classmethod\n")
	buffer.WriteString(fmt.Sprintf("    def from_json(cls, data: dict[str, %[1]s.Any]) -> %[1]s.Self:\n", typingPkg))

	buffer.WriteString(fmt.Sprintf("        args: dict[str, %s.Any] = {}\n", typingPkg))
	var assignments []string
	for _, field := range object.Type.AsStruct().Fields {
		fieldName := formatIdentifier(field.Name)
		value := fmt.Sprintf(`data["%s"]`, field.Name)

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
			decodingMap, decodingCall := jenny.disjunctionFromJSON(valueType.AsDisjunction(), "item")
			if decodingMap != "" {
				decodingMap += "\n            "
			}

			value = fmt.Sprintf(`[%[2]s for item in data["%[1]s"]]`, field.Name, decodingCall)

			assignment := fmt.Sprintf(`        if "%s" in data:
            %sargs["%s"] = %s`, field.Name, decodingMap, fieldName, value)
			assignments = append(assignments, assignment)
			continue
		} else if field.Type.IsDisjunction() {
			decodingMap, decodingCall := jenny.disjunctionFromJSON(field.Type.AsDisjunction(), fmt.Sprintf(`data["%s"]`, field.Name))
			if decodingMap != "" {
				decodingMap += "\n            "
			}

			assignment := fmt.Sprintf(`        if "%s" in data:
            %sargs["%s"] = %s`, field.Name, decodingMap, fieldName, decodingCall)
			assignments = append(assignments, assignment)
			continue
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

	return buffer.String(), nil
}

// TODO(kgz): this shouldn't be done by cog
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

// TODO(kgz): this shouldn't be done by cog
func (jenny RawTypes) generateDataqueryVariantConfigFunc(schema *ast.Schema, object ast.Object) string {
	cogruntime := jenny.importModule("cogruntime", "..cog", "runtime")
	objectName := tools.UpperCamelCase(object.Name)
	identifier := schema.Metadata.Identifier

	setup := ""

	// The `from_json_hook` needs to be generated differently if `object.Type` is a disjunction
	// since there is no "from_json" method to call
	fromJSONHook := fmt.Sprintf("%s.from_json", objectName)
	if object.Type.IsDisjunction() {
		decodingMap, decodingCall := jenny.disjunctionFromJSON(object.Type.AsDisjunction(), "data")
		fromJSONHook = "lambda data: " + decodingCall

		setup = decodingMap + "\n    "
	}

	return fmt.Sprintf(`def variant_config() -> %[2]s.DataqueryConfig:
    %[4]sreturn %[2]s.DataqueryConfig(
        identifier="%[3]s",
        from_json_hook=%[1]s,
    )`, fromJSONHook, cogruntime, identifier, setup)
}

func (jenny RawTypes) disjunctionFromJSON(disjunction ast.DisjunctionType, inputVar string) (string, string) {
	// this potentially generates incorrect code, but there isn't much we can do without more information.
	if disjunction.Discriminator == "" || disjunction.DiscriminatorMapping == nil {
		return "", inputVar
	}

	typingPkg := jenny.importPkg("typing", "typing")

	decodingMap := "{"
	branchTypes := make([]string, 0, len(disjunction.Branches))
	defaultBranch := ""
	discriminators := tools.Keys(disjunction.DiscriminatorMapping)
	sort.Strings(discriminators) // to ensure a deterministic output
	for _, discriminator := range discriminators {
		if discriminator == ast.DiscriminatorCatchAll {
			continue
		}

		objectRef := disjunction.DiscriminatorMapping[discriminator]
		decodingMap += fmt.Sprintf(`"%s": %s, `, discriminator, objectRef)
		branchTypes = append(branchTypes, fmt.Sprintf("%s.Type[%s]", typingPkg, objectRef))
	}

	decodingMap = strings.TrimSuffix(decodingMap, ", ") + "}"

	typeDecl := fmt.Sprintf("dict[str, %s.Union[%s]]", typingPkg, strings.Join(branchTypes, ", "))

	decodingMap = fmt.Sprintf("decoding_map: %s = %s", typeDecl, decodingMap)
	decodingCall := fmt.Sprintf(`decoding_map[%[2]s["%[1]s"]].from_json(%[2]s)`, disjunction.Discriminator, inputVar)

	if defaultBranchType, ok := disjunction.DiscriminatorMapping[ast.DiscriminatorCatchAll]; ok {
		defaultBranch = fmt.Sprintf(`, %s`, defaultBranchType)

		decodingCall = fmt.Sprintf(`decoding_map.get(%[3]s["%[1]s"]%[2]s).from_json(%[3]s)`, disjunction.Discriminator, defaultBranch, inputVar)
	}

	return decodingMap, decodingCall
}

func (jenny RawTypes) composableSlotFromJSON(context languages.Context, parentStruct ast.StructType, field ast.StructField) string {
	// TODO(kgz): this shouldn't be done by cog
	slot, _ := context.ResolveToComposableSlot(field.Type)
	if slot.AsComposableSlot().Variant != ast.SchemaVariantDataQuery {
		return "unknown composable slot variant"
	}

	cogruntime := jenny.importModule("cogruntime", "..cog", "runtime")

	// First: try to locate a field that would contain the type of datasource being used.
	// We're looking for a field defined as a reference to the `DataSourceRef` type.
	var hintField *ast.StructField
	for i, candidate := range parentStruct.Fields {
		if !candidate.Type.IsRef() {
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

	if field.Type.IsArray() {
		return fmt.Sprintf(`[%[3]s.dataquery_from_json(dataquery_json, %[2]s) for dataquery_json in data["%[1]s"]]`, field.Name, hintValue, cogruntime)
	}

	return fmt.Sprintf(`%[3]s.dataquery_from_json(data["%[1]s"], %[2]s)`, field.Name, hintValue, cogruntime)
}
