package python

import (
	"fmt"
	"strings"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/tools"
)

type pkgImporter func(alias string, pkg string) string
type moduleImporter func(alias string, pkg string, module string) string

type typeFormatter struct {
	importPkg    pkgImporter
	importModule moduleImporter

	forBuilder bool
	context    common.Context
}

func defaultTypeFormatter(importPkg pkgImporter, importModule moduleImporter) *typeFormatter {
	return &typeFormatter{
		importPkg:    importPkg,
		importModule: importModule,
	}
}

func builderTypeFormatter(context common.Context, importPkg pkgImporter, importModule moduleImporter) *typeFormatter {
	return &typeFormatter{
		importPkg:    importPkg,
		importModule: importModule,
		forBuilder:   true,
		context:      context,
	}
}

func (formatter *typeFormatter) formatObject(def ast.Object) (string, error) {
	var buffer strings.Builder

	defName := tools.UpperCamelCase(def.Name)

	if !def.Type.IsAnyOf(ast.KindStruct, ast.KindEnum) {
		buffer.WriteString(formatter.formatComments(def.Comments))
	}

	switch def.Type.Kind {
	case ast.KindEnum:
		buffer.WriteString(formatter.formatEnum(def))
	case ast.KindStruct:
		return formatter.formatStruct(def), nil
	case ast.KindDisjunction:
		typingPkg := formatter.importPkg("typing", "typing")
		buffer.WriteString(fmt.Sprintf("%s = %s.Union[%s]", defName, typingPkg, formatter.formatType(def.Type)))
	case ast.KindMap, ast.KindRef, ast.KindArray, ast.KindIntersection, ast.KindScalar:
		buffer.WriteString(fmt.Sprintf("%s = %s", defName, formatter.formatType(def.Type)))
	default:
		return "", fmt.Errorf("unhandled object of kind: %s", def.Type.Kind)
	}

	return buffer.String(), nil
}

