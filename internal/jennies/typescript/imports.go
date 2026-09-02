package typescript

import (
	"fmt"
	"strings"

	"github.com/grafana/cog/internal/tools"
	"github.com/grafana/cog/pkg/imports"
)

func NewImportMap(packagesImportMap map[string]string) *imports.DirectImportMap {
	return imports.NewDirectImportMap(
		imports.WithPackagesImportMap[imports.DirectImportMap](packagesImportMap),
		imports.WithAliasSanitizer[imports.DirectImportMap](formatPackageName),
		imports.WithImportPathSanitizer[imports.DirectImportMap](func(importPath string) string {
			parts := strings.Split(importPath, "/")

			return strings.Join(tools.Map(parts, func(input string) string {
				if input == ".." {
					return input
				}

				return formatPackageName(input)
			}), "/")
		}),
		imports.WithFormatter(func(importMap imports.DirectImportMap) string {
			if importMap.Imports.Len() == 0 {
				return ""
			}

			statements := make([]string, 0, importMap.Imports.Len())
			importMap.Imports.Iterate(func(alias string, importPath string) {
				statements = append(statements, fmt.Sprintf(`import * as %s from '%s';`, alias, importPath))
			})

			return strings.Join(statements, "\n") + "\n"
		}),
	)
}
