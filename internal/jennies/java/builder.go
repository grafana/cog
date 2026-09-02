package java

import (
	"fmt"
	"path/filepath"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/tools"
	"github.com/grafana/cog/pkg/apiref"
	"github.com/grafana/cog/pkg/ir"
	"github.com/grafana/cog/pkg/languages"
	"github.com/grafana/cog/pkg/template"
)

type Builder struct {
	config        Config
	tmpl          *template.Template
	imports       *common.DirectImportMap
	typeFormatter *typeFormatter

	apiRefCollector *apiref.APIReferenceCollector
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

		filename := filepath.Join(jenny.config.ProjectPath, formatPackageName(builder.Package), fmt.Sprintf("%sBuilder.java", jenny.getBuilderName(builder)))
		files = append(files, *codejen.NewFile(filename, output, jenny))
	}

	return files, nil
}

func (jenny Builder) genBuilder(context languages.Context, builder ir.Builder) ([]byte, error) {
	jenny.imports = NewImportMap(jenny.config.PackagePath)

	packageMapper := func(pkg string, class string) string {
		if jenny.imports.IsIdentical(pkg, builder.Package) {
			return ""
		}

		return jenny.imports.Add(class, pkg)
	}

	jenny.typeFormatter = createFormatter(context, jenny.config).withPackageMapper(packageMapper)

	object, _ := context.GetObject(builder.For.SelfRef.ReferredPkg, builder.For.SelfRef.ReferredType)
	tmpl := BuilderTemplate{
		Package:              jenny.config.formatPackage(builder.Package),
		RawPackage:           builder.Package,
		Imports:              jenny.imports,
		ObjectName:           tools.UpperCamelCase(object.Name),
		BuilderName:          jenny.getBuilderName(builder),
		BuilderSignatureType: jenny.getBuilderSignature(builder.Package, object),
		Constructor:          builder.Constructor,
		Options:              builder.Options,
		Properties:           builder.Properties,
		ImportAlias:          jenny.config.PackagePath,
		IsGenericPanel:       builder.IsGeneric,
		Comments:             builder.For.Comments,
		DeprecationMessage:   builder.DeprecationMessage,
	}

	jenny.apiRefCollector.BuilderMethod(builder, apiref.MethodReference{
		Name: "build",
		Comments: []string{
			"Builds the object.",
		},
		Return: tools.UpperCamelCase(builder.Name),
	})

	return jenny.tmpl.Funcs(map[string]any{
		"formatBuilderFieldType": jenny.typeFormatter.formatBuilderFieldType,
		"emptyValueForType": func(def ir.Type) string {
			return jenny.typeFormatter.emptyValueForType(def, true)
		},
		"typeHasBuilder":           jenny.typeFormatter.typeHasBuilder,
		"resolvesToComposableSlot": jenny.typeFormatter.resolvesToComposableSlot,
		"formatAssignmentPath":     jenny.typeFormatter.formatAssignmentPath,
		"formatPath":               jenny.typeFormatter.formatFieldPath,
		"formatRefType":            jenny.typeFormatter.formatRefType,
		"formatType":               jenny.typeFormatter.formatFieldType,
		"formatPathIndex":          jenny.typeFormatter.formatPathIndex,
	}).RenderAsBytes("builders/builder.tmpl", tmpl)
}

func (jenny Builder) getBuilderName(builder ir.Builder) string {
	if builder.For.SelfRef.ReferredPkg != builder.Package {
		return fmt.Sprintf("%s%s", tools.UpperCamelCase(builder.Package), tools.UpperCamelCase(builder.Name))
	}

	return tools.UpperCamelCase(builder.Name)
}

func (jenny Builder) getBuilderSignature(pkg string, obj ir.Object) string {
	if pkg != obj.SelfRef.ReferredPkg {
		jenny.imports.Add(obj.SelfRef.ReferredType, obj.SelfRef.ReferredPkg)
	}

	if !obj.Type.IsDataqueryVariant() {
		return obj.Name
	}

	return fmt.Sprintf("%s.%s", jenny.config.formatPackage("cog.variants"), tools.UpperCamelCase(obj.Type.ImplementedVariant()))
}
