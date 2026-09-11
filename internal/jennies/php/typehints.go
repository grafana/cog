package php

import (
	"fmt"
	"strings"

	"github.com/grafana/cog/internal/tools"
	"github.com/grafana/cog/pkg/ir"
	"github.com/grafana/cog/pkg/languages"
)

type typehints struct {
	config          Config
	context         languages.Context
	resolveBuilders bool
}

func (generator *typehints) requiresHint(def ir.Type) bool {
	if def.IsAny() {
		return true
	}

	if generator.resolveBuilders && def.IsRef() && generator.context.ResolveToBuilder(def) {
		return true
	}

	return !def.IsAnyOf(ir.KindScalar, ir.KindStruct, ir.KindRef, ir.KindEnum, ir.KindConstantRef)
}

func (generator *typehints) paramAnnotationForType(paramName string, def ir.Type) string {
	hintText := generator.forType(def, generator.resolveBuilders)
	if hintText == "" {
		return ""
	}

	return fmt.Sprintf("@param %s $%s", hintText, formatArgName(paramName))
}

func (generator *typehints) varAnnotationForType(def ir.Type) string {
	hintText := generator.forType(def, generator.resolveBuilders)
	if hintText == "" {
		return ""
	}

	return "@var " + hintText
}

func (generator *typehints) forType(def ir.Type, resolveBuilders bool) string {
	hint := ""

	switch {
	case def.IsArray():
		hint = generator.arrayHint(def, resolveBuilders)
	case def.IsMap():
		hint = generator.mapHint(def, resolveBuilders)
	case def.IsScalar():
		hint = scalarHint(def)
	case def.IsRef():
		hint = generator.refHint(def, resolveBuilders)
	case def.IsComposableSlot():
		hint = generator.composableSlotHint(def, resolveBuilders)
	case def.IsDisjunction():
		hint = generator.disjunctionHint(def, resolveBuilders)
	case def.IsConstantRef():
		hint = generator.constantRefHint(def)
	}

	if hint == "" {
		return ""
	}

	if def.Nullable {
		hint += "|null"
	}

	return hint
}

func (generator *typehints) arrayHint(def ir.Type, resolveBuilders bool) string {
	valueType := generator.forType(def.Array.ValueType, resolveBuilders)

	return fmt.Sprintf("array<%s>", valueType)
}

func (generator *typehints) mapHint(def ir.Type, resolveBuilders bool) string {
	indexType := generator.forType(def.Map.IndexType, resolveBuilders)
	valueType := generator.forType(def.Map.ValueType, resolveBuilders)

	return fmt.Sprintf("array<%s, %s>", indexType, valueType)
}

func scalarHint(def ir.Type) string {
	scalarKind := def.AsScalar().ScalarKind
	/*
		if def.HasHint(ast.HintStringFormatDateTime) {
			scalarKind = "time.Time" // TODO
		}
	*/

	switch scalarKind {
	case ir.KindNull:
		return "null"
	case ir.KindAny:
		return "mixed"

	case ir.KindBytes:
		return "string"
	case ir.KindString:
		return "string"

	case ir.KindFloat32, ir.KindFloat64:
		return "float"
	case ir.KindUint8, ir.KindUint16, ir.KindUint32, ir.KindUint64:
		return "int"
	case ir.KindInt8, ir.KindInt16, ir.KindInt32, ir.KindInt64:
		return "int"

	case ir.KindBool:
		return "bool"
	default:
		return string(scalarKind)
	}
}

func (generator *typehints) refHint(def ir.Type, resolveBuilders bool) string {
	referredPkg := formatPackageName(def.AsRef().ReferredPkg)
	typeName := formatObjectName(def.AsRef().ReferredType)

	fqcn := generator.config.fullNamespaceRef(referredPkg + "\\" + typeName)

	if !resolveBuilders || !generator.context.ResolveToBuilder(def) {
		return fqcn
	}

	return fmt.Sprintf("%s<%s>", generator.config.fullNamespaceRef("Cog\\Builder"), fqcn)
}

func (generator *typehints) composableSlotHint(def ir.Type, resolveBuilders bool) string {
	fqcn := generator.config.fullNamespaceRef("Cog\\" + formatObjectName(string(def.ComposableSlot.Variant)))
	if !resolveBuilders {
		return fqcn
	}

	return fmt.Sprintf("%s<%s>", generator.config.fullNamespaceRef("Cog\\Builder"), fqcn)
}

func (generator *typehints) disjunctionHint(def ir.Type, resolveBuilders bool) string {
	branches := tools.Map(def.Disjunction.Branches, func(branch ir.Type) string {
		return generator.forType(branch, resolveBuilders)
	})

	return strings.Join(branches, "|")
}

func (generator *typehints) constantRefHint(def ir.Type) string {
	referredPkg := formatPackageName(def.AsConstantRef().ReferredPkg)
	typeName := formatObjectName(def.AsConstantRef().ReferredType)

	return generator.config.fullNamespaceRef(referredPkg + "\\" + typeName)
}
