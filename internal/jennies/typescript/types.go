package typescript

import (
	"fmt"
	"strings"

	"github.com/grafana/cog/internal/tools"
	"github.com/grafana/cog/pkg/ir"
	"github.com/grafana/cog/pkg/languages"
)

type enumFormatter interface {
	formatDeclaration(def ir.Object) string
	formatValue(enumObj ir.Object, val any) string
}

type packageMapper func(pkg string) string

type typeFormatter struct {
	packageMapper func(pkg string) string
	enums         enumFormatter
	forBuilder    bool
	context       languages.Context
}

func defaultTypeFormatter(config Config, context languages.Context, packageMapper packageMapper) *typeFormatter {
	return &typeFormatter{
		packageMapper: packageMapper,
		context:       context,
		enums:         config.enumFormatter(packageMapper),
	}
}

func builderTypeFormatter(config Config, context languages.Context, packageMapper packageMapper) *typeFormatter {
	return &typeFormatter{
		packageMapper: packageMapper,
		forBuilder:    true,
		context:       context,
		enums:         config.enumFormatter(packageMapper),
	}
}

func (formatter *typeFormatter) variantInterface(variant string) string {
	referredPkg := formatter.packageMapper("cog")

	return fmt.Sprintf("%s.%s", referredPkg, tools.UpperCamelCase(variant))
}

func (formatter *typeFormatter) formatTypeDeclaration(def ir.Object) string {
	var buffer strings.Builder

	buffer.WriteString("export ")

	objectName := formatObjectName(def.Name)

	switch def.Type.Kind {
	case ir.KindStruct:
		buffer.WriteString(fmt.Sprintf("interface %s ", objectName))
		buffer.WriteString(formatter.formatStructFields(def.Type))
		buffer.WriteString("\n")
	case ir.KindEnum:
		buffer.WriteString(formatter.enums.formatDeclaration(def))
		buffer.WriteString("\n")
	case ir.KindDisjunction, ir.KindMap, ir.KindArray, ir.KindRef:
		buffer.WriteString(fmt.Sprintf("type %s = %s;\n", objectName, formatter.formatType(def.Type)))
	case ir.KindScalar:
		scalarType := def.Type.AsScalar()
		typeValue := formatValue(scalarType.Value)

		if !scalarType.IsConcrete() || def.Type.Hints["kind"] == "type" {
			if !scalarType.IsConcrete() {
				typeValue = formatter.formatScalarKind(scalarType.ScalarKind)
			}

			buffer.WriteString(fmt.Sprintf("type %s = %s;\n", objectName, typeValue))
		} else {
			buffer.WriteString(fmt.Sprintf("const %s = %s;\n", objectName, typeValue))
		}
	case ir.KindIntersection:
		buffer.WriteString(fmt.Sprintf("interface %s ", objectName))
		buffer.WriteString(formatter.formatType(def.Type))
		buffer.WriteString("\n")
	case ir.KindComposableSlot:
		buffer.WriteString(fmt.Sprintf("interface %s %s\n", objectName, formatter.variantInterface(string(def.Type.AsComposableSlot().Variant))))
	default:
		return fmt.Sprintf("unhandled object of type: %s", def.Type.Kind)
	}

	return buffer.String()
}

func (formatter *typeFormatter) formatType(def ir.Type) string {
	return formatter.doFormatType(def, formatter.forBuilder)
}

