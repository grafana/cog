package python

import (
	"fmt"
	"strings"

	"github.com/grafana/cog/internal/tools"
	"github.com/grafana/cog/pkg/ir"
	"github.com/grafana/cog/pkg/languages"
)

type pkgImporter func(alias string, pkg string) string
type moduleImporter func(alias string, pkg string, module string) string

type typeFormatter struct {
	importPkg    pkgImporter
	importModule moduleImporter

	forBuilder bool
	context    languages.Context
}

func defaultTypeFormatter(context languages.Context, importPkg pkgImporter, importModule moduleImporter) *typeFormatter {
	return &typeFormatter{
		context:      context,
		importPkg:    importPkg,
		importModule: importModule,
	}
}

func builderTypeFormatter(context languages.Context, importPkg pkgImporter, importModule moduleImporter) *typeFormatter {
	return &typeFormatter{
		importPkg:    importPkg,
		importModule: importModule,
		forBuilder:   true,
		context:      context,
	}
}

func (formatter *typeFormatter) formatObject(object ir.Object) (string, error) {
	var buffer strings.Builder

	defName := formatObjectName(object.Name)

	if !object.Type.IsAnyOf(ir.KindStruct, ir.KindEnum) {
		buffer.WriteString(formatter.formatComments(object.Comments))
	}

	if object.Type.IsConcreteScalar() {
		fmt.Fprintf(&buffer, "%s: %s = %s", defName, formatter.formatType(object.Type), formatValue(object.Type.AsScalar().Value))

		return buffer.String(), nil
	}

	switch object.Type.Kind {
	case ir.KindEnum:
		buffer.WriteString(formatter.formatEnum(object))
	case ir.KindStruct:
		return formatter.formatStruct(object), nil
	default:
		typingPkg := formatter.importPkg("typing", "typing")
		fmt.Fprintf(&buffer, "%s: %s.TypeAlias = %s", defName, typingPkg, formatter.formatType(object.Type))
	}

	return buffer.String(), nil
}

func (formatter *typeFormatter) formatType(def ir.Type) string {
	result := "unknown"

	if def.IsComposableSlot() {
		formatted := tools.UpperCamelCase(string(def.AsComposableSlot().Variant))
		cogVariants := formatter.importModule("cogvariants", "..cog", "variants")

		result = fmt.Sprintf("%s.%s", cogVariants, formatted)
	}

	if def.IsArray() {
		result = formatter.formatArray(def.AsArray())
	}

	if def.IsMap() {
		result = formatter.formatMap(def.AsMap())
	}

	if def.IsScalar() {
		// This scalar actually refers to a constant
		if def.AsScalar().IsConcrete() {
			typingPkg := formatter.importPkg("typing", "typing")
			result = fmt.Sprintf("%s.Literal[%s]", typingPkg, formatValue(def.AsScalar().Value))
		} else {
			result = formatter.formatScalarKind(def.AsScalar().ScalarKind)
		}
	}

	if def.IsRef() {
		result = formatter.formatRef(def.AsRef())
	}

	// anonymous enum
	if def.IsEnum() {
		result = formatter.formatAnonymousEnum(def)
	}

	if def.IsIntersection() {
		panic("formatting intersection type is not implemented for python")
	}

	if def.IsDisjunction() {
		result = formatter.formatDisjunction(def.AsDisjunction())
	}

	if formatter.forBuilder && (def.IsComposableSlot() || (def.IsRef() && formatter.context.ResolveToBuilder(def))) {
		cogBuilder := formatter.importModule("cogbuilder", "..cog", "builder")
		result = fmt.Sprintf("%s.Builder[%s]", cogBuilder, result)
	} else if def.Nullable {
		typingPkg := formatter.importPkg("typing", "typing")
		result = fmt.Sprintf("%s.Optional[%s]", typingPkg, result)
	}

	if def.IsConstantRef() {
		result = formatter.formatConstantReference(def.AsConstantRef(), false)
	}

	return result
}

func (formatter *typeFormatter) formatEnum(def ir.Object) string {
	var buffer strings.Builder

	enumPkg := formatter.importPkg("enum", "enum")

	enumName := formatObjectName(def.Name)
	enumType := def.Type.AsEnum()

	enumKind := enumPkg + ".IntEnum"
	if enumType.Values[0].Type.AsScalar().ScalarKind == ir.KindString {
		enumKind = enumPkg + ".StrEnum"
	}
	fmt.Fprintf(&buffer, "class %s(%s):\n", enumName, enumKind)
	buffer.WriteString(formatter.formatClassComments(def.Comments))

	for i, val := range enumType.Values {
		memberName := tools.UpperSnakeCase(val.Name)
		fmt.Fprintf(&buffer, "    %s = %#v", memberName, val.Value)

		if i != len(enumType.Values)-1 {
			buffer.WriteString("\n")
		}
	}

	return buffer.String()
}

func (formatter *typeFormatter) formatAnonymousEnum(typeDef ir.Type) string {
	typingPkg := formatter.importPkg("typing", "typing")
	literalValues := tools.Map(typeDef.AsEnum().Values, func(val ir.EnumValue) string {
		return formatValue(val.Value)
	})

	return fmt.Sprintf("%s.Literal[%s]", typingPkg, strings.Join(literalValues, ", "))
}

