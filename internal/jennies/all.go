package jennies

import (
	"fmt"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/ast/compiler"
	"github.com/grafana/cog/internal/jennies/golang"
	"github.com/grafana/cog/internal/jennies/typescript"
	"github.com/grafana/cog/internal/veneers/rewrite"
)

type LanguageTarget struct {
	Jennies        *codejen.JennyList[[]*ast.Schema]
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

func All(veneerRewriter *rewrite.Rewriter) LanguageTargets {
	targets := map[string]LanguageTarget{
		golang.LanguageRef: {
			Jennies:        golang.Jennies(veneerRewriter),
			CompilerPasses: golang.CompilerPasses(),
		},
		typescript.LanguageRef: {
			Jennies:        typescript.Jennies(veneerRewriter),
			CompilerPasses: typescript.CompilerPasses(),
		},
	}

	return targets
}
