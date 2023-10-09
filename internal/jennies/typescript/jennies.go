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

type BuilderContext struct {
	Schemas  ast.Schemas
	Builders ast.Builders
}

func (context *BuilderContext) BuilderForType(t ast.Type) (ast.Builder, bool) {
	if t.Kind != ast.KindRef {
		return ast.Builder{}, false
	}

	ref := t.AsRef()
	return context.Builders.LocateByObject(ref.ReferredPkg, ref.ReferredType)
}

func (context *BuilderContext) LocateObject(pkg string, name string) (ast.Object, bool) {
	return context.Schemas.LocateObject(pkg, name)
}

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
		codejen.AdaptOneToMany[BuilderContext, []*ast.Schema](
			&Builder{},
			func(schemas []*ast.Schema) BuilderContext {
				var err error

				generator := &ast.BuilderGenerator{}
				builders := generator.FromAST(schemas)

				// apply common veneers
				builders, err = veneers.Common().ApplyTo(builders)
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

				return BuilderContext{
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
		&compiler.Unspec{},
	}
}

func Veneers() *veneers.Rewriter {
	return veneers.NewRewrite(
		[]builder.RewriteRule{},

		[]option.RewriteRule{},
	)
}
