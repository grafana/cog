package jennies

import (
	"fmt"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast/compiler"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/jennies/golang"
	"github.com/grafana/cog/internal/jennies/java"
	"github.com/grafana/cog/internal/jennies/python"
	"github.com/grafana/cog/internal/jennies/typescript"
	"github.com/spf13/cobra"
)

type LanguageJenny interface {
	Jennies(config common.Config) *codejen.JennyList[common.Context]
	CompilerPasses() compiler.Passes
	RegisterCliFlags(cmd *cobra.Command)
}

type LanguageJennies map[string]LanguageJenny

func (languageJennies LanguageJennies) ForLanguages(languages []string) (LanguageJennies, error) {
	if languages == nil {
		return languageJennies, nil
	}

	filtered := make(LanguageJennies)
	for _, language := range languages {
		if target, exists := languageJennies[language]; exists {
			filtered[language] = target
			continue
		}

		return nil, fmt.Errorf("unknown language '%s'", language)
	}

	return filtered, nil
}

func (languageJennies LanguageJennies) AsLanguageRefs() []string {
	result := make([]string, 0, len(languageJennies))
	for language := range languageJennies {
		result = append(result, language)
	}
	return result
}

func All() LanguageJennies {
	return LanguageJennies{
		golang.LanguageRef:     golang.New(),
		python.LanguageRef:     python.New(),
		typescript.LanguageRef: typescript.New(),
		java.LanguageRef:       java.New(),
	}
}
