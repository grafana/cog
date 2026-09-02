package php

import (
	"io/fs"
	"log/slog"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/tools"
	"github.com/grafana/cog/pkg/apiref"
	"github.com/grafana/cog/pkg/ir"
	"github.com/grafana/cog/pkg/ir/transforms"
	"github.com/grafana/cog/pkg/languages"
	"github.com/grafana/cog/pkg/template"
)

const LanguageRef = "php"

type Config struct {
	debug bool

	converters bool

	NamespaceRoot string `yaml:"namespace_root"`

	// GenerateEqual controls the generation of `equals()` methods on types.
	GenerateEqual bool `yaml:"generate_equal"`

	// GenerateJSONMarshaller controls the generation of `fromArray()` and
	// `jsonSerialize()` methods on types.
	GenerateJSONMarshaller bool `yaml:"generate_json_marshaller"`

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

	// BuilderFactoriesClassMap allows to choose the name of the class that
	// will be generated to hold "builder factories".
	// By default, this class name is equal to the package name in which
	// factories are defined.
	// BuilderFactoriesClassMap associates these package names with a class
	// name.
	BuilderFactoriesClassMap map[string]string `yaml:"builder_factories_class_map"`
}

func (config *Config) InterpolateParameters(interpolator func(input string) string) {
	config.NamespaceRoot = interpolator(config.NamespaceRoot)
	config.OverridesTemplatesDirectories = tools.Map(config.OverridesTemplatesDirectories, interpolator)
	config.ExtraFilesTemplatesDirectories = tools.Map(config.ExtraFilesTemplatesDirectories, interpolator)
}

func (config Config) builderFactoryClassForPackage(pkg string) string {
	if config.BuilderFactoriesClassMap != nil && config.BuilderFactoriesClassMap[pkg] != "" {
		return config.BuilderFactoriesClassMap[pkg]
	}

	return pkg
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
	newConfig.converters = global.Converters

	return newConfig
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
	rawTypesJenny := RawTypes{config: config, tmpl: tmpl, apiRefCollector: language.apiRefCollector}

	jenny := codejen.JennyListWithNamer(func(_ languages.Context) string {
		return LanguageRef
	})
	jenny.AppendOneToMany(
		Runtime{config: config, tmpl: tmpl},
		common.If(globalConfig.Types, rawTypesJenny),
		common.If(globalConfig.Builders, &Builder{config: config, tmpl: tmpl, apiRefCollector: language.apiRefCollector}),
		common.If(globalConfig.Builders, &Factory{config: config, tmpl: tmpl, apiRefCollector: language.apiRefCollector}),
		common.If(globalConfig.Builders && globalConfig.Converters, &Converter{config: config, tmpl: tmpl, nullableConfig: language.NullableKinds()}),

		common.If(globalConfig.APIReference, apiref.APIReference{
			Collector: language.apiRefCollector,
			Language:  LanguageRef,
			Formatter: apiReferenceFormatter(tmpl, config),
			Tmpl:      tmpl,
		}),

		common.DynamicFiles{
			Tmpl: tmpl,
			Data: map[string]any{
				"Config": map[string]any{
					"Converters": config.converters,
				},
			},
			FuncsProvider: func(context languages.Context) template.FuncMap {
				return template.FuncMap{
					"unmarshalDisjunctionFunc": func(typeDef ir.Type) string {
						return rawTypesJenny.unmarshalDisjunctionFunc(context, typeDef.AsDisjunction())
					},
					"convertDisjunctionFunc": func(typeDef ir.Type) string {
						return rawTypesJenny.convertDisjunctionFunc(typeDef.AsDisjunction())
					},
				}
			},
		},

		common.CustomTemplates{
			TemplateDirectories: config.ExtraFilesTemplatesDirectories,
			Data: map[string]any{
				"Debug":         config.debug,
				"NamespaceRoot": config.NamespaceRoot,
			},
			ExtraData: config.ExtraFilesTemplatesData,
			TmplFuncs: formattingTemplateFuncs(),
		},
	)
	jenny.AddPostprocessors(common.GeneratedCommentHeader(globalConfig))

	return jenny
}

func (language *Language) Transform(schemas ir.Schemas) (ir.Schemas, error) {
	passes := transforms.Transforms{
		&transforms.AnonymousStructsToNamed{},
		&transforms.NotRequiredFieldAsNullableType{},
		&transforms.DisjunctionWithNullToOptional{},
		&transforms.DisjunctionOfConstantsToEnum{},
		&transforms.AnonymousEnumToExplicitType{},
		&transforms.SanitizeEnumMemberNames{},
		&transforms.FlattenDisjunctions{},
		&transforms.DisjunctionInferMapping{},
		&transforms.UndiscriminatedDisjunctionToAny{},
		&transforms.RemoveIntersections{},
		&transforms.DisjunctionPropagateVariant{},
		&transforms.InlineObjectsWithTypes{
			InlineTypes: []ir.Kind{ir.KindScalar, ir.KindArray, ir.KindMap, ir.KindDisjunction},
		},
	}

	return passes.Process(language.logger, schemas)
}

func (language *Language) NullableKinds() languages.NullableConfig {
	return languages.NullableConfig{
		ProtectArrayAppend: true,
		AnyIsNullable:      true,
	}
}
