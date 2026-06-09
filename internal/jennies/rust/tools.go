package rust

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/tools"
)

// formatTypeName converts an IR object/type name into an idiomatic Rust type
// name (PascalCase), e.g. "someStruct" -> "SomeStruct".
func formatTypeName(name string) string {
	return tools.UpperCamelCase(name)
}

// formatConstName converts an IR name into an idiomatic Rust constant name
// (SCREAMING_SNAKE_CASE), e.g. "constTypeString" -> "CONST_TYPE_STRING".
func formatConstName(name string) string {
	return tools.UpperSnakeCase(name)
}

// formatFieldName converts an IR field name into an idiomatic Rust struct field
// name (snake_case), escaping Rust keywords using raw identifiers where needed.
func formatFieldName(name string) string {
	// Assumption: field names are non-empty wire keys. TrimLeft is a cutset trim,
	// so a leading run of `$`/`_` is stripped; for a name made entirely of those
	// characters this would collapse to an empty identifier. Guard against that
	// by falling back to the snake-cased original name, which never emits `pub :`.
	trimmed := strings.TrimLeft(name, "$_")
	snake := tools.SnakeCase(trimmed)
	if snake == "" {
		snake = tools.SnakeCase(name)
	}
	return escapeRustKeyword(snake)
}

// rustFieldNeedsRename reports whether the idiomatic Rust field identifier
// differs from the original IR field name (which is the JSON wire key), in which
// case a #[serde(rename = "...")] attribute is required to preserve wire
// compatibility. The raw-identifier prefix is stripped before comparing because
// it is a Rust syntax artifact, not part of the logical identifier.
func rustFieldNeedsRename(originalName string) bool {
	rustIdent := strings.TrimPrefix(formatFieldName(originalName), "r#")
	return rustIdent != originalName
}

// escapeRustKeyword wraps reserved Rust keywords in the raw-identifier syntax
// (r#keyword) so they can be used as identifiers. A small set of keywords cannot
// be raw identifiers and are suffixed instead.
func escapeRustKeyword(name string) string {
	if isNonRawableKeyword(name) {
		return name + "_"
	}

	if isReservedRustKeyword(name) {
		return "r#" + name
	}

	return name
}

func isNonRawableKeyword(name string) bool {
	// These keywords are not valid as raw identifiers.
	switch name {
	case "crate", "self", "super", "Self":
		return true
	default:
		return false
	}
}

func isReservedRustKeyword(name string) bool {
	// See: https://doc.rust-lang.org/reference/keywords.html
	switch name {
	case "as", "break", "const", "continue", "dyn", "else", "enum", "extern",
		"false", "fn", "for", "if", "impl", "in", "let", "loop", "match", "mod",
		"move", "mut", "pub", "ref", "return", "static", "struct", "trait",
		"true", "type", "unsafe", "use", "where", "while", "async", "await",
		"abstract", "become", "box", "do", "final", "macro", "override", "priv",
		"try", "typeof", "unsized", "virtual", "yield", "union":
		return true
	default:
		return false
	}
}

// formatScalarKind maps an IR scalar kind to its idiomatic Rust type.
func formatScalarKind(kind ast.ScalarKind) string {
	switch kind {
	case ast.KindNull:
		// A bare null only appears as a disjunction branch; the surrounding
		// Optional handles it, so it never reaches type emission on its own.
		return "()"
	case ast.KindAny:
		return "serde_json::Value"
	case ast.KindBytes:
		return "Vec<u8>"
	case ast.KindString:
		return "String"
	case ast.KindFloat32:
		return "f32"
	case ast.KindFloat64:
		return "f64"
	case ast.KindUint8:
		return "u8"
	case ast.KindUint16:
		return "u16"
	case ast.KindUint32:
		return "u32"
	case ast.KindUint64:
		return "u64"
	case ast.KindInt8:
		return "i8"
	case ast.KindInt16:
		return "i16"
	case ast.KindInt32:
		return "i32"
	case ast.KindInt64:
		return "i64"
	case ast.KindBool:
		return "bool"
	default:
		// Unhandled kinds fall back to the catch-all JSON value, matching
		// formatInnerType's default and keeping emitted Rust compilable.
		return "serde_json::Value"
	}
}

// formatValue renders a Go value (as found in IR defaults/constants) into the
// equivalent Rust literal.
func formatValue(val any) string {
	switch v := val.(type) {
	case nil:
		return "None"
	case bool:
		return strconv.FormatBool(v)
	case string:
		return strconv.Quote(v)
	case float32:
		return formatFloat(float64(v))
	case float64:
		return formatFloat(v)
	case int:
		return strconv.Itoa(v)
	case int64:
		return strconv.FormatInt(v, 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// formatFloat renders a float so that whole numbers still carry a decimal point,
// keeping the literal unambiguously floating-point for the Rust compiler.
func formatFloat(v float64) string {
	formatted := strconv.FormatFloat(v, 'f', -1, 64)
	if !strings.ContainsAny(formatted, ".eE") {
		formatted += ".0"
	}
	return formatted
}
