package golang

import (
	"fmt"
	"strings"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/languages"
	"github.com/grafana/cog/internal/orderedmap"
	"github.com/grafana/cog/internal/tools"
)

type typeFormatter struct {
	packageMapper func(pkg string) string
	config        Config

	forBuilder bool
	context    languages.Context
}

func defaultTypeFormatter(config Config, context languages.Context, packageMapper func(pkg string) string) *typeFormatter {
	return &typeFormatter{
		packageMapper: packageMapper,
		config:        config,
		context:       context,
	}
}

func builderTypeFormatter(config Config, context languages.Context, packageMapper func(pkg string) string) *typeFormatter {
	return &typeFormatter{
		packageMapper: packageMapper,
		config:        config,
		forBuilder:    true,
		context:       context,
	}
}

func (formatter *typeFormatter) formatType(def ast.Type) string {
	return formatter.doFormatType(def, formatter.forBuilder)
}

func (formatter *typeFormatter) doFormatType(def ast.Type, resolveBuilders bool) string {
	actualFormatter := func() string {
		if def.IsAny() {
			return "any"
		}

		if def.IsComposableSlot() {
			formatted := formatter.variantInterface(string(def.AsComposableSlot().Variant))

			if !resolveBuilders {
				return formatted
			}

			cogAlias := formatter.packageMapper("cog")

			return fmt.Sprintf("%s.Builder[%s]", cogAlias, formatted)
		}

		if def.IsArray() {
			return formatter.formatArray(def.AsArray(), resolveBuilders)
		}

		if def.IsMap() {
			return formatter.formatMap(def.AsMap())
		}

		if def.IsScalar() {
			typeName := def.AsScalar().ScalarKind
			if def.HasHint(ast.HintStringFormatDateTime) {
				typeName = "time.Time"
			}
			if def.Nullable {
				typeName = "*" + typeName
			}
			if def.AsScalar().ScalarKind == ast.KindBytes {
				typeName = "[]byte"
			}

			return string(typeName)
		}

		if def.IsRef() {
			return formatter.formatRef(def, resolveBuilders)
		}

		// anonymous struct or struct body
		if def.IsStruct() {
			output := formatter.formatStructBody(def.AsStruct())
			if def.Nullable {
				output = "*" + output
			}

			return output
		}

		if def.IsIntersection() {
			return formatter.formatIntersection(def.AsIntersection())
		}

		// FIXME: we should never be here
		return "unknown"
	}

	passesTrail := ""
	if formatter.config.debug && len(def.PassesTrail) != 0 {
		passesTrail = fmt.Sprintf(" /* %s */", strings.Join(def.PassesTrail, ", "))
	}

	return actualFormatter() + passesTrail
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

	comments := def.Comments
	if formatter.config.debug {
		passesTrail := tools.Map(def.PassesTrail, func(trail string) string {
			return fmt.Sprintf("Modified by compiler pass '%s'", trail)
		})
		comments = append(comments, passesTrail...)
	}

	for _, commentLine := range comments {
		buffer.WriteString(fmt.Sprintf("// %s\n", commentLine))
	}

	jsonOmitEmpty := ""
	if !def.Required {
		jsonOmitEmpty = ",omitempty"
	}

	fieldType := def.Type

	// if the field's type is a reference to a constant,
	// we need to use the constant's type instead.
	// ie: `SomeField string` instead of `SomeField MyStringConstant`
	if def.Type.IsRef() {
		referredType, found := formatter.context.LocateObject(def.Type.AsRef().ReferredPkg, def.Type.AsRef().ReferredType)
		if found && referredType.Type.IsConcreteScalar() {
			fieldType = referredType.Type
		}
	}

	buffer.WriteString(fmt.Sprintf(
		"%s %s `json:\"%s%s\"`\n",
		tools.UpperCamelCase(def.Name),
		formatter.doFormatType(fieldType, false),
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

func formatDefaultReferenceStructForBuilder(refPkg string, name string, isBuilder bool, def ast.StructType, structMap *orderedmap.Map[string, interface{}]) string {
	starter, format, sep, lastSep, ending := fmt.Sprintf("%s {\n", name), "%s: %v", ",\n", ",\n", "}"
	if isBuilder {
		starter, format, sep, lastSep, ending = fmt.Sprintf("New%sBuilder().\n", name), "%s(%v)", ".\n", ",\n", ""
	}

	if refPkg != "" {
		starter = fmt.Sprintf("%s.%s", refPkg, starter)
	}

	var buffer strings.Builder
	count := 0
	structMap.Iterate(func(key string, value interface{}) {
		field, _ := def.FieldByName(key)
		if name != "" {
			key = tools.UpperCamelCase(key)
		}

		switch x := value.(type) {
		case map[string]interface{}:
			buffer.WriteString(fmt.Sprintf(format, key, formatDefaultReferenceStructForBuilder(refPkg, name, isBuilder, field.Type.AsStruct(), orderedmap.FromMap(x))))
		case nil:
			buffer.WriteString(fmt.Sprintf(format, key, formatScalar([]any{})))
		default:
			val := formatScalar(x)
			if !isBuilder && field.Type.IsScalar() && field.Type.Nullable {
				val = fmt.Sprintf("cog.ToPtr[%s](%v)", field.Type.AsScalar().ScalarKind, value)
			}
			buffer.WriteString(fmt.Sprintf(format, key, val))
		}

		if count != structMap.Len()-1 {
			buffer.WriteString(sep)
		} else {
			buffer.WriteString(lastSep)
		}
		count++
	})

	return fmt.Sprintf("%s%s%s", starter, buffer.String(), ending)
}

func formatAnonymousDefaultStruct(def ast.StructType, structMap *orderedmap.Map[string, interface{}]) string {
	return fmt.Sprintf("struct %s {\n %s }", defineAnonymousFields(def), defineAnonymousDefaults(def, structMap))
}

func defineAnonymousFields(def ast.StructType) string {
	var structDefinition strings.Builder

	for _, f := range def.Fields {
		key := tools.UpperCamelCase(f.Name)

		switch f.Type.Kind {
		case ast.KindScalar:
			structDefinition.WriteString(fmt.Sprintf("%s %v `json:\"%s\"`\n", key, f.Type.AsScalar().ScalarKind, tools.LowerCamelCase(key)))
		case ast.KindStruct:
			structFields := defineAnonymousFields(f.Type.AsStruct())
			structDefinition.WriteString(fmt.Sprintf("%s struct %v `json:\"%s\"`\n", key, structFields, tools.LowerCamelCase(key)))
		case ast.KindArray:
			array := f.Type.AsArray()
			if array.ValueType.IsScalar() {
				structDefinition.WriteString(fmt.Sprintf("%s []%v `json:\"%s\"`\n", key, array.ValueType.AsScalar().ScalarKind, tools.LowerCamelCase(key)))
			}
		// TODO: Map rest of array cases
		default:
			// TODO: Map rest of the cases when necessary. By default it sets any
			structDefinition.WriteString(fmt.Sprintf("%s any `json:\"%s\"`\n", key, strings.ToLower(key)))
		}
	}

	return fmt.Sprintf("{\n %s }", structDefinition.String())
}

func defineAnonymousDefaults(def ast.StructType, structMap *orderedmap.Map[string, interface{}]) string {
	var buffer strings.Builder
	structMap.Iterate(func(key string, value interface{}) {
		name := tools.UpperCamelCase(key)
		switch x := value.(type) {
		case map[string]interface{}:
			// FIXME: Set a default not defined shouldn't happen..
			field, _ := def.FieldByName(key)
			def = field.Type.AsStruct()
			buffer.WriteString(fmt.Sprintf("%s: struct %v {\n %v},\n", name, defineAnonymousFields(def), defineAnonymousDefaults(def, orderedmap.FromMap(x))))
		case []interface{}, interface{}:
			buffer.WriteString(fmt.Sprintf("%s: %v,\n", name, formatScalar(x)))
		}
	})

	return buffer.String()
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

func makePathFormatter(typeFormatter *typeFormatter) func(path ast.Path) string {
	return func(fieldPath ast.Path) string {
		parts := make([]string, len(fieldPath))

		for i := range fieldPath {
			output := fieldPath[i].Identifier
			if !fieldPath[i].Root {
				output = tools.UpperCamelCase(output)
			}

			// don't generate type hints if:
			// * there isn't one defined
			// * the type isn't "any"
			// * as a trailing element in the path
			if !fieldPath[i].Type.IsAny() || fieldPath[i].TypeHint == nil || i == len(fieldPath)-1 {
				parts[i] = output
				continue
			}

			formattedTypeHint := typeFormatter.formatType(*fieldPath[i].TypeHint)
			parts[i] = output + fmt.Sprintf(".(*%s)", formattedTypeHint)
		}

		return strings.Join(parts, ".")
	}
}
