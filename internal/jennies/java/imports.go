package java

import (
	"fmt"
	"strings"

	"github.com/grafana/cog/internal/jennies/common"
)

var ignorePaths = map[string]bool{
	"java.util":             true,
	"com.fasterxml.jackson": true,
}

func NewImportMap(pkgPrefix string) *common.DirectImportMap {
	return common.NewDirectImportMap(
		common.WithAliasSanitizer[common.DirectImportMap](func(alias string) string {
			return strings.ReplaceAll(alias, "/", "")
		}),
		common.WithFormatter(func(importMap common.DirectImportMap) string {
			if importMap.Imports.Len() == 0 {
				return ""
			}

			statements := make([]string, 0, importMap.Imports.Len())
			importMap.Imports.Iterate(func(class string, importPath string) {
				statements = append(statements, fmt.Sprintf("import %s.%s;", setPathPrefix(pkgPrefix, importPath), class))
			})

			return strings.Join(statements, "\n") + "\n"
		}),
	)
}

func setPathPrefix(prefix string, importPath string) string {
	if _, ok := ignorePaths[importPath]; ok || prefix == "" {
		return importPath
	}

	return fmt.Sprintf("%s.%s", prefix, importPath)
}
