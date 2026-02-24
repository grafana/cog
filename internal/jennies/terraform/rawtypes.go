package terraform

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/languages"
	"github.com/grafana/cog/internal/tools"
)

type RawTypes struct {
	config Config

	packageMapper func(pkg string) string
	typeFormatter *typeFormatter
}

func (jenny RawTypes) JennyName() string {
	return "TerraformRawTypes"
}

func (jenny RawTypes) Generate(context languages.Context) (codejen.Files, error) {
	files := make(codejen.Files, 0, len(context.Schemas))

	for _, schema := range context.Schemas {
		output, err := jenny.generateSchema(context, schema)
		if err != nil {
			return nil, err
		}

		filename := filepath.Join(
			formatPackageName(schema.Package),
			"types_gen.go",
		)

		files = append(files, *codejen.NewFile(filename, output, jenny))
	}
	return files, nil
}

func (jenny RawTypes) generateSchema(context languages.Context, schema *ast.Schema) ([]byte, error) {

	imports := NewImportMap(jenny.config.PackageRoot)
	jenny.packageMapper = func(pkg string) string {
		if imports.IsIdentical(pkg, schema.Package) {
			return ""
		}

		return imports.Add(pkg, jenny.config.importPath(pkg))
	}
	jenny.typeFormatter = defaultTypeFormatter(jenny.config, context, imports, jenny.packageMapper)

	// TF base types
	imports.Add("", "github.com/hashicorp/terraform-plugin-framework/types")

	var buffer strings.Builder
	var err error

	schema.Objects.Iterate(func(_ string, object ast.Object) {
		innerErr := jenny.formatObject(&buffer, schema, object)
		if innerErr != nil {
			err = innerErr
			return
		}

		buffer.WriteString("\n")
	})

	if err != nil {
		return nil, err
	}

	importStatements := imports.String()
	if importStatements != "" {
		importStatements += "\n\n"
	}

	return []byte(fmt.Sprintf(`package %[1]s

%[2]s%[3]s`, formatPackageName(schema.Package), importStatements, buffer.String())), nil
}

func (jenny RawTypes) formatObject(buffer *strings.Builder, schema *ast.Schema, object ast.Object) error {
	comments := object.Comments
	if jenny.config.debug {
		passesTrail := tools.Map(object.PassesTrail, func(trail string) string {
			return fmt.Sprintf("Modified by compiler pass '%s'", trail)
		})
		comments = append(comments, passesTrail...)
	}

	for _, commentLine := range comments {
		buffer.WriteString(fmt.Sprintf("// %s\n", commentLine))
	}

	buffer.WriteString(jenny.typeFormatter.formatTypeDeclaration(object))
	buffer.WriteString("\n")

	switch object.Type.Kind {

	}

	return nil
}
