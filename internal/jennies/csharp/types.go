package csharp

import (
	"fmt"
	"strings"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/languages"
	"github.com/grafana/cog/internal/tools"
)

// typeFormatter renders an ast.Type as a C# type expression and tracks
// the namespaces that need to be imported as a side effect.
type typeFormatter struct {
	config  Config
	context languages.Context
	imports *importMap
}

func newTypeFormatter(ctx languages.Context, config Config, imports *importMap) *typeFormatter {
	return &typeFormatter{
		context: ctx,
		config:  config,
		imports: imports,
	}
}

// withImports returns a copy of the formatter that records imports on the
// given map. Used when the formatter is reused across files in the same
// schema (one importMap per output file).
func (tf *typeFormatter) withImports(imports *importMap) *typeFormatter {
	clone := *tf
	clone.imports = imports
	return &clone
}

// formatFieldType renders the C# type expression for a field/property.
func (tf *typeFormatter) formatFieldType(def ast.Type) string {
	switch def.Kind {
	case ast.KindScalar:
		return formatScalarType(def.AsScalar())
	case ast.KindRef:
		return tf.formatReference(def.AsRef())
	case ast.KindArray:
		return tf.formatArray(def.AsArray())
	case ast.KindComposableSlot:
		return tf.formatComposable(def.AsComposableSlot())
	case ast.KindMap:
		return tf.formatMap(def.AsMap())
	case ast.KindStruct:
		// Anonymous structs should have been lifted by the
		// AnonymousStructsToNamed compiler pass; if one slips through,
		// fall back to `object`.
		return "object"
	case ast.KindConstantRef:
		return tf.formatConstantReference(def.AsConstantRef())
	}
	return "object"
}

// formatComposable renders a composable-slot type. Slots resolve to
// their variant interface (e.g. `Cog.Variants.Dataquery`) which lives
// under `<NamespaceRoot>.Cog.Variants`. The fully-qualified form is
// emitted so no extra `using` is required at the call site.
func (tf *typeFormatter) formatComposable(def ast.ComposableSlotType) string {
	return fmt.Sprintf("Cog.Variants.%s", formatObjectName(string(def.Variant)))
}

func (tf *typeFormatter) formatReference(def ast.RefType) string {
	object, found := tf.context.LocateObjectByRef(def)
	if found {
		switch object.Type.Kind {
		case ast.KindScalar:
			return formatScalarType(object.Type.AsScalar())
		case ast.KindMap:
			return tf.formatMap(object.Type.AsMap())
		case ast.KindArray:
			return tf.formatArray(object.Type.AsArray())
		}
	}
	tf.imports.addPackage(def.ReferredPkg)
	return formatObjectName(def.ReferredType)
}

func (tf *typeFormatter) formatConstantReference(def ast.ConstantReferenceType) string {
	object, found := tf.context.LocateObject(def.ReferredPkg, def.ReferredType)
	if !found {
		return "object"
	}
	if object.Type.IsEnum() {
		tf.imports.addPackage(def.ReferredPkg)
		return formatObjectName(def.ReferredType)
	}
	if object.Type.IsScalar() {
		return formatScalarType(object.Type.AsScalar())
	}
	return "object"
}

func (tf *typeFormatter) formatArray(def ast.ArrayType) string {
	tf.imports.addNamespace("System.Collections.Generic")
	return fmt.Sprintf("List<%s>", tf.formatFieldType(def.ValueType))
}

func (tf *typeFormatter) formatMap(def ast.MapType) string {
	tf.imports.addNamespace("System.Collections.Generic")
	keyType := formatScalarType(def.IndexType.AsScalar())
	return fmt.Sprintf("Dictionary<%s, %s>", keyType, tf.formatFieldType(def.ValueType))
}

// emptyValueForType returns a sensible default-value expression used in
// the parameterless constructor when no explicit default is provided.
//
// When useBuilders is true, refs that have a registered builder fall
// back to invoking the builder's `Build()` method (used by Builder
// generation when a nested-builder slot needs an empty placeholder).
func (tf *typeFormatter) emptyValueForType(def ast.Type) string {
	return tf.emptyValueForTypeOpts(def, false)
}

