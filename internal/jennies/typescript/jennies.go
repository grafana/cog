package typescript

import (
	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/tools"
	"github.com/grafana/cog/internal/veneers"
)

func Jennies() *codejen.JennyList[[]*ast.File] {
	targets := codejen.JennyListWithNamer[[]*ast.File](func(f []*ast.File) string {
		return "typescript"
	})
	targets.AppendOneToOne(
		TypescriptOptionsBuilder{},
	)
	targets.AppendManyToMany(
		tools.Foreach[*ast.File](TypescriptRawTypes{}),
	)
	targets.AppendOneToMany(
		codejen.AdaptOneToMany[[]ast.Builder, []*ast.File](
			&TypescriptBuilder{},
			func(files []*ast.File) []ast.Builder {
				generator := &ast.BuilderGenerator{}
				builders := generator.FromAST(files)

				return veneers.Engine().ApplyTo(builders)
			},
		),
	)

	return targets
}
