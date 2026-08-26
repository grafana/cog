package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"os"
	"reflect"
	"sort"
	"strings"

	"github.com/grafana/cog/internal/tools"
	"github.com/grafana/cog/internal/yaml"
)

//nolint:gosec
const irTransformsTypesSourceDir = "./internal/ir/transforms"

//nolint:gosec
const yamlCompilerPassTypesSourceDir = "./internal/yaml"

const outputFile = "./docs/reference/schema_transformations.md"

type transformParam struct {
	Name          string
	Type          string
	Documentation string
}

type transformationDocEntry struct {
	YamlName      string
	Parameters    []transformParam
	Documentation string
}

func yamlNameForField(field reflect.StructField) string {
	yamlName := field.Tag.Get("yaml")
	if yamlName == "" {
		yamlName = tools.SnakeCase(field.Name)
	}

	return yamlName
}

func transformParameters(yamlTransformComments map[string]string, transformType reflect.Type) []transformParam {
	params := make([]transformParam, 0, transformType.NumField())
	for i := 0; i < transformType.NumField(); i++ {
		field := transformType.Field(i)
		fieldType := field.Type.Name()
		if field.Type.Name() == "" {
			fieldType = field.Type.String()
		}

		params = append(params, transformParam{
			Name:          yamlNameForField(field),
			Type:          fieldType,
			Documentation: yamlTransformComments[fmt.Sprintf("%s.%s", transformType.Name(), field.Name)],
		})
	}

	return params
}

func buildYamlTransformTypesCommentsMap(yamlTypesInputDir string) (map[string]string, error) {
	commentsMap := make(map[string]string)

	packages, err := parser.ParseDir(token.NewFileSet(), yamlTypesInputDir, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	for _, packageAst := range packages {
		packageDocs := doc.New(packageAst, "./", doc.AllDecls)

		for _, t := range packageDocs.Types {
			if t.Decl == nil {
				continue
			}

			if len(t.Decl.Specs) != 1 {
				continue
			}

			if _, ok := t.Decl.Specs[0].(*ast.TypeSpec); !ok {
				continue
			}

			typeSpec := t.Decl.Specs[0].(*ast.TypeSpec)
			if _, ok := typeSpec.Type.(*ast.StructType); !ok {
				continue
			}
			structType := typeSpec.Type.(*ast.StructType)

			for _, fields := range structType.Fields.List {
				for _, field := range fields.Names {
					fieldRef := fmt.Sprintf("%s.%s", t.Name, field.Name)
					commentsMap[fieldRef] = fields.Doc.Text()
				}
			}
		}
	}

	return commentsMap, nil
}

func transformationsDocEntries(transformationsComments map[string]string, yamlTransformationsComments map[string]string) []transformationDocEntry {
	var entries []transformationDocEntry
	yamlTransformTypeOf := reflect.TypeFor[yaml.Transform]()

	for i := 0; i < yamlTransformTypeOf.NumField(); i++ {
		field := yamlTransformTypeOf.Field(i)

		asTransformMethod, found := field.Type.MethodByName("AsTransform")
		if !found {
			continue
		}

		transformType := asTransformMethod.Type.Out(0)
		transformTypeName := transformType.Elem().Name()

		entries = append(entries, transformationDocEntry{
			YamlName:      yamlNameForField(field),
			Parameters:    transformParameters(yamlTransformationsComments, field.Type.Elem()),
			Documentation: transformationsComments[transformTypeName],
		})
	}

	sort.SliceStable(entries, func(i, j int) bool {
		return entries[i].YamlName < entries[j].YamlName
	})

	return entries
}

func buildTransformTypesCommentsMap(typesInputDir string) (map[string]string, error) {
	transformComments := make(map[string]string)

	packages, err := parser.ParseDir(token.NewFileSet(), typesInputDir, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	for _, packageAst := range packages {
		packageDocs := doc.New(packageAst, "./", doc.AllDecls)

		for _, t := range packageDocs.Types {
			if t.Doc == "" {
				continue
			}

			transformComments[t.Name] = t.Doc
		}
	}

	return transformComments, nil
}

func docEntriesToMarkdown(entries []transformationDocEntry) []byte {
	var markdown bytes.Buffer

	markdown.WriteString("<!-- Generated with `make docs` -->\n")

	markdown.WriteString("# Schema transformations\n\n")
	for _, entry := range entries {
		markdown.WriteString(fmt.Sprintf("## `%s`\n", entry.YamlName))
		markdown.WriteString("\n")

		if entry.Documentation == "" {
			markdown.WriteString("N/A\n")
		} else {
			markdown.WriteString(entry.Documentation)
		}
		markdown.WriteString("\n")

		markdown.WriteString("### Usage\n\n")

		markdown.WriteString("```yaml\n")
		markdown.WriteString(fmt.Sprintf("%s:", entry.YamlName))

		if len(entry.Parameters) == 0 {
			markdown.WriteString(" {}")
		}

		markdown.WriteString("\n")

		for _, param := range entry.Parameters {
			if param.Documentation != "" {
				markdown.WriteString("  # " + strings.TrimSuffix(param.Documentation, "\n") + "\n")
			}
			markdown.WriteString(fmt.Sprintf("  %s: %s\n", param.Name, param.Type))
		}

		markdown.WriteString("```\n\n")
	}

	return markdown.Bytes()
}

func main() {
	transformComments, err := buildTransformTypesCommentsMap(irTransformsTypesSourceDir)
	if err != nil {
		panic(err)
	}

	yamlTransformsComments, err := buildYamlTransformTypesCommentsMap(yamlCompilerPassTypesSourceDir)
	if err != nil {
		panic(err)
	}

	docEntries := transformationsDocEntries(transformComments, yamlTransformsComments)

	if err := os.WriteFile(outputFile, docEntriesToMarkdown(docEntries), 0600); err != nil {
		panic(err)
	}
}
