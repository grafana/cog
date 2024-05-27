package java

import (
	"fmt"
	"strings"

	"github.com/grafana/cog/internal/jennies/common"
)

func NewImportMap() *common.DirectImportMap {
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
				statements = append(statements, fmt.Sprintf("import %s.%s;", setPackagePath(importPath), class))
			})

			return strings.Join(statements, "\n") + "\n"
		}),
	)
}

func setPackagePath(importPath string) string {
	ignorePaths := map[string]bool{
		"java.util": true,
	}
	if _, ok := ignorePaths[importPath]; ok {
		return importPath
	}

	return fmt.Sprintf("%s.%s", packagePath, importPath)
}
