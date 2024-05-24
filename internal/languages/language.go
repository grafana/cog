package languages

import (
	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/ast/compiler"
	"github.com/grafana/cog/internal/jennies/common"
)

type Language interface {
	Name() string
	Jennies(config common.Config) *codejen.JennyList[common.Context]
	CompilerPasses() compiler.Passes
}

type NullableConfig struct {
	Kinds              []ast.Kind
	ProtectArrayAppend bool
	AnyIsNullable      bool
}

type NullableKindsProvider interface {
	NullableKinds() NullableConfig
}

type Languages map[string]Language

func (languages Languages) AsLanguageRefs() []string {
	result := make([]string, 0, len(languages))
	for language := range languages {
		result = append(result, language)
	}
	return result
}
