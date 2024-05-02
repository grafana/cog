package jennies

import (
	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast/compiler"
	"github.com/grafana/cog/internal/jennies/common"
)

type LanguageJenny interface {
	Jennies(config common.Config) *codejen.JennyList[common.Context]
	CompilerPasses() compiler.Passes
}

type LanguageJennies map[string]LanguageJenny

func (languageJennies LanguageJennies) AsLanguageRefs() []string {
	result := make([]string, 0, len(languageJennies))
	for language := range languageJennies {
		result = append(result, language)
	}
	return result
}
