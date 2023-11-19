package typescript

import (
	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast/compiler"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/spf13/cobra"
)

const LanguageRef = "typescript"

type Language struct {
}

func New() *Language {
	return &Language{}
}

func (language *Language) RegisterCliFlags(_ *cobra.Command) {
}

func (language *Language) Jennies(targets common.Targets) *codejen.JennyList[common.Context] {
	jenny := codejen.JennyListWithNamer[common.Context](func(_ common.Context) string {
		return LanguageRef
	})
	jenny.AppendOneToMany(
		Runtime{},

		common.If[common.Context](targets.Types, RawTypes{}),
		common.If[common.Context](targets.Builders, &Builder{}),

		Index{Targets: targets},
	)
	jenny.AddPostprocessors(common.GeneratedCommentHeader())

	return jenny
}

func (language *Language) CompilerPasses() compiler.Passes {
	return compiler.Passes{
		&compiler.RenameNumericEnumValues{},
	}
}
