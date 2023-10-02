package golang

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
	targets := codejen.JennyListWithNamer[[]*ast.File](func(files []*ast.File) string {
		return "golang"
	})
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

				// apply Go-specific veneers
				return Veneers().ApplyTo(builders)
			},
		),
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

func Veneers() *veneers.Rewriter {
	return veneers.NewRewrite(
		[]builder.RewriteRule{},

		[]option.RewriteRule{
			/********************************************
			 * Dashboards
			 ********************************************/

			// Time(from, to) instead of time(struct {From string `json:"from"`, To   string `json:"to"`}{From: "lala", To: "lala})
			option.StructFieldsAsArguments(
				option.ByName("Dashboard", "time"),
			),
		},
	)
}
