package terraform

import (
	"io/fs"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast/compiler"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/jennies/golang"
	"github.com/grafana/cog/internal/languages"
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
	config Config
}

func New(config Config) *Language {
	return &Language{
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

func (language *Language) CompilerPasses() compiler.Passes {
	return compiler.Passes{
		&compiler.AnonymousStructsToNamed{},
		&compiler.NotRequiredFieldAsNullableType{},
		&compiler.DisjunctionWithNullToOptional{},
		&compiler.DisjunctionOfConstantsToEnum{},
		&compiler.AnonymousEnumToExplicitType{},
		&compiler.PrefixEnumValues{},
		&compiler.FlattenDisjunctions{},
		&compiler.DisjunctionOfAnonymousStructsToExplicit{},
		&compiler.DisjunctionInferMapping{},
		&compiler.UndiscriminatedDisjunctionToAny{},
		&compiler.DisjunctionToType{},
	}
}

func (language *Language) NullableKinds() languages.NullableConfig {
	return languages.NullableConfig{}
}
