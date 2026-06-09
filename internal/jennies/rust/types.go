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
	// packageName is the sanitized Rust module name of the schema currently being
	// emitted. A ref whose (sanitized) referred package matches this is rendered as
	// a bare short name; any other ref is rendered cross-module via a `use`.
	packageName string
}

func newTypeFormatter(context languages.Context, imports *importMap, packageName string) *typeFormatter {
	return &typeFormatter{
		context:     context,
		imports:     imports,
		packageName: formatPackageName(packageName),
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
		return formatter.formatRef(def.AsRef())
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

// formatRef renders a reference to a named type. A reference whose target lives
// in the same schema package is rendered as a bare short name (the type is
// emitted in this same Rust module). A reference into another package records a
// `use crate::types::<package>::<Type>;` through the import collector and then
// uses the bare short name, mirroring how the serde and HashMap imports are
// collected and deduped at the top of the module.
func (formatter *typeFormatter) formatRef(ref ast.RefType) string {
	typeName := formatTypeName(ref.ReferredType)
	referredPkg := formatPackageName(ref.ReferredPkg)

	if referredPkg != "" && referredPkg != formatter.packageName {
		formatter.imports.Add(fmt.Sprintf("crate::types::%s::%s", referredPkg, typeName))
	}

	return typeName
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

// formatConstantRef renders the storage type of a constant reference: the type
// of the object it refers to. A reference to a top-level scalar constant renders
// as that scalar's Rust kind; a reference to an enum renders as the enum type
// (cross-module references are recorded just like a plain ref). The constant
// value itself is pinned separately via the struct's Default impl.
func (formatter *typeFormatter) formatConstantRef(def ast.ConstantReferenceType) string {
	referredObject, found := formatter.context.LocateObject(def.ReferredPkg, def.ReferredType)
	if !found {
		return "serde_json::Value"
	}

	if referredObject.Type.IsScalar() {
		return formatScalarKind(referredObject.Type.AsScalar().ScalarKind)
	}

	if referredObject.Type.IsEnum() {
		return formatter.formatRef(ast.RefType{ReferredPkg: def.ReferredPkg, ReferredType: def.ReferredType})
	}

	return "serde_json::Value"
}
