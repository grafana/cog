package typescript

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

const LanguageRef = "typescript"

type Config struct {
	PathPrefix string `yaml:"path_prefix"`

	// SkipRuntime disables runtime-related code generation when enabled.
	// Note: builders can NOT be generated with this flag turned on, as they
	// rely on the runtime to function.
	SkipRuntime bool `yaml:"skip_runtime"`

	// SkipIndex disables the generation of `index.ts` files.
	SkipIndex bool `yaml:"skip_index"`

	// BuilderTemplatesDirectories holds a list of directories containing templates
	// to be used to override parts of builders.
	BuilderTemplatesDirectories []string `yaml:"builder_templates"`
}

func (config *Config) InterpolateParameters(interpolator func(input string) string) {
	config.BuilderTemplatesDirectories = tools.Map(config.BuilderTemplatesDirectories, interpolator)
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
	tmpl := initTemplates(language.config.BuilderTemplatesDirectories)

	jenny := codejen.JennyListWithNamer[languages.Context](func(_ languages.Context) string {
		return LanguageRef
	})
	jenny.AppendOneToMany(
		common.If[languages.Context](!language.config.SkipRuntime, Runtime{}),

		common.If[languages.Context](globalConfig.Types, RawTypes{}),
		common.If[languages.Context](!language.config.SkipRuntime && globalConfig.Builders, &Builder{tmpl: tmpl}),

		common.If[languages.Context](!language.config.SkipIndex, Index{Targets: globalConfig}),

		common.If[languages.Context](globalConfig.APIReference, common.APIReference{
			Language: "typescript",
			Formatter: common.APIReferenceFormatter{
				ObjectName: func(object ast.Object) string {
					return tools.CleanupNames(object.Name)
				},
				BuilderName: func(builder ast.Builder) string {
					return tools.UpperCamelCase(builder.Name) + "Builder"
				},
				ConstructorSignature: func(context languages.Context, builder ast.Builder) string {
					typesFormatter := builderTypeFormatter(context, func(pkg string) string {
						return pkg
					})
					args := tools.Map(builder.Constructor.Args, func(arg ast.Argument) string {
						return formatIdentifier(arg.Name) + ": " + typesFormatter.formatType(arg.Type)
					})

					return fmt.Sprintf("new %[1]s(%[2]s)", tools.UpperCamelCase(builder.Name)+"Builder", strings.Join(args, ", "))

				},
				OptionName: func(option ast.Option) string {
					return formatIdentifier(option.Name)
				},
				OptionSignature: func(context languages.Context, option ast.Option) string {
					typesFormatter := builderTypeFormatter(context, func(pkg string) string {
						return pkg
					})

					args := tools.Map(option.Args, func(arg ast.Argument) string {
						return formatIdentifier(arg.Name) + ": " + typesFormatter.formatType(arg.Type)
					})

					return fmt.Sprintf("%[1]s(%[2]s)", formatIdentifier(option.Name), strings.Join(args, ", "))
				},
			},
		}),
	)
	jenny.AddPostprocessors(common.GeneratedCommentHeader(globalConfig))

	return jenny
}

func (language *Language) CompilerPasses() compiler.Passes {
	return compiler.Passes{
		&compiler.RenameNumericEnumValues{},
	}
}

func (language *Language) NullableKinds() languages.NullableConfig {
	return languages.NullableConfig{
		Kinds:              []ast.Kind{ast.KindMap, ast.KindArray, ast.KindRef, ast.KindStruct},
		ProtectArrayAppend: true,
		AnyIsNullable:      true,
	}
}
