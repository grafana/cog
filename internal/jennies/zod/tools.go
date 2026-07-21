package zod

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

func formatPackageName(pkg string) string {
	parts := strings.Split(pkg, "/")
	if len(parts) > 1 {
		pkg = parts[len(parts)-1]
	}
	rgx := regexp.MustCompile("[^a-zA-Z0-9_]+")
	return strings.ToLower(rgx.ReplaceAllString(pkg, ""))
}

func (config *Config) pathWithPrefix(parts ...string) string {
	prefix := "src"
	if config.PathPrefix != nil {
		prefix = *config.PathPrefix
	}
	return filepath.Join(append([]string{prefix}, parts...)...)
}

// outputFilePath returns the generated file path for a (formatted) package
// name, honoring PathPrefix, FlatLayout and Filename.
func (config *Config) outputFilePath(pkg string) string {
	filename := "schemas.gen.ts"
	if config.Filename != nil && *config.Filename != "" {
		filename = *config.Filename
	}
	if config.FlatLayout {
		return config.pathWithPrefix(filename)
	}
	return config.pathWithPrefix(pkg, filename)
}

// formatLiteral renders a Go value as a TypeScript literal expression for
// z.literal(...) and .default(...) arguments.
func formatLiteral(val any) string {
	if val == nil {
		return "null"
	}
	switch v := val.(type) {
	case string:
		b, err := json.Marshal(v)
		if err != nil {
			return "null"
		}
		return string(b)
	case bool:
		return fmt.Sprintf("%t", v)
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64:
		return fmt.Sprintf("%v", v)
	case []any:
		parts := make([]string, 0, len(v))
		for _, item := range v {
			parts = append(parts, formatLiteral(item))
		}
		return "[" + strings.Join(parts, ", ") + "]"
	case map[string]any:
		keys := make([]string, 0, len(v))
		for k := range v {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		parts := make([]string, 0, len(keys))
		for _, k := range keys {
			parts = append(parts, fmt.Sprintf("%s: %s", quoteFieldName(k), formatLiteral(v[k])))
		}
		return "{" + strings.Join(parts, ", ") + "}"
	}

	// Fallback covers numeric types from yaml decoding, json.Number, etc.
	b, err := json.Marshal(val)
	if err != nil {
		return "null"
	}
	return string(b)
}

var jsIdentifier = regexp.MustCompile(`^[A-Za-z_$][A-Za-z0-9_$]*$`)

func quoteFieldName(name string) string {
	if jsIdentifier.MatchString(name) {
		return name
	}
	b, err := json.Marshal(name)
	if err != nil {
		return `""`
	}
	return string(b)
}

// joinDescription folds a comment block into a single .describe() argument.
// Returns "" for empty/whitespace-only blocks so callers can skip .describe().
func joinDescription(comments []string) string {
	trimmed := make([]string, 0, len(comments))
	for _, c := range comments {
		c = strings.TrimRight(c, " \t")
		trimmed = append(trimmed, c)
	}
	joined := strings.TrimSpace(strings.Join(trimmed, "\n"))
	return joined
}
