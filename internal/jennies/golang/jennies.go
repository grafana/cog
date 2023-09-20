package golang

import (
	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/ast/compiler"
	"github.com/grafana/cog/internal/jennies/tools"
	"github.com/grafana/cog/internal/veneers"
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

				return veneers.Engine().ApplyTo(builders)
			},
		),
	)
	targets.AddPostprocessors(PostProcessFile)

	return targets
}

func CompilerPasses() []compiler.Pass {
	return []compiler.Pass{
		&compiler.AnonymousEnumToExplicitType{},
		&compiler.DisjunctionToType{},
	}
}
