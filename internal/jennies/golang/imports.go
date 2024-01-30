package golang

import (
	"fmt"
	"strings"

	"github.com/grafana/cog/internal/jennies/common"
)

func NewImportMap() *common.DirectImportMap {
	return common.NewDirectImportMap(
		common.WithAliasSanitizer[common.DirectImportMap](func(alias string) string {
			return formatPackageName(strings.ReplaceAll(alias, "/", ""))
		}),
		common.WithFormatter(func(importMap common.DirectImportMap) string {
			if importMap.Imports.Len() == 0 {
				return ""
			}

			statements := make([]string, 0, importMap.Imports.Len())
			importMap.Imports.Iterate(func(alias string, importPath string) {
				statements = append(statements, fmt.Sprintf(`	%s "%s"`, alias, importPath))
			})

			return fmt.Sprintf(`import (
%[1]s
)`, strings.Join(statements, "\n"))
		}),
	)
}
