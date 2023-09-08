package builder

import (
	"github.com/grafana/cog/internal/ast"
)

type Selector func(builder ast.Builder) bool

func ByName(objectName string) Selector {
	return func(builder ast.Builder) bool {
		return builder.For.Name == objectName
	}
}

func EveryBuilder() Selector {
	return func(builder ast.Builder) bool {
		return true
	}
}
