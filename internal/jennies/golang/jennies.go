package golang

import (
	"fmt"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/template"
	"strings"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast/compiler"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/spf13/cobra"
)

const LanguageRef = "go"

type Config struct {
	Debug bool

	// GenerateGoMod indicates whether a go.mod file should be generated.
	// If enabled, PackageRoot is used as module path.
	GenerateGoMod bool

	// Root path for imports.
	// Ex: github.com/grafana/cog/generated
	PackageRoot string
}

func (config Config) MergeWithGlobal(global common.Config) Config {
	newConfig := config
	newConfig.Debug = global.Debug

	return newConfig
}

func (config Config) importPath(suffix string) string {
	root := strings.TrimSuffix(config.PackageRoot, "/")
	return fmt.Sprintf("%s/%s", root, suffix)
}

type Language struct {
	config Config
}

func New() *Language {
	return &Language{
		config: Config{},
	}
}

func (language *Language) RegisterCliFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&language.config.PackageRoot, "go-package-root", "github.com/grafana/cog/generated", "Go package root.")
	cmd.Flags().BoolVar(&language.config.GenerateGoMod, "go-mod", false, "Generate a go.mod file. If enabled, 'go-package-root' is used as module path.")
}

func (language *Language) Jennies(globalConfig common.Config) *codejen.JennyList[common.Context] {
	config := language.config.MergeWithGlobal(globalConfig)

	im := NewImportMap()
	tf := NewFormatter(config.PackageRoot, im)

	cfg := template.Config{
		Debug:         globalConfig.Debug,
		FileExtension: "go",
		TemplateConfig: template.FunctionsConfig{
			Name: "golang",
			FuncMap: map[string]any{
				"maybeAsPointer": func(intoNullable bool, variableName string) string {
					if intoNullable {
						return "&" + variableName
					}
					return variableName
				},
				"isNullableNonArray": func(typeDef ast.Type) bool {
					return typeDef.Nullable && typeDef.Kind != ast.KindArray
				},
			},
		},
		ImportMapper:    im,
		BuilderTemplate: GoTemplate{formatter: tf},
		Formatter:       tf,
	}

	jenny := codejen.JennyListWithNamer[common.Context](func(_ common.Context) string {
		return LanguageRef
	})

	jenny.AppendOneToMany(
		Runtime{Config: config},
		VariantsPlugins{Config: config},

		common.If[common.Context](language.config.GenerateGoMod, GoMod{Config: config}),

		common.If[common.Context](globalConfig.Types, RawTypes{Config: config}),
		common.If[common.Context](globalConfig.Types, JSONMarshalling{Config: config}),

		common.If[common.Context](globalConfig.Builders, template.NewExecutor(cfg)),
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
		&compiler.DisjunctionToType{},
	}
}
