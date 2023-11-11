package golang

import (
	"fmt"
	"strings"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast/compiler"
	"github.com/grafana/cog/internal/jennies/context"
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

func Jennies(config Config) *codejen.JennyList[context.Builders] {
	targets := codejen.JennyListWithNamer[context.Builders](func(_ context.Builders) string {
		return LanguageRef
	})
	targets.AppendOneToMany(
		Runtime{Config: config},
		VariantsPlugins{Config: config},

		RawTypes{Config: config},
		JSONMarshalling{Config: config},
		&Builder{Config: config},
	)
	targets.AddPostprocessors(PostProcessFile)

	return targets
}

func CompilerPasses() []compiler.Pass {
	return []compiler.Pass{
		&compiler.AnonymousEnumToExplicitType{},
		&compiler.PrefixEnumValues{},
		&compiler.NotRequiredFieldAsNullableType{},
		&compiler.DisjunctionToType{},
	}
}
