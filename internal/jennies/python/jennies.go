package python

import (
	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/ast/compiler"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/languages"
)

const LanguageRef = "python"

type Config struct {
	PathPrefix string `yaml:"path_prefix"`
}

func (config *Config) InterpolateParameters(interpolator func(input string) string) {
	config.PathPrefix = interpolator(config.PathPrefix)
}

type Language struct {
	config Config
}

func New(config Config) *Language {
	return &Language{
		config: config,
	}
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
		&compiler.RenameNumericEnumValues{},
	}
}

func (language *Language) NullableKinds() languages.NullableConfig {
	return languages.NullableConfig{
		Kinds:              []ast.Kind{ast.KindMap, ast.KindArray, ast.KindRef, ast.KindStruct},
		ProtectArrayAppend: true,
		AnyIsNullable:      true,
	}
}
