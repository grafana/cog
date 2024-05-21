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
	isPanel       map[string]bool
}

func parseBuilders(context common.Context, formatter *typeFormatter) Builders {
	b := make(map[string]map[string]ast.Builder)
	panels := make(map[string]bool)
	for _, builder := range context.Builders {
		if _, ok := b[builder.Package]; !ok {
			b[builder.Package] = map[string]ast.Builder{}
		}
		b[builder.Package][builder.Name] = builder
		panels[builder.Package] = builder.Name == "Panel" && builder.Package != "dashboard" // TODO: Ugh! Maybe a compiler pass??
	}

	return Builders{
		context:       context,
		builders:      b,
		typeFormatter: formatter,
		isPanel:       panels,
	}
}

func (b Builders) genBuilder(pkg string, name string) (template.Builder, bool) {
	builder, ok := b.getBuilder(pkg, name)
	if !ok {
		return template.Builder{}, false
	}

	object, _ := b.context.LocateObject(builder.For.SelfRef.ReferredPkg, builder.For.SelfRef.ReferredType)
	return template.Builder{
		Package:     builder.Package,
		ObjectName:  tools.UpperCamelCase(object.Name),
		BuilderName: builder.Name,
		Constructor: builder.Constructor,
		Options:     builder.Options,
		Properties:  builder.Properties,
	}, true
}

func (b Builders) genPanelBuilder(pkg string) (template.Builder, bool) {
	if !b.isPanel[pkg] {
		return template.Builder{}, false
	}

	return b.genBuilder(pkg, "Panel")
}

func (b Builders) getBuilder(pkg string, name string) (ast.Builder, bool) {
	builderMap, ok := b.builders[pkg]
	if !ok {
		return ast.Builder{}, false
	}

	builder, ok := builderMap[name]
	return builder, ok
}
