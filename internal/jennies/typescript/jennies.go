package typescript

import (
	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/ast/compiler"
	"github.com/grafana/cog/internal/jennies/tools"
	"github.com/grafana/cog/internal/veneers"
	"github.com/grafana/cog/internal/veneers/builder"
	"github.com/grafana/cog/internal/veneers/option"
)

func Jennies() *codejen.JennyList[[]*ast.File] {
	targets := codejen.JennyListWithNamer[[]*ast.File](func(f []*ast.File) string {
		return "typescript"
	})
	targets.AppendOneToOne(
		OptionsBuilder{},
	)
	targets.AppendManyToMany(
		tools.Foreach[*ast.File](RawTypes{}),
	)
	targets.AppendOneToMany(
		codejen.AdaptOneToMany[[]ast.Builder, []*ast.File](
			&Builder{},
			func(files []*ast.File) []ast.Builder {
				generator := &ast.BuilderGenerator{}
				builders := generator.FromAST(files)

				// apply common veneers
				builders = veneers.Common().ApplyTo(builders)

				// apply TS-specific veneers
				return Veneers().ApplyTo(builders)
			},
		),
	)

	return targets
}

func CompilerPasses() []compiler.Pass {
	return nil
}

func Veneers() *veneers.Rewriter {
	return veneers.NewRewrite(
		[]builder.RewriteRule{},

		[]option.RewriteRule{},
	)
}
