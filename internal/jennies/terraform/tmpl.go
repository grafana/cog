package terraform

import (
	"fmt"

	"github.com/grafana/cog/pkg/ir"
	"github.com/grafana/cog/pkg/languages"
	"github.com/grafana/cog/pkg/template"
)

func initTemplates(config Config) *template.Template {
	tmpl, err := template.New("terraform",
		template.Funcs(template.TypeResolvingHelpers(languages.Context{})),
		template.Funcs(template.TypesHelpers(languages.Context{})),
		template.Funcs(template.FuncMap{
			// placeholder — overridden per-schema in RawTypes.generateSchema
			"importStdPkg": func(_ string) string {
				panic("importStdPkg() needs to be overridden by a jenny")
			},
			"toTfType": func(_ ir.Type) string {
				panic("toTfType() needs to be overridden by a jenny")
			},
			"toGoType": func(_ ir.Type) string {
				panic("toGoType() needs to be overridden by a jenny")
			},
			"toTfModel": func(_ ir.Type) string {
				panic("toTfModel() needs to be overridden by a jenny")
			},
			"toTfModelWithRefs": func(_ ir.Type) string {
				panic("toTfModelWithRefs() needs to be overridden by a jenny")
			},

			"formatScalar": formatScalar,

			// tfValueOf returns the Terraform value-getter method for a scalar field.
			// Nullable types use the pointer variant (e.g. "ValueStringPointer()" for *string).
			"tfValueOf": func(typeDef ir.Type, intoNullable bool) string {
				if !typeDef.IsScalar() {
					return ""
				}

				ptr := ""
				if intoNullable {
					ptr = "Pointer"
				}

				switch typeDef.Scalar.ScalarKind {
				case ir.KindString:
					return "ValueString" + ptr + "()"
				case ir.KindBool:
					return "ValueBool" + ptr + "()"
				case ir.KindFloat32, ir.KindFloat64:
					return "ValueFloat64" + ptr + "()"
				case ir.KindInt8, ir.KindUint8, ir.KindInt16, ir.KindUint16, ir.KindInt32, ir.KindUint32, ir.KindInt64, ir.KindUint64:
					return "ValueInt64" + ptr + "()"
				default:
					return fmt.Sprintf("unsupported scalar kind '%s'", typeDef.Scalar.ScalarKind)
				}
			},

			"tfTypeNullValueOf": func(typeDef ir.Type) string {
				if !typeDef.IsScalar() {
					return ""
				}

				switch typeDef.Scalar.ScalarKind {
				case ir.KindString:
					return "types.StringNull"
				case ir.KindBool:
					return "types.BoolNull"
				case ir.KindFloat32, ir.KindFloat64:
					return "types.Float64Null"
				case ir.KindInt8, ir.KindUint8, ir.KindInt16, ir.KindUint16, ir.KindInt32, ir.KindUint32, ir.KindInt64, ir.KindUint64:
					return "types.Int64Null"
				default:
					return fmt.Sprintf("unsupported tfTypeNullValueOf kind '%s'", typeDef.Scalar.ScalarKind)
				}
			},

			// tfTypeValueOf returns the Terraform constructor for converting a native Go value
			// to a Terraform SDK type. Nullable types use the pointer variant
			// (e.g. "types.StringPointerValue" for *string).
			"tfTypeValueOf": func(typeDef ir.Type, intoNullable bool) string {
				if !typeDef.IsScalar() {
					return ""
				}

				ptr := ""
				if intoNullable {
					ptr = "Pointer"
				}

				switch typeDef.Scalar.ScalarKind {
				case ir.KindString:
					return "types.String" + ptr + "Value"
				case ir.KindBool:
					return "types.Bool" + ptr + "Value"
				case ir.KindFloat32, ir.KindFloat64:
					return "types.Float64" + ptr + "Value"
				case ir.KindInt8, ir.KindUint8, ir.KindInt16, ir.KindUint16, ir.KindInt32, ir.KindUint32, ir.KindInt64, ir.KindUint64:
					return "types.Int64" + ptr + "Value"
				default:
					return fmt.Sprintf("unsupported scalar kind '%s'", typeDef.Scalar.ScalarKind)
				}
			},
		}),
		template.Funcs(config.OverridesTemplateFuncs),
		template.ParseDirectories(config.OverridesTemplatesDirectories...),
		template.ParseFS(config.OverridesTemplatesFS, "custom"),
	)

	if err != nil {
		panic(fmt.Errorf("could not initialize templates: %w", err))
	}
	return tmpl
}
