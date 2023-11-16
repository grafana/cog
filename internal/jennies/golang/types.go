package golang

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

func (formatter *typeFormatter) formatType(def ast.Type) string {
	return formatter.doFormatType(def, formatter.forBuilder)
}

func (formatter *typeFormatter) doFormatType(def ast.Type, resolveBuilders bool) string {
	if def.IsAny() {
		return "any"
	}

	if def.Kind == ast.KindComposableSlot {
		formatted := formatter.variantInterface(string(def.AsComposableSlot().Variant))

		if !resolveBuilders {
			return formatted
		}

		cogAlias := formatter.packageMapper("cog")

		return fmt.Sprintf("%s.Builder[%s]", cogAlias, formatted)
	}

	if def.Kind == ast.KindArray {
		return formatter.formatArray(def.AsArray(), resolveBuilders)
	}

	if def.Kind == ast.KindMap {
		return formatter.formatMap(def.AsMap())
	}

	if def.Kind == ast.KindScalar {
		typeName := def.AsScalar().ScalarKind
		if def.Nullable {
			typeName = "*" + typeName
		}

		return string(typeName)
	}

	if def.Kind == ast.KindRef {
		return formatter.formatRef(def, resolveBuilders)
	}

	// anonymous struct or struct body
	if def.Kind == ast.KindStruct {
		output := formatter.formatStructBody(def.AsStruct())
		if def.Nullable {
			output = "*" + output
		}

		return output
	}

	if def.Kind == ast.KindIntersection {
		return formatter.formatIntersection(def.AsIntersection())
	}

	// FIXME: we should never be here
	return "unknown"
}

func (formatter *typeFormatter) variantInterface(variant string) string {
	referredPkg := formatter.packageMapper("cog/variants")

	return fmt.Sprintf("%s.%s", referredPkg, tools.UpperCamelCase(variant))
}

func (formatter *typeFormatter) formatStructBody(def ast.StructType) string {
	var buffer strings.Builder

	buffer.WriteString("struct {\n")

	for _, fieldDef := range def.Fields {
		buffer.WriteString("\t" + formatter.formatField(fieldDef))
	}

	buffer.WriteString("}")

	return buffer.String()
}

func (formatter *typeFormatter) formatField(def ast.StructField) string {
	var buffer strings.Builder

	for _, commentLine := range def.Comments {
		buffer.WriteString(fmt.Sprintf("// %s\n", commentLine))
	}

	jsonOmitEmpty := ""
	if !def.Required {
		jsonOmitEmpty = ",omitempty"
	}

	buffer.WriteString(fmt.Sprintf(
		"%s %s `json:\"%s%s\"`\n",
		tools.UpperCamelCase(def.Name),
		formatter.doFormatType(def.Type, false),
		def.Name,
		jsonOmitEmpty,
	))

	return buffer.String()
}

func (formatter *typeFormatter) formatArray(def ast.ArrayType, resolveBuilders bool) string {
	subTypeString := formatter.doFormatType(def.ValueType, resolveBuilders)

	return fmt.Sprintf("[]%s", subTypeString)
}

func (formatter *typeFormatter) formatMap(def ast.MapType) string {
	keyTypeString := formatter.doFormatType(def.IndexType, false)
	valueTypeString := formatter.doFormatType(def.ValueType, false)

	return fmt.Sprintf("map[%s]%s", keyTypeString, valueTypeString)
}

func formatScalar(val any) string {
	if list, ok := val.([]any); ok {
		items := make([]string, 0, len(list))

		for _, item := range list {
			items = append(items, formatScalar(item))
		}

		// FIXME: this is wrong, we can't just assume a list of strings.
		return fmt.Sprintf("[]string{%s}", strings.Join(items, ", "))
	}

	return fmt.Sprintf("%#v", val)
}

func (formatter *typeFormatter) formatRef(def ast.Type, resolveBuilders bool) string {
	referredPkg := formatter.packageMapper(def.AsRef().ReferredPkg)
	typeName := tools.UpperCamelCase(def.AsRef().ReferredType)

	if referredPkg != "" {
		typeName = referredPkg + "." + typeName
	}

	if resolveBuilders && formatter.context.ResolveToBuilder(def) {
		cogAlias := formatter.packageMapper("cog")

		return fmt.Sprintf("%s.Builder[%s]", cogAlias, typeName)
	}

	if def.Nullable {
		typeName = "*" + typeName
	}

	return typeName
}

func (formatter *typeFormatter) formatIntersection(def ast.IntersectionType) string {
	var buffer strings.Builder

	buffer.WriteString("struct {\n")

	refs := make([]ast.Type, 0)
	rest := make([]ast.Type, 0)
	for _, b := range def.Branches {
		if b.Ref != nil {
			refs = append(refs, b)
			continue
		}

		rest = append(rest, b)
	}

	for _, ref := range refs {
		buffer.WriteString("\t" + formatter.formatRef(ref, false) + "\n")
	}

	if len(refs) > 0 {
		buffer.WriteString("\n")
	}

	for _, r := range rest {
		if r.Struct != nil {
			for _, fieldDef := range r.AsStruct().Fields {
				buffer.WriteString("\t" + formatter.formatField(fieldDef))
			}
			continue
		}
		buffer.WriteString("\t" + formatter.doFormatType(r, false) + "\n")
	}

	buffer.WriteString("}")

	return buffer.String()
}
