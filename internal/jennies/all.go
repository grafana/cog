package jennies

import (
	"fmt"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast/compiler"
	"github.com/grafana/cog/internal/jennies/context"
	"github.com/grafana/cog/internal/jennies/golang"
	"github.com/grafana/cog/internal/jennies/typescript"
)

type LanguageTarget struct {
	Jennies        *codejen.JennyList[context.Builders]
	CompilerPasses []compiler.Pass
}

type LanguageTargets map[string]LanguageTarget

func (languageTargets LanguageTargets) ForLanguages(languages []string) (LanguageTargets, error) {
	if languages == nil {
		return languageTargets, nil
	}

	filtered := make(LanguageTargets)

	for _, language := range languages {
		if target, exists := languageTargets[language]; exists {
			filtered[language] = target
			continue
		}

		return nil, fmt.Errorf("unknown language '%s'", language)
	}

	return filtered, nil
}

type Config struct {
	Go golang.Config
}

func All(config Config) LanguageTargets {
	targets := map[string]LanguageTarget{
		golang.LanguageRef: {
			Jennies:        golang.Jennies(config.Go),
			CompilerPasses: golang.CompilerPasses(),
		},
		typescript.LanguageRef: {
			Jennies:        typescript.Jennies(),
			CompilerPasses: typescript.CompilerPasses(),
		},
	}

	return targets
}
