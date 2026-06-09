// Package rust implements the Rust code generation target for cog.
package rust

import (
	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/ast/compiler"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/languages"
)

const LanguageRef = "rust"

var _ languages.Language = (*Language)(nil)

type Config struct {
	PathPrefix string `yaml:"path_prefix"`

	// SkipRuntime disables runtime-related code generation when enabled.
	SkipRuntime bool `yaml:"skip_runtime"`
}

func (config *Config) InterpolateParameters(interpolator func(input string) string) {
	config.PathPrefix = interpolator(config.PathPrefix)
}

type Language struct {
	config          Config
	apiRefCollector *common.APIReferenceCollector
}

func New(config Config) *Language {
	return &Language{
		config:          config,
		apiRefCollector: common.NewAPIReferenceCollector(),
	}
}

func (language *Language) Name() string {
	return LanguageRef
}

func (language *Language) Jennies(globalConfig languages.Config) *codejen.JennyList[languages.Context] {
	tmpl := initTemplates(language.config, language.apiRefCollector)

	jenny := codejen.JennyListWithNamer(func(_ languages.Context) string {
		return LanguageRef
	})
	jenny.AppendOneToMany(
		common.If(!language.config.SkipRuntime, Runtime{tmpl: tmpl}),
		RawTypes{config: language.config, apiRefCollector: language.apiRefCollector},
		common.If(!language.config.SkipRuntime && globalConfig.Builders, Builder{config: language.config, apiRefCollector: language.apiRefCollector}),
	)
	jenny.AddPostprocessors(common.GeneratedCommentHeader(globalConfig))
	jenny.AddPostprocessors(FormatRustFiles)

	if language.config.PathPrefix != "" {
		jenny.AddPostprocessors(common.PathPrefixer(language.config.PathPrefix))
	}

	return jenny
}

func (language *Language) CompilerPasses() compiler.Passes {
	return compiler.Passes{
		&compiler.AnonymousStructsToNamed{},
		&compiler.AnonymousEnumToExplicitType{},
		&compiler.NotRequiredFieldAsNullableType{},
		&compiler.DisjunctionWithNullToOptional{},
		&compiler.DisjunctionOfConstantsToEnum{},
		&compiler.FlattenDisjunctions{},
		&compiler.DisjunctionInferMapping{},
		&compiler.RenameNumericEnumValues{},
		&compiler.DisjunctionPropagateVariant{},
	}
}

func (language *Language) NullableKinds() languages.NullableConfig {
	return languages.NullableConfig{
		Kinds:              []ast.Kind{ast.KindMap, ast.KindArray, ast.KindRef, ast.KindStruct},
		ProtectArrayAppend: true,
		AnyIsNullable:      true,
	}
}
