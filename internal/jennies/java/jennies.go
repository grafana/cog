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
	ProjectPath   string `yaml:"-"`
	PackagePath   string `yaml:"package_path"`
	SkipGradleDev bool   `yaml:"skip_gradle_dev"`

	// BuilderTemplatesDirectories holds a list of directories containing templates
	// to be used to override parts of builders.
	BuilderTemplatesDirectories []string `yaml:"builder_templates"`

	// SkipRuntime disables runtime-related code generation when enabled.
	// Note: builders can NOT be generated with this flag turned on, as they
	// rely on the runtime to function.
	SkipRuntime        bool `yaml:"skip_runtime"`
	generateBuilders   bool
	generateConverters bool
}

func (config *Config) InterpolateParameters(interpolator func(input string) string) {
	config.PackagePath = interpolator(config.PackagePath)
	config.BuilderTemplatesDirectories = tools.Map(config.BuilderTemplatesDirectories, interpolator)
	config.ProjectPath = fmt.Sprintf("src/main/java/%s", strings.ReplaceAll(config.PackagePath, ".", "/"))
}

func (config Config) MergeWithGlobal(global languages.Config) Config {
	newConfig := config
	newConfig.generateBuilders = global.Builders
	newConfig.generateConverters = global.Converters

	return newConfig
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

	tmpl := initTemplates(language.config.BuilderTemplatesDirectories)

	jenny := codejen.JennyListWithNamer[languages.Context](func(_ languages.Context) string {
		return LanguageRef
	})
	jenny.AppendOneToMany(
		common.If[languages.Context](!config.SkipRuntime, Runtime{config: config, tmpl: tmpl}),
		common.If[languages.Context](!config.SkipRuntime, Registry{config: config, tmpl: tmpl}),
		common.If[languages.Context](!config.SkipRuntime, &Deserializers{config: config, tmpl: tmpl}),
		common.If[languages.Context](!config.SkipRuntime, &Serializers{config: config, tmpl: tmpl}),
		RawTypes{config: config, tmpl: tmpl},
		common.If[languages.Context](!config.SkipGradleDev, Gradle{config: config, tmpl: tmpl}),
		common.If[languages.Context](!config.SkipRuntime && globalConfig.Builders && globalConfig.Converters, &Converter{config: config, tmpl: tmpl}),
	)
	jenny.AddPostprocessors(common.GeneratedCommentHeader(globalConfig))

	return jenny
}

func (language *Language) CompilerPasses() compiler.Passes {
	return compiler.Passes{
		&compiler.AnonymousEnumToExplicitType{},
		&compiler.AnonymousStructsToNamed{},
		&compiler.NotRequiredFieldAsNullableType{},
		&compiler.FlattenDisjunctions{},
		&compiler.DisjunctionWithNullToOptional{},
		&compiler.DisjunctionInferMapping{},
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
