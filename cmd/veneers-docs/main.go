package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"reflect"
	"sort"
	"strings"

	"github.com/grafana/cog/internal/tools"
	"github.com/grafana/cog/internal/yaml"
)

const builderRulesSource = "./internal/veneers/builder/rules.go"
const builderRulesOuputFile = "./docs/reference/builders_transformations.md"

const optionRulesSource = "./internal/veneers/option/rules.go"
const optionRulesOuputFile = "./docs/reference/options_transformations.md"

type Rule struct {
	Name string
	Doc  string
}

type veneerParam struct {
	Name          string
	Type          string
	Documentation string
}

type veneerDocEntry struct {
	YamlName      string
	Parameters    []veneerParam
	Documentation string
}

func yamlNameForField(field reflect.StructField) string {
	yamlName := field.Tag.Get("yaml")
	if yamlName == "" {
		yamlName = tools.SnakeCase(field.Name)
	}

	return yamlName
}

func ruleDocForYamlType(field reflect.StructField, rules map[string]Rule) string {
	ruleName := field.Tag.Get("rule_name")
	if ruleName == "" {
		ruleName = strings.TrimPrefix(field.Type.String(), "*yaml.")
	}

	return rules[ruleName].Doc
}

func ruleParameters(def reflect.Type) []veneerParam {
	var params []veneerParam

	for i := 0; i < def.NumField(); i++ {
		field := def.Field(i)
		fieldType := field.Type.Name()
		if field.Type.Name() == "" {
			fieldType = field.Type.String()
		}

		if fieldType == "BuilderSelector" || fieldType == "OptionSelector" {
			continue
		}

		params = append(params, veneerParam{
			Name: yamlNameForField(field),
			Type: fieldType,
		})
	}

	return params
}

func parseYamlRules(entrypointTypeOf reflect.Type, rules map[string]Rule) []veneerDocEntry {
	var entries []veneerDocEntry

	for i := 0; i < entrypointTypeOf.NumField(); i++ {
		field := entrypointTypeOf.Field(i)

		entries = append(entries, veneerDocEntry{
			YamlName:      yamlNameForField(field),
			Documentation: ruleDocForYamlType(field, rules),
			Parameters:    ruleParameters(field.Type.Elem()),
		})
	}

	sort.SliceStable(entries, func(i, j int) bool {
		return entries[i].YamlName < entries[j].YamlName
	})

	return entries
}

func parseRules(sourceFile string) (map[string]Rule, error) {
	fileAst, err := parser.ParseFile(token.NewFileSet(), sourceFile, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	rules := make(map[string]Rule)
	for _, decl := range fileAst.Decls {
		// We're looking for functions
		if _, ok := decl.(*ast.FuncDecl); !ok {
			continue
		}

		funcDecl := decl.(*ast.FuncDecl)
		funcName := funcDecl.Name.String()

		// That are public
		firstChar := funcName[0]
		if firstChar < 65 || firstChar > 90 {
			continue
		}

		// and that return a `Rule`
		results := funcDecl.Type.Results.List
		if len(results) != 1 {
			continue
		}

		var returnType *ast.Ident
		if expr, ok := results[0].Type.(*ast.StarExpr); ok {
			if ret, ok := expr.X.(*ast.Ident); !ok {
				continue
			} else {
				returnType = ret
			}
		} else {
			if expr, ok := results[0].Type.(*ast.Ident); !ok {
				continue
			} else {
				returnType = expr
			}
		}

		if returnType.String() != "Rule" {
			continue
		}

		rules[funcName] = Rule{
			Name: funcName,
			Doc:  strings.TrimSpace(funcDecl.Doc.Text()),
		}
	}

	return rules, nil
}

func builderDocEntriesToMarkdown(entries []veneerDocEntry) []byte {
	var markdown bytes.Buffer

	markdown.WriteString("<!-- Generated with `make docs` -->\n")

	markdown.WriteString("# Builder transformations\n\n")

	selectorTypeOf := reflect.TypeFor[yaml.BuilderSelector]()
	selectorParams := ruleParameters(selectorTypeOf)

	markdown.WriteString(`Each builder transformation requires the use of one of the following selectors, explicitly and unambiguously stating on which builder(s) the transformation should apply.`)
	markdown.WriteString("\n```yaml\n")
	for _, param := range selectorParams {
		if param.Documentation != "" {
			markdown.WriteString("  # " + strings.TrimSuffix(param.Documentation, "\n") + "\n")
		}
		markdown.WriteString(fmt.Sprintf("%s: %s\n", param.Name, strings.TrimPrefix(param.Type, "*")))
	}
	markdown.WriteString("```\n\n")

	markdown.WriteString("Example:\n```yaml\n")
	markdown.WriteString(`- rename:
    by_object: RowPanel
    as: Row
`)
	markdown.WriteString("```\n\n")

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

func optionDocEntriesToMarkdown(entries []veneerDocEntry) []byte {
	var markdown bytes.Buffer

	markdown.WriteString(`---
weight: 10
---
`)
	markdown.WriteString("<!-- Generated with `make docs` -->\n")

	markdown.WriteString("# Option transformations\n\n")

	selectorTypeOf := reflect.TypeFor[yaml.OptionSelector]()
	selectorParams := ruleParameters(selectorTypeOf)

	markdown.WriteString(`Each option transformation requires the use of one of the following selectors, explicitly and unambiguously stating on which option(s) the transformation should apply.`)
	markdown.WriteString("\n```yaml\n")
	for _, param := range selectorParams {
		if param.Documentation != "" {
			markdown.WriteString("  # " + strings.TrimSuffix(param.Documentation, "\n") + "\n")
		}
		markdown.WriteString(fmt.Sprintf("%s: %s\n", param.Name, strings.TrimPrefix(param.Type, "*")))
	}
	markdown.WriteString("```\n\n")

	markdown.WriteString("Example:\n```yaml\n")
	markdown.WriteString(`# H() → Height()
- rename:
    by_name: Panel.h
    as: height
`)
	markdown.WriteString("```\n\n")

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
	builderRules, err := parseRules(builderRulesSource)
	if err != nil {
		panic(err)
	}

	builderDocEntries := parseYamlRules(reflect.TypeFor[yaml.BuilderRule](), builderRules)
	if err := os.WriteFile(builderRulesOuputFile, builderDocEntriesToMarkdown(builderDocEntries), 0600); err != nil {
		panic(err)
	}

	optionRules, err := parseRules(optionRulesSource)
	if err != nil {
		panic(err)
	}

	optionDocEntries := parseYamlRules(reflect.TypeFor[yaml.OptionRule](), optionRules)
	if err := os.WriteFile(optionRulesOuputFile, optionDocEntriesToMarkdown(optionDocEntries), 0600); err != nil {
		panic(err)
	}
}
