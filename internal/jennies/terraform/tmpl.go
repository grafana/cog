package terraform

import (
	"fmt"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/jennies/template"
	"github.com/grafana/cog/internal/languages"
)

func initTemplates(config Config) *template.Template {
	tmpl, err := template.New("terraform",
		template.Funcs(common.TypeResolvingTemplateHelpers(languages.Context{})),
		template.Funcs(common.TypesTemplateHelpers(languages.Context{})),
		template.Funcs(template.FuncMap{
			// placeholder — overridden per-schema in RawTypes.generateSchema
			"importStdPkg": func(_ string) string {
				panic("importStdPkg() needs to be overridden by a jenny")
			},
			"toTfType": func(_ ast.Type) string {
				panic("toTfType() needs to be overridden by a jenny")
			},
			"toGoType": func(_ ast.Type) string {
				panic("toGoType() needs to be overridden by a jenny")
			},
			"toTfModel": func(_ ast.Type) string {
				panic("toTfModel() needs to be overridden by a jenny")
			},
			"toTfModelWithRefs": func(_ ast.Type) string {
				panic("toTfModelWithRefs() needs to be overridden by a jenny")
			},

			"formatScalar": formatScalar,

			// tfValueOf returns the Terraform value-getter method for a scalar field.
			// Nullable types use the pointer variant (e.g. "ValueStringPointer()" for *string).
			"tfValueOf": func(typeDef ast.Type, intoNullable bool) string {
				if !typeDef.IsScalar() {
					return ""
				}

				ptr := ""
				if intoNullable {
					ptr = "Pointer"
				}

				switch typeDef.Scalar.ScalarKind {
				case ast.KindString:
					return "ValueString" + ptr + "()"
				case ast.KindBool:
					return "ValueBool" + ptr + "()"
				case ast.KindFloat32, ast.KindFloat64:
					return "ValueFloat64" + ptr + "()"
				case ast.KindInt8, ast.KindUint8, ast.KindInt16, ast.KindUint16, ast.KindInt32, ast.KindUint32, ast.KindInt64, ast.KindUint64:
					return "ValueInt64" + ptr + "()"
				default:
					return fmt.Sprintf("unsupported scalar kind '%s'", typeDef.Scalar.ScalarKind)
				}
			},

			"tfTypeNullValueOf": func(typeDef ast.Type) string {
				if !typeDef.IsScalar() {
					return ""
				}

				switch typeDef.Scalar.ScalarKind {
				case ast.KindString:
					return "types.StringNull"
				case ast.KindBool:
					return "types.BoolNull"
				case ast.KindFloat32, ast.KindFloat64:
					return "types.Float64Null"
				case ast.KindInt8, ast.KindUint8, ast.KindInt16, ast.KindUint16, ast.KindInt32, ast.KindUint32, ast.KindInt64, ast.KindUint64:
					return "types.Int64Null"
				default:
					return fmt.Sprintf("unsupported tfTypeNullValueOf kind '%s'", typeDef.Scalar.ScalarKind)
				}
			},

			// tfTypeValueOf returns the Terraform constructor for converting a native Go value
			// to a Terraform SDK type. Nullable types use the pointer variant
			// (e.g. "types.StringPointerValue" for *string).
			"tfTypeValueOf": func(typeDef ast.Type, intoNullable bool) string {
				if !typeDef.IsScalar() {
					return ""
				}

				ptr := ""
				if intoNullable {
					ptr = "Pointer"
				}

				switch typeDef.Scalar.ScalarKind {
				case ast.KindString:
					return "types.String" + ptr + "Value"
				case ast.KindBool:
					return "types.Bool" + ptr + "Value"
				case ast.KindFloat32, ast.KindFloat64:
					return "types.Float64" + ptr + "Value"
				case ast.KindInt8, ast.KindUint8, ast.KindInt16, ast.KindUint16, ast.KindInt32, ast.KindUint32, ast.KindInt64, ast.KindUint64:
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
