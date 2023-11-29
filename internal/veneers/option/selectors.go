package option

import (
	"strings"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/tools"
)

type Selector func(builder ast.Builder, option ast.Option) bool

// ByName matches options by their name, defined for builders for the given
// object (referred to by its package and name).
// Note: the comparison on object and options names is case-insensitive.
func ByName(pkg string, objectName string, optionNames ...string) Selector {
	return func(builder ast.Builder, option ast.Option) bool {
		return builder.For.SelfRef.ReferredPkg == pkg &&
			strings.EqualFold(builder.For.Name, objectName) &&
			tools.StringInListEqualFold(option.Name, optionNames)
	}
}
