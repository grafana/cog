package python

import (
	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/ast/compiler"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/languages"
	"github.com/grafana/cog/internal/tools"
)

const LanguageRef = "python"

type Config struct {
	PathPrefix string `yaml:"path_prefix"`

	// SkipRuntime disables runtime-related code generation when enabled.
	// Note: builders can NOT be generated with this flag turned on, as they
	// rely on the runtime to function.
	SkipRuntime bool `yaml:"skip_runtime"`

	// BuilderTemplatesDirectories holds a list of directories containing templates
	// to be used to override parts of builders.
	BuilderTemplatesDirectories []string `yaml:"builder_templates"`
}

func (config *Config) InterpolateParameters(interpolator func(input string) string) {
	config.PathPrefix = interpolator(config.PathPrefix)
	config.BuilderTemplatesDirectories = tools.Map(config.BuilderTemplatesDirectories, interpolator)
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
	tmpl := initTemplates(language.config.BuilderTemplatesDirectories)

	jenny := codejen.JennyListWithNamer[languages.Context](func(_ languages.Context) string {
		return LanguageRef
	})
	jenny.AppendOneToMany(
		ModuleInit{},
		common.If[languages.Context](!language.config.SkipRuntime, Runtime{tmpl: tmpl}),

		common.If[languages.Context](globalConfig.Types, RawTypes{tmpl: tmpl}),
		common.If[languages.Context](!language.config.SkipRuntime && globalConfig.Builders, &Builder{tmpl: tmpl}),

		common.If[languages.Context](globalConfig.APIReference, common.APIReference{
			Collector: language.apiRefCollector,
			Language:  LanguageRef,
			Formatter: apiReferenceFormatter(),
			Tmpl:      tmpl,
		}),
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
