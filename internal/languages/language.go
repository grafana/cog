package languages

import (
	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/ast/compiler"
)

type Language interface {
	Name() string
	Jennies(config Config) *codejen.JennyList[Context]
	CompilerPasses() compiler.Passes
}

type NullableConfig struct {
	Kinds              []ast.Kind
	ProtectArrayAppend bool
	AnyIsNullable      bool
}

type DefaultConfig struct {
	FormatScalarFunc func(t ast.Type, v any) string
	FormatListFunc   func(t ast.Type, v any) string
	FormatEnumFunc   func(name string, t ast.Type, v any) string
	FormatStructFunc func(name string, t ast.Type, v any) string
	FormatMapFunc    func(t ast.Type, v any) string
}

type NullableKindsProvider interface {
	NullableKinds() NullableConfig
}

type DefaultKindsProvider interface {
	DefaultKinds() DefaultConfig
}

type Languages map[string]Language

func (languages Languages) AsLanguageRefs() []string {
	result := make([]string, 0, len(languages))
	for language := range languages {
		result = append(result, language)
	}
	return result
}
