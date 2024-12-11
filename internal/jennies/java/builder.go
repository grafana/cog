package java

import (
	"fmt"
	"path/filepath"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/jennies/template"
	"github.com/grafana/cog/internal/languages"
	"github.com/grafana/cog/internal/tools"
)

type Builder struct {
	config        Config
	tmpl          *template.Template
	imports       *common.DirectImportMap
	typeFormatter *typeFormatter
}

func (jenny Builder) JennyName() string {
	return "Builder"
}

func (jenny Builder) Generate(context languages.Context) (codejen.Files, error) {
	files := make(codejen.Files, 0)

	for _, builder := range context.Builders {
		output, err := jenny.genBuilder(context, builder)
		if err != nil {
			return nil, err
		}

		filename := filepath.Join(jenny.config.ProjectPath, formatPackageName(builder.Package), fmt.Sprintf("%sBuilder.java", tools.UpperCamelCase(builder.Name)))
		files = append(files, *codejen.NewFile(filename, output, jenny))
	}

	return files, nil
}

func (jenny Builder) genBuilder(context languages.Context, builder ast.Builder) ([]byte, error) {
	jenny.imports = NewImportMap(jenny.config.PackagePath)

	packageMapper := func(pkg string, class string) string {
		if jenny.imports.IsIdentical(pkg, builder.Package) {
			return ""
		}

		return jenny.imports.Add(class, pkg)
	}

	jenny.typeFormatter = createFormatter(context, jenny.config).withPackageMapper(packageMapper)

	object, _ := context.LocateObject(builder.For.SelfRef.ReferredPkg, builder.For.SelfRef.ReferredType)
	tmpl := BuilderTemplate{
		Package:              jenny.typeFormatter.formatPackage(builder.Package),
		Imports:              jenny.imports,
		ObjectName:           tools.UpperCamelCase(object.Name),
		BuilderName:          tools.UpperCamelCase(builder.Name),
		BuilderSignatureType: jenny.getBuilderSignature(builder.Package, object),
		Constructor:          builder.Constructor,
		Options:              builder.Options,
		Properties:           builder.Properties,
		Defaults:             jenny.genDefaults(context, builder.Options),
		ImportAlias:          jenny.config.PackagePath,
		IsGenericPanel:       jenny.isGenericPanel(builder),
	}

	return jenny.tmpl.Funcs(map[string]any{
		"formatBuilderFieldType":   jenny.typeFormatter.formatBuilderFieldType,
		"emptyValueForType":        jenny.typeFormatter.defaultValueFor,
		"shouldCastNilCheck":       jenny.typeFormatter.shouldCastNilCheck,
		"formatCastValue":          jenny.typeFormatter.formatCastValue,
		"typeHasBuilder":           jenny.typeFormatter.typeHasBuilder,
		"resolvesToComposableSlot": jenny.typeFormatter.resolvesToComposableSlot,
		"formatAssignmentPath":     jenny.typeFormatter.formatAssignmentPath,
		"formatPath":               jenny.typeFormatter.formatFieldPath,
		"formatRefType":            jenny.typeFormatter.formatRefType,
		"formatType":               jenny.typeFormatter.formatFieldType,
	}).RenderAsBytes("builders/builder.tmpl", tmpl)
}

func (jenny Builder) getBuilderSignature(pkg string, obj ast.Object) string {
	if pkg != obj.SelfRef.ReferredPkg {
		jenny.imports.Add(obj.SelfRef.ReferredType, obj.SelfRef.ReferredPkg)
	}

	if !obj.Type.IsDataqueryVariant() {
		return obj.Name
	}

	return fmt.Sprintf("%s.%s", jenny.typeFormatter.formatPackage("cog.variants"), tools.UpperCamelCase(obj.Type.ImplementedVariant()))
}

func (jenny Builder) isGenericPanel(builder ast.Builder) bool {
	if builder.Package != builder.For.SelfRef.ReferredPkg {
		return false
	}

	return builder.Name == "Panel"
}

func (jenny Builder) genDefaults(context languages.Context, options []ast.Option) []OptionCall {
	calls := make([]OptionCall, 0)
	for _, opt := range options {
		if opt.Default == nil || len(opt.Args) == 0 {
			continue
		}

		calls = append(calls, OptionCall{
			Initializers: jenny.formatInitializers(context, opt.Args),
			OptionName:   tools.UpperCamelCase(opt.Name),
			Args:         jenny.formatDefaultValues(context, opt.Args),
		})
	}

	return calls
}

// formatInitializers initialises objects with their defaults before set the value in the corresponding setter.
// TODO: It could have conflicts if we have different fields with the same kind of argument.
// TODO: It means that we need to initialize the objects with different names in that case.
func (jenny Builder) formatInitializers(context languages.Context, args []ast.Argument) []string {
	initializers := make([]string, 0)
	for _, arg := range args {
		if !arg.Type.IsRef() {
			return nil
		}

		ref := arg.Type.AsRef()
		object, _ := context.LocateObject(ref.ReferredPkg, ref.ReferredType)
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
		if jenny.typeFormatter.typeHasBuilder(arg.Type) {
			constructorFormat = "%sBuilder %sResource = new %sBuilder();"
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

func (jenny Builder) formatDefaultValues(context languages.Context, args []ast.Argument) []string {
	argumentList := make([]string, 0, len(args))
	for _, arg := range args {
		switch arg.Type.Kind {
		case ast.KindRef:
			argumentList = append(argumentList, jenny.formatDefaultReference(context, arg.Type.AsRef(), arg.Type.Default))
		case ast.KindScalar:
			scalar := arg.Type.AsScalar()
			argumentList = append(argumentList, formatType(scalar.ScalarKind, arg.Type.Default))
		case ast.KindArray:
			array := arg.Type.AsArray()
			if array.IsArrayOfScalars() {
				argumentList = append(argumentList, fmt.Sprintf("List.of(%s)", jenny.typeFormatter.formatScalar(arg.Type.Default)))
			}
		case ast.KindStruct:
			// TODO: Java is using veneers to avoid anonymous structs but it should be detailed if we need it at any moment.
			argumentList = append(argumentList, "new Object()")
		}
		// TODO: Implement the rest of types if any
	}

	return argumentList
}

func (jenny Builder) formatDefaultReference(context languages.Context, ref ast.RefType, defValue any) string {
	object, _ := context.LocateObject(ref.ReferredPkg, ref.ReferredType)
	switch object.Type.Kind {
	case ast.KindEnum:
		for _, v := range object.Type.AsEnum().Values {
			if defValue == v.Value {
				return fmt.Sprintf("%s.%s", object.Name, tools.UpperSnakeCase(v.Name))
			}
		}
	case ast.KindStruct:
		if jenny.typeFormatter.typeHasBuilder(object.Type) {
			// TODO: Builder could have arguments 🙃
			builder := fmt.Sprintf("new %sBuilder()", tools.UpperCamelCase(object.Name))
			structType := object.Type.AsStruct()
			defValues := defValue.(map[string]interface{})
			for _, field := range structType.Fields {
				if f, ok := defValues[field.Name]; ok {
					builder = fmt.Sprintf("%s.%s(%#v)", builder, tools.LowerCamelCase(field.Name), f)
				}
			}
			return builder + ".build()"
		}

		jenny.imports.Add(ref.ReferredType+"Builder", ref.ReferredPkg)

		return fmt.Sprintf("%sResource", tools.LowerCamelCase(object.Name))
	}

	return ""
}
