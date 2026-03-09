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
	validators    *validators
}

func defaultTypeFormatter(config Config, context languages.Context, imports *common.DirectImportMap, packageMapper func(pkg string) string) *typeFormatter {
	tf := &typeFormatter{
		config:        config,
		context:       context,
		imports:       imports,
		packageMapper: packageMapper,
	}

	tf.validators = newValidators(tf)
	return tf
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

func (formatter *typeFormatter) formatTypeAttributeForObject(obj ast.Object, comments string) string {
	return fmt.Sprintf("\"%s\": %s", tools.SnakeCase(obj.Name), formatter.formatTypeAttribute(obj.Type, comments))
}

func (formatter *typeFormatter) formatTypeAttribute(field ast.Type, comments string) string {
	switch field.Kind {
	case ast.KindScalar:
		return formatter.formatScalarAttribute(field)
	case ast.KindStruct:
		return formatter.formatStructAttributes(field, comments)
	case ast.KindArray:
		return formatter.formatArrayAttributes(field)
	case ast.KindMap:
		return formatter.formatMapAttributes(field)
	case ast.KindRef:
		return formatter.formatReferenceAttribute(field)
	case ast.KindEnum:
		return formatter.formatEnumAttribute(field)
	case ast.KindConstantRef:
		return formatter.formatConstantReferenceAttribute(field)
	default:
		return ""
	}
}

func (formatter *typeFormatter) formatStructAttributes(def ast.Type, comments string) string {
	var buffer strings.Builder
	buffer.WriteString("schema.SingleNestedAttribute{\n")

	if def.Nullable {
		buffer.WriteString("Optional: true,\n")
	} else {
		buffer.WriteString("Required: true,\n")
	}

	if comments != "" {
		buffer.WriteString(fmt.Sprintf("Description: %s,\n", comments))
	}

	buffer.WriteString("Attributes: map[string]schema.Attribute{\n")
	buffer.WriteString(formatter.formatFieldAttributes(def.AsStruct().Fields))
	buffer.WriteString("},\n},\n")
	return buffer.String()
}

func (formatter *typeFormatter) formatFieldAttributes(fields []ast.StructField) string {
	var buffer strings.Builder
	for _, field := range fields {
		if field.Type.IsIntersection() {
			continue
		}
		buffer.WriteString(fmt.Sprintf("\"%s\": %s\n", tools.SnakeCase(field.Name), formatter.formatTypeAttribute(field.Type, "")))
	}

	return buffer.String()
}

func (formatter *typeFormatter) formatArrayAttributes(def ast.Type) string {
	var buffer strings.Builder

	defVal := ""
	if def.Default != nil {
		formatter.packageMapper("github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault")
		formatter.packageMapper("github.com/hashicorp/terraform-plugin-framework/attr")
		defVal = fmt.Sprintf("Default: listdefault.StaticValue(%s),\n", formatter.parseArrayOrMapDefaults(def.AsArray().ValueType, def.Default, ListDefault))
	}

	switch def.AsArray().ValueType.Kind {
	case ast.KindRef:
		ref := def.AsArray().ValueType.AsRef()
		obj, ok := formatter.context.LocateObject(ref.ReferredPkg, ref.ReferredType)
		if !ok {
			return "unknown"
		}

		if !obj.Type.IsEnum() {
			buffer.WriteString("schema.ListNestedAttribute{\n")
			buffer.WriteString(fmt.Sprintf("NestedObject: schema.NestedAttributeObject {\n"))
			buffer.WriteString("Attributes: map[string]schema.Attribute {\n")
		}

		switch obj.Type.Kind {
		case ast.KindEnum:
			buffer.WriteString("schema.ListAttribute{\n")
			enumType := "types.StringType"
			if obj.Type.AsEnum().Values[0].Type.AsScalar().IsNumeric() {
				enumType = "types.Int64Type"
			}
			buffer.WriteString(fmt.Sprintf("ElementType: %s,\n", enumType))
		case ast.KindStruct:
			buffer.WriteString(formatter.formatFieldAttributes(obj.Type.AsStruct().Fields))
		default:
			buffer.WriteString(formatter.formatTypeAttributeForObject(obj, formatComments(obj.Comments)))
		}

		if defVal != "" {
			buffer.WriteString(defVal)
		}

		if !obj.Type.IsEnum() {
			buffer.WriteString("},\n},\n")
		}

	default:
		buffer.WriteString("schema.ListAttribute{\n ")
		buffer.WriteString(fmt.Sprintf("ElementType: %s,\n", formatter.formatElementType(def.AsArray().ValueType)))
		if defVal != "" {
			buffer.WriteString(defVal)
		}
	}
	buffer.WriteString(fmt.Sprintf("},\n"))

	return buffer.String()
}

func (formatter *typeFormatter) formatMapAttributes(def ast.Type) string {
	var buffer strings.Builder

	// TODO: It seems that we aren't parsing default values for maps, so we are not generating the default value for the map.
	defVal := ""
	if def.Default != nil {
		formatter.packageMapper("github.com/hashicorp/terraform-plugin-framework/resource/schema/mapdefault")
		formatter.packageMapper("github.com/hashicorp/terraform-plugin-framework/attr")
		defVal = fmt.Sprintf("Default: mapdefault.StaticValue(%s),\n", formatter.parseArrayOrMapDefaults(def.AsMap().ValueType, def.Default, MapDefault))
	}

	switch def.AsMap().ValueType.Kind {
	case ast.KindRef:
		ref := def.AsMap().ValueType.AsRef()
		obj, ok := formatter.context.LocateObject(ref.ReferredPkg, ref.ReferredType)
		if !ok {
			return "unknown"
		}

		if !obj.Type.IsEnum() {
			buffer.WriteString("schema.MapNestedAttribute{\n")
			buffer.WriteString(fmt.Sprintf("NestedObject: schema.NestedAttributeObject {\n"))
			buffer.WriteString("Attributes: map[string]schema.Attribute {\n")
		}

		switch obj.Type.Kind {
		case ast.KindEnum:
			buffer.WriteString("schema.MapAttribute{\n")
			enumType := "types.StringType"
			if obj.Type.AsEnum().Values[0].Type.AsScalar().IsNumeric() {
				enumType = "types.Int64Type"
			}
			buffer.WriteString(fmt.Sprintf("ElementType: %s,\n", enumType))
		case ast.KindStruct:
			buffer.WriteString(formatter.formatFieldAttributes(obj.Type.AsStruct().Fields))
		default:
			buffer.WriteString(formatter.formatTypeAttributeForObject(obj, formatComments(obj.Comments)))
		}

		if defVal != "" {
			buffer.WriteString(defVal)
		}

		if !obj.Type.IsEnum() {
			buffer.WriteString("},\n},\n")
		}
	default:
		buffer.WriteString("schema.MapAttribute{\n ")
		buffer.WriteString(fmt.Sprintf("ElementType: %s,\n", formatter.formatElementType(def.AsMap().ValueType)))
	}
	buffer.WriteString(fmt.Sprintf("},\n"))

	return buffer.String()
}

func (formatter *typeFormatter) formatScalarAttribute(def ast.Type) string {
	var attrs = map[ast.ScalarKind]attributeFormatter{
		ast.KindString:  {name: "StringAttribute", defImport: "stringdefault", defVal: "StaticString"},
		ast.KindBytes:   {name: "StringAttribute", defImport: "stringdefault", defVal: "StaticString"},
		ast.KindNull:    {name: "StringAttribute", defImport: "stringdefault", defVal: "StaticString"},
		ast.KindBool:    {name: "BoolAttribute", defImport: "booldefault", defVal: "StaticBool"},
		ast.KindInt32:   {name: "Int32Attribute", defImport: "int32default", defVal: "StaticInt32"},
		ast.KindUint32:  {name: "Int32Attribute", defImport: "int32default", defVal: "StaticInt32"},
		ast.KindInt64:   {name: "Int64Attribute", defImport: "int64default", defVal: "StaticInt64"},
		ast.KindUint64:  {name: "Int64Attribute", defImport: "int64default", defVal: "StaticInt64"},
		ast.KindFloat32: {name: "Float32Attribute", defImport: "float32default", defVal: "StaticFloat32"},
		ast.KindFloat64: {name: "Float64Attribute", defImport: "float64default", defVal: "StaticFloat64"},
		ast.KindAny:     {name: "ObjectAttribute", defImport: "objectdefault", defVal: "StaticValue"},
		ast.KindInt8:    {name: "NumberAttribute", defImport: "numberdefault", defVal: "StaticBigFloat"},
		ast.KindUint8:   {name: "NumberAttribute", defImport: "numberdefault", defVal: "StaticBigFloat"},
		ast.KindInt16:   {name: "NumberAttribute", defImport: "numberdefault", defVal: "StaticBigFloat"},
		ast.KindUint16:  {name: "NumberAttribute", defImport: "numberdefault", defVal: "StaticBigFloat"},
	}

	attr, ok := attrs[def.AsScalar().ScalarKind]
	if ok {
		required := fmt.Sprintf("Required: %v,\n", !def.Nullable)
		if def.Nullable {
			required = "Optional: true,\n"
		}
		defaultVal := ""
		if def.Default != nil {
			formatter.packageMapper(fmt.Sprintf("github.com/hashicorp/terraform-plugin-framework/resource/schema/%s", attr.defImport))
			defaultVal = fmt.Sprintf("Default: %s.%s(%s),\n", attr.defImport, attr.defVal, formatScalar(def.Default))
		} else if def.AsScalar().Value != nil {
			formatter.packageMapper(fmt.Sprintf("github.com/hashicorp/terraform-plugin-framework/resource/schema/%s", attr.defImport))
			defaultVal = fmt.Sprintf("Default: %s.%s(%s),\n", attr.defImport, attr.defVal, formatScalar(def.AsScalar().Value))
		}

		customType := ""
		if def.HasHint(ast.HintStringFormatDateTime) {
			formatter.packageMapper("github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes")
			customType = "CustomType: timetypes.RFC3339Type{},\n"
		}

		return fmt.Sprintf("schema.%s{\n %s%s%s%s},\n", attr.name, required, defaultVal, customType, formatter.validators.validate(def.AsScalar().ScalarKind, def.AsScalar().Constraints))
	}

	return "unknown"
}

type attributeFormatter struct {
	name      string
	defImport string
	defVal    string
}

func (formatter *typeFormatter) formatReferenceAttribute(def ast.Type) string {
	obj, ok := formatter.context.LocateObject(def.AsRef().ReferredPkg, def.AsRef().ReferredType)
	if !ok {
		return "unknown"
	}

	return formatter.formatTypeAttribute(obj.Type, formatComments(obj.Comments))
}

func (formatter *typeFormatter) formatEnumAttribute(def ast.Type) string {
	return formatter.formatScalarAttribute(def.AsEnum().Values[0].Type)
}

func (formatter *typeFormatter) formatConstantReferenceAttribute(def ast.Type) string {
	obj, ok := formatter.context.LocateObject(def.AsConstantRef().ReferredPkg, def.AsConstantRef().ReferredType)
	if !ok {
		return "unknown"
	}

	return formatter.formatTypeAttribute(obj.Type, formatComments(obj.Comments))
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
		return "types.DynamicType"
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

type defaultFormatter struct {
	name         string
	defFunc      string
	extraDefFunc string
}

type defaultType string

const MapDefault = "Map"
const ListDefault = "List"

func (formatter *typeFormatter) parseArrayOrMapDefaults(def ast.Type, defVal any, t defaultType) string {
	var scalarDef = map[ast.ScalarKind]defaultFormatter{
		ast.KindString:  {name: "StringType", defFunc: "StringValue"},
		ast.KindBytes:   {name: "StringType", defFunc: "StringValue"},
		ast.KindNull:    {name: "StringType", defFunc: "StringValue"},
		ast.KindBool:    {name: "BoolType", defFunc: "BoolValue"},
		ast.KindInt32:   {name: "Int32Type", defFunc: "Int32Value"},
		ast.KindUint32:  {name: "Int32Type", defFunc: "Int32Value"},
		ast.KindInt64:   {name: "Int64Type", defFunc: "Int64Value"},
		ast.KindUint64:  {name: "Int64Type", defFunc: "Int64Value"},
		ast.KindFloat32: {name: "Float32Type", defFunc: "Float32Value"},
		ast.KindFloat64: {name: "Float64Type", defFunc: "Float64Value"},
		ast.KindInt8:    {name: "NumberType", defFunc: "NumberValue", extraDefFunc: "big.NewFloat"},
		ast.KindUint8:   {name: "NumberType", defFunc: "NumberValue", extraDefFunc: "big.NewFloat"},
		ast.KindInt16:   {name: "NumberType", defFunc: "NumberValue", extraDefFunc: "big.NewFloat"},
		ast.KindUint16:  {name: "NumberType", defFunc: "NumberValue", extraDefFunc: "big.NewFloat"},
	}

	var enumDef = map[ast.ScalarKind]defaultFormatter{
		ast.KindString: {name: "StringType", defFunc: "StringValue"},
		ast.KindInt64:  {name: "StringType", defFunc: "StringValue"},
	}

	v := defVal.([]interface{})

	attrValue := "[]attr.Value"
	mustValue := "ListValueMust"
	if t == MapDefault {
		attrValue = "map[string]attr.Value"
		mustValue = "MapValueMust"
	}

	var buffer strings.Builder
	buffer.WriteString(fmt.Sprintf("types.%s(\n", mustValue))
	switch def.Kind {
	case ast.KindScalar:
		if scalar, ok := scalarDef[def.AsScalar().ScalarKind]; ok {
			buffer.WriteString(fmt.Sprintf("types.%s, %s{\n", scalar.name, attrValue))
			for _, val := range v {
				defaultValue := formatScalar(val)
				if scalar.extraDefFunc != "" {
					defaultValue = fmt.Sprintf("%s(%s)", scalar.extraDefFunc, defaultValue)
				}
				buffer.WriteString(fmt.Sprintf("types.%s(%s),\n", scalar.defFunc, defaultValue))
			}
			buffer.WriteString("},\n")
		}

	case ast.KindRef:
		obj, ok := formatter.context.LocateObject(def.AsRef().ReferredPkg, def.AsRef().ReferredType)
		if !ok {
			return "unknown"
		}
		if obj.Type.IsEnum() {
			if enum, ok := enumDef[obj.Type.AsEnum().Values[0].Type.AsScalar().ScalarKind]; ok {
				buffer.WriteString(fmt.Sprintf("types.%s, %s{\n", enum.name, attrValue))
				for _, val := range v {
					buffer.WriteString(fmt.Sprintf("types.%s(%s),\n", enum.defFunc, formatScalar(val)))
				}
				buffer.WriteString("},\n")
			}
		} else {
			return "unknown"
		}
	default:
		fmt.Println("Unknown type", def.Kind)
	}

	buffer.WriteString(")")
	return buffer.String()
}
