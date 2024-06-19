package terraform

import (
	"fmt"
	"strings"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast/compiler"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/languages"
)

const LanguageRef = "terraform"

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

func (config Config) MergeWithGlobal(global languages.Config) Config {
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

func (language *Language) Name() string {
	return LanguageRef
}

func (language *Language) Jennies(globalConfig languages.Config) *codejen.JennyList[languages.Context] {
	config := language.config.MergeWithGlobal(globalConfig)

	jenny := codejen.JennyListWithNamer[languages.Context](func(_ languages.Context) string {
		return LanguageRef
	})
	jenny.AppendOneToMany(
		common.If[languages.Context](config.GenerateGoMod, GoMod{Config: config}),

		common.If[languages.Context](globalConfig.Types, RawTypes{Config: config}),
	)
	//@TODO Re-enable this when ready
	//jenny.AddPostprocessors(PostProcessFile, common.GeneratedCommentHeader(globalConfig))

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
