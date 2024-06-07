package typescript

import (
	"github.com/grafana/cog/internal/tools"
)

func escapeIdentifier(name string) string {
	if isReservedTypescriptKeyword(name) {
		return name + "Val"
	}

	return name
}

func isReservedTypescriptKeyword(input string) bool {
	// see: https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Lexical_grammar#reserved_words
	switch input {
	case "break", "case", "catch", "class", "const", "continue", "debugger", "default", "delete",
		"do", "else", "export", "extends", "false", "finally", "for", "function", "if", "import",
		"in", "instanceof", "new", "null", "return", "super", "switch", "this", "throw", "true",
		"try", "typeof", "var", "void", "while", "with", "let", "static", "yield", "await":
		return true

	default:
		return false
	}
}

var formatPackageName = tools.LowerCamelCase
