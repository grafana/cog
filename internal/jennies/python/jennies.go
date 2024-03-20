package python

import (
	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast/compiler"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/spf13/cobra"
)

const LanguageRef = "python"

type Config struct {
	PathPrefix string
}

type Language struct {
	config Config
}

func New() *Language {
	return &Language{
		config: Config{},
	}
}

func (language *Language) RegisterCliFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&language.config.PathPrefix, "python-path-prefix", "", "Python path prefix.")
}

func (language *Language) Jennies(globalConfig common.Config) *codejen.JennyList[common.Context] {
	jenny := codejen.JennyListWithNamer[common.Context](func(_ common.Context) string {
		return LanguageRef
	})
	jenny.AppendOneToMany(
		ModuleInit{},
		Runtime{},

		common.If[common.Context](globalConfig.Types, RawTypes{}),
		common.If[common.Context](globalConfig.Builders, &Builder{}),
	)
	jenny.AddPostprocessors(common.GeneratedCommentHeader(globalConfig))

	if language.config.PathPrefix != "" {
		jenny.AddPostprocessors(common.PathPrefixer(language.config.PathPrefix))
	}

	return jenny
}

func (language *Language) CompilerPasses() compiler.Passes {
	return compiler.Passes{
		&compiler.FlattenDisjunctions{},
		&compiler.DisjunctionWithNullToOptional{},
		&compiler.DisjunctionInferMapping{},
		&compiler.NotRequiredFieldAsNullableType{},
		&compiler.AnonymousStructsToNamed{},
	}
}
