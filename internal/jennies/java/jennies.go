package java

import (
	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast/compiler"
	"github.com/grafana/cog/internal/jennies/common"
)

const LanguageRef = "java"

type Config struct {
}

type Language struct {
	config Config
}

func New() *Language {
	return &Language{config: Config{}}
}

func (language *Language) Jennies(globalConfig common.Config) *codejen.JennyList[common.Context] {
	jenny := codejen.JennyListWithNamer[common.Context](func(_ common.Context) string {
		return LanguageRef
	})

	jenny.AppendOneToMany(
		Runtime{},
		RawTypes{},
	)
	jenny.AddPostprocessors(common.GeneratedCommentHeader(globalConfig))

	return jenny
}

func (language *Language) CompilerPasses() compiler.Passes {
	return compiler.Passes{
		&compiler.AnonymousEnumToExplicitType{},
		&compiler.AnonymousStructsToNamed{},
		&compiler.NotRequiredFieldAsNullableType{},
		&compiler.FlattenDisjunctions{},
		&compiler.DisjunctionWithNullToOptional{},
		&compiler.DisjunctionInferMapping{},
		&compiler.DisjunctionToType{},
	}
}
