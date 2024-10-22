package golang

import (
	"fmt"
	"strings"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/languages"
	"github.com/grafana/cog/internal/tools"
)

func apiReferenceFormatter(config Config) common.APIReferenceFormatter {
	return common.APIReferenceFormatter{
		ObjectName: func(object ast.Object) string {
			return tools.UpperCamelCase(object.Name)
		},
		ObjectDefinition: func(context languages.Context, object ast.Object) string {
			dummyImports := NewImportMap("")
			typesFormatter := defaultTypeFormatter(config, context, dummyImports, func(pkg string) string {
				return pkg
			})
			return typesFormatter.formatTypeDeclaration(object)
		},

		MethodName: func(method common.MethodReference) string {
			return tools.UpperCamelCase(method.Name)
		},
		MethodSignature: func(context languages.Context, object ast.Object, method common.MethodReference) string {
			args := tools.Map(method.Arguments, func(arg common.ArgumentReference) string {
				return fmt.Sprintf("%s %s", arg.Name, arg.Type)
			})

			returnType := ""
			if method.Return != "" {
				returnType = " " + method.Return
			}

			receiverName := tools.LowerCamelCase(object.Name)
			objectName := tools.UpperCamelCase(object.Name)
			methodName := tools.UpperCamelCase(method.Name)

			return fmt.Sprintf("func (%[1]s *%[2]s) %[3]s(%[4]s)%[5]s", receiverName, objectName, methodName, strings.Join(args, ", "), returnType)
		},

		BuilderName: func(builder ast.Builder) string {
			return tools.UpperCamelCase(builder.Name) + "Builder"
		},
		ConstructorSignature: func(context languages.Context, builder ast.Builder) string {
			dummyImports := NewImportMap("")
			typesFormatter := builderTypeFormatter(config, context, dummyImports, func(pkg string) string {
				return pkg
			})
			args := tools.Map(builder.Constructor.Args, func(arg ast.Argument) string {
				argType := typesFormatter.formatType(arg.Type)
				if argType != "" {
					argType = " " + argType
				}

				return formatArgName(arg.Name) + argType
			})

			return fmt.Sprintf("%[1]s(%[2]s)", tools.UpperCamelCase(builder.Name)+"Builder", strings.Join(args, ", "))
		},
		OptionName: func(option ast.Option) string {
			return tools.UpperCamelCase(option.Name)
		},
		OptionSignature: func(context languages.Context, builder ast.Builder, option ast.Option) string {
			dummyImports := NewImportMap("")
			typesFormatter := builderTypeFormatter(config, context, dummyImports, func(pkg string) string {
				return pkg
			})
			args := tools.Map(option.Args, func(arg ast.Argument) string {
				argType := typesFormatter.formatType(arg.Type)
				if argType != "" {
					argType = " " + strings.TrimPrefix(argType, "*")
				}

				return formatArgName(arg.Name) + argType
			})

			builderName := tools.UpperCamelCase(builder.Name) + "Builder"
			optionName := tools.UpperCamelCase(option.Name)

			return fmt.Sprintf("func (builder *%[1]s) %[2]s(%[3]s) *%[1]s", builderName, optionName, strings.Join(args, ", "))
		},
	}
}
