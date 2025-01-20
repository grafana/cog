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
const compilerPassTypesSourceDir = "./internal/ast/compiler"

//nolint:gosec
const yamlCompilerPassTypesSourceDir = "./internal/yaml"

const outputFile = "./docs/reference/schema_transformations.md"

type compilerPassParam struct {
	Name          string
	Type          string
	Documentation string
}

type compilerPassDocEntry struct {
	YamlName      string
	Parameters    []compilerPassParam
	Documentation string
}

func yamlNameForField(field reflect.StructField) string {
	yamlName := field.Tag.Get("yaml")
	if yamlName == "" {
		yamlName = tools.SnakeCase(field.Name)
	}

	return yamlName
}

func compilerPassParameters(yamlCompilerPassComments map[string]string, compilerPassType reflect.Type) []compilerPassParam {
	params := make([]compilerPassParam, 0, compilerPassType.NumField())
	for i := 0; i < compilerPassType.NumField(); i++ {
		field := compilerPassType.Field(i)
		fieldType := field.Type.Name()
		if field.Type.Name() == "" {
			fieldType = field.Type.String()
		}

		params = append(params, compilerPassParam{
			Name:          yamlNameForField(field),
			Type:          fieldType,
			Documentation: yamlCompilerPassComments[fmt.Sprintf("%s.%s", compilerPassType.Name(), field.Name)],
		})
	}

	return params
}

func buildYamlCompilerPassTypesCommentsMap(yamlTypesInputDir string) (map[string]string, error) {
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

func compilerPassDocEntries(compilerPassComments map[string]string, yamlCompilerPassComments map[string]string) []compilerPassDocEntry {
	var entries []compilerPassDocEntry
	yamlCompilerPassTypeOf := reflect.TypeOf(yaml.CompilerPass{})

	for i := 0; i < yamlCompilerPassTypeOf.NumField(); i++ {
		field := yamlCompilerPassTypeOf.Field(i)

		compilerPassMethod, found := field.Type.MethodByName("AsCompilerPass")
		if !found {
			continue
		}

		compilerPassType := compilerPassMethod.Type.Out(0)
		compilerPassTypeName := compilerPassType.Elem().Name()

		entries = append(entries, compilerPassDocEntry{
			YamlName:      yamlNameForField(field),
			Parameters:    compilerPassParameters(yamlCompilerPassComments, field.Type.Elem()),
			Documentation: compilerPassComments[compilerPassTypeName],
		})
	}

	sort.SliceStable(entries, func(i, j int) bool {
		return entries[i].YamlName < entries[j].YamlName
	})

	return entries
}

func buildCompilerPassTypesCommentsMap(typesInputDir string) (map[string]string, error) {
	compilerPassComments := make(map[string]string)

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

			compilerPassComments[t.Name] = t.Doc
		}
	}

	return compilerPassComments, nil
}

func docEntriesToMarkdown(entries []compilerPassDocEntry) []byte {
	var markdown bytes.Buffer

	markdown.WriteString(`---
weight: 10
---
`)
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
	compilerPassComments, err := buildCompilerPassTypesCommentsMap(compilerPassTypesSourceDir)
	if err != nil {
		panic(err)
	}

	yamlCompilerPassComments, err := buildYamlCompilerPassTypesCommentsMap(yamlCompilerPassTypesSourceDir)
	if err != nil {
		panic(err)
	}

	docEntries := compilerPassDocEntries(compilerPassComments, yamlCompilerPassComments)

	if err := os.WriteFile(outputFile, docEntriesToMarkdown(docEntries), 0600); err != nil {
		panic(err)
	}
}
