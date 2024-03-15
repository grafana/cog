package openapi

import (
	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast/compiler"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/spf13/cobra"
)

const LanguageRef = "openapi"

type Config struct {
	Debug bool
}

func (config Config) MergeWithGlobal(global common.Config) Config {
	newConfig := config
	newConfig.Debug = global.Debug

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
		Schema{Config: config},
	)

	return jenny
}

func (language *Language) CompilerPasses() compiler.Passes {
	return compiler.Passes{
		// should be a superset of the compiler passes defined for jsonschema jennies
		&compiler.DisjunctionWithNullToOptional{},
	}
}
