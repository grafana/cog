package typescript

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
		statements = append(statements, fmt.Sprintf(`import * as %s from "%s";`, alias, importPath))
	}

	return strings.Join(statements, "\n")
}
