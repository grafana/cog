package typescript

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
			importMap.Imports.Iterate(func(alias string, importPath string) {
				statements = append(statements, fmt.Sprintf(`import * as %s from '%s';`, alias, importPath))
			})

			return strings.Join(statements, "\n") + "\n"
		}),
	)
}
