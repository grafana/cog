package golang

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/jennies/template"
	"github.com/grafana/cog/internal/languages"
	"github.com/grafana/cog/internal/tools"
)

type RawTypes struct {
	Config          Config
	Tmpl            *template.Template
	apiRefCollector *common.APIReferenceCollector

	typeFormatter *typeFormatter
}

func (jenny RawTypes) JennyName() string {
	return "GoRawTypes"
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
	var buffer strings.Builder
	var err error

	imports := NewImportMap(jenny.Config.PackageRoot)
	packageMapper := func(pkg string) string {
		if imports.IsIdentical(pkg, schema.Package) {
			return ""
		}

		return imports.Add(pkg, jenny.Config.importPath(pkg))
	}
	jenny.typeFormatter = defaultTypeFormatter(jenny.Config, context, imports, packageMapper)
	unmarshallerGenerator := NewJSONMarshalling(jenny.Config, jenny.Tmpl, imports, packageMapper, jenny.typeFormatter, jenny.apiRefCollector)
	strictUnmarshallerGenerator := newStrictJSONUnmarshal(jenny.Tmpl, imports, packageMapper, jenny.typeFormatter, jenny.apiRefCollector)
	equalityMethodsGenerator := newEqualityMethods(jenny.Tmpl, jenny.apiRefCollector)
	validationMethodsGenerator := newValidationMethods(jenny.Tmpl, packageMapper, jenny.apiRefCollector)

	schema.Objects.Iterate(func(_ string, object ast.Object) {
		objectOutput, innerErr := jenny.formatObject(schema, object)
		if innerErr != nil {
			err = innerErr
			return
		}

		buffer.Write(objectOutput)
		buffer.WriteString("\n")

		innerErr = unmarshallerGenerator.generateForObject(&buffer, context, schema, object)
		if innerErr != nil {
			err = innerErr
			return
		}

		innerErr = strictUnmarshallerGenerator.generateForObject(&buffer, context, object)
		if innerErr != nil {
			err = innerErr
			return
		}

		innerErr = equalityMethodsGenerator.generateForObject(&buffer, context, object, imports)
		if innerErr != nil {
			err = innerErr
			return
		}

		innerErr = validationMethodsGenerator.generateForObject(&buffer, context, object, imports)
		if innerErr != nil {
			err = innerErr
			return
		}
	})
	if err != nil {
		return nil, err
	}

	if err := unmarshallerGenerator.generateForSchema(&buffer, schema); err != nil {
		return nil, err
	}

	importStatements := imports.String()
	if importStatements != "" {
		importStatements += "\n\n"
	}

	return []byte(fmt.Sprintf(`package %[1]s

%[2]s%[3]s`, formatPackageName(schema.Package), importStatements, buffer.String())), nil
}

func (jenny RawTypes) formatObject(schema *ast.Schema, def ast.Object) ([]byte, error) {
	var buffer strings.Builder

	defName := tools.UpperCamelCase(def.Name)

	comments := def.Comments
	if jenny.Config.debug {
		passesTrail := tools.Map(def.PassesTrail, func(trail string) string {
			return fmt.Sprintf("Modified by compiler pass '%s'", trail)
		})
		comments = append(comments, passesTrail...)
	}

	for _, commentLine := range comments {
		buffer.WriteString(fmt.Sprintf("// %s\n", commentLine))
	}

	buffer.WriteString(jenny.typeFormatter.formatTypeDeclaration(def))
	buffer.WriteString("\n")

	if def.Type.ImplementsVariant() && !def.Type.IsRef() {
		variant := tools.UpperCamelCase(def.Type.ImplementedVariant())

		buffer.WriteString(fmt.Sprintf("func (resource %s) Implements%sVariant() {}\n", defName, variant))
		buffer.WriteString("\n")

		if def.Type.ImplementedVariant() == string(ast.SchemaVariantDataQuery) {
			buffer.WriteString(fmt.Sprintf("func (resource %s) DataqueryType() string {\n", defName))
			buffer.WriteString(fmt.Sprintf("\treturn \"%s\"\n", strings.ToLower(schema.Metadata.Identifier)))
			buffer.WriteString("}\n")
		}
	}

	return []byte(buffer.String()), nil
}
