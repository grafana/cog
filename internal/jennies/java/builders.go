package java

import (
	"fmt"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/template"
	"github.com/grafana/cog/internal/languages"
	"github.com/grafana/cog/internal/tools"
)

type Builders struct {
	config           Config
	context          languages.Context
	typeFormatter    *typeFormatter
	builders         map[string]map[string]ast.Builder
	externalBuilders map[string]map[string]bool
}

func parseBuilders(config Config, context languages.Context, formatter *typeFormatter) Builders {
	if !config.generateBuilders || config.SkipRuntime {
		return Builders{
			builders:         make(map[string]map[string]ast.Builder),
			externalBuilders: make(map[string]map[string]bool),
		}
	}
	b := make(map[string]map[string]ast.Builder)
	externalBuilders := make(map[string]map[string]bool)
	for _, builder := range context.Builders {
		if _, ok := b[builder.Package]; !ok {
			b[builder.Package] = map[string]ast.Builder{}
		}

		// This should be better 🤮
		if builder.For.SelfRef.ReferredType == "VariableModel" || (builder.Name == "Panel" && builder.Package != "dashboard") {
			if _, ok := externalBuilders[builder.Package]; !ok {
				externalBuilders[builder.Package] = map[string]bool{}
			}
			b[builder.Package][builder.Name] = builder
			externalBuilders[builder.Package][builder.Name] = true
			continue
		}

		b[builder.Package][builder.For.SelfRef.ReferredType] = builder
	}

	return Builders{
		config:           config,
		context:          context,
		builders:         b,
		typeFormatter:    formatter,
		externalBuilders: externalBuilders,
	}
}

func (b Builders) genBuilder(pkg string, name string) (template.Builder, bool) {
	builder, ok := b.getBuilder(pkg, name)
	if !ok {
		return template.Builder{}, false
	}

	if builder.For.SelfRef.ReferredPkg != pkg {
		b.typeFormatter.packageMapper(builder.For.SelfRef.ReferredPkg, builder.For.SelfRef.ReferredType)
	}

	object, _ := b.context.LocateObject(builder.For.SelfRef.ReferredPkg, builder.For.SelfRef.ReferredType)
	return template.Builder{
		Package:     b.typeFormatter.formatPackage(pkg),
		ObjectName:  tools.UpperCamelCase(object.Name),
		BuilderName: builder.Name,
		Constructor: builder.Constructor,
		Options:     builder.Options,
		Properties:  builder.Properties,
		Defaults:    b.genDefaults(builder.Options),
	}, true
}

func (b Builders) genExternalBuilders(pkg string) (map[string]template.Builder, bool) {
	builders := b.externalBuilders[pkg]
	if len(builders) == 0 {
		return nil, false
	}

	tmpl := make(map[string]template.Builder, len(builders))
	for name := range builders {
		generatedTmpl, ok := b.genBuilder(pkg, name)
		if !ok {
			return nil, false
		}
		tmpl[name] = generatedTmpl
	}

	return tmpl, true
}

func (b Builders) getBuilder(pkg string, name string) (ast.Builder, bool) {
	builderMap, ok := b.builders[pkg]
	if !ok {
		return ast.Builder{}, false
	}

	builder, ok := builderMap[name]
	return builder, ok
}

func (b Builders) genDefaults(options []ast.Option) []template.OptionCall {
	calls := make([]template.OptionCall, 0)
	for _, opt := range options {
		if opt.Default == nil || len(opt.Args) == 0 {
			continue
		}

		calls = append(calls, template.OptionCall{
			OptionName: tools.UpperCamelCase(opt.Name),
			Args:       b.formatDefaultValues(opt.Args),
		})
	}

	return calls
}

func (b Builders) formatDefaultValues(args []ast.Argument) []string {
	argumentList := make([]string, 0, len(args))
	for _, arg := range args {
		switch arg.Type.Kind {
		case ast.KindRef:
			argumentList = append(argumentList, b.formatDefaultReference(arg.Type.AsRef(), arg.Type.Default))
		case ast.KindScalar:
			scalar := arg.Type.AsScalar()
			if scalar.ScalarKind == ast.KindFloat32 || scalar.ScalarKind == ast.KindFloat64 {
				val := arg.Type.Default
				if v, ok := val.(int64); ok {
					val = float64(v)
				} else {
					val = val.(float64)
				}
				argumentList = append(argumentList, fmt.Sprintf("%.1f", val))
			} else {
				argumentList = append(argumentList, fmt.Sprintf("%#v", arg.Type.Default))
			}
		case ast.KindArray:
			array := arg.Type.AsArray()
			if array.IsArrayOfScalars() {
				argumentList = append(argumentList, fmt.Sprintf("List.of(%s)", b.typeFormatter.formatScalar(arg.Type.Default)))
			}
		case ast.KindStruct:
			// TODO: Java is using veneers to avoid anonymous structs but it should be detailed if we need it at any moment.
			argumentList = append(argumentList, "new Object()")
		}
		// TODO: Implement the rest of types if any
	}

	return argumentList
}

func (b Builders) formatDefaultReference(ref ast.RefType, defValue any) string {
	object, _ := b.context.LocateObject(ref.ReferredPkg, ref.ReferredType)
	switch object.Type.Kind {
	case ast.KindEnum:
		for _, v := range object.Type.AsEnum().Values {
			if defValue == v.Value {
				return fmt.Sprintf("%s.%s", object.Name, tools.UpperSnakeCase(v.Name))
			}
		}
	case ast.KindStruct:
		// TODO: Builder could have arguments 🙃
		builder := fmt.Sprintf("new %s.Builder()", tools.UpperCamelCase(object.Name))
		structType := object.Type.AsStruct()
		defValues := defValue.(map[string]interface{})
		for _, field := range structType.Fields {
			if f, ok := defValues[field.Name]; ok {
				builder = fmt.Sprintf("%s.set%s(%#v)", builder, tools.UpperCamelCase(field.Name), f)
			}
		}
		return builder + ".build()"
	}

	return ""
}
