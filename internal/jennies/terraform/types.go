package terraform

import (
	"fmt"
	"strings"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/languages"
	"github.com/grafana/cog/internal/tools"
)

type typeFormatter struct {
	config        Config
	context       languages.Context
	imports       *common.DirectImportMap
	packageMapper func(pkg string) string
}

func defaultTypeFormatter(config Config, context languages.Context, imports *common.DirectImportMap, packageMapper func(pkg string) string) *typeFormatter {
	return &typeFormatter{
		config:        config,
		context:       context,
		imports:       imports,
		packageMapper: packageMapper,
	}
}

func (formatter *typeFormatter) formatTypeDeclaration(object ast.Object) string {
	objectName := formatObjectName(object.Name)
	switch object.Type.Kind {
	case ast.KindScalar:
		if object.Type.AsScalar().Value != nil {
			return fmt.Sprintf("const %s = %s", objectName, formatScalar(object.Type.AsScalar().Value))
		}
		return fmt.Sprintf("type %s %s", objectName, formatter.formatScalar(object.Type.AsScalar()))
	case ast.KindArray:
		return fmt.Sprintf("type %s %s", objectName, formatter.formatArray(object.Type.AsArray()))
	case ast.KindStruct:
		return fmt.Sprintf("type %s struct {\n %s }", objectName, formatter.formatStruct(object.Type.AsStruct()))
	case ast.KindRef:
		return fmt.Sprintf("type %s = %s", objectName, formatter.formatReference(object.Type.AsRef()))
	case ast.KindMap:
		return fmt.Sprintf("type %s types.Map", objectName)
	default:
		return ""
	}
}

func (formatter *typeFormatter) formatType(def ast.Type) string {
	switch def.Kind {
	case ast.KindScalar:
		if def.HasHint(ast.HintStringFormatDateTime) {
			formatter.packageMapper("github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes")
			return "timetypes.RFC3339"
		}
		return formatter.formatScalar(def.AsScalar())
	case ast.KindMap:
		return "types.Map"
	case ast.KindArray:
		return formatter.formatArray(def.AsArray())
	case ast.KindStruct:
		return "types.Object"
	case ast.KindRef:
		return formatter.formatReference(def.AsRef())
	case ast.KindConstantRef:
		return formatter.formatConstantReference(def.AsConstantRef())
	case ast.KindEnum:
		return formatter.formatType(def.AsEnum().Values[0].Type)
	default:
		return "unknown"
	}
}

func (formatter *typeFormatter) formatScalar(scalar ast.ScalarType) string {
	switch scalar.ScalarKind {
	case ast.KindString, ast.KindBytes, ast.KindNull:
		return "types.String"
	case ast.KindBool:
		return "types.Bool"
	case ast.KindInt32, ast.KindUint32:
		return "types.Int32"
	case ast.KindInt64, ast.KindUint64:
		return "types.Int64"
	case ast.KindFloat32:
		return "types.Float32"
	case ast.KindFloat64:
		return "types.Float64"
	case ast.KindAny:
		return "types.Object"
	case ast.KindInt8, ast.KindUint8, ast.KindInt16, ast.KindUint16:
		return "types.Number" // types.Number can be converted into any numeric type https://developer.hashicorp.com/terraform/plugin/framework/handling-data/types/number#setting-values
	default:
		return "unknown"
	}
}

func (formatter *typeFormatter) formatStruct(s ast.StructType) string {
	var buffer strings.Builder
	for _, field := range s.Fields {
		for _, comment := range field.Comments {
			buffer.WriteString(fmt.Sprintf("// %s\n", comment))
		}
		buffer.WriteString(fmt.Sprintf("%s %s `tfsdk:\"%s\"`", tools.UpperCamelCase(field.Name), formatter.formatType(field.Type), field.Name))
		buffer.WriteString("\n")
	}
	return buffer.String()
}

func (formatter *typeFormatter) formatReference(ref ast.RefType) string {
	obj, ok := formatter.context.LocateObject(ref.ReferredPkg, ref.ReferredType)
	if !ok {
		return "unknown" // We don't find the referenced object, so we assume it's a generic object
	}

	if obj.Type.IsEnum() {
		return formatter.formatType(obj.Type.AsEnum().Values[0].Type)
	}

	if obj.Type.IsConstantRef() {
		return formatter.formatConstantReference(obj.Type.AsConstantRef())
	}

	pkg := formatter.packageMapper(ref.ReferredPkg)
	if pkg != "" {
		return pkg + "." + formatObjectName(ref.ReferredType)
	}

	return formatObjectName(ref.ReferredType)
}