func (formatter *typeFormatter) formatType(def ast.Type) string {
	result := "unknown"

	if def.Kind == ast.KindComposableSlot {
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

	if def.Kind == ast.KindIntersection {
		panic("formatting intersection type is not implemented for python")
	}

	if def.Kind == ast.KindDisjunction {
		result = formatter.formatDisjunction(def.AsDisjunction())
	}

	if formatter.forBuilder && (def.Kind == ast.KindComposableSlot || (def.Kind == ast.KindRef && formatter.context.ResolveToBuilder(def))) {
		cogBuilder := formatter.importModule("cogbuilder", "..cog", "builder")
		result = fmt.Sprintf("%s.Builder[%s]", cogBuilder, result)
	} else if def.Nullable {
		typingPkg := formatter.importPkg("typing", "typing")
		result = fmt.Sprintf("%s.Optional[%s]", typingPkg, result)
	}

	return result
}

func (formatter *typeFormatter) formatEnum(def ast.Object) string {
	var buffer strings.Builder

	enumPkg := formatter.importPkg("enum", "enum")

	enumName := tools.UpperCamelCase(def.Name)
	enumType := def.Type.AsEnum()

	enumKind := enumPkg + ".IntEnum"
	if enumType.Values[0].Type.AsScalar().ScalarKind == ast.KindString {
		enumKind = enumPkg + ".StrEnum"
	}
	buffer.WriteString(fmt.Sprintf("class %s(%s):\n", enumName, enumKind))
	buffer.WriteString(formatter.formatClassComments(def.Comments))

	for i, val := range enumType.Values {
		memberName := tools.UpperSnakeCase(val.Name)
		buffer.WriteString(fmt.Sprintf("    %s = %#v", memberName, val.Value))

		if i != len(enumType.Values)-1 {
			buffer.WriteString("\n")
		}
	}

	return buffer.String()
}

func (formatter *typeFormatter) formatAnonymousEnum(typeDef ast.Type) string {
	typingPkg := formatter.importPkg("typing", "typing")
	literalValues := tools.Map(typeDef.AsEnum().Values, func(val ast.EnumValue) string {
		return formatValue(val.Value)
	})

	return fmt.Sprintf("%s.Literal[%s]", typingPkg, strings.Join(literalValues, ", "))
}

func (formatter *typeFormatter) formatStruct(def ast.Object) string {
	var buffer strings.Builder

	classBases := ""
	if def.Type.Kind == ast.KindStruct && def.Type.ImplementsVariant() {
		cogVariants := formatter.importModule("cogvariants", "..cog", "variants")
		variant := tools.UpperCamelCase(def.Type.ImplementedVariant())

		classBases = fmt.Sprintf("(%s.%s)", cogVariants, variant)
	}

	buffer.WriteString(fmt.Sprintf("class %s%s:\n", tools.UpperCamelCase(def.Name), classBases))
	buffer.WriteString(formatter.formatClassComments(def.Comments))

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

func (formatter *typeFormatter) formatStructField(def ast.StructField) string {
	var buffer strings.Builder

	for _, commentLine := range def.Comments {
		buffer.WriteString(fmt.Sprintf("    # %s\n", commentLine))
	}

	field := formatter.formatType(def.Type)

	buffer.WriteString(fmt.Sprintf("    %s: %s", formatIdentifier(def.Name), field))

	return buffer.String()
}

func (formatter *typeFormatter) formatArray(def ast.ArrayType) string {
	return fmt.Sprintf("list[%s]", formatter.formatType(def.ValueType))
}

func (formatter *typeFormatter) formatMap(def ast.MapType) string {
	keyTypeString := formatter.formatType(def.IndexType)
	valueTypeString := formatter.formatType(def.ValueType)

	return fmt.Sprintf("dict[%s, %s]", keyTypeString, valueTypeString)
}

func (formatter *typeFormatter) formatRef(def ast.RefType) string {
	return formatter.formatFullyQualifiedRef(def, !formatter.forBuilder)
}

func (formatter *typeFormatter) formatFullyQualifiedRef(def ast.RefType, escapeForwardRef bool) string {
	formatted := tools.UpperCamelCase(def.ReferredType)

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

func (formatter *typeFormatter) formatDisjunction(def ast.DisjunctionType) string {
	branches := tools.Map(def.Branches, func(branch ast.Type) string {
		return formatter.formatType(branch)
	})

	typingPkg := formatter.importPkg("typing", "typing")
	return fmt.Sprintf("%s.Union[%s]", typingPkg, strings.Join(branches, ", "))
}

func (formatter *typeFormatter) formatEnumValue(enumObj ast.Object, val any) string {
	referredPkg := enumObj.SelfRef.ReferredPkg
	referredPkg = formatter.importModule(referredPkg, "..models", referredPkg)

	enumName := tools.UpperSnakeCase(enumObj.Type.AsEnum().Values[0].Name)
	for _, enumValue := range enumObj.Type.AsEnum().Values {
		if enumValue.Value == val {
			enumName = tools.UpperSnakeCase(enumValue.Name)
			break
		}
	}

	if referredPkg == "" {
		return fmt.Sprintf("%s.%s", enumObj.Name, enumName)
	}

	return fmt.Sprintf("%s.%s.%s", referredPkg, enumObj.Name, enumName)
}

func (formatter *typeFormatter) formatScalarKind(kind ast.ScalarKind) string {
	switch kind {
	case ast.KindNull:
		return "None"
	case ast.KindAny:
		return "object"

	case ast.KindBytes:
		return "bytes"
	case ast.KindString:
		return "str"

	case ast.KindFloat32, ast.KindFloat64:
		return "float"
	case ast.KindUint8, ast.KindUint16, ast.KindUint32, ast.KindUint64:
		return "int"
	case ast.KindInt8, ast.KindInt16, ast.KindInt32, ast.KindInt64:
		return "int"

	case ast.KindBool:
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
		buffer.WriteString(fmt.Sprintf("    %s\n", commentLine))
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
