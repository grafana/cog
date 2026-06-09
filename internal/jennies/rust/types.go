package rust

import (
	"fmt"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/languages"
)

// typeFormatter renders IR types into idiomatic Rust type syntax, recording any
// imports the rendered type requires into the supplied importMap.
type typeFormatter struct {
	context languages.Context
	imports *importMap
}

func newTypeFormatter(context languages.Context, imports *importMap) *typeFormatter {
	return &typeFormatter{
		context: context,
		imports: imports,
	}
}

// formatType renders an arbitrary IR type. A nullable (optional) type is wrapped
// in Option<T>, except for arrays and maps: Vec and HashMap carry their own
// "empty" representation (an empty collection), so a nullable collection is
// rendered as the bare Vec/HashMap and its absence is modelled by emptiness.
// This mirrors the Go target, which likewise never pointer-wraps slices or maps.
func (formatter *typeFormatter) formatType(def ast.Type) string {
	inner := formatter.formatInnerType(def)

	if def.Nullable && !def.IsArray() && !def.IsMap() {
		return fmt.Sprintf("Option<%s>", inner)
	}

	return inner
}

func (formatter *typeFormatter) formatInnerType(def ast.Type) string {
	switch {
	case def.IsScalar():
		return formatter.formatScalar(def.AsScalar())
	case def.IsArray():
		return fmt.Sprintf("Vec<%s>", formatter.formatType(def.AsArray().ValueType))
	case def.IsMap():
		return formatter.formatMap(def.AsMap())
	case def.IsRef():
		return formatTypeName(def.AsRef().ReferredType)
	case def.IsConstantRef():
		return formatter.formatConstantRef(def.AsConstantRef())
	case def.IsEnum():
		// Inline enums are hoisted to named top-level enum objects by the
		// AnonymousEnumToExplicitType compiler pass before emission, so a field
		// referencing one arrives as a ref. This branch is a defensive fallback
		// that renders any residual inline enum as its underlying scalar type.
		return formatScalarKind(def.AsEnum().Values[0].Type.AsScalar().ScalarKind)
	default:
		return "serde_json::Value"
	}
}

func (formatter *typeFormatter) formatScalar(scalar ast.ScalarType) string {
	// Both concrete (constant) and non-concrete scalars map to the same storage
	// type: a concrete scalar's value is fixed by a Default impl, so the field
	// type is still just the underlying scalar kind. `any` maps to
	// serde_json::Value, referenced by its fully qualified path (no import).
	return formatScalarKind(scalar.ScalarKind)
}

func (formatter *typeFormatter) formatMap(def ast.MapType) string {
	formatter.imports.Add("std::collections::HashMap")
	return fmt.Sprintf("HashMap<%s, %s>", formatter.formatType(def.IndexType), formatter.formatType(def.ValueType))
}

func (formatter *typeFormatter) formatConstantRef(def ast.ConstantReferenceType) string {
	referredObject, found := formatter.context.LocateObject(def.ReferredPkg, def.ReferredType)
	if found && referredObject.Type.IsScalar() {
		return formatScalarKind(referredObject.Type.AsScalar().ScalarKind)
	}
	return "serde_json::Value"
}
