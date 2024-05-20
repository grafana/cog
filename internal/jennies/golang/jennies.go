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
	debug bool

	// GenerateGoMod indicates whether a go.mod file should be generated.
	// If enabled, PackageRoot is used as module path.
	GenerateGoMod bool `yaml:"go_mod"`

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

func (language *Language) Jennies(globalConfig common.Config) *codejen.JennyList[common.Context] {
	config := language.config.MergeWithGlobal(globalConfig)
	identifiersFormatter := language.IdentifiersFormatter()

	jenny := codejen.JennyListWithNamer[common.Context](func(_ common.Context) string {
		return LanguageRef
	})
	jenny.AppendOneToMany(
		Runtime{Config: config},
		VariantsPlugins{Config: config},

		common.If[common.Context](config.GenerateGoMod, GoMod{Config: config}),

		common.If[common.Context](globalConfig.Types, RawTypes{Config: config, IdentifiersFormatter: identifiersFormatter}),
		common.If[common.Context](globalConfig.Types, JSONMarshalling{Config: config, IdentifiersFormatter: identifiersFormatter}),

		common.If[common.Context](globalConfig.Builders, &Builder{Config: config}),
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
		&compiler.DisjunctionToType{},
	}
}

func (language *Language) IdentifiersFormatter() *languages.IdentifierFormatter {
	return languages.NewIdentifierFormatter(
		languages.PackageFormatter(formatPackageName),
		languages.ObjectFormatter(tools.UpperCamelCase),
		languages.ObjectPublicFieldFormatter(tools.UpperCamelCase),
		languages.ObjectPrivateFieldFormatter(tools.LowerCamelCase),
		languages.ConstantFormatter(tools.UpperCamelCase),
		languages.OptionFormatter(tools.UpperCamelCase),
		languages.EnumFormatter(tools.UpperCamelCase),
		languages.EnumMemberFormatter(func(s string) string {
			return tools.StripNonAlphaNumeric(tools.UpperCamelCase(s))
		}),
		languages.VariableFormatter(func(name string) string {
			return escapeVarName(tools.LowerCamelCase(name))
		}),
	)
}

func (language *Language) NullableKinds() []ast.Kind {
	return []ast.Kind{
		ast.KindMap,
		ast.KindArray,
	}
}
