package typescript

import (
	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/ast/compiler"
	"github.com/grafana/cog/internal/jennies/context"
	"github.com/grafana/cog/internal/jennies/tools"
	"github.com/grafana/cog/internal/veneers/builder"
	"github.com/grafana/cog/internal/veneers/option"
	"github.com/grafana/cog/internal/veneers/rewrite"
)

func Jennies(veneers *rewrite.Rewriter) *codejen.JennyList[[]*ast.Schema] {
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
		codejen.AdaptOneToMany[context.Builders, []*ast.Schema](
			&Builder{},
			func(schemas []*ast.Schema) context.Builders {
				var err error

				generator := &ast.BuilderGenerator{}
				builders := generator.FromAST(schemas)

				// apply given veneers
				builders, err = veneers.ApplyTo(builders)
				if err != nil {
					// FIXME: codejen.AdaptOneToMany() doesn't let us return an error
					panic(err)
				}

				// apply TS-specific veneers
				builders, err = Veneers().ApplyTo(builders)
				if err != nil {
					// FIXME: codejen.AdaptOneToMany() doesn't let us return an error
					panic(err)
				}

				return context.Builders{
					Schemas:  schemas,
					Builders: builders,
				}
			},
		),
	)

	return targets
}

func CompilerPasses() []compiler.Pass {
	return []compiler.Pass{
		&compiler.PrefixEnumValues{},
	}
}

func Veneers() *rewrite.Rewriter {
	return rewrite.NewRewrite(
		[]builder.RewriteRule{},

		[]option.RewriteRule{},
	)
}
