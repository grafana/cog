package typescript

import (
	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast/compiler"
	"github.com/grafana/cog/internal/jennies/context"
)

const LanguageRef = "typescript"

func Jennies() *codejen.JennyList[context.Builders] {
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

func CompilerPasses() []compiler.Pass {
	return []compiler.Pass{
		&compiler.PrefixEnumValues{},
	}
}
