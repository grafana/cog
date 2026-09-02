package terraform

import (
	"fmt"
	"strings"

	"github.com/grafana/cog/pkg/imports"
)

func NewImportMap(packageRoot string) *imports.DirectImportMap {
	return imports.NewDirectImportMap(
		imports.WithAliasSanitizer[imports.DirectImportMap](formatPackageName),
		imports.WithImportPathSanitizer[imports.DirectImportMap](strings.ToLower),
		imports.WithFormatter(func(importMap imports.DirectImportMap) string {
			if importMap.Imports.Len() == 0 {
				return ""
			}

			statements := make([]string, 0, importMap.Imports.Len())
			importMap.Imports.Iterate(func(alias string, importPath string) {
				if strings.HasPrefix(importPath, packageRoot) {
					statements = append(statements, fmt.Sprintf(`	%s "%s"`, alias, importPath))
				} else { // stdlib import
					statements = append(statements, fmt.Sprintf(`	"%s"`, importPath))
				}
			})

			return fmt.Sprintf(`import (
%[1]s
)`, strings.Join(statements, "\n"))
		}),
	)
}