func (tf *typeFormatter) emptyValueForTypeOpts(def ast.Type, useBuilders bool) string {
	switch def.Kind {
	case ast.KindArray:
		tf.imports.addNamespace("System.Collections.Generic")
		return fmt.Sprintf("new %s()", tf.formatArray(def.AsArray()))
	case ast.KindMap:
		tf.imports.addNamespace("System.Collections.Generic")
		return fmt.Sprintf("new %s()", tf.formatMap(def.AsMap()))
	case ast.KindRef:
		referred, found := tf.context.LocateObjectByRef(def.AsRef())
		if found && referred.Type.IsEnum() {
			defaultMember := referred.Type.AsEnum().Values[0]
			tf.imports.addPackage(def.AsRef().ReferredPkg)
			return fmt.Sprintf("%s.%s", formatObjectName(referred.Name), formatObjectName(defaultMember.Name))
		}
		if useBuilders && tf.typeHasBuilder(def) {
			tf.imports.addPackage(def.AsRef().ReferredPkg)
			return fmt.Sprintf("new %sBuilder().Build()", formatObjectName(def.AsRef().ReferredType))
		}
		tf.imports.addPackage(def.AsRef().ReferredPkg)
		return fmt.Sprintf("new %s()", formatObjectName(def.AsRef().ReferredType))
	case ast.KindStruct:
		return "new object()"
	case ast.KindScalar:
		switch def.AsScalar().ScalarKind {
		case ast.KindBool:
			return "false"
		case ast.KindFloat32:
			return "0f"
		case ast.KindFloat64:
			return "0d"
		case ast.KindInt8, ast.KindInt16, ast.KindInt32, ast.KindUint8, ast.KindUint16, ast.KindUint32:
			return "0"
		case ast.KindInt64:
			return "0L"
		case ast.KindUint64:
			return "0UL"
		case ast.KindString:
			return `""`
		case ast.KindBytes:
			return "(byte) 0"
		case ast.KindAny:
			return "new object()"
		}
	}
	return "default!"
}

// ---------------------------------------------------------------------
// Builder-jenny helpers (used only by the Builder template).
// ---------------------------------------------------------------------

// typeHasBuilder reports whether a builder is registered for the given
// type. Used by the Builder template to decide whether option arguments
// should accept an `IBuilder<T>` instead of the raw `T`.
func (tf *typeFormatter) typeHasBuilder(def ast.Type) bool {
	return tf.context.ResolveToBuilder(def)
}

// resolvesToComposableSlot reports whether the type ultimately resolves
// to a composable slot (e.g. `Dataquery`). Composable slots are
// rendered as their variant interface and surfaced via `IBuilder<T>`.
func (tf *typeFormatter) resolvesToComposableSlot(def ast.Type) bool {
	_, found := tf.context.ResolveToComposableSlot(def)
	return found
}

// formatBuilderFieldType renders the C# type expression used for a
// builder option's parameter slot. When the underlying type has a
// builder (or resolves to a composable slot), it is wrapped in
// `Cog.IBuilder<T>` (or `List<…>` / `Dictionary<string, …>` for
// collections of builders).
func (tf *typeFormatter) formatBuilderFieldType(def ast.Type) string {
	if tf.resolvesToComposableSlot(def) || tf.typeHasBuilder(def) {
		switch def.Kind {
		case ast.KindArray:
			return tf.formatBuilderCollectionFields(def.AsArray().ValueType, "List", "List<")
		case ast.KindMap:
			keyType := formatScalarType(def.AsMap().IndexType.AsScalar())
			return tf.formatBuilderCollectionFields(def.AsMap().ValueType, "Dictionary", fmt.Sprintf("Dictionary<%s, ", keyType))
		default:
			tf.imports.addPackage("cog")
			return fmt.Sprintf("Cog.IBuilder<%s>", tf.formatFieldType(def))
		}
	}

	return tf.formatFieldType(def)
}

