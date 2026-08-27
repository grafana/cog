package terraform

import (
	"io/fs"
	"log/slog"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/jennies/golang"
	"github.com/grafana/cog/pkg/ir"
	"github.com/grafana/cog/pkg/ir/transforms"
	"github.com/grafana/cog/pkg/languages"
)

const LanguageRef = "terraform"

type ValidatorsConfig struct {
	// Name of a validator function that verifies that exactly n of the
	// specified attributes are configured.
	// Expected signature: `func(n int, attributeNames ...string) validator.Object`
	AttributeCountExactly string

	// Name of a validator function that verifies that required attributes are
	// set only when a block is configured.
	// Expected signature: `func(names ...string) validator.Object`
	RequireAttrsWhenPresent string
}

type Config struct {
	debug bool

	// Root path for imports.
	// Ex: github.com/grafana/cog/generated
	PackageRoot string `yaml:"package_root"`

	// OverridesTemplatesDirectories holds a list of directories containing templates
	// defining blocks used to override parts of builders/types/....
	OverridesTemplatesDirectories []string `yaml:"overrides_templates"`
	// OverridesTemplatesFS holds an embedded filesystem containing templates.
	OverridesTemplatesFS fs.FS `yaml:"-"`
	// OverridesTemplateFuncs holds additional template functions to inject into override templates.
	OverridesTemplateFuncs map[string]any `yaml:"-"`

	// SkipPostFormatting disables formatting of Go files done with go imports
	// after code generation.
	SkipPostFormatting bool `yaml:"skip_post_formatting"`

	// SkipGeneratedHeader disables the addition of a
	// "Code generated - EDITING IS FUTILE. DO NOT EDIT." comment in generated
	// files headers.
	SkipGeneratedHeader bool `yaml:"skip_generated_header"`

	Validators ValidatorsConfig `yaml:"-"`
}

type Language struct {
	logger *slog.Logger
	config Config
}

func New(logger *slog.Logger, config Config) *Language {
	return &Language{
		logger: logger,
		config: config,
	}
}

func (config *Config) MergeWithGlobal(global languages.Config) Config {
	newConfig := config
	newConfig.debug = global.Debug

	return *newConfig
}

func (config *Config) InterpolateParameters(interpolator func(input string) string) {
	config.PackageRoot = interpolator(config.PackageRoot)
}

func (language *Language) Name() string {
	return LanguageRef
}

func (language *Language) Jennies(globalConfig languages.Config) *codejen.JennyList[languages.Context] {
	config := language.config.MergeWithGlobal(globalConfig)
	tmpl := initTemplates(config)

	jenny := codejen.JennyListWithNamer(func(_ languages.Context) string {
		return LanguageRef
	})

	jenny.AppendOneToMany(
		common.If(globalConfig.Types, RawTypes{config: config, tmpl: tmpl}))

	if !config.SkipGeneratedHeader {
		jenny.AddPostprocessors(common.GeneratedCommentHeader(globalConfig))
	}
	if !config.SkipPostFormatting {
		jenny.AddPostprocessors(golang.FormatGoFiles)
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
		&transforms.UndiscriminatedDisjunctionToAny{},
		&transforms.DisjunctionToType{},
	}

	return passes.Process(language.logger, schemas)
}

func (language *Language) NullableKinds() languages.NullableConfig {
	return languages.NullableConfig{}
}
