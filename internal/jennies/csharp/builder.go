package csharp

import (
	"fmt"
	"path/filepath"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/jennies/template"
	"github.com/grafana/cog/internal/languages"
)

// Builder emits one fluent builder class per ast.Builder in the input
// context. Each builder lives next to its target type at
// <ProjectPath>/<Pkg>/<Name>Builder.cs and implements
// `Cog.IBuilder<TargetType>`. The implementation closely mirrors the
// Java jenny so the generated APIs stay aligned across targets.
type Builder struct {
	config        Config
	tmpl          *template.Template
	imports       *importMap
	typeFormatter *typeFormatter

	apiRefCollector *common.APIReferenceCollector
}

func (jenny Builder) JennyName() string {
	return "CSharpBuilder"
}

func (jenny Builder) Generate(context languages.Context) (codejen.Files, error) {
	files := make(codejen.Files, 0, len(context.Builders))
	for _, builder := range context.Builders {
		out, err := jenny.genBuilder(context, builder)
		if err != nil {
			return nil, err
		}
		filename := filepath.Join(
			jenny.config.ProjectPath,
			formatPackageName(builder.Package),
			fmt.Sprintf("%sBuilder.cs", jenny.builderName(builder)),
		)
		files = append(files, *codejen.NewFile(filename, out, jenny))
	}
	return files, nil
}

// builderTemplate carries the data passed into the builder template.
type builderTemplate struct {
	Namespace            string
	Imports              fmt.Stringer
	ObjectName           string
	BuilderName          string
	BuilderSignatureType string
	Constructor          ast.Constructor
	Options              []ast.Option
	Properties           []ast.StructField
}

func (jenny Builder) genBuilder(context languages.Context, builder ast.Builder) ([]byte, error) {
	namespace := jenny.config.formatNamespace(builder.Package)

	jenny.imports = newImportMap(jenny.config.namespaceRoot(), namespace)
	jenny.typeFormatter = newTypeFormatter(context, jenny.config, jenny.imports)

	object, _ := context.LocateObject(builder.For.SelfRef.ReferredPkg, builder.For.SelfRef.ReferredType)

	data := builderTemplate{
		Namespace:            namespace,
		Imports:              jenny.imports,
		ObjectName:           formatObjectName(object.Name),
		BuilderName:          jenny.builderName(builder),
		BuilderSignatureType: jenny.builderSignature(builder.Package, object),
		Constructor:          builder.Constructor,
		Options:              builder.Options,
		Properties:           builder.Properties,
	}

	jenny.apiRefCollector.BuilderMethod(builder, common.MethodReference{
		Name:     "Build",
		Comments: []string{"Builds the object."},
		Return:   formatObjectName(builder.Name),
	})

	// Cog.IBuilder<T> resolves via parent-namespace lookup
	// (`<NamespaceRoot>.Cog` is a sibling of every generated package),
	// so no explicit `using Grafana.Foundation.Cog;` is needed.

	return jenny.tmpl.Funcs(map[string]any{
		"formatBuilderFieldType":   jenny.typeFormatter.formatBuilderFieldType,
		"emptyValueForType":        func(def ast.Type) string { return jenny.typeFormatter.emptyValueForTypeOpts(def, true) },
		"typeHasBuilder":           jenny.typeFormatter.typeHasBuilder,
		"resolvesToComposableSlot": jenny.typeFormatter.resolvesToComposableSlot,
		"formatAssignmentPath":     jenny.typeFormatter.formatAssignmentPath,
		"formatPath":               jenny.typeFormatter.formatFieldPath,
		"formatRefType":            jenny.typeFormatter.formatRefType,
		"formatType":               jenny.typeFormatter.formatFieldType,
		"formatPathIndex":          jenny.typeFormatter.formatPathIndex,
	}).RenderAsBytes("builders/builder.tmpl", data)
}

// builderName picks the C# class name for a builder. When the builder
// targets a type from a different package we prefix it with the source
// package name to avoid clashes (mirrors Java's foreign-builder
// naming).
func (jenny Builder) builderName(builder ast.Builder) string {
	if builder.For.SelfRef.ReferredPkg != builder.Package {
		return formatObjectName(builder.Package) + formatObjectName(builder.Name)
	}
	return formatObjectName(builder.Name)
}

// builderSignature picks the type that ends up in
// `Cog.IBuilder<…>`. For composable-slot variants this becomes the
// variant interface (e.g. `Cog.Variants.Dataquery`); otherwise it is
// the target object's name. Imports are added as a side effect.
func (jenny Builder) builderSignature(pkg string, obj ast.Object) string {
	if obj.Type.IsDataqueryVariant() {
		return fmt.Sprintf("Cog.Variants.%s", formatObjectName(obj.Type.ImplementedVariant()))
	}
	if pkg != obj.SelfRef.ReferredPkg {
		jenny.imports.addPackage(obj.SelfRef.ReferredPkg)
	}
	return formatObjectName(obj.Name)
}
