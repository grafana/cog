package python

import (
	"fmt"
	"strings"

	"github.com/grafana/cog/internal/orderedmap"
	"github.com/grafana/cog/pkg/imports"
)

func NewImportMap() *ModuleImportMap {
	return NewModuleImportMap(
		imports.WithAliasSanitizer[ModuleImportMap](func(alias string) string {
			return strings.ReplaceAll(alias, "/", "")
		}),
		imports.WithFormatter(func(importMap ModuleImportMap) string {
			if importMap.Imports.Len() == 0 {
				return ""
			}

			statements := make([]string, 0, importMap.Imports.Len())
			importMap.Imports.Iterate(func(alias string, stmt ImportStmt) {
				var importStmt string
				if stmt.Module == "" {
					if stmt.Package == alias {
						importStmt = fmt.Sprintf(`import %s`, stmt.Package)
					} else {
						importStmt = fmt.Sprintf(`import %s as %s`, stmt.Package, alias)
					}
				} else {
					if stmt.Module == alias {
						importStmt = fmt.Sprintf(`from %s import %s`, stmt.Package, stmt.Module)
					} else {
						importStmt = fmt.Sprintf(`from %s import %s as %s`, stmt.Package, stmt.Module, alias)
					}
				}

				statements = append(statements, importStmt)
			})

			return strings.Join(statements, "\n")
		}),
	)
}

type ImportStmt struct {
	Package string
	Module  string
}

type ModuleImportMap struct {
	// alias → ImportStmt
	Imports *orderedmap.Map[string, ImportStmt]
	config  imports.ImportMapConfig[ModuleImportMap]
}

func NewModuleImportMap(opts ...imports.ImportMapOption[ModuleImportMap]) *ModuleImportMap {
	config := imports.ImportMapConfig[ModuleImportMap]{
		Formatter: func(importMap ModuleImportMap) string {
			return fmt.Sprintf("%#v\n", importMap.Imports)
		},
		AliasSanitizer:      imports.NoopImportSanitizer,
		ImportPathSanitizer: imports.NoopImportSanitizer,
	}

	for _, opt := range opts {
		opt(&config)
	}

	return &ModuleImportMap{
		Imports: orderedmap.New[string, ImportStmt](),
		config:  config,
	}
}

func (im ModuleImportMap) AddPackage(alias string, pkg string) string {
	sanitizedAlias := im.config.AliasSanitizer(alias)

	im.Imports.Set(sanitizedAlias, ImportStmt{
		Package: pkg,
	})

	return sanitizedAlias
}

func (im ModuleImportMap) AddModule(alias string, pkg string, module string) string {
	sanitizedAlias := im.config.AliasSanitizer(alias)

	im.Imports.Set(sanitizedAlias, ImportStmt{
		Package: pkg,
		Module:  module,
	})

	return sanitizedAlias
}

func (im ModuleImportMap) Sort() {
	im.Imports.Sort(orderedmap.SortStrings)
}

func (im ModuleImportMap) String() string {
	return im.config.Formatter(im)
}
