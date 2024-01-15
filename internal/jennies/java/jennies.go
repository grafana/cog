package java

import (
	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast/compiler"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/spf13/cobra"
)

const LanguageRef = "java"

type Config struct {
	Debug bool

	GenGettersAndSetters bool
}

type Language struct {
	config Config
}

func (config Config) MergeWithGlobal(global common.Config) Config {
	newConfig := config
	newConfig.Debug = global.Debug

	return newConfig
}

func New() *Language {
	return &Language{config: Config{}}
}

func (language *Language) RegisterCliFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&language.config.GenGettersAndSetters, "java-getters-and-setters", true, "Generate getters and setters for types")
}

func (language *Language) Jennies(globalConfig common.Config) *codejen.JennyList[common.Context] {
	config := language.config.MergeWithGlobal(globalConfig)

	jenny := codejen.JennyListWithNamer[common.Context](func(_ common.Context) string {
		return LanguageRef
	})

	jenny.AppendOneToMany(
		Runtime{},
		common.If[common.Context](globalConfig.Types, RawTypes{config: config}),
	)
	jenny.AddPostprocessors(common.GeneratedCommentHeader(globalConfig))

	return jenny
}

func (language *Language) CompilerPasses() compiler.Passes {
	return compiler.Passes{
		&compiler.AnonymousEnumToExplicitType{},
		&compiler.FlattenDisjunctions{},
		&compiler.DisjunctionToType{},
	}
}
