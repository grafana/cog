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
	debug            bool
	generateBuilders bool

	// GenerateGoMod indicates whether a go.mod file should be generated.
	// If enabled, PackageRoot is used as module path.
	GenerateGoMod bool `yaml:"go_mod"`

	// SkipRuntime disables runtime-related code generation when enabled.
	// Note: builders can NOT be generated with this flag turned on, as they
	// rely on the runtime to function.
	SkipRuntime bool `yaml:"skip_runtime"`

	// Root path for imports.
	// Ex: github.com/grafana/cog/generated
	PackageRoot string `yaml:"package_root"`
}

func (config *Config) InterpolateParameters(interpolator func(input string) string) {
	config.PackageRoot = interpolator(config.PackageRoot)
}

func (config Config) MergeWithGlobal(global common.Config) Config {
	newConfig := config
	newConfig.debug = global.Debug
	newConfig.generateBuilders = global.Builders

	return newConfig
}

func (config Config) importPath(suffix string) string {
	root := strings.TrimSuffix(config.PackageRoot, "/")
	return fmt.Sprintf("%s/%s", root, suffix)
}

type Language struct {
	config Config
}

func New(config Config) *Language {
	return &Language{
		config: config,
	}
}

func (language *Language) Name() string {
	return LanguageRef
}

func (language *Language) Jennies(globalConfig common.Config) *codejen.JennyList[common.Context] {
	config := language.config.MergeWithGlobal(globalConfig)

	jenny := codejen.JennyListWithNamer[common.Context](func(_ common.Context) string {
		return LanguageRef
	})
	jenny.AppendOneToMany(
		common.If[common.Context](!config.SkipRuntime, Runtime{Config: config}),
		common.If[common.Context](!config.SkipRuntime, VariantsPlugins{Config: config}),

		common.If[common.Context](config.GenerateGoMod, GoMod{Config: config}),

		common.If[common.Context](globalConfig.Types, RawTypes{Config: config}),

		common.If[common.Context](!config.SkipRuntime && globalConfig.Builders, &Builder{Config: config}),
	)
	jenny.AddPostprocessors(PostProcessFile, common.GeneratedCommentHeader(globalConfig))

	return jenny
}

func (language *Language) CompilerPasses() compiler.Passes {
	return compiler.Passes{
		&compiler.AnonymousEnumToExplicitType{},
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

func (language *Language) IdentifiersFormatter() *ast.IdentifierFormatter {
	return ast.NewIdentifierFormatter(
		ast.PackageFormatter(formatPackageName),
		ast.ClassFormatter(tools.UpperCamelCase),
		ast.ClassFieldFormatter(tools.UpperCamelCase),
		ast.EnumFormatter(tools.UpperCamelCase),
		ast.EnumMemberFormatter(func(s string) string {
			return tools.CleanupNames(tools.UpperCamelCase(s))
		}),
		ast.ConstantFormatter(tools.UpperCamelCase),
		ast.VariableFormatter(formatArgName),
	)
}
