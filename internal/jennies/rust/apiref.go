package rust

import (
	"fmt"
	"strings"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/languages"
	"github.com/grafana/cog/internal/tools"
)

// apiReferenceFormatter renders the Rust-flavored signatures shown in the
// generated markdown API reference, mirroring the Python and Go targets. The
// type formatters use a throwaway import map: the reference shows short type
// names, never `use` statements.
func apiReferenceFormatter() common.APIReferenceFormatter {
	return common.APIReferenceFormatter{
		KindName: func(kind ast.Kind) string {
			return string(kind)
		},

		FunctionName: func(function common.FunctionReference) string {
			return formatArgName(function.Name)
		},
		FunctionSignature: func(_ languages.Context, function common.FunctionReference) string {
			args := tools.Map(function.Arguments, func(arg common.ArgumentReference) string {
				return fmt.Sprintf("%s: %s", formatArgName(arg.Name), arg.Type)
			})

			returnType := ""
			if function.Return != "" {
				returnType = " -> " + function.Return
			}

			return fmt.Sprintf("pub fn %[1]s(%[2]s)%[3]s", formatArgName(function.Name), strings.Join(args, ", "), returnType)
		},

		ObjectName: func(object ast.Object) string {
			return formatTypeName(object.Name)
		},
		ObjectDefinition: apiRefObjectDefinition,

		MethodName: func(method common.MethodReference) string {
			return formatArgName(method.Name)
		},
		MethodSignature: func(_ languages.Context, method common.MethodReference) string {
			args := tools.Map(method.Arguments, func(arg common.ArgumentReference) string {
				return fmt.Sprintf("%s: %s", formatArgName(arg.Name), arg.Type)
			})
			receiver := "&self"
			if method.Static {
				receiver = ""
			} else if len(args) > 0 {
				receiver += ", "
			}

			returnType := ""
			if method.Return != "" {
				returnType = " -> " + method.Return
			}

			return fmt.Sprintf("pub fn %[1]s(%[2]s%[3]s)%[4]s", formatArgName(method.Name), receiver, strings.Join(args, ", "), returnType)
		},

		BuilderName: func(builder ast.Builder) string {
			return formatTypeName(builder.Name) + "Builder"
		},
		ConstructorSignature: func(context languages.Context, builder ast.Builder) string {
			formatter := newTypeFormatter(context, newImportMap(), builder.Package)
			args := tools.Map(builder.Constructor.Args, func(arg ast.Argument) string {
				return fmt.Sprintf("%s: %s", formatArgName(arg.Name), formatter.formatType(arg.Type))
			})

			return fmt.Sprintf("%sBuilder::new(%s)", formatTypeName(builder.Name), strings.Join(args, ", "))
		},
		OptionName: func(option ast.Option) string {
			return formatFieldName(option.Name)
		},
		OptionSignature: func(context languages.Context, builder ast.Builder, option ast.Option) string {
			formatter := newTypeFormatter(context, newImportMap(), builder.Package)
			args := tools.Map(option.Args, func(arg ast.Argument) string {
				argType := arg.Type.DeepCopy()
				argType.Nullable = false

				return fmt.Sprintf("%s: %s", formatArgName(arg.Name), formatter.formatType(argType))
			})

			return fmt.Sprintf("pub fn %[1]s(self, %[2]s) -> Self", formatFieldName(option.Name), strings.Join(args, ", "))
		},
	}
}

// apiRefObjectDefinition renders a compact Rust-flavored definition of an
// object for the API reference: struct fields, enum members or the aliased
// type. It intentionally omits derives and serde attributes; the reference
// documents the shape, not the wire format.
func apiRefObjectDefinition(context languages.Context, object ast.Object) string {
	formatter := newTypeFormatter(context, newImportMap(), object.SelfRef.ReferredPkg)

	switch {
	case object.Type.IsStruct():
		var buffer strings.Builder
		fmt.Fprintf(&buffer, "pub struct %s {\n", formatTypeName(object.Name))
		for _, field := range object.Type.AsStruct().Fields {
			fmt.Fprintf(&buffer, "    pub %s: %s,\n", formatFieldName(field.Name), formatter.formatType(field.Type))
		}
		buffer.WriteString("}")
		return buffer.String()
	case object.Type.IsEnum():
		var buffer strings.Builder
		fmt.Fprintf(&buffer, "pub enum %s {\n", formatTypeName(object.Name))
		for _, member := range object.Type.AsEnum().Values {
			fmt.Fprintf(&buffer, "    %s,\n", formatTypeName(member.Name))
		}
		buffer.WriteString("}")
		return buffer.String()
	case object.Type.IsConcreteScalar():
		return fmt.Sprintf("pub const %s: %s = %s;", formatConstName(object.Name), formatter.formatType(object.Type), formatValue(object.Type.AsScalar().Value))
	default:
		return fmt.Sprintf("pub type %s = %s;", formatTypeName(object.Name), formatter.formatType(object.Type))
	}
}
