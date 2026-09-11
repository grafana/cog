package languages

import (
	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ir"
	"github.com/grafana/cog/internal/ir/transforms"
)

type Language interface {
	Name() string
	Jennies(config Config) *codejen.JennyList[Context]
	CompilerPasses() transforms.Transforms
}

type NullableConfig struct {
	Kinds              []ir.Kind
	ProtectArrayAppend bool
	AnyIsNullable      bool
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