func (tf *typeFormatter) formatBuilderCollectionFields(def ast.Type, collection string, prefix string) string {
	tf.imports.addNamespace("System.Collections.Generic")
	_ = collection
	if def.Kind == ast.KindArray || def.Kind == ast.KindMap {
		return fmt.Sprintf("%s%s>", prefix, tf.formatBuilderFieldType(def))
	}

	tf.imports.addPackage("cog")
	return fmt.Sprintf("%sCog.IBuilder<%s>>", prefix, tf.formatFieldType(def))
}

// formatFieldPath renders a dotted field-access expression for read
// positions (e.g. `parent.child.field`). Path elements after an `any`
// boundary are dropped because the value beyond that point is opaque.
func (tf *typeFormatter) formatFieldPath(fieldPath ast.Path) string {
	parts := make([]string, 0)
	for i, part := range fieldPath {
		output := formatFieldName(part.Identifier)

		if i > 0 && fieldPath[i-1].Type.IsAny() {
			return output
		}

		parts = append(parts, output)
	}

	return strings.Join(parts, ".")
}

// formatPathIndex renders a map/array indexing expression (the value
// inside the `[…]`). Constants become C# literals; argument-driven
// indices reference the formatted argument name.
func (tf *typeFormatter) formatPathIndex(pathIndex *ast.PathIndex) string {
	if pathIndex.Constant != nil {
		switch v := pathIndex.Constant.(type) {
		case string:
			return fmt.Sprintf("%q", v)
		default:
			return fmt.Sprintf("%v", v)
		}
	}
	return formatArgName(pathIndex.Argument.Name)
}

// formatAssignmentPath renders the LHS of an assignment, walking the
// AST path from the resource root (typically `this.@internal`) down to
// the field being assigned. Map indexing nodes produce `[key]` access;
// the trailing element keeps its dot syntax even when annotated as
// `any` so Phase 4 stays out of the cast-emission business.
func (tf *typeFormatter) formatAssignmentPath(resourceRoot string, fieldPath ast.Path) string {
	path := resourceRoot

	for i := range fieldPath {
		output := fieldPath[i].Identifier
		if !fieldPath[i].Root {
			output = formatFieldName(output)
		}

		if fieldPath[i].Index != nil {
			path += output + "[" + tf.formatPathIndex(fieldPath[i].Index) + "]"
			continue
		}

		// Without a TypeHint (or as the trailing element) we just
		// append the field; the cast path used by Java to navigate
		// `any` values isn't needed for the C# scenarios we currently
		// support and would obscure generated code.
		if !fieldPath[i].Type.IsAny() || fieldPath[i].TypeHint == nil || i == len(fieldPath)-1 {
			path += "." + output
			continue
		}

		path = fmt.Sprintf("((%s) %s.%s)", tf.formatReference(fieldPath[i].TypeHint.AsRef()), path, output)
	}

	return path
}

// formatRefType renders a literal value of the destination type. It is
// used by constant assignments (e.g. discriminator fields). For enum
// destinations the corresponding member is emitted; otherwise we fall
// back to a Go-style %#v rendering.
func (tf *typeFormatter) formatRefType(destinationType ast.Type, value any) string {
	if !destinationType.IsRef() {
		if value == nil {
			return "null"
		}
		switch v := value.(type) {
		case string:
			return fmt.Sprintf("%q", v)
		case bool:
			if v {
				return "true"
			}
			return "false"
		default:
			return fmt.Sprintf("%v", v)
		}
	}

	referred, found := tf.context.LocateObjectByRef(destinationType.AsRef())
	if !found {
		return fmt.Sprintf("%#v", value)
	}
	if referred.Type.IsEnum() {
		member, _ := referred.Type.AsEnum().MemberForValue(value)
		tf.imports.addPackage(referred.SelfRef.ReferredPkg)
		return fmt.Sprintf("%s.%s", formatObjectName(referred.Name), formatEnumMemberName(member.Name))
	}
	return fmt.Sprintf("%#v", value)
}

// _unused keeps the tools import alive when the file is trimmed during
// edits.
var _ = tools.UpperCamelCase
