package tools

// UniqueFormatted applies format to each item in items and returns unique results, preserving order.
// Two formatted results are considered equal if they produce the same string.
func UniqueFormatted[T any](items []T, format func(T) string) []string {
	seen := make(map[string]struct{})
	result := make([]string, 0, len(items))
	for _, item := range items {
		formatted := format(item)
		if _, exists := seen[formatted]; !exists {
			seen[formatted] = struct{}{}
			result = append(result, formatted)
		}
	}
	return result
}

// UniqueFormattedBy applies format to each item, deduplicates by the string produced by keyFn,
// and returns unique formatted results, preserving order.
func UniqueFormattedBy[T any, F any](items []T, format func(T) F, keyFn func(F) string) []F {
	seen := make(map[string]struct{})
	result := make([]F, 0, len(items))
	for _, item := range items {
		formatted := format(item)
		key := keyFn(formatted)
		if _, exists := seen[key]; !exists {
			seen[key] = struct{}{}
			result = append(result, formatted)
		}
	}
	return result
}
