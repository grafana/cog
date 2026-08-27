package terraform

import (
	"fmt"
	"strings"

	"github.com/grafana/cog/internal/ir"
	"github.com/grafana/cog/internal/languages"
)

type typeFormatter struct {
	context       languages.Context
	packageMapper func(pkg string) string
}

func defaultTypeFormatter(context languages.Context, packageMapper func(pkg string) string) *typeFormatter {
	return &typeFormatter{
		context:       context,
		packageMapper: packageMapper,
	}
}

func (formatter *typeFormatter) formatDeclaration(object ir.Object) string {
	return fmt.Sprintf("var %s = %s", formatTypeName(object.SelfRef), formatter.formatType(object.Type))
}

func (formatter *typeFormatter) formatType(def ir.Type) string {
	switch def.Kind {
	case ir.KindScalar:
		if def.HasHint(ir.HintStringFormatDateTime) {
			formatter.packageMapper("github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes")
			return "timetypes.RFC3339"
		}
		if def.HasHint(ir.HintStringFormatDuration) {
			formatter.packageMapper("github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes")
			return "timetypes.GoDurationType"
		}
		return formatter.formatScalarType(def.AsScalar())
	case ir.KindMap:
		return formatter.formatMapType(def.AsMap())
	case ir.KindArray:
		return formatter.formatArrayType(def.AsArray())
	case ir.KindStruct:
		return formatter.formatStructType(def.AsStruct())
	case ir.KindRef:
		return formatter.formatReference(def.AsRef())
	case ir.KindEnum:
		return formatter.formatType(def.AsEnum().Values[0].Type)
	default:
		return "unknown"
	}
}

func (formatter *typeFormatter) formatStructType(s ir.StructType) string {
	var buffer strings.Builder

	formatter.packageMapper("github.com/hashicorp/terraform-plugin-framework/attr")

	buffer.WriteString("types.ObjectType{\n")
	buffer.WriteString("\tAttrTypes: map[string]attr.Type{\n")

	for _, field := range s.Fields {
		// constant refs shouldn't be exposed to users as their value can be set directly by the provider
		if field.Type.IsConstantRef() {
			continue
		}

		fmt.Fprintf(&buffer, "\t\t\"%s\": %s,\n", formatTfSDKAttrName(field.Name), formatter.formatType(field.Type))
	}

	buffer.WriteString("\t},\n")
	buffer.WriteString("}")

	return buffer.String()
}

func (formatter *typeFormatter) formatScalarType(scalar ir.ScalarType) string {
	switch scalar.ScalarKind {
	case ir.KindString, ir.KindBytes, ir.KindNull:
		return "types.StringType"
	case ir.KindBool:
		return "types.BoolType"
	case ir.KindInt32, ir.KindUint32:
		return "types.Int32Type"
	case ir.KindInt64, ir.KindUint64:
		return "types.Int64Type"
	case ir.KindFloat32:
		return "types.Float32Type"
	case ir.KindFloat64:
		return "types.Float64Type"
	case ir.KindAny:
		// `any` should be represented as a string holding a JSON payload
		return "types.StringType"
	case ir.KindInt8, ir.KindUint8, ir.KindInt16, ir.KindUint16:
		return "types.NumberType"
	default:
		return fmt.Sprintf("unsupported scalar kind '%s'", scalar.ScalarKind)
	}
}

func (formatter *typeFormatter) formatReference(ref ir.RefType) string {
	obj, ok := formatter.context.GetObject(ref.ReferredPkg, ref.ReferredType)
	if !ok {
		return "unknown" // We don't find the referenced object, so we assume it's a generic object
	}

	if obj.Type.IsEnum() {
		return formatter.formatType(obj.Type.AsEnum().Values[0].Type)
	}

	pkg := formatter.packageMapper(ref.ReferredPkg)
	if pkg != "" {
		return pkg + "." + formatTypeName(ref)
	}

	return formatTypeName(ref)
}

func (formatter *typeFormatter) formatArrayType(array ir.ArrayType) string {
	return fmt.Sprintf("types.ListType{\n\tElemType: %s,\n}", formatter.formatType(array.ValueType))
}

func (formatter *typeFormatter) formatMapType(mapType ir.MapType) string {
	return fmt.Sprintf("types.MapType{\n\tElemType: %s,\n}", formatter.formatType(mapType.ValueType))
}
