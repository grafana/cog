package python

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/common"
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
	jenny.typeFormatter = defaultTypeFormatter(jenny.importPkg, jenny.importModule)

	for i, object := range schema.Objects {
		objectOutput, err := jenny.typeFormatter.formatObject(object)
		if err != nil {
			return nil, err
		}

		buffer.WriteString(objectOutput)

		if object.Type.Kind == ast.KindStruct {
			buffer.WriteString("\n\n")
			buffer.Write([]byte(jenny.generateToInitMethod(context.Schemas, object)))
		}

		// we want two blank lines between objects, except at the end of the file
		if i != len(schema.Objects)-1 {
			buffer.WriteString("\n\n\n")
		}
	}

	buffer.WriteString("\n")

	importStatements := imports.String()
	if importStatements != "" {
		importStatements += "\n\n\n"
	}

	return []byte(importStatements + buffer.String()), nil
}

func (jenny RawTypes) generateToInitMethod(schemas ast.Schemas, object ast.Object) string {
	var buffer strings.Builder

	var args []string
	var assignments []string

	for _, field := range object.Type.AsStruct().Fields {
		fieldName := formatFieldName(field.Name)
		fieldType := jenny.typeFormatter.formatType(field.Type)
		defaultValue := defaultValueForType(schemas, field.Type, jenny.importModule)

		if field.Type.IsScalar() && field.Type.AsScalar().IsConcrete() {
			assignments = append(assignments, fmt.Sprintf("        self.%s = %s", fieldName, formatValue(field.Type.AsScalar().Value)))
			continue
		} else if field.Type.IsAnyOf(ast.KindStruct, ast.KindRef, ast.KindEnum, ast.KindMap, ast.KindArray) || field.Type.IsAny() {
			if !field.Type.Nullable {
				typingPkg := jenny.importPkg("typing", "typing")
				fieldType = fmt.Sprintf("%s.Optional[%s]", typingPkg, fieldType)
			}

			args = append(args, fmt.Sprintf("%s: %s = None", fieldName, fieldType))
			assignments = append(assignments, fmt.Sprintf("        self.%[1]s = %[1]s if %[1]s is not None else %[2]s", fieldName, formatValue(defaultValue)))
			continue
		}

		args = append(args, fmt.Sprintf("%s: %s = %s", fieldName, fieldType, formatValue(defaultValue)))
		assignments = append(assignments, fmt.Sprintf("        self.%[1]s = %[1]s", fieldName))
	}

	buffer.WriteString(fmt.Sprintf("    def __init__(self, %s):\n", strings.Join(args, ", ")))
	buffer.WriteString(strings.Join(assignments, "\n"))

	return strings.TrimSuffix(buffer.String(), "\n")
}
