package golang

import (
	"fmt"
	"strings"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast/compiler"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/jennies/context"
	"github.com/spf13/cobra"
)

const LanguageRef = "go"

type Config struct {
	// Root path for imports.
	// Ex: github.com/grafana/cog/generated
	PackageRoot string
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
}

func (language *Language) Jennies(targets common.Targets) *codejen.JennyList[context.Builders] {
	jenny := codejen.JennyListWithNamer[context.Builders](func(_ context.Builders) string {
		return LanguageRef
	})
	jenny.AppendOneToMany(
		Runtime{Config: language.config},
		VariantsPlugins{Config: language.config},

		common.If[context.Builders](targets.Types, RawTypes{Config: language.config}),
		common.If[context.Builders](targets.Types, JSONMarshalling{Config: language.config}),

		common.If[context.Builders](targets.Builders, &Builder{Config: language.config}),
	)
	jenny.AddPostprocessors(PostProcessFile)

	return jenny
}

func (language *Language) CompilerPasses() compiler.Passes {
	return compiler.Passes{
		&compiler.AnonymousEnumToExplicitType{},
		&compiler.PrefixEnumValues{},
		&compiler.NotRequiredFieldAsNullableType{},
		&compiler.DisjunctionToType{},
	}
}
