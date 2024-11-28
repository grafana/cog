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
		object, _ := b.context.LocateObjectByRef(builder.For.SelfRef)
		return Builder{
			Package:              b.typeFormatter.formatPackage(pkg),
			ObjectName:           tools.UpperCamelCase(object.Name),
			BuilderName:          builder.Name,
			BuilderSignatureType: b.getBuilderSignature(object),
			Constructor:          builder.Constructor,
			Options:              builder.Options,
			Properties:           builder.Properties,
			Defaults:             b.genDefaults(builder.Options),
			ImportAlias:          b.config.PackagePath,
			IsGeneric:            builder.For.SelfRef.ReferredPkg == "dashboard" && object.Name == "Panel",
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

func (b Builders) getBuilderSignature(obj ast.Object) string {
	if !obj.Type.IsDataqueryVariant() {
		return obj.Name
	}

	return fmt.Sprintf("%s.%s", b.typeFormatter.formatPackage("cog.variants"), tools.UpperCamelCase(obj.Type.ImplementedVariant()))
}

func (b Builders) genDefaults(options []ast.Option) []OptionCall {
	calls := make([]OptionCall, 0)
	for _, opt := range options {
		if opt.Default == nil || len(opt.Args) == 0 {
			continue
		}

		calls = append(calls, OptionCall{
			Initializers: b.formatInitializers(opt.Args),
			OptionName:   tools.UpperCamelCase(opt.Name),
			Args:         b.formatDefaultValues(opt.Args),
		})
	}

	return calls
}

// formatInitializers initialises objects with their defaults before set the value in the corresponding setter.
// TODO: It could have conflicts if we have different fields with the same kind of argument.
// TODO: It means that we need to initialize the objects with different names in that case.
func (b Builders) formatInitializers(args []ast.Argument) []string {
	initializers := make([]string, 0)
	for _, arg := range args {
		if !arg.Type.IsRef() {
			return nil
		}

		ref := arg.Type.AsRef()
		object, _ := b.context.LocateObjectByRef(ref)
		if !object.Type.IsStruct() {
			return nil
		}

		structType := object.Type.AsStruct()
		defValues := arg.Type.Default.(map[string]interface{})

		constructorFormat := "%s %sResource = new %s();"
		setterFormat := "%sResource.%s = %s;"
		fieldNameFunc := func(s string) string {
			return s
		}
		if b.typeFormatter.typeHasBuilder(arg.Type) {
			constructorFormat = "%s.Builder %sResource = new %s.Builder();"
			setterFormat = "%sResource.%s(%s);"
			fieldNameFunc = tools.LowerCamelCase
		}

		initializers = append(initializers, fmt.Sprintf(constructorFormat, ref.ReferredType, tools.LowerCamelCase(ref.ReferredType), ref.ReferredType))
		for _, field := range structType.Fields {
			if defVal, ok := defValues[field.Name]; ok {
				if field.Type.IsScalar() {
					initializers = append(initializers, fmt.Sprintf(setterFormat, tools.LowerCamelCase(ref.ReferredType), fieldNameFunc(field.Name), formatType(field.Type.AsScalar().ScalarKind, defVal)))
				}
				// TODO: Implement lists if needed
			}
		}
	}

	return initializers
}

func (b Builders) formatDefaultValues(args []ast.Argument) []string {
	argumentList := make([]string, 0, len(args))
	for _, arg := range args {
		switch arg.Type.Kind {
		case ast.KindRef:
			argumentList = append(argumentList, b.formatDefaultReference(arg.Type.AsRef(), arg.Type.Default))
		case ast.KindScalar:
			scalar := arg.Type.AsScalar()
			argumentList = append(argumentList, formatType(scalar.ScalarKind, arg.Type.Default))
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
		if b.typeFormatter.typeHasBuilder(object.Type) {
			// TODO: Builder could have arguments 🙃
			builder := fmt.Sprintf("new %s.Builder()", tools.UpperCamelCase(object.Name))
			structType := object.Type.AsStruct()
			defValues := defValue.(map[string]interface{})
			for _, field := range structType.Fields {
				if f, ok := defValues[field.Name]; ok {
					builder = fmt.Sprintf("%s.%s(%#v)", builder, tools.LowerCamelCase(field.Name), f)
				}
			}
			return builder + ".build()"
		}

		return fmt.Sprintf("%sResource", tools.LowerCamelCase(object.Name))
	}

	return ""
}
