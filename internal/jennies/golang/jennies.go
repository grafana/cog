package golang

import (
	"fmt"
	"strings"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/ast/compiler"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/languages"
	"github.com/grafana/cog/internal/tools"
)

const LanguageRef = "go"

type Config struct {
	debug              bool
	generateBuilders   bool
	generateConverters bool

	// GenerateStrictUnmarshaller controls the generation of
	// `UnmarshalJSONStrict()` methods on types.
	GenerateStrictUnmarshaller bool `yaml:"generate_strict_unmarshaller"`

	// GenerateEqual controls the generation of `Equal()` methods on types.
	GenerateEqual bool `yaml:"generate_equal"`

	// GenerateValidate controls the generation of `Validate()` methods on types.
	GenerateValidate bool `yaml:"generate_validate"`

	// SkipRuntime disables runtime-related code generation when enabled.
	// Note: builders can NOT be generated with this flag turned on, as they
	// rely on the runtime to function.
	SkipRuntime bool `yaml:"skip_runtime"`

	// OverridesTemplatesDirectories holds a list of directories containing templates
	// defining blocks used to override parts of builders/types/....
	OverridesTemplatesDirectories []string `yaml:"overrides_templates"`

	// ExtraFilesTemplatesDirectories holds a list of directories containing
	// templates describing files to be added to the generated output.
	ExtraFilesTemplatesDirectories []string `yaml:"extra_files_templates"`

	// ExtraFilesTemplatesData holds additional data to be injected into the
	// templates described in ExtraFilesTemplatesDirectories.
	ExtraFilesTemplatesData map[string]string `yaml:"-"`

	// Root path for imports.
	// Ex: github.com/grafana/cog/generated
	PackageRoot string `yaml:"package_root"`
}

func (config *Config) InterpolateParameters(interpolator func(input string) string) {
	config.PackageRoot = interpolator(config.PackageRoot)
	config.OverridesTemplatesDirectories = tools.Map(config.OverridesTemplatesDirectories, interpolator)
	config.ExtraFilesTemplatesDirectories = tools.Map(config.ExtraFilesTemplatesDirectories, interpolator)
}

func (config Config) MergeWithGlobal(global languages.Config) Config {
	newConfig := config
	newConfig.debug = global.Debug
	newConfig.generateBuilders = global.Builders
	newConfig.generateConverters = global.Converters

	return newConfig
}

func (config Config) importPath(suffix string) string {
	root := strings.TrimSuffix(config.PackageRoot, "/")
	return fmt.Sprintf("%s/%s", root, suffix)
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
	config := language.config.MergeWithGlobal(globalConfig)

	tmpl := initTemplates(language.config.OverridesTemplatesDirectories)

	jenny := codejen.JennyListWithNamer[languages.Context](func(_ languages.Context) string {
		return LanguageRef
	})
	jenny.AppendOneToMany(
		common.If[languages.Context](!config.SkipRuntime, Runtime{Config: config, Tmpl: tmpl}),

		common.If[languages.Context](globalConfig.Types, RawTypes{Config: config, Tmpl: tmpl, apiRefCollector: language.apiRefCollector}),

		common.If[languages.Context](!config.SkipRuntime && globalConfig.Builders, &Builder{Config: config, Tmpl: tmpl, apiRefCollector: language.apiRefCollector}),
		common.If[languages.Context](!config.SkipRuntime && globalConfig.Builders && globalConfig.Converters, &Converter{
			Config:          config,
			Tmpl:            tmpl,
			NullableConfig:  language.NullableKinds(),
			apiRefCollector: language.apiRefCollector,
		}),

		common.If[languages.Context](globalConfig.APIReference, common.APIReference{
			Collector: language.apiRefCollector,
			Language:  LanguageRef,
			Formatter: apiReferenceFormatter(config),
			Tmpl:      tmpl,
		}),

		common.CustomTemplates{
			TemplateDirectories: config.ExtraFilesTemplatesDirectories,
			Data: map[string]any{
				"Debug":       config.debug,
				"PackageRoot": config.PackageRoot,
			},
			ExtraData: config.ExtraFilesTemplatesData,
			TmplFuncs: formattingTemplateFuncs(),
		},
	)
	jenny.AddPostprocessors(formatGoFiles, common.GeneratedCommentHeader(globalConfig))

	return jenny
}

func (language *Language) CompilerPasses() compiler.Passes {
	return compiler.Passes{
		&compiler.AnonymousEnumToExplicitType{},
		&compiler.AnonymousStructsToNamed{},
		&compiler.PrefixEnumValues{},
		&compiler.NotRequiredFieldAsNullableType{},
		&compiler.FlattenDisjunctions{},
		&compiler.DisjunctionWithNullToOptional{},
		&compiler.DisjunctionOfAnonymousStructsToExplicit{},
		&compiler.DisjunctionInferMapping{},
		&compiler.UndiscriminatedDisjunctionToAny{},
		&compiler.DisjunctionToType{},
	}
}

func (language *Language) NullableKinds() languages.NullableConfig {
	return languages.NullableConfig{
		Kinds:              []ast.Kind{ast.KindMap, ast.KindArray},
		ProtectArrayAppend: false,
		AnyIsNullable:      true,
	}
}
