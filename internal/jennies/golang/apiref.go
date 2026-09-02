package golang

import (
	"fmt"
	"strings"

	"github.com/grafana/cog/internal/tools"
	"github.com/grafana/cog/pkg/apiref"
	"github.com/grafana/cog/pkg/ir"
	"github.com/grafana/cog/pkg/languages"
)

func apiReferenceFormatter(config Config) apiref.APIReferenceFormatter {
	builderName := func(builder ir.Builder) string {
		return formatObjectName(builder.Name) + "Builder"
	}
	methodSignature := func(context languages.Context, method apiref.MethodReference) string {
		args := tools.Map(method.Arguments, func(arg apiref.ArgumentReference) string {
			return fmt.Sprintf("%s %s", arg.Name, arg.Type)
		})

		returnType := ""
		if method.Return != "" {
			returnType = " " + method.Return
		}

		var receiverName, objectName string
		if method.ReceiverObject != nil {
			receiverName = formatArgName(method.ReceiverObject.Name)
			objectName = formatObjectName(method.ReceiverObject.Name)
		} else {
			receiverName = "builder"
			objectName = builderName(*method.ReceiverBuilder)
		}
		methodName := formatFunctionName(method.Name)

		return fmt.Sprintf("func (%[1]s *%[2]s) %[3]s(%[4]s)%[5]s", receiverName, objectName, methodName, strings.Join(args, ", "), returnType)
	}

	functionSignature := func(context languages.Context, function apiref.FunctionReference) string {
		args := tools.Map(function.Arguments, func(arg apiref.ArgumentReference) string {
			return fmt.Sprintf("%s %s", arg.Name, arg.Type)
		})

		returnType := ""
		if function.Return != "" {
			returnType = " " + function.Return
		}

		return fmt.Sprintf("func %[1]s(%[2]s)%[3]s", formatFunctionName(function.Name), strings.Join(args, ", "), returnType)
	}

	return apiref.APIReferenceFormatter{
		KindName: func(kind ir.Kind) string {
			return string(kind)
		},

		FunctionName: func(function apiref.FunctionReference) string {
			return formatFunctionName(function.Name)
		},
		FunctionSignature: functionSignature,

		ObjectName: func(object ir.Object) string {
			return formatObjectName(object.Name)
		},
		ObjectDefinition: func(context languages.Context, object ir.Object) string {
			dummyImports := NewImportMap("")
			typesFormatter := defaultTypeFormatter(config, context, dummyImports, func(pkg string) string {
				return pkg
			})
			return typesFormatter.formatTypeDeclaration(object)
		},

		MethodName: func(method apiref.MethodReference) string {
			return formatFunctionName(method.Name)
		},
		MethodSignature: methodSignature,

		BuilderName: builderName,
		ConstructorSignature: func(context languages.Context, builder ir.Builder) string {
			dummyImports := NewImportMap("")
			typesFormatter := builderTypeFormatter(config, context, dummyImports, func(pkg string) string {
				return pkg
			})
			args := tools.Map(builder.Constructor.Args, func(arg ir.Argument) apiref.ArgumentReference {
				return apiref.ArgumentReference{
					Name: formatArgName(arg.Name),
					Type: strings.TrimPrefix(typesFormatter.formatType(arg.Type), "*"),
				}
			})

			return functionSignature(context, apiref.FunctionReference{
				Name:      "New" + builderName(builder),
				Arguments: args,
				Return:    "*" + builderName(builder),
			})
		},
		OptionName: func(option ir.Option) string {
			return formatFunctionName(option.Name)
		},
		OptionSignature: func(context languages.Context, builder ir.Builder, option ir.Option) string {
			dummyImports := NewImportMap("")
			typesFormatter := builderTypeFormatter(config, context, dummyImports, func(pkg string) string {
				return pkg
			})
			args := tools.Map(option.Args, func(arg ir.Argument) apiref.ArgumentReference {
				return apiref.ArgumentReference{
					Name: formatArgName(arg.Name),
					Type: strings.TrimPrefix(typesFormatter.formatType(arg.Type), "*"),
				}
			})

			return methodSignature(context, apiref.MethodReference{
				ReceiverBuilder: &builder,
				Name:            option.Name,
				Comments:        option.Comments,
				Arguments:       args,
				Return:          "*" + builderName(builder),
			})
		},
	}
}