func (formatter *typeFormatter) doFormatType(def ir.Type, resolveBuilders bool) string {
	switch def.Kind {
	case ir.KindDisjunction:
		return formatter.formatDisjunction(def.AsDisjunction(), resolveBuilders)
	case ir.KindRef:
		return formatter.formatRef(def.AsRef(), resolveBuilders)
	case ir.KindArray:
		return formatter.formatArray(def.AsArray(), resolveBuilders)
	case ir.KindStruct:
		return formatter.formatStructFields(def)
	case ir.KindMap:
		return formatter.formatMap(def.AsMap(), resolveBuilders)
	case ir.KindEnum:
		return formatter.formatAnonymousEnum(def.AsEnum())
	case ir.KindScalar:
		// This scalar actually refers to a constant
		if def.AsScalar().Value != nil {
			return formatValue(def.AsScalar().Value)
		}

		return formatter.formatScalarKind(def.AsScalar().ScalarKind)
	case ir.KindIntersection:
		return formatter.formatIntersection(def.AsIntersection())
	case ir.KindComposableSlot:
		formatted := formatter.variantInterface(string(def.AsComposableSlot().Variant))

		if !resolveBuilders {
			return formatted
		}

		cogAlias := formatter.packageMapper("cog")

		return fmt.Sprintf("%s.Builder<%s>", cogAlias, formatted)
	case ir.KindConstantRef:
		return formatter.formatConstantReferences(def.AsConstantRef())
	default:
		return string(def.Kind)
	}
}

func (formatter *typeFormatter) formatRef(refType ir.RefType, resolveBuilders bool) string {
	formatted := tools.CleanupNames(refType.ReferredType)

	referredPkg := formatter.packageMapper(refType.ReferredPkg)
	if referredPkg != "" {
		formatted = referredPkg + "." + formatted
	}

	if resolveBuilders && formatter.context.ResolveToBuilder(refType.AsType()) {
		cogAlias := formatter.packageMapper("cog")

		return fmt.Sprintf("%s.Builder<%s>", cogAlias, formatted)
	}

	// if the field's type is a reference to a constant,
	// we need to use the constant's value instead.
	// ie: `SomeField: "foo"` instead of `SomeField: MyStringConstant`
	referredType, found := formatter.context.GetObjectByRef(refType)
	if found && referredType.Type.IsConcreteScalar() {
		return formatter.doFormatType(referredType.Type, resolveBuilders)
	}

	return formatted
}

func (formatter *typeFormatter) formatStructFields(structType ir.Type) string {
	var buffer strings.Builder

	buffer.WriteString("{\n")

	for _, fieldDef := range structType.AsStruct().Fields {
		fieldDefGen := formatter.formatField(fieldDef)

		buffer.WriteString(
			strings.TrimSuffix(
				prefixLinesWith(fieldDefGen, "\t"),
				"\t",
			),
		)
	}

	if structType.ImplementsVariant() {
		variant := tools.UpperCamelCase(structType.ImplementedVariant())
		buffer.WriteString(fmt.Sprintf("\t_implements%sVariant(): void;\n", variant))
	}

	buffer.WriteString("}")

	return buffer.String()
}

func (formatter *typeFormatter) formatField(def ir.StructField) string {
	var buffer strings.Builder

	for _, commentLine := range def.Comments {
		buffer.WriteString(fmt.Sprintf("// %s\n", commentLine))
	}

	required := ""
	if !def.Required {
		required = "?"
	}

	formattedType := formatter.doFormatType(def.Type, false)

	buffer.WriteString(fmt.Sprintf(
		"%s%s: %s;\n",
		def.Name,
		required,
		formattedType,
	))

	return buffer.String()
}

func (formatter *typeFormatter) formatScalarKind(kind ir.ScalarKind) string {
	switch kind {
	case ir.KindNull:
		return "null"
	case ir.KindAny:
		return "any"

	case ir.KindBytes, ir.KindString:
		return "string"

	case ir.KindFloat32, ir.KindFloat64:
		return "number"
	case ir.KindUint8, ir.KindUint16, ir.KindUint32, ir.KindUint64:
		return "number"
	case ir.KindInt8, ir.KindInt16, ir.KindInt32, ir.KindInt64:
		return "number"

	case ir.KindBool:
		return "boolean"
	default:
		return string(kind)
	}
}

func (formatter *typeFormatter) formatArray(def ir.ArrayType, resolveBuilders bool) string {
	subTypeString := formatter.doFormatType(def.ValueType, resolveBuilders)

	if def.ValueType.IsDisjunction() {
		return fmt.Sprintf("(%s)[]", subTypeString)
	}

	return fmt.Sprintf("%s[]", subTypeString)
}

