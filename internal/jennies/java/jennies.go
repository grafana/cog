package java

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

const LanguageRef = "java"

type Config struct {
	ProjectPath string `yaml:"-"`
	PackagePath string `yaml:"package_path"`

	// OverridesTemplatesDirectories holds a list of directories containing templates
	// defining blocks used to override parts of builders/types/....
	OverridesTemplatesDirectories []string `yaml:"overrides_templates"`

	// ExtraFilesTemplatesDirectories holds a list of directories containing
	// templates describing files to be added to the generated output.
	ExtraFilesTemplatesDirectories []string `yaml:"extra_files_templates"`

	// ExtraFilesTemplatesData holds additional data to be injected into the
	// templates described in ExtraFilesTemplatesDirectories.
	ExtraFilesTemplatesData map[string]string `yaml:"-"`

	// SkipRuntime disables runtime-related code generation when enabled.
	// Note: builders can NOT be generated with this flag turned on, as they
	// rely on the runtime to function.
	SkipRuntime        bool `yaml:"skip_runtime"`
	GenerateBuilders   bool `yaml:"-"`
	GenerateConverters bool `yaml:"-"`
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
	config Config
}

func New(config Config) *Language {
	return &Language{config}
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
		common.If[languages.Context](!config.SkipRuntime, Runtime{config: config, tmpl: tmpl}),
		common.If[languages.Context](!config.SkipRuntime, &Deserializers{config: config, tmpl: tmpl}),
		common.If[languages.Context](!config.SkipRuntime, &Serializers{config: config, tmpl: tmpl}),
		RawTypes{config: config, tmpl: tmpl},
		common.If[languages.Context](config.GenerateBuilders, Builder{config: config, tmpl: tmpl}),
		common.If[languages.Context](!config.SkipRuntime && config.GenerateBuilders && config.GenerateConverters, &Converter{config: config, tmpl: tmpl}),

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
	jenny.AddPostprocessors(common.GeneratedCommentHeader(globalConfig))

	return jenny
}

func (language *Language) CompilerPasses() compiler.Passes {
	return compiler.Passes{
		&compiler.AnonymousStructsToNamed{},
		&compiler.NotRequiredFieldAsNullableType{},
		&compiler.DisjunctionWithNullToOptional{},
		&compiler.DisjunctionOfConstantsToEnum{},
		&compiler.AnonymousEnumToExplicitType{},
		&compiler.FlattenDisjunctions{},
		&compiler.DisjunctionInferMapping{},
		&compiler.UndiscriminatedDisjunctionToAny{},
		&compiler.DisjunctionToType{},
		&compiler.RemoveIntersections{},
	}
}

func (language *Language) NullableKinds() languages.NullableConfig {
	return languages.NullableConfig{
		Kinds:              []ast.Kind{ast.KindMap, ast.KindArray, ast.KindRef, ast.KindStruct},
		ProtectArrayAppend: true,
		AnyIsNullable:      true,
	}
}
