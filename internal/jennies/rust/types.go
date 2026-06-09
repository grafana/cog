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
	// rawPackage is the unsanitized IR package of the schema currently being
	// emitted, used when synthesizing same-package refs (e.g. hoisted disjunction
	// enums) so they render as bare short names.
	rawPackage string
	// importSamePackageRefs forces refs whose referred package matches the current
	// package to still be imported (from crate::types::<pkg>). The builder jenny
	// sets this: a builder module is distinct from its types module, so even a
	// same-package type it names (an enum-typed option argument, say) must be
	// brought in explicitly rather than assumed to be in scope.
	importSamePackageRefs bool
}

func newTypeFormatter(context languages.Context, imports *importMap, packageName string) *typeFormatter {
	return &typeFormatter{
		context:     context,
		imports:     imports,
		packageName: formatPackageName(packageName),
		rawPackage:  packageName,
	}
}

// formatType renders an arbitrary IR type. A nullable (optional) type is wrapped
// in Option<T>, except for arrays and maps: Vec and HashMap carry their own
// "empty" representation (an empty collection), so a nullable collection is
// rendered as the bare Vec/HashMap and its absence is modelled by emptiness.
// This mirrors the Go target, which likewise never pointer-wraps slices or maps.
func (formatter *typeFormatter) formatType(def ast.Type) string {
	inner := formatter.formatInnerType(def)

	if isOptionWrapped(def) {
		return fmt.Sprintf("Option<%s>", inner)
	}

	return inner
}

// isBareCollection reports whether a type is rendered as a bare collection (Vec
// or HashMap) rather than wrapped in Option, even when it is nullable. A
// collection carries its own "empty" representation, so absence is modelled by
// emptiness, mirroring the Go target which never pointer-wraps slices or maps.
func isBareCollection(def ast.Type) bool {
	return def.IsArray() || def.IsMap()
}

// isOptionWrapped reports whether a type is rendered as Option<T>. A nullable
// non-collection type is the optional case; collections stay bare (see
// isBareCollection). This is the single predicate the type formatter, the serde
// attribute generator and the constant-default decision all share, so the serde
// attribute can never drift from what formatType actually emits.
func isOptionWrapped(def ast.Type) bool {
	return def.Nullable && !isBareCollection(def)
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
		// An inline (field-level) disjunction reaches this branch: the Rust target
		// does not run the DisjunctionToType pass (it would rewrite the top-level
		// untagged-enum disjunctions emitted in earlier phases into Go-style
		// struct-of-options), so a disjunction that was never promoted to a named
		// top-level object is rendered as the permissive serde_json::Value. This
		// round-trips any branch and stays wire-compatible; a typed union is only
		// emitted for disjunctions that exist as named top-level objects.
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
	// A reference to a concrete scalar object resolves to its underlying Rust
	// scalar type. Such an object is emitted as a `const` item (see
	// RawTypes.formatConstant), not a type alias, so its name is not a usable type;
	// the field instead takes the underlying scalar kind, matching the Go target
	// (which renders the field as the bare scalar). The pinned value, if used as a
	// field default, is restored separately through the owning struct's Default.
	if referred, found := formatter.context.LocateObject(ref.ReferredPkg, ref.ReferredType); found && referred.Type.IsConcreteScalar() {
		return formatScalarKind(referred.Type.AsScalar().ScalarKind)
	}

	typeName := formatTypeName(ref.ReferredType)
	referredPkg := formatPackageName(ref.ReferredPkg)

	if referredPkg != "" && (formatter.importSamePackageRefs || referredPkg != formatter.packageName) {
		formatter.imports.Add(fmt.Sprintf("crate::types::%s::%s", referredPkg, typeName))
	}

	return typeName
}

// formatBuilderArgType renders a builder-valued option argument's type. A ref
// whose target has a builder becomes `impl cog::Builder<T>` (the idiomatic Rust
// equivalent of Go's `cog.Builder[T]`); a collection of builders becomes
// Vec<...>/HashMap<_, ...> wrapping the same bound recursively, matching Go's
// `[]cog.Builder[T]` / `map[K]cog.Builder[T]`. The referenced type is imported so
// the bound resolves. Only types for which ResolveToBuilder is true reach here.
func (formatter *typeFormatter) formatBuilderArgType(def ast.Type) string {
	switch {
	case def.IsArray():
		return fmt.Sprintf("Vec<%s>", formatter.formatBuilderArgType(def.AsArray().ValueType))
	case def.IsMap():
		m := def.AsMap()
		formatter.imports.Add("std::collections::HashMap")
		return fmt.Sprintf("HashMap<%s, %s>", formatter.formatType(m.IndexType), formatter.formatBuilderArgType(m.ValueType))
	case def.IsRef():
		// formatRef records the cross-module import for the referred type.
		return fmt.Sprintf("impl cog::Builder<%s>", formatter.formatRef(def.AsRef()))
	default:
		// Defensive: a non-ref, non-collection that resolved to a builder should not
		// occur in the in-scope fixtures. Fall back to the plain rendering.
		return formatter.formatType(def)
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
