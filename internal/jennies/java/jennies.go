package java

import (
	"fmt"
	"strings"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/ast/compiler"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/languages"
)

const LanguageRef = "java"

type Config struct {
	ProjectPath   string `yaml:"-"`
	PackagePath   string `yaml:"package_path"`
	SkipGradleDev bool   `yaml:"skip_gradle_dev"`

	// SkipRuntime disables runtime-related code generation when enabled.
	// Note: builders can NOT be generated with this flag turned on, as they
	// rely on the runtime to function.
	SkipRuntime      bool `yaml:"skip_runtime"`
	generateBuilders bool
}

func (config *Config) InterpolateParameters(interpolator func(input string) string) {
	config.PackagePath = interpolator(config.PackagePath)
	config.ProjectPath = fmt.Sprintf("src/main/java/%s", strings.ReplaceAll(config.PackagePath, ".", "/"))
}

func (config Config) MergeWithGlobal(global languages.Config) Config {
	newConfig := config
	newConfig.generateBuilders = global.Builders

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

	jenny := codejen.JennyListWithNamer[languages.Context](func(_ languages.Context) string {
		return LanguageRef
	})
	jenny.AppendOneToMany(
		common.If[languages.Context](!config.SkipRuntime, Runtime{config: language.config}),
		common.If[languages.Context](!config.SkipRuntime, Registry{config: language.config}),
		common.If[languages.Context](!config.SkipRuntime, &Deserializers{config: language.config}),
		RawTypes{config: config},
		common.If[languages.Context](!config.SkipGradleDev, Gradle{config: config}),
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