func (formatter *typeFormatter) formatStruct(def ir.Object) string {
	var buffer strings.Builder

	classBases := ""
	if def.Type.IsStruct() && def.Type.ImplementsVariant() {
		cogVariants := formatter.importModule("cogvariants", "..cog", "variants")
		variant := tools.UpperCamelCase(def.Type.ImplementedVariant())

		classBases = fmt.Sprintf("(%s.%s)", cogVariants, variant)
	}

	fmt.Fprintf(&buffer, "class %s%s:\n", formatObjectName(def.Name), classBases)
	buffer.WriteString(formatter.formatClassComments(def.Comments))

	if def.DeprecationMessage != "" {
		formatter.importPkg("warnings", "warnings")
		fmt.Fprintf(&buffer, "    warnings.warn(%s, DeprecationWarning)\n", formatValue(def.DeprecationMessage))
	}

	fields := def.Type.AsStruct().Fields

	// shouldn't happen, but we never know.
	if len(fields) == 0 {
		buffer.WriteString("    pass")
	}

	for i, fieldDef := range def.Type.AsStruct().Fields {
		buffer.WriteString(formatter.formatStructField(fieldDef))

		if i != len(fields)-1 {
			buffer.WriteString("\n")
		}
	}

	return buffer.String()
}

func (formatter *typeFormatter) formatStructField(def ir.StructField) string {
	var buffer strings.Builder

	for _, commentLine := range def.Comments {
		fmt.Fprintf(&buffer, "    # %s\n", commentLine)
	}

	fmt.Fprintf(&buffer, "    %s: %s", formatIdentifier(def.Name), formatter.formatType(def.Type))

	return buffer.String()
}

func (formatter *typeFormatter) formatArray(def ir.ArrayType) string {
	return fmt.Sprintf("list[%s]", formatter.formatType(def.ValueType))
}

func (formatter *typeFormatter) formatMap(def ir.MapType) string {
	keyTypeString := formatter.formatType(def.IndexType)
	valueTypeString := formatter.formatType(def.ValueType)

	return fmt.Sprintf("dict[%s, %s]", keyTypeString, valueTypeString)
}

func (formatter *typeFormatter) formatRef(def ir.RefType) string {
	return formatter.formatFullyQualifiedRef(def, !formatter.forBuilder)
}

func (formatter *typeFormatter) formatFullyQualifiedRef(def ir.RefType, escapeForwardRef bool) string {
	referredObject, found := formatter.context.GetObject(def.ReferredPkg, def.ReferredType)
	if found && referredObject.Type.IsConcreteScalar() {
		return formatter.formatType(referredObject.Type)
	}

	formatted := formatObjectName(def.ReferredType)

	referredPkg := def.ReferredPkg
	referredPkg = formatter.importModule(referredPkg, "..models", referredPkg)
	if referredPkg != "" {
		formatted = referredPkg + "." + formatted
	}

	if !escapeForwardRef || referredPkg != "" {
		return formatted
	}

	// The quotes are important to allow for forward-references.
	return fmt.Sprintf("'%s'", formatted)
}

func (formatter *typeFormatter) formatDisjunction(def ir.DisjunctionType) string {
	typingPkg := formatter.importPkg("typing", "typing")
	branches := tools.UniqueFormatted(def.Branches, formatter.formatType)

	return fmt.Sprintf("%s.Union[%s]", typingPkg, strings.Join(branches, ", "))
}

func (formatter *typeFormatter) formatEnumValue(enumObj ir.Object, val any) string {
	referredPkg := enumObj.SelfRef.ReferredPkg
	referredPkg = formatter.importModule(referredPkg, "..models", referredPkg)

	member, _ := enumObj.Type.AsEnum().MemberForValue(val)
	memberName := tools.UpperSnakeCase(member.Name)

	if referredPkg == "" {
		return fmt.Sprintf("%s.%s", enumObj.Name, memberName)
	}

	return fmt.Sprintf("%s.%s.%s", referredPkg, enumObj.Name, memberName)
}

func (formatter *typeFormatter) formatScalarKind(kind ir.ScalarKind) string {
	switch kind {
	case ir.KindNull:
		return "None"
	case ir.KindAny:
		return "object"

	case ir.KindBytes:
		return "bytes"
	case ir.KindString:
		return "str"

	case ir.KindFloat32, ir.KindFloat64:
		return "float"
	case ir.KindUint8, ir.KindUint16, ir.KindUint32, ir.KindUint64:
		return "int"
	case ir.KindInt8, ir.KindInt16, ir.KindInt32, ir.KindInt64:
		return "int"

	case ir.KindBool:
		return "bool"
	default:
		return string(kind)
	}
}

func (formatter *typeFormatter) formatClassComments(comments []string) string {
	if len(comments) == 0 {
		return ""
	}

	var buffer strings.Builder

	buffer.WriteString(`    """` + "\n")
	for _, commentLine := range comments {
		fmt.Fprintf(&buffer, "    %s\n", commentLine)
	}
	buffer.WriteString(`    """` + "\n\n")

	return buffer.String()
}

func (formatter *typeFormatter) formatComments(comments []string) string {
	if len(comments) == 0 {
		return ""
	}

	var buffer strings.Builder

	for _, commentLine := range comments {
		buffer.WriteString(strings.TrimRight(fmt.Sprintf("# %s\n", commentLine), " "))
	}

	return buffer.String()
}

func (formatter *typeFormatter) formatConstantReference(def ir.ConstantReferenceType, shouldSetValue bool) string {
	referredObject, found := formatter.context.GetObject(def.ReferredPkg, def.ReferredType)
	if !found {
		return "unknown"
	}

	if referredObject.Type.IsScalar() {
		if shouldSetValue {
			return def.ReferredType
		}

		return formatter.formatScalarKind(referredObject.Type.AsScalar().ScalarKind)
	}

	if !referredObject.Type.IsEnum() {
		return "unknown"
	}

	if shouldSetValue {
		return formatter.formatEnumValue(referredObject, def.ReferenceValue)
	}

	t := referredObject.Type.AsEnum().Values[0].Type
	if t.AsScalar().ScalarKind == ir.KindString {
		return "str"
	}

	return "int"
}
