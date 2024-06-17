package php

import (
	"fmt"
	"strings"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/tools"
)

type typehints struct {
	config Config
}

func (generator *typehints) requiresHint(def ast.Type) bool {
	if def.IsAny() {
		return true
	}

	return !def.IsAnyOf(ast.KindScalar, ast.KindStruct, ast.KindRef, ast.KindEnum)
}

func (generator *typehints) paramAnnotationForType(paramName string, def ast.Type) string {
	hintText := generator.forType(def)
	if hintText == "" {
		return ""
	}

	return fmt.Sprintf("@param %s $%s", hintText, paramName)
}

func (generator *typehints) varAnnotationForType(def ast.Type) string {
	hintText := generator.forType(def)
	if hintText == "" {
		return ""
	}

	return "@var " + hintText
}

func (generator *typehints) forType(def ast.Type) string {
	hint := ""

	switch {
	case def.IsArray():
		hint = generator.arrayHint(def)
	case def.IsMap():
		hint = generator.mapHint(def)
	case def.IsScalar():
		hint = generator.scalarHint(def)
	case def.IsRef():
		hint = generator.refHint(def)
	case def.IsComposableSlot():
		hint = generator.composableSlotHint(def)
	case def.IsDisjunction():
		hint = generator.disjunctionHint(def)
	}

	if hint == "" {
		return ""
	}

	if def.Nullable {
		hint += "|null"
	}

	return hint
}

func (generator *typehints) arrayHint(def ast.Type) string {
	valueType := generator.forType(def.Array.ValueType)

	return fmt.Sprintf("array<%s>", valueType)
}

func (generator *typehints) mapHint(def ast.Type) string {
	indexType := generator.forType(def.Map.IndexType)
	valueType := generator.forType(def.Map.ValueType)

	return fmt.Sprintf("array<%s, %s>", indexType, valueType)
}

func (generator *typehints) scalarHint(def ast.Type) string {
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
		return "mixed"

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

func (generator *typehints) refHint(def ast.Type) string {
	referredPkg := formatPackageName(def.AsRef().ReferredPkg)
	typeName := formatObjectName(def.AsRef().ReferredType)

	return generator.config.fullNamespaceRef("Types\\" + referredPkg + "\\" + typeName)
}

func (generator *typehints) composableSlotHint(def ast.Type) string {
	return generator.config.fullNamespaceRef("Runtime\\Variants\\" + formatObjectName(string(def.ComposableSlot.Variant)))
}

func (generator *typehints) disjunctionHint(def ast.Type) string {
	branches := tools.Map(def.Disjunction.Branches, func(branch ast.Type) string {
		return generator.forType(branch)
	})

	return strings.Join(branches, "|")
}
