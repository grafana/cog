package java

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

const LanguageRef = "java"

type Config struct {
	ProjectPath string `yaml:"-"`
	PackagePath string `yaml:"package_path"`

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

	// SkipRuntime disables runtime-related code generation when enabled.
	// Note: builders can NOT be generated with this flag turned on, as they
	// rely on the runtime to function.
	SkipRuntime bool `yaml:"skip_runtime"`

	// GenerateEqual controls the generation of `equals()` and `hashCode()` methods on types.
	GenerateEqual bool `yaml:"generate_equal"`

	// GenerateJSONMarshaller controls the generation of `MarshalJSON()` and
	// `UnmarshalJSON()` methods on types.
	GenerateJSONMarshaller bool `yaml:"generate_json_marshaller"`

	GenerateBuilders   bool `yaml:"-"`
	GenerateConverters bool `yaml:"-"`

	// BuilderFactoriesClassMap allows to choose the name of the class that
	// will be generated to hold "builder factories".
	// By default, this class name is equal to the package name in which
	// factories are defined.
	// BuilderFactoriesClassMap associates these package names with a class
	// name.
	BuilderFactoriesClassMap map[string]string `yaml:"builder_factories_class_map"`
}

func (config *Config) builderFactoryClassForPackage(pkg string) string {
	if config.BuilderFactoriesClassMap != nil && config.BuilderFactoriesClassMap[pkg] != "" {
		return config.BuilderFactoriesClassMap[pkg]
	}

	return pkg
}

func (config *Config) formatPackage(pkg string) string {
	if config.PackagePath != "" {
		return fmt.Sprintf("%s.%s", config.PackagePath, pkg)
	}

	return pkg
}

func (config *Config) InterpolateParameters(interpolator func(input string) string) {
	config.PackagePath = interpolator(config.PackagePath)
	config.OverridesTemplatesDirectories = tools.Map(config.OverridesTemplatesDirectories, interpolator)
	config.ExtraFilesTemplatesDirectories = tools.Map(config.ExtraFilesTemplatesDirectories, interpolator)
	config.ProjectPath = fmt.Sprintf("src/main/java/%s", strings.ReplaceAll(config.PackagePath, ".", "/"))
}

func (config *Config) MergeWithGlobal(global languages.Config) Config {
	newConfig := config
	newConfig.GenerateBuilders = global.Builders
	// newConfig.GenerateConverters = global.Converters
	newConfig.GenerateConverters = false

	return *newConfig
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

	apiRef := APIRef{config: config, tmpl: tmpl}

	jenny := codejen.JennyListWithNamer(func(_ languages.Context) string {
		return LanguageRef
	})
	jenny.AppendOneToMany(
		jennies.If(!config.SkipRuntime, Runtime{config: config, tmpl: tmpl}),
		jennies.If(!config.SkipRuntime && config.GenerateJSONMarshaller, &Deserializers{config: config, tmpl: tmpl}),
		jennies.If(!config.SkipRuntime && config.GenerateJSONMarshaller, &Serializers{config: config, tmpl: tmpl}),
		RawTypes{config: config, tmpl: tmpl},
		jennies.If(config.GenerateBuilders, Builder{config: config, tmpl: tmpl, apiRefCollector: language.apiRefCollector}),
		jennies.If(globalConfig.Builders, &Factory{config: config, tmpl: tmpl}),
		jennies.If(!config.SkipRuntime && config.GenerateBuilders && config.GenerateConverters, &Converter{config: config, tmpl: tmpl}),

		jennies.If(globalConfig.APIReference, apiref.APIReference{
			Collector: language.apiRefCollector,
			Language:  LanguageRef,
			Formatter: apiRef.apiReferenceFormatter(),
			Tmpl:      tmpl,
		}),

		common.CustomTemplates{
			TmplFuncs:           formattingTemplateFuncs(),
			TemplateDirectories: config.ExtraFilesTemplatesDirectories,
			Data: map[string]any{
				"Debug":  globalConfig.Debug,
				"Config": config,
			},
			ExtraData: config.ExtraFilesTemplatesData,
		},
	)
	jenny.AddPostprocessors(jennies.GeneratedCommentHeader(globalConfig.Debug))

	return jenny
}

func (language *Language) Transform(schemas ir.Schemas) (ir.Schemas, error) {
	passes := transforms.Transforms{
		&transforms.AnonymousStructsToNamed{},
		&transforms.NotRequiredFieldAsNullableType{},
		&transforms.DisjunctionWithNullToOptional{},
		&transforms.DisjunctionOfConstantsToEnum{},
		&transforms.AnonymousEnumToExplicitType{},
		&transforms.FlattenDisjunctions{},
		&transforms.DisjunctionInferMapping{},
		&transforms.UndiscriminatedDisjunctionToAny{},
		&transforms.DisjunctionToType{},
		&transforms.RemoveIntersections{},
		&transforms.InlineObjectsWithTypes{InlineTypes: []ir.Kind{ir.KindScalar, ir.KindMap, ir.KindArray}},
	}

	return passes.Process(language.logger, schemas)
}

func (language *Language) NullableKinds() languages.NullableConfig {
	return languages.NullableConfig{
		Kinds:              []ir.Kind{ir.KindMap, ir.KindArray, ir.KindRef, ir.KindStruct},
		ProtectArrayAppend: true,
		AnyIsNullable:      true,
	}
}
