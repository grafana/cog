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

				return veneers.Engine().ApplyTo(builders)
			},
		),
	)

	return targets
}
