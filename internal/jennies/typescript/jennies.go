package typescript

import (
	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast/compiler"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/jennies/context"
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

func (language *Language) Jennies(targets common.Targets) *codejen.JennyList[context.Builders] {
	jenny := codejen.JennyListWithNamer[context.Builders](func(_ context.Builders) string {
		return LanguageRef
	})
	jenny.AppendOneToMany(
		Runtime{},

		common.If[context.Builders](targets.Types, RawTypes{}),
		common.If[context.Builders](targets.Builders, &Builder{}),

		Index{Targets: targets},
	)

	return jenny
}

func (language *Language) CompilerPasses() compiler.Passes {
	return compiler.Passes{
		&compiler.RenameNumericEnumValues{},
	}
}
