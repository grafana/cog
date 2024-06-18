package php

import (
	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/ast/compiler"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/languages"
)

const LanguageRef = "php"

type Config struct {
	debug bool

	NamespaceRoot string `yaml:"namespace_root"`
}

func (config *Config) InterpolateParameters(interpolator func(input string) string) {
	config.NamespaceRoot = interpolator(config.NamespaceRoot)
}

func (config Config) fullNamespace(typeName string) string {
	return config.NamespaceRoot + "\\" + typeName
}

func (config Config) fullNamespaceRef(typeName string) string {
	return "\\" + config.fullNamespace(typeName)
}

func (config Config) MergeWithGlobal(global languages.Config) Config {
	newConfig := config
	newConfig.debug = global.Debug

	return newConfig
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
	config := language.config.MergeWithGlobal(globalConfig)

	jenny := codejen.JennyListWithNamer[languages.Context](func(_ languages.Context) string {
		return LanguageRef
	})
	jenny.AppendOneToMany(
		VariantsPlugins{config: config},
		common.If[languages.Context](globalConfig.Types, RawTypes{config: config}),
		common.If[languages.Context](globalConfig.Builders, &Builder{config: config}),
	)
	jenny.AddPostprocessors(common.GeneratedCommentHeader(globalConfig))

	return jenny
}

func (language *Language) CompilerPasses() compiler.Passes {
	return compiler.Passes{
		&compiler.InlineObjectsWithTypes{
			InlineTypes: []ast.Kind{ast.KindScalar, ast.KindArray, ast.KindMap, ast.KindDisjunction},
		},
		&compiler.AnonymousEnumToExplicitType{},
		&compiler.SanitizeEnumMemberNames{},
		&compiler.NotRequiredFieldAsNullableType{},
		&compiler.FlattenDisjunctions{},
		&compiler.DisjunctionWithNullToOptional{},
		&compiler.AnonymousStructsToNamed{},
		&compiler.DisjunctionInferMapping{},
		&compiler.UndiscriminatedDisjunctionToAny{},
	}
}

func (language *Language) NullableKinds() languages.NullableConfig {
	return languages.NullableConfig{
		ProtectArrayAppend: true,
		AnyIsNullable:      true,
	}
}
