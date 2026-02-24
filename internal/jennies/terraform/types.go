package terraform

import (
	"fmt"
	"strings"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/languages"
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
		buffer.WriteString(fmt.Sprintf("%s %s `tfsdk:\"%s\"`", formatFieldName(field.Name), formatter.formatType(field.Type), field.Name))
		buffer.WriteString("\n")
	}
	return buffer.String()
}

func (formatter *typeFormatter) formatReference(ref ast.RefType) string {
	obj, ok := formatter.context.LocateObject(ref.ReferredPkg, ref.ReferredType)
	if !ok {
		return "types.Object" // We don't find the referenced object, so we assume it's a generic object
	}

	if obj.Type.IsEnum() {
		return formatter.formatType(obj.Type.AsEnum().Values[0].Type)
	}

	if obj.Type.IsConstantRef() {
		return formatter.formatConstantReference(obj.Type.AsConstantRef())
	}

	pkg := formatter.packageMapper(ref.ReferredPkg)
	if pkg != "" {
		return ref.ReferredPkg + "." + formatObjectName(ref.ReferredType)
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
		return "types.Object"
	}

	if obj.Type.IsScalar() {
		return formatter.formatScalar(obj.Type.AsScalar())
	}

	if obj.Type.IsEnum() {
		return formatter.formatType(obj.Type.AsEnum().Values[0].Type)
	}

	return "unknown"
}
