package java

import (
	"fmt"
	"strings"

	"github.com/grafana/cog/pkg/imports"
)

func NewImportMap(packagePath string) *imports.DirectImportMap {
	return imports.NewDirectImportMap(
		imports.WithAliasSanitizer[imports.DirectImportMap](func(alias string) string {
			return strings.ReplaceAll(alias, "/", "")
		}),
		imports.WithFormatter(func(importMap imports.DirectImportMap) string {
			if importMap.Imports.Len() == 0 {
				return ""
			}

			statements := make([]string, 0, importMap.Imports.Len())
			importMap.Imports.Iterate(func(class string, importPath string) {
				statements = append(statements, fmt.Sprintf("import %s.%s;", setPackagePath(packagePath, importPath), class))
			})

			return strings.Join(statements, "\n") + "\n"
		}),
	)
}

func setPackagePath(packagePath string, importPath string) string {
	ignorePaths := map[string]bool{
		"com.fasterxml.jackson": true,
		"java.lang":             true,
		"java.util":             true,
	}
	if _, ok := ignorePaths[importPath]; ok || packagePath == "" {
		return importPath
	}

	return fmt.Sprintf("%s.%s", packagePath, importPath)
}
