package java

import (
	"fmt"
	"strings"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/languages"
	"github.com/grafana/cog/internal/tools"
)

func apiReferenceFormatter(config Config) common.APIReferenceFormatter {
	pkgMapper := func(pkg string, class string) string {
		return pkg
	}

	return common.APIReferenceFormatter{
		KindName: func(kind ast.Kind) string {
			if kind == ast.KindStruct {
				return "class"
			}

			return string(kind)
		},

		FunctionName: func(function common.FunctionReference) string {
			return tools.LowerCamelCase(function.Name)
		},
		FunctionSignature: func(context languages.Context, function common.FunctionReference) string {
			args := tools.Map(function.Arguments, func(arg common.ArgumentReference) string {
				return fmt.Sprintf("%s %s", arg.Type, arg.Name)
			})

			return fmt.Sprintf("public %[1]s(%[2]s)", tools.LowerCamelCase(function.Name), strings.Join(args, ", "))
		},

		ObjectName: func(object ast.Object) string {
			return formatObjectName(object.Name)
		},
		ObjectDefinition: func(context languages.Context, object ast.Object) string {
			typesFormatter := createFormatter(context, config).withPackageMapper(pkgMapper)
			return typesFormatter.formatFieldType(object.Type)
		},

		MethodName: func(method common.MethodReference) string {
			return tools.LowerCamelCase(method.Name)
		},
		MethodSignature: func(context languages.Context, method common.MethodReference) string {
			args := tools.Map(method.Arguments, func(arg common.ArgumentReference) string {
				return fmt.Sprintf("%s %s", arg.Type, arg.Name)
			})

			return fmt.Sprintf("public %[1]s %[2]s(%[3]s)", method.Return, tools.LowerCamelCase(method.Name), strings.Join(args, ", "))
		},

		BuilderName: func(builder ast.Builder) string {
			return formatObjectName(builder.Name) + "Builder"
		},
		ConstructorSignature: func(context languages.Context, builder ast.Builder) string {
			typesFormatter := createFormatter(context, config).withPackageMapper(func(pkg string, class string) string {
				return pkg
			})
			args := tools.Map(builder.Constructor.Args, func(arg ast.Argument) string {
				argType := typesFormatter.formatFieldType(arg.Type)
				return argType + " " + formatArgName(arg.Name)
			})

			return fmt.Sprintf("new %[1]s(%[2]s)", formatObjectName(builder.Name)+"Builder", strings.Join(args, ", "))
		},
		OptionName: func(option ast.Option) string {
			return tools.LowerCamelCase(option.Name)
		},
		OptionSignature: func(context languages.Context, builder ast.Builder, option ast.Option) string {
			typesFormatter := createFormatter(context, config).withPackageMapper(pkgMapper)
			args := tools.Map(option.Args, func(arg ast.Argument) string {
				argType := typesFormatter.formatFieldType(arg.Type)
				if argType != "" {
					argType += " "
				}

				return argType + formatArgName(arg.Name)
			})

			return fmt.Sprintf("public %[1]sBuilder %[2]s(%[3]s)", builder.Name, tools.LowerCamelCase(option.Name), strings.Join(args, ", "))
		},
	}
}
