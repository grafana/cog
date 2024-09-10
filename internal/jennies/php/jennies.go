package php

import (
	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/ast/compiler"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/languages"
	"github.com/grafana/cog/internal/tools"
)

const LanguageRef = "php"

type Config struct {
	debug bool

	NamespaceRoot string `yaml:"namespace_root"`

	// BuilderTemplatesDirectories holds a list of directories containing templates
	// to be used to override parts of builders.
	BuilderTemplatesDirectories []string `yaml:"builder_templates"`
}

func (config *Config) InterpolateParameters(interpolator func(input string) string) {
	config.NamespaceRoot = interpolator(config.NamespaceRoot)
	config.BuilderTemplatesDirectories = tools.Map(config.BuilderTemplatesDirectories, interpolator)
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

	tmpl := initTemplates(language.config.BuilderTemplatesDirectories)

	jenny := codejen.JennyListWithNamer[languages.Context](func(_ languages.Context) string {
		return LanguageRef
	})
	jenny.AppendOneToMany(
		VariantsPlugins{config: config, tmpl: tmpl},
		Runtime{config: config, tmpl: tmpl},
		common.If[languages.Context](globalConfig.Types, RawTypes{config: config, tmpl: tmpl}),
		common.If[languages.Context](globalConfig.Builders, &Builder{config: config, tmpl: tmpl}),
	)
	jenny.AddPostprocessors(common.GeneratedCommentHeader(globalConfig))

	return jenny
}

func (language *Language) CompilerPasses() compiler.Passes {
	return compiler.Passes{
		&compiler.AnonymousEnumToExplicitType{},
		&compiler.SanitizeEnumMemberNames{},
		&compiler.NotRequiredFieldAsNullableType{},
		&compiler.FlattenDisjunctions{},
		&compiler.DisjunctionWithNullToOptional{},
		&compiler.AnonymousStructsToNamed{},
		&compiler.DisjunctionInferMapping{},
		&compiler.UndiscriminatedDisjunctionToAny{},
		&compiler.InlineObjectsWithTypes{
			InlineTypes: []ast.Kind{ast.KindScalar, ast.KindArray, ast.KindMap, ast.KindDisjunction},
		},
	}
}

func (language *Language) NullableKinds() languages.NullableConfig {
	return languages.NullableConfig{
		ProtectArrayAppend: true,
		AnyIsNullable:      true,
	}
}