func (formatter *typeFormatter) formatArray(array ast.ArrayType) string {
	if array.IsArrayOf(ast.KindArray, ast.KindMap, ast.KindRef, ast.KindStruct) {
		return fmt.Sprintf("[]%s", formatter.formatType(array.ValueType))
	}

	return "types.List"
}

func (formatter *typeFormatter) formatConstantReference(constantRef ast.ConstantReferenceType) string {
	obj, ok := formatter.context.LocateObject(constantRef.ReferredPkg, constantRef.ReferredType)
	if !ok {
		return "unknown"
	}

	if obj.Type.IsScalar() {
		return formatter.formatScalar(obj.Type.AsScalar())
	}

	if obj.Type.IsEnum() {
		return formatter.formatType(obj.Type.AsEnum().Values[0].Type)
	}

	return "unknown"
}

// Schema attributes

func (formatter *typeFormatter) formatTypeAttribute(object ast.Object) string {
	comments := ""
	if len(object.Comments) > 0 {
		comments += "`\n"

		for _, comment := range object.Comments {
			comments += comment + "\n"
		}
		comments += "`"
	}

	switch object.Type.Kind {
	case ast.KindScalar:
		return formatter.formatScalarAttribute(object.Type)
	case ast.KindStruct:
		return formatter.formatStructAttributes(object.Type, comments)
	case ast.KindArray:
		return formatter.formatArrayAttributes(object.Type)
	case ast.KindMap:
		return formatter.formatMapAttributes(object.Type)
	default:
		return ""
	}
}

func (formatter *typeFormatter) formatStructAttributes(def ast.Type, comments string) string {
	var buffer strings.Builder
	buffer.WriteString("schema.ObjectAttribute{\n")

	if def.Nullable {
		buffer.WriteString("Optional: true,\n")
	} else {
		buffer.WriteString("Required: true,\n")
	}

	if comments != "" {
		buffer.WriteString(fmt.Sprintf("Description: %s,\n", comments))
	}

	formatter.packageMapper("github.com/hashicorp/terraform-plugin-framework/attr")
	buffer.WriteString("AttributeTypes: map[string]attr.Type{\n")
	for _, field := range def.AsStruct().Fields {
		if field.Type.IsIntersection() {
			continue
		}
		buffer.WriteString(fmt.Sprintf("\"%s\": %s,\n", tools.SnakeCase(field.Name), formatter.formatElementType(field.Type)))
	}

	buffer.WriteString("},\n},\n")
	return buffer.String()
}

func (formatter *typeFormatter) formatArrayAttributes(def ast.Type) string {
	var buffer strings.Builder
	buffer.WriteString("schema.ListAttribute{\n ")
	buffer.WriteString(fmt.Sprintf("ElementType: %s,\n", formatter.formatElementType(def.AsArray().ValueType)))
	buffer.WriteString(fmt.Sprintf("},\n"))

	return buffer.String()
}

func (formatter *typeFormatter) formatMapAttributes(def ast.Type) string {
	var buffer strings.Builder
	buffer.WriteString("schema.ListMapAttribute{\n ")
	buffer.WriteString(fmt.Sprintf("ElementType: %s,\n", formatter.formatElementType(def.AsMap().ValueType)))
	buffer.WriteString(fmt.Sprintf("},\n"))

	return buffer.String()
}

func (formatter *typeFormatter) formatScalarAttribute(def ast.Type) string {
	required := fmt.Sprintf("Required: %v,", !def.Nullable)
	if def.Nullable {
		required = "Optional: true,"
	}

	customType := "\n"
	if def.HasHint(ast.HintStringFormatDateTime) {
		formatter.packageMapper("github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes")
		customType = "\nCustomType: timetypes.RFC3339{},"
	}

	switch def.AsScalar().ScalarKind {
	case ast.KindString, ast.KindBytes, ast.KindNull:
		return fmt.Sprintf("schema.StringAttribute{\n %s%s\n},\n", required, customType)
	case ast.KindBool:
		return fmt.Sprintf("schema.BoolAttribute{\n %s \n},\n", required)
	case ast.KindInt32, ast.KindUint32:
		return fmt.Sprintf("schema.Int32Attribute{\n %s \n},\n", required)
	case ast.KindInt64, ast.KindUint64:
		return fmt.Sprintf("schema.Int64Attribute{\n %s \n},\n", required)
	case ast.KindFloat32:
		return fmt.Sprintf("schema.Float32Attribute{\n %s \n},\n", required)
	case ast.KindFloat64:
		return fmt.Sprintf("schema.Float64Attribute{\n %s \n},\n", required)
	case ast.KindAny:
		return fmt.Sprintf("schema.ObjectAttribute{\n %s \n},\n", required)
	case ast.KindInt8, ast.KindUint8, ast.KindInt16, ast.KindUint16:
		return fmt.Sprintf("schema.NumberAttribute{\n %s \n},\n", required) // types.Number can be converted into any numeric type https://developer.hashicorp.com/terraform/plugin/framework/handling-data/types/number#setting-values
	default:
		return "unknown"
	}
}

