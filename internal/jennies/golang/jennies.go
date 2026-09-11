package golang

import (
	"fmt"
	"io/fs"
	"log/slog"
	"strings"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/tools"
	"github.com/grafana/cog/pkg/apiref"
	"github.com/grafana/cog/pkg/ir"
	"github.com/grafana/cog/pkg/ir/transforms"
	"github.com/grafana/cog/pkg/jennies"
	"github.com/grafana/cog/pkg/languages"
)

const LanguageRef = "go"

type Config struct {
	debug              bool
	generateBuilders   bool
	GenerateConverters bool `yaml:"-"`

	// GenerateJSONMarshaller controls the generation of `MarshalJSON()` and
	// `UnmarshalJSON()` methods on types.
	GenerateJSONMarshaller bool `yaml:"generate_json_marshaller"`

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

	// SkipPostFormatting disables formatting of Go files done with go imports
	// after code generation.
	SkipPostFormatting bool `yaml:"skip_post_formatting"`

	// OverridesTemplatesDirectories holds a list of directories containing templates
	// defining blocks used to override parts of builders/types/....
	OverridesTemplatesDirectories []string `yaml:"overrides_templates"`
	// OverridesTemplatesFS holds an embedded filesystem containing templates
	OverridesTemplatesFS fs.FS `yaml:"-"`
	// OverridesTemplateFuncs holds additional template functions to be injected into the override templates.
	OverridesTemplateFuncs map[string]any `yaml:"-"`

	// ExtraFilesTemplatesDirectories holds a list of directories containing
	// templates describing files to be added to the generated output.
	ExtraFilesTemplatesDirectories []string `yaml:"extra_files_templates"`

	// ExtraFilesTemplatesData holds additional data to be injected into the
	// templates described in ExtraFilesTemplatesDirectories.
	ExtraFilesTemplatesData map[string]string `yaml:"-"`

	// Root path for imports.
	// Ex: github.com/grafana/cog/generated
	PackageRoot string `yaml:"package_root"`

	// AnyAsInterface instructs this jenny to emit `interface{}` instead of `any`.
	AnyAsInterface bool `yaml:"any_as_interface"`

	// GenerateUndiscriminatedDisjunctions controls whether disjunctions of refs
	// without a discriminator field are generated as a union struct type.
	// When false (default), such disjunctions are emitted as interface{}/any.
	GenerateUndiscriminatedDisjunctions bool `yaml:"generate_undiscriminated_disjunctions"`
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
	newConfig.GenerateConverters = global.Converters

	return newConfig
}

func (config Config) importPath(suffix string) string {
	root := strings.TrimSuffix(config.PackageRoot, "/")
	return fmt.Sprintf("%s/%s", root, suffix)
}

type Language struct {
	logger          *slog.Logger
	config          Config
	apiRefCollector *apiref.APIReferenceCollector
}

func New(logger *slog.Logger, config Config) *Language {
	return &Language{
		logger:          logger,
		config:          config,
		apiRefCollector: apiref.NewAPIReferenceCollector(),
	}
}

func (language *Language) Name() string {
	return LanguageRef
}

func (language *Language) Jennies(globalConfig languages.Config) *codejen.JennyList[languages.Context] {
	config := language.config.MergeWithGlobal(globalConfig)

	tmpl := initTemplates(config, language.apiRefCollector)

	jenny := codejen.JennyListWithNamer(func(_ languages.Context) string {
		return LanguageRef
	})
	jenny.AppendOneToMany(
		jennies.If(!config.SkipRuntime, Runtime{Config: config, Tmpl: tmpl}),

		jennies.If(globalConfig.Types, RawTypes{config: config, tmpl: tmpl, apiRefCollector: language.apiRefCollector}),

		jennies.If(!config.SkipRuntime && globalConfig.Builders, &Builder{Config: config, Tmpl: tmpl, apiRefCollector: language.apiRefCollector}),
		jennies.If(!config.SkipRuntime && globalConfig.Builders && globalConfig.Converters, &Converter{
			Config:          config,
			Tmpl:            tmpl,
			NullableConfig:  language.NullableKinds(),
			apiRefCollector: language.apiRefCollector,
		}),

		jennies.If(globalConfig.APIReference, apiref.APIReference{
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
			TmplFuncs: formattingTemplateFuncs(config),
		},
	)
	jenny.AddPostprocessors(jennies.GeneratedCommentHeader(globalConfig.Debug))
	if !config.SkipPostFormatting {
		jenny.AddPostprocessors(FormatGoFiles)
	}

	return jenny
}

func (language *Language) Transform(schemas ir.Schemas) (ir.Schemas, error) {
	passes := transforms.Transforms{
		&transforms.AnonymousStructsToNamed{},
		&transforms.NotRequiredFieldAsNullableType{},
		&transforms.DisjunctionWithNullToOptional{},
		&transforms.DisjunctionOfConstantsToEnum{},
		&transforms.AnonymousEnumToExplicitType{},
		&transforms.PrefixEnumValues{},
		&transforms.FlattenDisjunctions{},
		&transforms.DisjunctionOfAnonymousStructsToExplicit{},
		&transforms.DisjunctionInferMapping{},
		&transforms.UndiscriminatedDisjunctionToAny{
			GenerateUndiscriminatedDisjunctions: language.config.GenerateUndiscriminatedDisjunctions,
		},
		&transforms.DisjunctionToType{
			GenerateUndiscriminatedDisjunctions: language.config.GenerateUndiscriminatedDisjunctions,
		},
	}

	return passes.Process(language.logger, schemas)
}

func (language *Language) NullableKinds() languages.NullableConfig {
	return languages.NullableConfig{
		Kinds:              []ir.Kind{ir.KindMap, ir.KindArray},
		ProtectArrayAppend: false,
		AnyIsNullable:      true,
	}
}
