package option

import (
	"github.com/grafana/cog/internal/ast"
)

type Selector func(builder ast.Builder, option ast.Option) bool

func ByName(objectName string, optionName string) Selector {
	return func(builder ast.Builder, option ast.Option) bool {
		return builder.For.Name == objectName && option.Name == optionName
	}
}

func EveryOption() Selector {
	return func(builder ast.Builder, option ast.Option) bool {
		return true
	}
}
