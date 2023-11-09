package template

import (
	"strings"
)

type ImportMap map[string]string

func NewImportMap() ImportMap {
	return make(map[string]string)
}

func (im ImportMap) Add(alias string, importPath string) string {
	sanitizedAlias := strings.ReplaceAll(alias, "/", "")
	im[sanitizedAlias] = importPath

	return sanitizedAlias
}
