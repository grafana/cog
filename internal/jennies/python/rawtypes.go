package python

import (
	"path/filepath"
	"strings"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/common"
)

type RawTypes struct {
	typeFormatter *typeFormatter
	importModule  moduleImporter
}

func (jenny RawTypes) JennyName() string {
	return "PythonRawTypes"
}

func (jenny RawTypes) Generate(context common.Context) (codejen.Files, error) {
	files := make(codejen.Files, 0, len(context.Schemas))

	for _, schema := range context.Schemas {
		output, err := jenny.generateSchema(schema)
		if err != nil {
			return nil, err
		}

		filename := filepath.Join("models", schema.Package+".py")

		files = append(files, *codejen.NewFile(filename, output, jenny))
	}

	return files, nil
}

func (jenny RawTypes) generateSchema(schema *ast.Schema) ([]byte, error) {
	var buffer strings.Builder

	imports := NewImportMap()
	jenny.importModule = func(alias string, pkg string, module string) string {
		if module == schema.Package {
			return ""
		}

		return imports.AddModule(alias, pkg, module)
	}
	jenny.typeFormatter = defaultTypeFormatter(func(alias string, pkg string) string {
		if strings.TrimPrefix(pkg, ".") == schema.Package {
			return ""
		}

		return imports.AddPackage(alias, pkg)
	}, jenny.importModule)

	for i, object := range schema.Objects {
		objectOutput, err := jenny.typeFormatter.formatObject(object)
		if err != nil {
			return nil, err
		}

		buffer.WriteString(objectOutput)

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
