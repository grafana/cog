package php

import (
	"fmt"
	"strings"

	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/tools"
	"github.com/grafana/cog/pkg/ir"
	"github.com/grafana/cog/pkg/languages"
	"github.com/grafana/cog/pkg/template"
)

func apiReferenceFormatter(tmpl *template.Template, config Config) common.APIReferenceFormatter {
	return common.APIReferenceFormatter{
		KindName: func(kind ir.Kind) string {
			if kind == ir.KindStruct {
				return "class"
			}

			return string(kind)
		},

		FunctionName: func(function common.FunctionReference) string {
			return formatOptionName(function.Name)
		},
		FunctionSignature: func(context languages.Context, function common.FunctionReference) string {
			args := tools.Map(function.Arguments, func(arg common.ArgumentReference) string {
				return fmt.Sprintf("%s $%s", arg.Type, arg.Name)
			})

			return fmt.Sprintf("%[1]s(%[2]s)", formatOptionName(function.Name), strings.Join(args, ", "))
		},

		ObjectName: func(object ir.Object) string {
			return formatObjectName(object.Name)
		},
		ObjectDefinition: func(context languages.Context, object ir.Object) string {
			typesFormatter := defaultTypeFormatter(config, context)
			return typesFormatter.formatTypeDeclaration(tmpl, context, object)
		},

		MethodName: func(method common.MethodReference) string {
			return formatOptionName(method.Name)
		},
		MethodSignature: func(context languages.Context, method common.MethodReference) string {
			args := tools.Map(method.Arguments, func(arg common.ArgumentReference) string {
				return fmt.Sprintf("%s $%s", arg.Type, arg.Name)
			})

			signature := fmt.Sprintf("%[1]s(%[2]s)", formatOptionName(method.Name), strings.Join(args, ", "))
			if method.Static {
				signature = "static " + signature
			}

			return signature
		},

		BuilderName: func(builder ir.Builder) string {
			return formatObjectName(builder.Name) + "Builder"
		},
		ConstructorSignature: func(context languages.Context, builder ir.Builder) string {
			typesFormatter := builderTypeFormatter(config, context)
			args := tools.Map(builder.Constructor.Args, func(arg ir.Argument) string {
				argType := typesFormatter.formatType(arg.Type)
				if argType != "" {
					argType += " "
				}

				return argType + "$" + formatArgName(arg.Name)
			})

			return fmt.Sprintf("new %[1]s(%[2]s)", formatObjectName(builder.Name)+"Builder", strings.Join(args, ", "))
		},
		OptionName: func(option ir.Option) string {
			return formatOptionName(option.Name)
		},
		OptionSignature: func(context languages.Context, builder ir.Builder, option ir.Option) string {
			typesFormatter := builderTypeFormatter(config, context)
			args := tools.Map(option.Args, func(arg ir.Argument) string {
				argType := typesFormatter.formatType(arg.Type)
				if argType != "" {
					argType += " "
				}

				return argType + "$" + formatArgName(arg.Name)
			})

			return fmt.Sprintf("%[1]s(%[2]s)", formatOptionName(option.Name), strings.Join(args, ", "))
		},
	}
}
