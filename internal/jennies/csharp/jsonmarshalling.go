package csharp

import (
	"fmt"

	"github.com/grafana/cog/internal/ast"
)

// jsonMarshaller centralises System.Text.Json related code generation
// hooks. Activated only when Config.GenerateJSONConverters is true; all
// methods return the empty string / no-op otherwise so callers can use
// them unconditionally.
type jsonMarshaller struct {
	config Config
}

func newJSONMarshaller(config Config) jsonMarshaller {
	return jsonMarshaller{config: config}
}

// enabled reports whether JSON-marshalling code should be emitted.
func (j jsonMarshaller) enabled() bool {
	return j.config.GenerateJSONConverters
}

// fieldPropertyAttribute returns the `[JsonPropertyName("origName")]`
// attribute for a struct field, ensuring round-tripping with the
// original (typically camelCase) JSON name even though our generated
// fields are PascalCased.
func (j jsonMarshaller) fieldPropertyAttribute(field ast.StructField) string {
	if !j.enabled() {
		return ""
	}
	return fmt.Sprintf("[JsonPropertyName(%q)]", field.Name)
}

// classConverterAttribute returns `[JsonConverter(typeof(<T>JsonConverter))]`
// for disjunction-derived structs that need a custom converter, or the
// empty string otherwise. The Phase 4 builder layer will emit the
// matching JsonConverter<T> class via the Converters jenny.
func (j jsonMarshaller) classConverterAttribute(object ast.Object) string {
	if !j.enabled() || !needsCustomConverter(object) {
		return ""
	}
	return fmt.Sprintf("[JsonConverter(typeof(%sJsonConverter))]", formatObjectName(object.Name))
}

// enumConverterAttribute returns the `[JsonConverter]` attribute for
// string enums so that values are serialized as their original string
// payload (declared via [JsonStringEnumMemberName("...")]).
func (j jsonMarshaller) enumConverterAttribute(name string) string {
	if !j.enabled() {
		return ""
	}
	return fmt.Sprintf("[JsonConverter(typeof(JsonStringEnumConverter<%s>))]", name)
}

// needsCustomConverter reports whether the type requires a generated
// JsonConverter<T> alongside its class definition.
//
// We mirror the Java jenny's set of disjunction hints: scalar-only,
// discriminated-refs, and the scalars+refs combination.
func needsCustomConverter(object ast.Object) bool {
	if !object.Type.IsStruct() {
		return false
	}
	t := object.Type
	return t.HasHint(ast.HintDisjunctionOfScalars) ||
		t.HasHint(ast.HintDiscriminatedDisjunctionOfRefs) ||
		t.HasHint(ast.HintDisjunctionOfScalarsAndRefs)
}

// disjunctionFieldType wraps a value-type scalar in `T?` so a "branch
// not selected" can be distinguished from "branch holds the zero value".
// Reference types are already nullable under NRT and need no change.
func disjunctionFieldType(typeExpr string, fieldType ast.Type) string {
	if !isCSharpValueType(fieldType) {
		return typeExpr
	}
	return typeExpr + "?"
}

// isCSharpValueType reports whether the C# type emitted for `t` is a
// value type (and therefore needs an explicit `?` suffix to be
// nullable).
func isCSharpValueType(t ast.Type) bool {
	if !t.IsScalar() {
		return false
	}
	switch t.AsScalar().ScalarKind {
	case ast.KindString, ast.KindAny, ast.KindBytes:
		return false
	default:
		// All numeric scalars (incl. bool) are value types.
		return true
	}
}
