package typescript

import (
	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast/compiler"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/spf13/cobra"
)

const LanguageRef = "typescript"

type Config struct {
	Debug bool

	RenameOutputFunc func(pkg string) string
	ImportMapper     func(pkg string) string
}

func (config Config) MergeWithGlobal(global common.Config) Config {
	newConfig := config
	newConfig.Debug = global.Debug
	newConfig.RenameOutputFunc = global.RenameOutputFunc
	newConfig.ImportMapper = global.TSConfig.ImportMapper

	return newConfig
}

type Language struct {
	config Config
}

func New() *Language {
	return &Language{
		config: Config{},
	}
}

func (language *Language) RegisterCliFlags(_ *cobra.Command) {
}

func (language *Language) Jennies(globalConfig common.Config) *codejen.JennyList[common.Context] {
	config := language.config.MergeWithGlobal(globalConfig)

	jenny := codejen.JennyListWithNamer[common.Context](func(_ common.Context) string {
		return LanguageRef
	})
	jenny.AppendOneToMany(
		common.If[common.Context](globalConfig.TSConfig.GenRuntime, Runtime{RuntimePath: globalConfig.TSConfig.RuntimePath}),

		common.If[common.Context](globalConfig.Types, RawTypes{Config: config}),
		common.If[common.Context](globalConfig.Builders, &Builder{Config: config}),

		common.If[common.Context](globalConfig.TSConfig.GenTSIndex, Index{Targets: globalConfig}),
	)
	jenny.AddPostprocessors(common.GeneratedCommentHeader(globalConfig))

	return jenny
}

func (language *Language) CompilerPasses() compiler.Passes {
	return compiler.Passes{
		&compiler.RenameNumericEnumValues{},
	}
}
