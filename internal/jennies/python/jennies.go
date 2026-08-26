package python

import (
	"io/fs"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ir"
	compiler2 "github.com/grafana/cog/internal/ir/transforms"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/languages"
	"github.com/grafana/cog/internal/tools"
)

const LanguageRef = "python"

type Config struct {
	PathPrefix string `yaml:"path_prefix"`

	// GenerateEqual controls the generation of `__eq__()` methods on types.
	GenerateEqual bool `yaml:"generate_equal"`

	// GenerateJSONMarshaller controls the generation of `to_json()` and
	// `from_json()` methods on types.
	GenerateJSONMarshaller bool `yaml:"generate_json_marshaller"`

	// SkipRuntime disables runtime-related code generation when enabled.
	// Note: builders can NOT be generated with this flag turned on, as they
	// rely on the runtime to function.
	SkipRuntime bool `yaml:"skip_runtime"`

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
}

func (config *Config) InterpolateParameters(interpolator func(input string) string) {
	config.PathPrefix = interpolator(config.PathPrefix)
	config.OverridesTemplatesDirectories = tools.Map(config.OverridesTemplatesDirectories, interpolator)
	config.ExtraFilesTemplatesDirectories = tools.Map(config.ExtraFilesTemplatesDirectories, interpolator)
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

	extraTemplatesJenny := common.CustomTemplates{
		TemplateDirectories: language.config.ExtraFilesTemplatesDirectories,
		Data: map[string]any{
			"Debug": globalConfig.Debug,
		},
		ExtraData: language.config.ExtraFilesTemplatesData,
	}

	jenny := codejen.JennyListWithNamer(func(_ languages.Context) string {
		return LanguageRef
	})
	jenny.AppendOneToMany(
		ModuleInit{},
		common.If(!language.config.SkipRuntime, Runtime{tmpl: tmpl}),

		common.If(globalConfig.Types, RawTypes{config: language.config, tmpl: tmpl, apiRefCollector: language.apiRefCollector}),
		common.If(!language.config.SkipRuntime && globalConfig.Builders, &Builder{tmpl: tmpl, apiRefCollector: language.apiRefCollector}),

		common.If(globalConfig.APIReference, common.APIReference{
			Collector: language.apiRefCollector,
			Language:  LanguageRef,
			Formatter: apiReferenceFormatter(),
			Tmpl:      tmpl,
		}),

		extraTemplatesJenny,
	)
	jenny.AddPostprocessors(common.GeneratedCommentHeader(globalConfig))

	if language.config.PathPrefix != "" {
		jenny.AddPostprocessors(common.PathPrefixer(
			language.config.PathPrefix,
			common.PrefixExcept("docs/"),
			common.ExcludeCreatedByJenny(extraTemplatesJenny.JennyName()),
		))
	}

	return jenny
}

func (language *Language) CompilerPasses() compiler2.Transforms {
	return compiler2.Transforms{
		&compiler2.AnonymousStructsToNamed{},
		&compiler2.NotRequiredFieldAsNullableType{},
		&compiler2.DisjunctionWithNullToOptional{},
		&compiler2.DisjunctionOfConstantsToEnum{},
		&compiler2.FlattenDisjunctions{},
		&compiler2.DisjunctionInferMapping{},
		&compiler2.RenameNumericEnumValues{},
		&compiler2.DisjunctionPropagateVariant{},
	}
}

func (language *Language) NullableKinds() languages.NullableConfig {
	return languages.NullableConfig{
		Kinds:              []ir.Kind{ir.KindMap, ir.KindArray, ir.KindRef, ir.KindStruct},
		ProtectArrayAppend: true,
		AnyIsNullable:      true,
	}
}
