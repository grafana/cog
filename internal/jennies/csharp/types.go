package csharp

import (
	"fmt"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/languages"
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
func (tf *typeFormatter) emptyValueForType(def ast.Type) string {
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
