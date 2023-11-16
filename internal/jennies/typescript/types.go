package typescript

import (
	"fmt"
	"strings"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/context"
	"github.com/grafana/cog/internal/tools"
)

type typeFormatter struct {
	packageMapper func(pkg string) string

	forBuilder bool
	context    context.Builders
}

func defaultTypeFormatter(packageMapper func(pkg string) string) *typeFormatter {
	return &typeFormatter{
		packageMapper: packageMapper,
		context:       context.Builders{},
	}
}

func builderTypeFormatter(context context.Builders, packageMapper func(pkg string) string) *typeFormatter {
	return &typeFormatter{
		packageMapper: packageMapper,
		forBuilder:    true,
		context:       context,
	}
}

func (formatter *typeFormatter) variantInterface(variant string) string {
	referredPkg := formatter.packageMapper("cog")

	return fmt.Sprintf("%s.%s", referredPkg, tools.UpperCamelCase(variant))
}

func (formatter *typeFormatter) formatType(def ast.Type) string {
	return formatter.doFormatType(def, formatter.forBuilder)
}

func (formatter *typeFormatter) doFormatType(def ast.Type, resolveBuilders bool) string {
	switch def.Kind {
	case ast.KindDisjunction:
		return formatter.formatDisjunction(def.AsDisjunction(), resolveBuilders)
	case ast.KindRef:
		formatted := def.AsRef().ReferredType

		referredPkg := formatter.packageMapper(def.AsRef().ReferredPkg)
		if referredPkg != "" {
			formatted = referredPkg + "." + formatted
		}

		if resolveBuilders && formatter.context.ResolveToBuilder(def) {
			cogAlias := formatter.packageMapper("cog")

			return fmt.Sprintf("%s.OptionsBuilder<%s>", cogAlias, formatted)
		}

		return formatted
	case ast.KindArray:
		return formatter.formatArray(def.AsArray(), resolveBuilders)
	case ast.KindStruct:
		return formatter.formatStructFields(def)
	case ast.KindMap:
		return formatter.formatMap(def.AsMap())
	case ast.KindEnum:
		return formatter.formatAnonymousEnum(def.AsEnum())
	case ast.KindScalar:
		// This scalar actually refers to a constant
		if def.AsScalar().Value != nil {
			return formatScalar(def.AsScalar().Value)
		}

		return formatter.formatScalarKind(def.AsScalar().ScalarKind)
	case ast.KindIntersection:
		return formatter.formatIntersection(def.AsIntersection())
	case ast.KindComposableSlot:
		formatted := formatter.variantInterface(string(def.AsComposableSlot().Variant))

		if !resolveBuilders {
			return formatted
		}

		cogAlias := formatter.packageMapper("cog")

		return fmt.Sprintf("%s.OptionsBuilder<%s>", cogAlias, formatted)
	default:
		return string(def.Kind)
	}
}

func (formatter *typeFormatter) formatStructFields(structType ast.Type) string {
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

func (formatter *typeFormatter) formatField(def ast.StructField) string {
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
func (formatter *typeFormatter) formatScalarKind(kind ast.ScalarKind) string {
	switch kind {
	case ast.KindNull:
		return "null"
	case ast.KindAny:
		return "any"

	case ast.KindBytes, ast.KindString:
		return "string"

	case ast.KindFloat32, ast.KindFloat64:
		return "number"
	case ast.KindUint8, ast.KindUint16, ast.KindUint32, ast.KindUint64:
		return "number"
	case ast.KindInt8, ast.KindInt16, ast.KindInt32, ast.KindInt64:
		return "number"

	case ast.KindBool:
		return "boolean"
	default:
		return string(kind)
	}
}

func (formatter *typeFormatter) formatArray(def ast.ArrayType, resolveBuilders bool) string {
	subTypeString := formatter.doFormatType(def.ValueType, resolveBuilders)

	if def.ValueType.Kind == ast.KindDisjunction {
		return fmt.Sprintf("(%s)[]", subTypeString)
	}

	return fmt.Sprintf("%s[]", subTypeString)
}

func (formatter *typeFormatter) formatDisjunction(def ast.DisjunctionType, resolveBuilders bool) string {
	subTypes := make([]string, 0, len(def.Branches))
	for _, subType := range def.Branches {
		subTypes = append(subTypes, formatter.doFormatType(subType, resolveBuilders))
	}

	return strings.Join(subTypes, " | ")
}

func (formatter *typeFormatter) formatMap(def ast.MapType) string {
	keyTypeString := formatter.doFormatType(def.IndexType, false)
	valueTypeString := formatter.doFormatType(def.ValueType, false)

	return fmt.Sprintf("Record<%s, %s>", keyTypeString, valueTypeString)
}

func (formatter *typeFormatter) formatAnonymousEnum(def ast.EnumType) string {
	values := make([]string, 0, len(def.Values))
	for _, value := range def.Values {
		values = append(values, fmt.Sprintf("%#v", value.Value))
	}

	enumeration := strings.Join(values, " | ")

	return enumeration
}

func (formatter *typeFormatter) formatIntersection(def ast.IntersectionType) string {
	var buffer strings.Builder

	refs := make([]ast.Type, 0)
	rest := make([]ast.Type, 0)
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

func formatScalar(val any) string {
	if list, ok := val.([]any); ok {
		items := make([]string, 0, len(list))

		for _, item := range list {
			items = append(items, formatScalar(item))
		}

		return fmt.Sprintf("[%s]", strings.Join(items, ", "))
	}

	return fmt.Sprintf("%#v", val)
}
