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
		VariantsPlugins{},
		common.If[languages.Context](globalConfig.Types, RawTypes{Config: config}),
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
		&AddTypehintsComments{},
	}
}

func (language *Language) NullableKinds() languages.NullableConfig {
	return languages.NullableConfig{
		Kinds:              []ast.Kind{ast.KindMap, ast.KindArray},
		ProtectArrayAppend: false,
		AnyIsNullable:      true,
	}
}
