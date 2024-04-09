package typescript

import (
	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/ast/compiler"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/languages"
)

const LanguageRef = "typescript"

type Config struct {
	PathPrefix string `yaml:"path_prefix"`

	// SkipRuntime disables runtime-related code generation when enabled.
	// Note: builders can NOT be generated with this flag turned on, as they
	// rely on the runtime to function.
	SkipRuntime bool `yaml:"skip_runtime"`

	// SkipIndex disables the generation of `index.ts` files.
	SkipIndex bool `yaml:"skip_index"`
}

type Language struct {
	config Config
}

func New(config Config) *Language {
	return &Language{
		config: config,
	}
}

func (language *Language) Name() string {
	return LanguageRef
}

func (language *Language) Jennies(globalConfig languages.Config) *codejen.JennyList[languages.Context] {
	jenny := codejen.JennyListWithNamer[languages.Context](func(_ languages.Context) string {
		return LanguageRef
	})
	jenny.AppendOneToMany(
		common.If[languages.Context](!language.config.SkipRuntime, Runtime{}),

		common.If[languages.Context](globalConfig.Types, RawTypes{}),
		common.If[languages.Context](!language.config.SkipRuntime && globalConfig.Builders, &Builder{}),
		common.If[languages.Context](globalConfig.Converters, &Converter{}),

		common.If[languages.Context](!language.config.SkipIndex, Index{Targets: globalConfig}),
	)
	jenny.AddPostprocessors(common.GeneratedCommentHeader(globalConfig))

	return jenny
}

func (language *Language) CompilerPasses() compiler.Passes {
	return compiler.Passes{
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
