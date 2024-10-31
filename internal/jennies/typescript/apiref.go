package typescript

import (
	"fmt"
	"strings"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/languages"
	"github.com/grafana/cog/internal/tools"
)

func apiReferenceFormatter() common.APIReferenceFormatter {
	return common.APIReferenceFormatter{
		ObjectName: func(object ast.Object) string {
			return tools.CleanupNames(object.Name)
		},
		ObjectDefinition: func(context languages.Context, object ast.Object) string {
			typesFormatter := defaultTypeFormatter(context, func(pkg string) string {
				return pkg
			})

			return typesFormatter.formatTypeDeclaration(object)
		},

		MethodName: func(method common.MethodReference) string {
			return formatIdentifier(method.Name)
		},
		MethodSignature: func(context languages.Context, method common.MethodReference) string {
			args := tools.Map(method.Arguments, func(arg common.ArgumentReference) string {
				return fmt.Sprintf("%s: %s", arg.Name, arg.Type)
			})

			return fmt.Sprintf("%[1]s(%[2]s)", formatIdentifier(method.Name), strings.Join(args, ", "))
		},

		BuilderName: func(builder ast.Builder) string {
			return tools.UpperCamelCase(builder.Name) + "Builder"
		},
		ConstructorSignature: func(context languages.Context, builder ast.Builder) string {
			typesFormatter := builderTypeFormatter(context, func(pkg string) string {
				return pkg
			})
			args := tools.Map(builder.Constructor.Args, func(arg ast.Argument) string {
				return formatIdentifier(arg.Name) + ": " + typesFormatter.formatType(arg.Type)
			})

			return fmt.Sprintf("new %[1]s(%[2]s)", tools.UpperCamelCase(builder.Name)+"Builder", strings.Join(args, ", "))
		},
		OptionName: func(option ast.Option) string {
			return formatIdentifier(option.Name)
		},
		OptionSignature: func(context languages.Context, option ast.Option) string {
			typesFormatter := builderTypeFormatter(context, func(pkg string) string {
				return pkg
			})

			args := tools.Map(option.Args, func(arg ast.Argument) string {
				return formatIdentifier(arg.Name) + ": " + typesFormatter.formatType(arg.Type)
			})

			return fmt.Sprintf("%[1]s(%[2]s)", formatIdentifier(option.Name), strings.Join(args, ", "))
		},
	}
}