func (formatter *typeFormatter) formatDisjunction(def ir.DisjunctionType, resolveBuilders bool) string {
	subTypes := tools.UniqueFormatted(def.Branches, func(subType ir.Type) string {
		return formatter.doFormatType(subType, resolveBuilders)
	})

	return strings.Join(subTypes, " | ")
}

func (formatter *typeFormatter) formatMap(def ir.MapType, resolveBuilders bool) string {
	keyTypeString := formatter.doFormatType(def.IndexType, resolveBuilders)
	valueTypeString := formatter.doFormatType(def.ValueType, resolveBuilders)

	return fmt.Sprintf("Record<%s, %s>", keyTypeString, valueTypeString)
}

func (formatter *typeFormatter) formatAnonymousEnum(def ir.EnumType) string {
	values := make([]string, 0, len(def.Values))
	for _, value := range def.Values {
		values = append(values, fmt.Sprintf("%#v", value.Value))
	}

	enumeration := strings.Join(values, " | ")

	return enumeration
}

func (formatter *typeFormatter) formatIntersection(def ir.IntersectionType) string {
	var buffer strings.Builder

	refs := make([]ir.Type, 0)
	rest := make([]ir.Type, 0)
	for _, b := range def.Branches {
		if b.Ref != nil {
			refs = append(refs, b)
			continue
		}
		rest = append(rest, b)
	}

	if len(refs) > 0 {
		buffer.WriteString("extends ")
	}

	for i, ref := range refs {
		if i != 0 && i < len(refs) {
			buffer.WriteString(", ")
		}

		buffer.WriteString(formatter.doFormatType(ref, false))
	}

	buffer.WriteString(" {\n")

	for _, r := range rest {
		if r.Struct != nil {
			for _, fieldDef := range r.AsStruct().Fields {
				buffer.WriteString("\t" + formatter.formatField(fieldDef))
			}
			continue
		}
		buffer.WriteString("\t" + formatter.doFormatType(r, false))
	}

	buffer.WriteString("}")

	return buffer.String()
}

func (formatter *typeFormatter) formatConstantReferences(def ir.ConstantReferenceType) string {
	referredType, found := formatter.context.GetObject(def.ReferredPkg, def.ReferredType)
	if !found {
		return "unknown"
	}

	if referredType.Type.IsEnum() {
		return formatter.enums.formatValue(referredType, def.ReferenceValue)
	}
	if referredType.Type.IsScalar() {
		return formatValue(def.ReferenceValue)
	}

	return "unknown"
}

type enumAsTypeFormatter struct {
	packageMapper func(pkg string) string
}

func (formatter *enumAsTypeFormatter) formatDeclaration(def ir.Object) string {
	var buffer strings.Builder
	objectName := formatObjectName(def.Name)

	buffer.WriteString(fmt.Sprintf("enum %s {\n", objectName))
	for _, val := range def.Type.AsEnum().Values {
		buffer.WriteString(fmt.Sprintf("\t%s = %s,\n", formatEnumMemberName(val.Name), formatValue(val.Value)))
	}
	buffer.WriteString("}")

	return buffer.String()
}

func (formatter *enumAsTypeFormatter) formatValue(enumObj ir.Object, val any) string {
	referredPkg := formatter.packageMapper(enumObj.SelfRef.ReferredPkg)
	pkgPrefix := ""
	if referredPkg != "" {
		pkgPrefix = referredPkg + "."
	}

	member, _ := enumObj.Type.AsEnum().MemberForValue(val)

	return fmt.Sprintf("%s%s.%s", pkgPrefix, enumObj.Name, formatEnumMemberName(member.Name))
}

type enumAsDisjunctionFormatter struct {
}

func (formatter *enumAsDisjunctionFormatter) formatDeclaration(def ir.Object) string {
	values := tools.Map(def.Type.Enum.Values, func(value ir.EnumValue) string {
		return formatValue(value.Value)
	})

	return fmt.Sprintf("type %s = %s;", formatObjectName(def.Name), strings.Join(values, " | "))
}

func (formatter *enumAsDisjunctionFormatter) formatValue(enumObj ir.Object, val any) string {
	if val == nil {
		return formatValue(enumObj.Type.Enum.Values[0].Value)
	}

	return formatValue(val)
}
