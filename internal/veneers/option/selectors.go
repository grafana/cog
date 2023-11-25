package option

import (
	"strings"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/tools"
)

type Selector func(builder ast.Builder, option ast.Option) bool

// Note: comparison on object name and option names is case insensitive
func ByName(pkg string, objectName string, optionNames ...string) Selector {
	return func(builder ast.Builder, option ast.Option) bool {
		return builder.For.SelfRef.ReferredPkg == pkg &&
			strings.EqualFold(builder.For.Name, objectName) &&
			tools.StringInListEqualFold(option.Name, optionNames)
	}
}
