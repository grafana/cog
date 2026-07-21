package zod

import (
	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast/compiler"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/languages"
	"github.com/grafana/cog/internal/tools"
)

const LanguageRef = "zod"

type Config struct {
	// PathPrefix is the directory prefix for generated files. Defaults to "src".
	PathPrefix *string `yaml:"path_prefix"`

	// Filename overrides the per-package output filename. Defaults to
	// "schemas.gen.ts".
	Filename *string `yaml:"filename"`

	// FlatLayout drops the per-package directory and writes directly into
	// PathPrefix. Useful when the consumer manages its own layout.
	FlatLayout bool `yaml:"flat_layout"`

	// OptionalKindLiterals turns required fields whose type is a concrete
	// string literal (e.g. `kind: "Panel"`) into `.optional().default("X")`,
	// so callers don't have to repeat the discriminator. The literal still
	// constrains values that are supplied. Implies VerboseDefaults so the
	// generated file uses one form for every defaulted field.
	OptionalKindLiterals bool `yaml:"optional_kind_literals"`

	// ObjectMode controls how unknown keys are handled at parse time:
	//   - "strip" (default): unknown keys are silently dropped (Zod's default).
	//   - "strict":          unknown keys cause parse errors.
	//   - "loose":           unknown keys are kept on the parsed output.
	ObjectMode string `yaml:"object_mode"`

	// PreferUnknownForAny emits `z.unknown()` instead of `z.any()` for CUE's
	// top type (`_`), forcing consumers to narrow before use.
	PreferUnknownForAny bool `yaml:"prefer_unknown_for_any"`

	// VerboseDefaults emits defaulted fields as `.optional().default(X)`
	// rather than `.default(X).optional()`. The verbose form fixes the
	// inferred output type, which would otherwise include `undefined` even
	// though the default always fires. Implied by OptionalKindLiterals.
	VerboseDefaults bool `yaml:"verbose_defaults"`
}

func (config *Config) InterpolateParameters(interpolator func(input string) string) {
	if config.PathPrefix != nil {
		config.PathPrefix = tools.ToPtr(interpolator(*config.PathPrefix))
	}
	if config.Filename != nil {
		config.Filename = tools.ToPtr(interpolator(*config.Filename))
	}
}

func (config *Config) applyDefaults() {
	if config.PathPrefix == nil {
		config.PathPrefix = tools.ToPtr("src")
	}
	if config.Filename == nil {
		config.Filename = tools.ToPtr("schemas.gen.ts")
	}
	if config.ObjectMode == "" {
		config.ObjectMode = "strip"
	}
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
	language.config.applyDefaults()

	jenny := codejen.JennyListWithNamer(func(_ languages.Context) string {
		return LanguageRef
	})

	jenny.AppendOneToMany(
		common.If(globalConfig.Types, RawTypes{config: language.config}),
	)

	jenny.AddPostprocessors(common.GeneratedCommentHeader(globalConfig))

	return jenny
}

func (language *Language) CompilerPasses() compiler.Passes {
	// Order matters: each pass assumes the shape produced by the previous one.
	return compiler.Passes{
		// `T | null` → `T?` so we can attach .nullable() at the call site.
		&compiler.DisjunctionWithNullToOptional{},
		// `"a" | "b" | "c"` → enum, so we emit z.enum instead of a union of literals.
		&compiler.DisjunctionOfConstantsToEnum{},
		// Flatten before DisjunctionInferMapping, which needs flat disjunctions
		// to spot a shared discriminator.
		&compiler.FlattenDisjunctions{},
		// Sets Discriminator on unions of struct refs that share a scalar field,
		// enabling z.discriminatedUnion instead of z.union.
		&compiler.DisjunctionInferMapping{},
	}
}

func (language *Language) NullableKinds() languages.NullableConfig {
	return languages.NullableConfig{}
}
