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

func Jennies() *codejen.JennyList[[]*ast.Schema] {
	targets := codejen.JennyListWithNamer[[]*ast.Schema](func(_ []*ast.Schema) string {
		return "typescript"
	})
	targets.AppendOneToOne(
		OptionsBuilder{},
	)
	targets.AppendManyToMany(
		tools.Foreach[*ast.Schema](RawTypes{}),
	)
	targets.AppendOneToMany(
		codejen.AdaptOneToMany[[]ast.Builder, []*ast.Schema](
			&Builder{},
			func(schemas []*ast.Schema) []ast.Builder {
				generator := &ast.BuilderGenerator{}
				builders := generator.FromAST(schemas)

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
