package python

import (
	"fmt"
	"strings"

	"github.com/grafana/cog/internal/tools"
	"github.com/grafana/cog/pkg/apiref"
	"github.com/grafana/cog/pkg/ir"
	"github.com/grafana/cog/pkg/languages"
)

func apiReferenceFormatter() apiref.APIReferenceFormatter {
	return apiref.APIReferenceFormatter{
		KindName: func(kind ir.Kind) string {
			if kind == ir.KindStruct {
				return "class"
			}

			return string(kind)
		},

		FunctionName: func(function apiref.FunctionReference) string {
			return formatFunctionName(function.Name)
		},
		FunctionSignature: func(context languages.Context, function apiref.FunctionReference) string {
			args := tools.Map(function.Arguments, func(arg apiref.ArgumentReference) string {
				return fmt.Sprintf("%s: %s", arg.Name, arg.Type)
			})

			returnType := ""
			if function.Return != "" {
				returnType = " -> " + function.Return
			}

			return fmt.Sprintf("def %[1]s(%[2]s)%[3]s", formatFunctionName(function.Name), strings.Join(args, ", "), returnType)
		},

		ObjectName: func(object ir.Object) string {
			return formatObjectName(object.Name)
		},
		ObjectDefinition: func(context languages.Context, object ir.Object) string {
			typesFormatter := defaultTypeFormatter(context, func(alias string, pkg string) string {
				return alias
			}, func(alias string, pkg string, module string) string {
				return alias
			})
			formatted, err := typesFormatter.formatObject(object)
			if err != nil {
				return err.Error()
			}

			return formatted
		},

		MethodName: func(method apiref.MethodReference) string {
			return formatIdentifier(method.Name)
		},
		MethodSignature: func(context languages.Context, method apiref.MethodReference) string {
			args := tools.Map(method.Arguments, func(arg apiref.ArgumentReference) string {
				return fmt.Sprintf("%s: %s", arg.Name, arg.Type)
			})

			returnType := ""
			if method.Return != "" {
				returnType = " -> " + method.Return
			}

			methodName := formatIdentifier(method.Name)
			classmethod := ""
			if method.Static {
				classmethod = "@classmethod\n"
			}

			return fmt.Sprintf("%[1]sdef %[2]s(%[3]s)%[4]s", classmethod, methodName, strings.Join(args, ", "), returnType)
		},

		BuilderName: func(builder ir.Builder) string {
			return formatObjectName(builder.Name)
		},
		ConstructorSignature: func(context languages.Context, builder ir.Builder) string {
			typesFormatter := builderTypeFormatter(context, func(alias string, pkg string) string {
				return alias
			}, func(alias string, pkg string, module string) string {
				return alias
			})
			args := tools.Map(builder.Constructor.Args, func(arg ir.Argument) string {
				argType := typesFormatter.formatType(arg.Type)
				if argType != "" {
					argType = ": " + argType
				}

				return formatIdentifier(arg.Name) + argType
			})

			return fmt.Sprintf("%[1]s(%[2]s)", formatObjectName(builder.Name), strings.Join(args, ", "))
		},
		OptionName: func(option ir.Option) string {
			return formatIdentifier(option.Name)
		},
		OptionSignature: func(context languages.Context, builder ir.Builder, option ir.Option) string {
			typesFormatter := builderTypeFormatter(context, func(alias string, pkg string) string {
				return alias
			}, func(alias string, pkg string, module string) string {
				return alias
			})
			args := tools.Map(option.Args, func(arg ir.Argument) string {
				newArgType := arg.Type.DeepCopy()
				newArgType.Nullable = false

				argType := typesFormatter.formatType(newArgType)
				if argType != "" {
					argType = ": " + argType
				}

				return formatIdentifier(arg.Name) + argType
			})

			optionName := formatIdentifier(option.Name)

			return fmt.Sprintf("def %[1]s(%[2]s) -> typing.Self", optionName, strings.Join(args, ", "))
		},
	}
}
