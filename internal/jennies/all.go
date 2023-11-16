package jennies

import (
	"fmt"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast/compiler"
	"github.com/grafana/cog/internal/jennies/context"
	"github.com/grafana/cog/internal/jennies/golang"
	"github.com/grafana/cog/internal/jennies/typescript"
	"github.com/spf13/cobra"
)

type LanguageJenny interface {
	Jennies() *codejen.JennyList[context.Builders]
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

func All() LanguageJennies {
	return LanguageJennies{
		golang.LanguageRef:     golang.New(),
		typescript.LanguageRef: typescript.New(),
	}
}
