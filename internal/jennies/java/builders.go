package java

import (
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/jennies/template"
	"github.com/grafana/cog/internal/tools"
)

type Builders struct {
	context       common.Context
	typeFormatter *typeFormatter
	builders      map[string]map[string]ast.Builder
}

func parseBuilders(context common.Context, formatter *typeFormatter) Builders {
	b := make(map[string]map[string]ast.Builder)
	for _, builder := range context.Builders {
		if _, ok := b[builder.Package]; !ok {
			b[builder.Package] = map[string]ast.Builder{}
		}
		b[builder.Package][builder.Name] = builder
	}

	return Builders{
		context:       context,
		builders:      b,
		typeFormatter: formatter,
	}
}

func (b Builders) genBuilder(pkg string, name string) (template.Builder, bool) {
	builder, ok := b.getBuilder(pkg, name)
	if !ok {
		return template.Builder{}, false
	}

	object, _ := b.context.LocateObject(builder.For.SelfRef.ReferredPkg, builder.For.SelfRef.ReferredType)
	return template.Builder{
		ObjectName:  tools.UpperCamelCase(object.Name),
		BuilderName: builder.Name,
		Constructor: builder.Constructor,
		Options:     tools.Map(builder.Options, b.genOption),
		Properties:  builder.Properties,
	}, true
}

func (b Builders) getBuilder(pkg string, name string) (ast.Builder, bool) {
	builderMap, ok := b.builders[pkg]
	if !ok {
		return ast.Builder{}, false
	}

	builder, ok := builderMap[name]
	return builder, ok
}

func (b Builders) genOption(opt ast.Option) ast.Option {
	opt.Name = tools.LowerCamelCase(opt.Name)
	opt.Args = tools.Map(opt.Args, func(arg ast.Argument) ast.Argument {
		arg.Name = tools.LowerCamelCase(arg.Name)
		return arg
	})

	return opt
}