// ElementTypes defines the type of the elements for the attributes

func (formatter *typeFormatter) formatElementType(def ast.Type) string {
	switch def.Kind {
	case ast.KindScalar:
		return formatter.formatScalarAsElementType(def.AsScalar())
	case ast.KindArray:
		return formatter.formatArrayAsElementType(def.AsArray())
	case ast.KindRef:
		return formatter.formatReferenceAsElementType(def.AsRef())
	case ast.KindMap:
		return formatter.formatMapAsElementType(def.AsMap())
	case ast.KindStruct:
		return formatter.formatStructAsElementType(def.AsStruct())
	case ast.KindConstantRef:
		return formatter.formatConstantReferenceAsElementType(def.AsConstantRef())
	case ast.KindEnum:
		return formatter.formatEnumAsElementType(def.AsEnum())
	default:
		return "unknown"
	}
}

func (formatter *typeFormatter) formatScalarAsElementType(def ast.ScalarType) string {
	switch def.ScalarKind {
	case ast.KindString, ast.KindBytes, ast.KindNull:
		return "types.StringType"
	case ast.KindBool:
		return "types.BoolType"
	case ast.KindInt32, ast.KindUint32:
		return "types.Int32Type"
	case ast.KindInt64, ast.KindUint64:
		return "types.Int64Type"
	case ast.KindFloat32:
		return "types.Float32Type"
	case ast.KindFloat64:
		return "types.Float64Type"
	case ast.KindAny:
		return "types.ObjectType{}"
	case ast.KindInt8, ast.KindUint8, ast.KindInt16, ast.KindUint16:
		return "types.NumberType" // types.Number can be converted into any numeric type https://developer.hashicorp.com/terraform/plugin/framework/handling-data/types/number#setting-values
	default:
		return "unknown"
	}
}

func (formatter *typeFormatter) formatArrayAsElementType(def ast.ArrayType) string {
	return fmt.Sprintf("types.ListType{\n ElemType: %s,\n}", formatter.formatElementType(def.ValueType))
}

func (formatter *typeFormatter) formatMapAsElementType(def ast.MapType) string {
	return fmt.Sprintf("types.MapType{\n ElemType: %s,\n}", formatter.formatElementType(def.ValueType))
}

func (formatter *typeFormatter) formatReferenceAsElementType(ref ast.RefType) string {
	var buffer strings.Builder

	obj, ok := formatter.context.LocateObject(ref.ReferredPkg, ref.ReferredType)
	if !ok {
		return "unknown"
	}

	buffer.WriteString(formatter.formatElementType(obj.Type))
	return buffer.String()
}

func (formatter *typeFormatter) formatStructAsElementType(s ast.StructType) string {
	formatter.packageMapper("github.com/hashicorp/terraform-plugin-framework/attr")
	var buffer strings.Builder

	buffer.WriteString("types.ObjectType{\n AttrTypes: map[string]attr.Type{\n")
	for _, field := range s.Fields {
		buffer.WriteString(fmt.Sprintf("\"%s\": %s,\n", tools.LowerCamelCase(field.Name), formatter.formatElementType(field.Type)))
	}
	buffer.WriteString("},\n}")
	return buffer.String()
}

func (formatter *typeFormatter) formatConstantReferenceAsElementType(ref ast.ConstantReferenceType) string {
	obj, ok := formatter.context.LocateObject(ref.ReferredPkg, ref.ReferredType)
	if !ok {
		return "unknown"
	}

	return formatter.formatElementType(obj.Type)
}

func (formatter *typeFormatter) formatEnumAsElementType(enum ast.EnumType) string {
	return formatter.formatScalarAsElementType(enum.Values[0].Type.AsScalar())
}
