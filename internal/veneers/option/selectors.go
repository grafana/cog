package option

import (
	"fmt"
	"strings"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/tools"
)

type Selector struct {
	description string
	matcher     func(builder ast.Builder, option ast.Option) bool
}

func (selector Selector) Matches(builder ast.Builder, option ast.Option) bool {
	return selector.matcher(builder, option)
}

func (selector Selector) String() string {
	return selector.description
}

// EveryOption accepts any given option.
func EveryOption() *Selector {
	return &Selector{
		description: "every_option",
		matcher: func(_ ast.Builder, _ ast.Option) bool {
			return true
		},
	}
}

// ByName matches options by their name, defined for builders for the given
// object (referred to by its package and name).
// Note: the comparison on object and options names is case-insensitive.
func ByName(pkg string, objectName string, optionNames ...string) *Selector {
	return &Selector{
		description: fmt.Sprintf("by_name[builder.for.pkg='%s', option.name=(%s)]", objectName, strings.Join(optionNames, ", ")),
		matcher: func(builder ast.Builder, option ast.Option) bool {
			return (builder.For.SelfRef.ReferredPkg == pkg || pkg == "*") &&
				strings.EqualFold(builder.For.Name, objectName) &&
				tools.StringInListEqualFold(option.Name, optionNames)
		},
	}
}

// ByBuilder matches options by their name and the name of the builder containing them..
// Note: the comparison on builder and options names is case-insensitive.
func ByBuilder(pkg string, builderName string, optionNames ...string) *Selector {
	return &Selector{
		description: fmt.Sprintf("by_builder[builder.pkg='%s', builder.name='%s', option.name=(%s)]", pkg, builderName, strings.Join(optionNames, ", ")),
		matcher: func(builder ast.Builder, option ast.Option) bool {
			return (builder.Package == pkg || pkg == "*") &&
				strings.EqualFold(builder.Name, builderName) &&
				tools.StringInListEqualFold(option.Name, optionNames)
		},
	}
}
