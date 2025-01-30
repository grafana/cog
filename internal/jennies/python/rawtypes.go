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
		Funcs(common.TypeResolvingTemplateHelpers(context)).
		Funcs(template.FuncMap{
			"formatFullyQualifiedRef": func(typeDef ast.RefType) string {
				return jenny.typeFormatter.formatFullyQualifiedRef(typeDef, false)
			},
			"importModule": jenny.importModule,
			"importPkg":    jenny.importPkg,
			"disjunctionFromJSON": func(typeDef ast.Type, inputVar string) disjunctionFromJSONCode {
				return jenny.disjunctionFromJSON(context, typeDef, inputVar)
			},
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

		customVariantTmpl := template.CustomObjectVariantBlock(object)
		if object.Type.ImplementsVariant() && jenny.tmpl.Exists(customVariantTmpl) {
			buffer.WriteString("\n\n\n")
			if innerErr := jenny.tmpl.RenderInBuffer(&buffer, customVariantTmpl, map[string]any{
				"Object": object,
				"Schema": schema,
			}); innerErr != nil {
				err = innerErr
				return
			}
		}

		// we want two blank lines between objects, except at the end of the file
		if i != schema.Objects.Len()-1 {
			buffer.WriteString("\n\n\n")
		}
	})
	if err != nil {
		return nil, err
	}

	customSchemaVariant := template.CustomSchemaVariantBlock(schema)
	if schema.Metadata.Kind == ast.SchemaKindComposable && jenny.tmpl.Exists(customSchemaVariant) {
		buffer.WriteString("\n\n\n")

		if err := jenny.tmpl.RenderInBuffer(&buffer, customSchemaVariant, map[string]any{
			"Schema": schema,
		}); err != nil {
			return nil, err
		}
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
			assignments = append(assignments, fmt.Sprintf("        self.%s = %s", fieldName, formatValue(field.Type.Value)))
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
	var err error

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
			value, err = jenny.composableSlotFromJSON(context, object, field)
			if err != nil {
				return "", err
			}
		} else if field.Type.IsRef() { //nolint:gocritic
			ref := field.Type.AsRef()
			referredObject, found := context.LocateObject(ref.ReferredPkg, ref.ReferredType)
			if found && referredObject.Type.IsStruct() {
				formattedRef := jenny.typeFormatter.formatFullyQualifiedRef(ref, false)

				value = fmt.Sprintf(`%s.from_json(data["%s"])`, formattedRef, field.Name)
			}
		} else if field.Type.IsArray() && field.Type.Array.ValueType.IsDisjunction() {
			valueType := field.Type.Array.ValueType
			code := jenny.disjunctionFromJSON(context, valueType, "item")
			if code.Setup != "" {
				code.Setup += "\n            "
			}

			value = fmt.Sprintf(`[%[2]s for item in data["%[1]s"]]`, field.Name, code.DecodingCall)

			assignment := fmt.Sprintf(`        if "%s" in data:
            %sargs["%s"] = %s`, field.Name, code.Setup, fieldName, value)
			assignments = append(assignments, assignment)
			continue
		} else if field.Type.IsDisjunction() {
			code := jenny.disjunctionFromJSON(context, field.Type, fmt.Sprintf(`data["%s"]`, field.Name))
			if code.Setup != "" {
				code.Setup += "\n            "
			}

			assignment := fmt.Sprintf(`        if "%s" in data:
            %sargs["%s"] = %s`, field.Name, code.Setup, fieldName, code.DecodingCall)
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

type disjunctionFromJSONCode struct {
	Setup        string
	DecodingCall string
}

func (jenny RawTypes) disjunctionFromJSON(context languages.Context, typeDef ast.Type, inputVar string) disjunctionFromJSONCode {
	disjunction := context.ResolveRefs(typeDef).AsDisjunction()

	// this potentially generates incorrect code, but there isn't much we can do without more information.
	if disjunction.Discriminator == "" || disjunction.DiscriminatorMapping == nil {
		return disjunctionFromJSONCode{DecodingCall: inputVar}
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

	return disjunctionFromJSONCode{
		Setup:        decodingMap,
		DecodingCall: decodingCall,
	}
}

func (jenny RawTypes) composableSlotFromJSON(context languages.Context, parentObject ast.Object, field ast.StructField) (string, error) {
	slot, _ := context.ResolveToComposableSlot(field.Type)
	variant := string(slot.AsComposableSlot().Variant)
	unmarshalVariantBlock := template.VariantFieldUnmarshalBlock(variant)
	if !jenny.tmpl.Exists(unmarshalVariantBlock) {
		return "", fmt.Errorf("can not generate custom unmarshal function for composable slot with variant '%s': template block %s not found", variant, unmarshalVariantBlock)
	}

	return jenny.tmpl.Render(unmarshalVariantBlock, map[string]any{
		"Object": parentObject,
		"Field":  field,
	})
}
