package terraform

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/jennies/golang"
	"github.com/grafana/cog/internal/jennies/template"
	"github.com/grafana/cog/internal/languages"
	"github.com/grafana/cog/internal/tools"
)

type RawTypes struct {
	config Config
	tmpl   *template.Template

	packageMapper func(pkg string) string

	typeFormatter  *typeFormatter
	modelFormatter *modelFormatter
	attributes     *attributes
}

func (jenny RawTypes) JennyName() string {
	return "TerraformRawTypes"
}

func (jenny RawTypes) Generate(context languages.Context) (codejen.Files, error) {
	files := make(codejen.Files, 0, len(context.Schemas))

	jenny.tmpl = jenny.tmpl.
		Funcs(common.TypeResolvingTemplateHelpers(context)).
		Funcs(common.TypesTemplateHelpers(context))

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

		return imports.Add(pkg, pkg)
	}
	goTypesFormatterHelper := golang.MakeTypeFormatterHelper(golang.Config{}, context, imports, jenny.packageMapper)
	jenny.typeFormatter = defaultTypeFormatter(context, jenny.packageMapper)
	jenny.modelFormatter = defaultModelFormatter(context, jenny.packageMapper)
	jenny.attributes = newAttributesGenerator(context, jenny.config, jenny.typeFormatter, jenny.packageMapper)

	// TF base types
	imports.Add("", "github.com/hashicorp/terraform-plugin-framework/types")

	// Per-schema template with importStdPkg wired to the local imports map.
	schemaTmpl := jenny.tmpl.Funcs(template.FuncMap{
		"importStdPkg": func(pkg string) string {
			return imports.Add(pkg, pkg)
		},
		"toGoType":          goTypesFormatterHelper,
		"toTfType":          jenny.typeFormatter.formatType,
		"toTfModel":         jenny.modelFormatter.formatModel,
		"toTfModelWithRefs": modelFormatterWithRefs(context, jenny.packageMapper).formatModel,
	})

	var buffer strings.Builder
	var err error

	jenny.attributes.identifyDisjunctionBranches(schema)

	schema.Objects.Iterate(func(_ string, object ast.Object) {
		innerErr := jenny.formatObject(&buffer, object)
		if innerErr != nil {
			err = innerErr
			return
		}

		customMethodsBlock := template.CustomObjectMethodsBlock(object)
		if schemaTmpl.Exists(customMethodsBlock) {
			innerErr = schemaTmpl.RenderInBuffer(&buffer, customMethodsBlock, map[string]any{
				"Schema": schema,
				"Object": object,
			})
			if innerErr != nil {
				err = innerErr
				return
			}
			buffer.WriteString("\n")
		}

		customAllBlock := template.CustomObjectMethodAllBlock()
		if schemaTmpl.Exists(customAllBlock) {
			innerErr = schemaTmpl.RenderInBuffer(&buffer, customAllBlock, map[string]any{
				"Schema": schema,
				"Object": object,
			})
			if innerErr != nil {
				err = innerErr
				return
			}
			buffer.WriteString("\n")
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

func (jenny RawTypes) formatObject(buffer *strings.Builder, object ast.Object) error {
	comments := object.Comments
	if jenny.config.debug {
		passesTrail := tools.Map(object.PassesTrail, func(trail string) string {
			return fmt.Sprintf("Modified by compiler pass '%s'", trail)
		})
		comments = append(comments, passesTrail...)
	}

	for _, commentLine := range comments {
		fmt.Fprintf(buffer, "// %s\n", commentLine)
	}

	buffer.WriteString(jenny.modelFormatter.formatDeclaration(object))
	buffer.WriteString("\n")

	buffer.WriteString(jenny.typeFormatter.formatDeclaration(object))
	buffer.WriteString("\n")

	buffer.WriteString(jenny.attributes.generateForObject(object))
	buffer.WriteString("\n")

	return nil
}
