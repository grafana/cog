package java

import (
	"fmt"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/languages"
	"github.com/grafana/cog/internal/tools"
)

type Builders struct {
	config        Config
	context       languages.Context
	typeFormatter *typeFormatter
	builders      map[string]map[string]ast.Builders
	isPanel       map[string]bool
}

func parseBuilders(config Config, context languages.Context, formatter *typeFormatter) Builders {
	if !config.generateBuilders || config.SkipRuntime {
		return Builders{
			builders: make(map[string]map[string]ast.Builders),
			isPanel:  make(map[string]bool),
		}
	}
	b := make(map[string]map[string]ast.Builders)
	panels := make(map[string]bool)
	for _, builder := range context.Builders {
		if _, ok := b[builder.Package]; !ok {
			b[builder.Package] = map[string]ast.Builders{}
		}

		b[builder.Package][builder.For.SelfRef.ReferredType] = append(b[builder.Package][builder.For.SelfRef.ReferredType], builder)
		panels[builder.Package] = builder.Name == "Panel" && builder.Package != "dashboard" // TODO: Ugh! Maybe a compiler pass??
	}

	return Builders{
		config:        config,
		context:       context,
		builders:      b,
		typeFormatter: formatter,
		isPanel:       panels,
	}
}

func (b Builders) genBuilders(pkg string, name string) ([]Builder, bool) {
	builders := b.getBuilders(pkg, name)
	if len(builders) == 0 {
		return nil, false
	}

	return tools.Map(builders, func(builder ast.Builder) Builder {
		object, _ := b.context.LocateObject(builder.For.SelfRef.ReferredPkg, builder.For.SelfRef.ReferredType)
		return Builder{
			Package:              b.typeFormatter.formatPackage(pkg),
			ObjectName:           tools.UpperCamelCase(object.Name),
			BuilderName:          builder.Name,
			BuilderSignatureType: b.getBuilderSignature(builder, object),
			Constructor:          builder.Constructor,
			Options:              builder.Options,
			Properties:           builder.Properties,
			ImportAlias:          b.config.PackagePath,
		}
	}), true
}

func (b Builders) genPanelBuilder(pkg string) (Builder, bool) {
	if !b.isPanel[pkg] {
		return Builder{}, false
	}

	b.typeFormatter.packageMapper("dashboard", "Panel")
	builderTmpl, found := b.genBuilders(pkg, "Panel")
	if !found {
		return Builder{}, false
	}

	return builderTmpl[0], true
}

func (b Builders) getBuilders(pkg string, name string) ast.Builders {
	builderMap, ok := b.builders[pkg]
	if !ok {
		return nil
	}

	return builderMap[name]
}

func (b Builders) getBuilderSignature(builder ast.Builder, obj ast.Object) string {
	if builder.Name != obj.Type.ImplementedVariant() {
		return obj.Name
	}

	return fmt.Sprintf("%s.%s", b.typeFormatter.formatPackage("cog.variants"), tools.UpperCamelCase(obj.Name))
}
