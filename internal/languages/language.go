package languages

import (
	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ir"
)

// Language describes the interface that language generators must implement.
type Language interface {
	// Name returns the name of the language.
	Name() string

	// Transform receives schemas and apply language-specific transformations to them.
	// For example, transformations may include eliminating disjunction (or union/sum types)
	// for languages that don't support them.
	// This method is called by the codegen pipeline before Jennies() is called.
	Transform(schemas ir.Schemas) (ir.Schemas, error)

	// Jennies returns a list of codegen jennies that do the actual codegen.
	Jennies(config Config) *codejen.JennyList[Context]
}

// NullableConfig describes some properties of nullable types for a given language.
type NullableConfig struct {
	// Kinds lists all the kinds that can be null (without involving pointers).
	// Example for go: maps and arrays/slices
	Kinds []ir.Kind

	// ProtectArrayAppend indicates whether appending a value to a null array is allowed
	// or whether the array needs to be initialized first.
	ProtectArrayAppend bool

	// AnyIsNullable indicates that the `any` type for the current language is
	// nullable.
	AnyIsNullable bool
}

type NullableKindsProvider interface {
	NullableKinds() NullableConfig
}

func (nullableConfig NullableConfig) TypeIsNullable(typeDef ir.Type) bool {
	return typeDef.Nullable ||
		(typeDef.IsAny() && nullableConfig.AnyIsNullable) ||
		typeDef.IsAnyOf(nullableConfig.Kinds...)
}

type Languages map[string]Language

func (languages Languages) AsLanguageRefs() []string {
	result := make([]string, 0, len(languages))
	for language := range languages {
		result = append(result, language)
	}
	return result
}
