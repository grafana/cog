package python

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/orderedmap"
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
			buffer.WriteString(jenny.generateToJSONMethod(context, object))

			buffer.WriteString("\n\n")
			buffer.WriteString(jenny.generateFromJSONMethod(context, object))
		}

		// we want two blank lines between objects, except at the end of the file
		if i != schema.Objects.Len()-1 {
			buffer.WriteString("\n\n\n")
		}
	})
	if err != nil {
		return nil, err
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

func (jenny RawTypes) generateToJSONMethod(context common.Context, object ast.Object) string {
	var buffer strings.Builder

	buffer.WriteString("    def to_json(self) -> dict[str, object]:\n")
	buffer.WriteString("        payload: dict[str, object] = {\n")

	fieldValue := func(field ast.StructField, nilCheck bool) string {
		fieldName := formatIdentifier(field.Name)

		if context.ResolveToStruct(field.Type) {
			if nilCheck {
				return fmt.Sprintf("None if self.%[1]s is None else self.%[1]s.to_json()", fieldName)
			}

			return fmt.Sprintf("self.%s.to_json()", fieldName)
		}

		return fmt.Sprintf("self.%s", fieldName)
	}

	for _, field := range object.Type.AsStruct().Fields {
		if !field.Required {
			continue
		}

		buffer.WriteString(fmt.Sprintf(`            "%s": %s,`+"\n", field.Name, fieldValue(field, true)))
	}

	buffer.WriteString("        }\n")

	for _, field := range object.Type.AsStruct().Fields {
		if field.Required {
			continue
		}

		fieldName := formatIdentifier(field.Name)

		buffer.WriteString(fmt.Sprintf("        if self.%s is not None:\n", fieldName))
		buffer.WriteString(fmt.Sprintf(`            payload["%s"] = %s`+"\n", field.Name, fieldValue(field, false)))
	}

	buffer.WriteString("        return payload")

	return buffer.String()
}

func (jenny RawTypes) generateFromJSONMethod(context common.Context, object ast.Object) string {
	var buffer strings.Builder

	typingPkg := jenny.importPkg("typing", "typing")

	buffer.WriteString("    @classmethod\n")
	buffer.WriteString(fmt.Sprintf("    def from_json(cls, data: dict[str, %[1]s.Any]) -> %[1]s.Self:\n", typingPkg))

	buffer.WriteString("        args = {\n")
	var optionalFields []string
	for _, field := range object.Type.AsStruct().Fields {
		fieldName := formatIdentifier(field.Name)
		value := fmt.Sprintf(`data["%s"]`, field.Name)

		// No need to unmarshal constant scalar fields since they're set in
		// the object's constructor
		if field.Type.IsConcreteScalar() {
			continue
		}

		if field.Type.IsRef() {
			ref := field.Type.AsRef()
			referredObject, found := context.LocateObject(ref.ReferredPkg, ref.ReferredType)
			if found && referredObject.Type.IsStruct() {
				formattedRef := jenny.typeFormatter.formatFullyQualifiedRef(ref, false)

				value = fmt.Sprintf(`%s.from_json(data["%s"])`, formattedRef, field.Name)
			}
		}

		if field.Required {
			buffer.WriteString(fmt.Sprintf(`            "%s": %s,`+"\n", fieldName, value))
		} else {
			assignment := fmt.Sprintf(`        if "%s" in data:
            args["%s"] = %s`, field.Name, fieldName, value)

			optionalFields = append(optionalFields, assignment)
		}
	}
	buffer.WriteString("        }\n")

	if len(optionalFields) != 0 {
		buffer.WriteString("        \n")
		buffer.WriteString(strings.Join(optionalFields, "\n"))
		buffer.WriteString("        \n\n")
	}

	buffer.WriteString("        return cls(**args)")

	return buffer.String()
}
