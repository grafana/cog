package java

import (
	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast/compiler"
	"github.com/grafana/cog/internal/jennies/common"
)

const LanguageRef = "java"

type Config struct {
	GenGettersAndSetters bool `yaml:"getters_and_setters"`
}

type Language struct {
	config Config
}

func New(config Config) *Language {
	return &Language{config: config}
}

func (language *Language) Jennies(globalConfig common.Config) *codejen.JennyList[common.Context] {
	jenny := codejen.JennyListWithNamer[common.Context](func(_ common.Context) string {
		return LanguageRef
	})

	jenny.AppendOneToMany(
		Runtime{},
		common.If[common.Context](globalConfig.Types, RawTypes{config: language.config}),
	)
	jenny.AddPostprocessors(common.GeneratedCommentHeader(globalConfig))

	return jenny
}

func (language *Language) CompilerPasses() compiler.Passes {
	return compiler.Passes{
		&compiler.AnonymousEnumToExplicitType{},
		&compiler.AnonymousStructsToNamed{},
		&compiler.FlattenDisjunctions{},
		&compiler.DisjunctionInferMapping{},
		&compiler.DisjunctionToType{},
	}
}
