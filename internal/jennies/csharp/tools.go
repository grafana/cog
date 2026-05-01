package csharp

import "strings"

// pascalCase converts a snake_case, kebab-case, or lower-case identifier
// into PascalCase suitable for C# namespace and type names.
//
// Examples:
//
//	pascalCase("dashboard")       -> "Dashboard"
//	pascalCase("library_panel")   -> "LibraryPanel"
//	pascalCase("library-panel")   -> "LibraryPanel"
//	pascalCase("HTTPRequest")     -> "HTTPRequest" (already capitalised)
func pascalCase(input string) string {
	if input == "" {
		return ""
	}

	// Split on common separators used by the schema source identifiers.
	parts := strings.FieldsFunc(input, func(r rune) bool {
		return r == '_' || r == '-' || r == '.' || r == ' '
	})

	var b strings.Builder
	for _, part := range parts {
		if part == "" {
			continue
		}
		b.WriteString(strings.ToUpper(part[:1]))
		b.WriteString(part[1:])
	}
	return b.String()
}
