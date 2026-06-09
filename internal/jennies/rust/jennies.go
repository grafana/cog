// Package rust implements the Rust code generation target for cog.
package rust

import (
	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/ast/compiler"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/languages"
)

const LanguageRef = "rust"

var _ languages.Language = (*Language)(nil)

type Config struct {
	PathPrefix string `yaml:"path_prefix"`

	// SkipRuntime disables runtime-related code generation when enabled.
	SkipRuntime bool `yaml:"skip_runtime"`

	// CrateName is the package name written to the generated Cargo.toml. It
	// defaults to defaultCrateName when left empty.
	CrateName string `yaml:"crate_name"`

	// CrateVersion is the package version written to the generated Cargo.toml.
	// It defaults to defaultCrateVersion when left empty.
	CrateVersion string `yaml:"crate_version"`

	// GenerateCargoToml controls whether a Cargo.toml manifest is emitted. The
	// lib.rs and mod.rs scaffolding is always emitted, but the manifest is gated
	// because consumers frequently vendor the generated modules into an existing
	// crate that already owns its Cargo.toml. This mirrors the golang target,
	// which gates go.mod generation rather than always overwriting it.
	GenerateCargoToml bool `yaml:"generate_cargo_toml"`
}

const (
	defaultCrateName    = "grafana_foundation_sdk"
	defaultCrateVersion = "0.0.1"
)

func (config *Config) InterpolateParameters(interpolator func(input string) string) {
	config.PathPrefix = interpolator(config.PathPrefix)
	config.CrateName = interpolator(config.CrateName)
	config.CrateVersion = interpolator(config.CrateVersion)
}

// crateName returns the configured crate name or the default.
func (config Config) crateName() string {
	if config.CrateName != "" {
		return config.CrateName
	}
	return defaultCrateName
}

// crateVersion returns the configured crate version or the default.
func (config Config) crateVersion() string {
	if config.CrateVersion != "" {
		return config.CrateVersion
	}
	return defaultCrateVersion
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

	jenny := codejen.JennyListWithNamer(func(_ languages.Context) string {
		return LanguageRef
	})
	jenny.AppendOneToMany(
		common.If(!language.config.SkipRuntime, Runtime{tmpl: tmpl}),
		common.If(!language.config.SkipRuntime, Plugins{config: language.config}),
		RawTypes{config: language.config, apiRefCollector: language.apiRefCollector},
		common.If(!language.config.SkipRuntime && globalConfig.Builders, Builder{config: language.config, apiRefCollector: language.apiRefCollector}),
		ModuleInit{config: language.config},
	)
	jenny.AddPostprocessors(common.GeneratedCommentHeader(globalConfig))
	jenny.AddPostprocessors(FormatRustFiles)

	if language.config.PathPrefix != "" {
		jenny.AddPostprocessors(common.PathPrefixer(language.config.PathPrefix))
	}

	return jenny
}

func (language *Language) CompilerPasses() compiler.Passes {
	return compiler.Passes{
		&compiler.AnonymousStructsToNamed{},
		&compiler.AnonymousEnumToExplicitType{},
		&compiler.PrefixEnumValues{},
		&compiler.NotRequiredFieldAsNullableType{},
		&compiler.DisjunctionWithNullToOptional{},
		&compiler.DisjunctionOfConstantsToEnum{},
		&compiler.FlattenDisjunctions{},
		&compiler.DisjunctionInferMapping{},
		&compiler.RenameNumericEnumValues{},
		&compiler.DisjunctionPropagateVariant{},
	}
}

func (language *Language) NullableKinds() languages.NullableConfig {
	return languages.NullableConfig{
		Kinds:              []ast.Kind{ast.KindMap, ast.KindArray, ast.KindRef, ast.KindStruct},
		ProtectArrayAppend: true,
		AnyIsNullable:      true,
	}
}
