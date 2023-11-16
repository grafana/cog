package typescript

import (
	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast/compiler"
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

func (language *Language) Jennies() *codejen.JennyList[context.Builders] {
	targets := codejen.JennyListWithNamer[context.Builders](func(_ context.Builders) string {
		return LanguageRef
	})
	targets.AppendOneToMany(
		Runtime{},
		RawTypes{},
		&Builder{},
		Index{},
	)

	return targets
}

func (language *Language) CompilerPasses() compiler.Passes {
	return compiler.Passes{
		&compiler.RenameNumericEnumValues{},
	}
}
