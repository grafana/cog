package golang

import (
	"fmt"
	"strings"
)

type importMap map[string]string

func newImportMap() importMap {
	return make(map[string]string)
}

func (im importMap) Add(alias string, importPath string) {
	im[alias] = importPath
}

func (im importMap) Format() string {
	if len(im) == 0 {
		return ""
	}

	statements := make([]string, 0, len(im))

	for alias, importPath := range im {
		statements = append(statements, fmt.Sprintf(`	%s "%s"`, alias, importPath))
	}

	return fmt.Sprintf(`import (
%[1]s
)`, strings.Join(statements, "\n"))
}
