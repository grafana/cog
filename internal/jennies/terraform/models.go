package terraform

import (
	"fmt"
	"strings"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/languages"
)

type modelFormatter struct {
	context             languages.Context
	packageMapper       func(pkg string) string
	refToStructAsObject bool
}

func defaultModelFormatter(context languages.Context, packageMapper func(pkg string) string) *modelFormatter {
	return &modelFormatter{
		context:             context,
		packageMapper:       packageMapper,
		refToStructAsObject: true,
	}
}

func modelFormatterWithRefs(context languages.Context, packageMapper func(pkg string) string) *modelFormatter {
	return &modelFormatter{
		context:             context,
		packageMapper:       packageMapper,
		refToStructAsObject: false,
	}
}

func (formatter *modelFormatter) formatDeclaration(object ast.Object) string {
	objectName := formatModelName(object.SelfRef)

	switch object.Type.Kind {
	case ast.KindScalar:
		if object.Type.AsScalar().Value != nil {
			return fmt.Sprintf("const %s = %s", objectName, formatScalar(object.Type.AsScalar().Value))
		}
		return fmt.Sprintf("type %s = %s", objectName, formatScalarAsModel(object.Type.AsScalar()))
	case ast.KindRef, ast.KindMap, ast.KindArray:
		return fmt.Sprintf("type %s = %s", objectName, formatter.formatModel(object.Type))
	default:
		return fmt.Sprintf("type %s %s", objectName, formatter.formatModel(object.Type))
	}
}

func (formatter *modelFormatter) formatModel(def ast.Type) string {
	switch def.Kind {
	case ast.KindScalar:
		if def.HasHint(ast.HintStringFormatDateTime) {
			formatter.packageMapper("github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes")
			return "timetypes.RFC3339"
		}
		if def.HasHint(ast.HintStringFormatDuration) {
			formatter.packageMapper("github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes")
			return "timetypes.GoDurationType"
		}
		return formatScalarAsModel(def.AsScalar())
	case ast.KindMap:
		return "types.Map"
	case ast.KindArray:
		return "types.List"
	case ast.KindRef:
		return formatter.formatReference(def)
	case ast.KindEnum:
		return formatter.formatModel(def.AsEnum().Values[0].Type)
	case ast.KindStruct:
		return formatter.formatStruct(def.AsStruct())
	default:
		return "unknown"
	}
}

func (formatter *modelFormatter) formatStruct(s ast.StructType) string {
	var buffer strings.Builder
	buffer.WriteString("struct {\n")
	for _, field := range s.Fields {
		for _, comment := range field.Comments {
			fmt.Fprintf(&buffer, "\t// %s\n", comment)
		}

		// constant refs shouldn't be exposed to users as their value can be set directly by the provider
		if field.Type.IsConstantRef() {
			continue
		}

		fmt.Fprintf(&buffer, "\t%s %s `tfsdk:\"%s\"`\n", formatModelFieldName(field.Name), formatter.formatModel(field.Type), formatTfSDKAttrName(field.Name))
	}
	buffer.WriteString("}")
	return buffer.String()
}

func (formatter *modelFormatter) formatReference(typeDef ast.Type) string {
	resolved := formatter.context.ResolveRefs(typeDef)
	if resolved.IsRef() {
		return fmt.Sprintf("could not resolve ref '%s'", typeDef.Ref.String())
	}

	if resolved.IsEnum() {
		return formatter.formatModel(resolved.AsEnum().Values[0].Type)
	}

	if resolved.IsStruct() {
		if formatter.refToStructAsObject {
			return "types.Object"
		}

		ref := typeDef.AsRef()
		pkg := formatter.packageMapper(ref.ReferredPkg)
		if pkg != "" {
			return pkg + "." + formatModelName(ref)
		}

		return formatModelName(ref)
	}

	return formatter.formatModel(resolved)
}
