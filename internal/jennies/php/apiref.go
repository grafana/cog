package php

import (
	"fmt"
	"strings"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/jennies/template"
	"github.com/grafana/cog/internal/languages"
	"github.com/grafana/cog/internal/tools"
)

func apiReferenceFormatter(tmpl *template.Template, config Config) common.APIReferenceFormatter {
	return common.APIReferenceFormatter{
		ObjectName: func(object ast.Object) string {
			return formatObjectName(object.Name)
		},
		ObjectDefinition: func(context languages.Context, object ast.Object) string {
			typesFormatter := defaultTypeFormatter(config, context)
			return typesFormatter.formatTypeDeclaration(tmpl, context, object)
		},
		BuilderName: func(builder ast.Builder) string {
			return formatObjectName(builder.Name) + "Builder"
		},
		ConstructorSignature: func(context languages.Context, builder ast.Builder) string {
			typesFormatter := builderTypeFormatter(config, context)
			args := tools.Map(builder.Constructor.Args, func(arg ast.Argument) string {
				argType := typesFormatter.formatType(arg.Type)
				if argType != "" {
					argType += " "
				}

				return argType + "$" + formatArgName(arg.Name)
			})

			return fmt.Sprintf("new %[1]s(%[2]s)", formatObjectName(builder.Name)+"Builder", strings.Join(args, ", "))

		},
		OptionName: func(option ast.Option) string {
			return formatOptionName(option.Name)
		},
		OptionSignature: func(context languages.Context, option ast.Option) string {
			typesFormatter := builderTypeFormatter(config, context)
			args := tools.Map(option.Args, func(arg ast.Argument) string {
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
