package python

import (
	"fmt"
	"strings"

	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/orderedmap"
)

func NewImportMap(opts ...common.ImportMapOption[ModuleImportMap]) *ModuleImportMap {
	allOpts := []common.ImportMapOption[ModuleImportMap]{
		common.WithAliasSanitizer[ModuleImportMap](func(alias string) string {
			return strings.ReplaceAll(alias, "/", "")
		}),
		common.WithFormatter(func(importMap ModuleImportMap) string {
			if importMap.imports.Len() == 0 {
				return ""
			}

			statements := make([]string, 0, importMap.imports.Len())
			importMap.imports.Iterate(func(alias string, stmt ImportStmt) {
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
	}

	return NewModuleImportMap(append(allOpts, opts...)...)
}

type ImportStmt struct {
	Package string
	Module  string
}

type ModuleImportMap struct {
	// alias → ImportStmt
	imports         *orderedmap.Map[string, ImportStmt]
	config          common.ImportMapConfig[ModuleImportMap]
	currentPkgGuard func(pkg string) bool
}

func NewModuleImportMap(opts ...common.ImportMapOption[ModuleImportMap]) *ModuleImportMap {
	config := common.ImportMapConfig[ModuleImportMap]{
		Formatter: func(importMap ModuleImportMap) string {
			return fmt.Sprintf("%#v\n", importMap.imports)
		},
		AliasSanitizer:      common.NoopImportSanitizer,
		ImportPathSanitizer: common.NoopImportSanitizer,
	}

	for _, opt := range opts {
		opt(&config)
	}

	return &ModuleImportMap{
		imports: orderedmap.New[string, ImportStmt](),
		config:  config,
		currentPkgGuard: func(pkg string) bool {
			if config.CurrentPkg == "" {
				return false
			}

			return strings.TrimPrefix(pkg, ".") == config.CurrentPkg
		},
	}
}

func (im ModuleImportMap) Import(pkg string) string {
	return im.ImportAs(pkg, pkg)
}

func (im ModuleImportMap) ImportAs(pkg string, alias string) string {
	if im.currentPkgGuard(pkg) {
		return ""
	}

	sanitizedAlias := im.config.AliasSanitizer(alias)

	im.imports.Set(sanitizedAlias, ImportStmt{
		Package: pkg,
	})

	return sanitizedAlias
}

func (im ModuleImportMap) FromImport(pkg string, module string) string {
	return im.FromImportAs(pkg, module, module)
}

func (im ModuleImportMap) FromImportAs(pkg string, module string, alias string) string {
	if im.currentPkgGuard(module) {
		return ""
	}

	sanitizedAlias := im.config.AliasSanitizer(alias)

	im.imports.Set(sanitizedAlias, ImportStmt{
		Package: pkg,
		Module:  module,
	})

	return sanitizedAlias
}

func (im ModuleImportMap) String() string {
	im.imports.Sort(orderedmap.SortStrings)

	return im.config.Formatter(im)
}
