package java

import (
	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast/compiler"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/spf13/cobra"
)

const LanguageRef = "java"

type Config struct {
	GenGettersAndSetters bool
}

type Language struct {
	config Config
}

func New() *Language {
	return &Language{config: Config{}}
}

func (language *Language) RegisterCliFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&language.config.GenGettersAndSetters, "java-getters-and-setters", true, "Generate getters and setters for types")
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
