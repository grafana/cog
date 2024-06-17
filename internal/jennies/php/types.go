package php

import (
	"fmt"
	"strings"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/languages"
	"github.com/grafana/cog/internal/tools"
)

type typeFormatter struct {
	config Config

	forBuilder bool
	context    languages.Context
}

func defaultTypeFormatter(config Config, context languages.Context) *typeFormatter {
	return &typeFormatter{
		config:  config,
		context: context,
	}
}

func (formatter *typeFormatter) formatType(def ast.Type) string {
	return formatter.doFormatType(def, formatter.forBuilder)
}

func (formatter *typeFormatter) doFormatType(def ast.Type, resolveBuilders bool) string {
	actualFormatter := func() string {
		if def.IsAny() {
			return ""
		}

		if def.IsComposableSlot() {
			formatted := formatter.variantInterface(string(def.AsComposableSlot().Variant))

			if !resolveBuilders {
				return formatted
			}

			// TODO
			cogAlias := "cog"

			return fmt.Sprintf("%s.Builder[%s]", cogAlias, formatted)
		}

		if def.IsArray() || def.IsMap() {
			return "array"
		}

		if def.IsScalar() {
			return formatter.formatScalar(def)
		}

		if def.IsRef() {
			return formatter.formatRef(def, resolveBuilders)
		}

		if def.IsDisjunction() {
			return ""
		}

		// FIXME: we should never be here
		return "unknown"
	}

	passesTrail := ""
	if formatter.config.debug && len(def.PassesTrail) != 0 {
		passesTrail = fmt.Sprintf(" /* %s */", strings.Join(def.PassesTrail, ", "))
	}

	formatted := actualFormatter()
	if def.Nullable && formatted != "" {
		formatted = "?" + formatted
	}

	return formatted + passesTrail
}

func (formatter *typeFormatter) variantInterface(variant string) string {
	return formatter.config.fullNamespaceRef("Runtime\\Variants\\" + formatObjectName(variant))
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

	buffer.WriteString(formatCommentsBlock(comments))

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

	formattedType := formatter.doFormatType(fieldType, false)
	if formattedType != "" {
		formattedType = " " + formattedType
	}

	buffer.WriteString(fmt.Sprintf(
		"public%s $%s;",
		formattedType,
		formatFieldName(def.Name),
	))

	return buffer.String()
}

func (formatter *typeFormatter) formatScalar(def ast.Type) string {
	scalarKind := def.AsScalar().ScalarKind
	/*
		if def.HasHint(ast.HintStringFormatDateTime) {
			scalarKind = "time.Time" // TODO
		}
	*/

	switch scalarKind {
	case ast.KindNull:
		return "null"
	case ast.KindAny:
		return ""

	case ast.KindBytes:
		return "string"
	case ast.KindString:
		return "string"

	case ast.KindFloat32, ast.KindFloat64:
		return "float"
	case ast.KindUint8, ast.KindUint16, ast.KindUint32, ast.KindUint64:
		return "int"
	case ast.KindInt8, ast.KindInt16, ast.KindInt32, ast.KindInt64:
		return "int"

	case ast.KindBool:
		return "bool"
	default:
		return string(scalarKind)
	}
}

func (formatter *typeFormatter) formatRef(def ast.Type, resolveBuilders bool) string {
	referredPkg := formatPackageName(def.AsRef().ReferredPkg)
	typeName := formatObjectName(def.AsRef().ReferredType)

	if referredPkg != "" {
		typeName = "Types\\" + referredPkg + "\\" + typeName
	}

	if resolveBuilders {
		typeName = "Builders\\" + referredPkg + "\\" + typeName
	}

	typeName = formatter.config.fullNamespaceRef(typeName)

	if resolveBuilders && formatter.context.ResolveToBuilder(def) {
		cogAlias := "cog"

		return fmt.Sprintf("%s.Builder[%s]", cogAlias, typeName)
	}

	return typeName
}
