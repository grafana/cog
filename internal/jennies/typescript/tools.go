package typescript

import (
	"strings"

	"github.com/grafana/cog/internal/tools"
)

func formatIdentifier(name string) string {
	return tools.LowerCamelCase(escapeIdentifier(name))
}

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

func prefixLinesWith(input string, prefix string) string {
	lines := strings.Split(input, "\n")
	prefixed := make([]string, 0, len(lines))

	for _, line := range lines {
		prefixed = append(prefixed, prefix+line)
	}

	return strings.Join(prefixed, "\n")
}
