package python

import (
	"path/filepath"
	"strings"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/jennies/template"
	"github.com/grafana/cog/internal/languages"
	"github.com/grafana/cog/internal/tools"
)

type Builder struct {
	tmpl            *template.Template
	apiRefCollector *common.APIReferenceCollector

	imports          *ModuleImportMap
	typeFormatter    *typeFormatter
	rawTypeFormatter *typeFormatter
	importModule     func(alias string, pkg string, module string) string
}

func (jenny *Builder) JennyName() string {
	return "PythonBuilder"
}

func (jenny *Builder) Generate(context languages.Context) (codejen.Files, error) {
	files := codejen.Files{}
	buildersByPackage := make(map[string][]ast.Builder)

	for _, builder := range context.Builders {
		buildersByPackage[strings.ToLower(builder.Package)] = append(buildersByPackage[strings.ToLower(builder.Package)], builder)
	}

	for pkg, builders := range buildersByPackage {
		var source strings.Builder

		jenny.imports = NewImportMap()
		jenny.importModule = func(alias string, pkg string, module string) string {
			return jenny.imports.AddModule(alias, pkg, module)
		}
		jenny.typeFormatter = builderTypeFormatter(context, func(alias string, pkg string) string {
			return jenny.imports.AddPackage(alias, pkg)
		}, jenny.importModule)
		jenny.rawTypeFormatter = defaultTypeFormatter(context, func(alias string, pkg string) string {
			return jenny.imports.AddPackage(alias, pkg)
		}, jenny.importModule)

		for _, builder := range builders {
			source.WriteString("\n\n")

			output, err := jenny.generateBuilder(context, builder)
			if err != nil {
				return nil, err
			}

			source.Write(output)
		}

		filename := filepath.Join("builders", pkg+".py")
		completeSource := jenny.imports.String() + "\n" + source.String()

		files = append(files, *codejen.NewFile(filename, []byte(completeSource), jenny))
	}

	return files, nil
}

func (jenny *Builder) generateBuilder(context languages.Context, builder ast.Builder) ([]byte, error) {
	// every builder uses the following imports
	jenny.imports.AddPackage("typing", "typing")
	jenny.importModule("cogbuilder", "..cog", "builder")

	fullObjectName := jenny.typeFormatter.formatRef(builder.For.SelfRef)
	buildObjectSignature := fullObjectName

	jenny.apiRefCollector.BuilderMethod(builder, common.MethodReference{
		Name: "build",
		Comments: []string{
			"Builds the object.",
		},
		Return: fullObjectName,
	})

	return jenny.tmpl.
		Funcs(common.TypeResolvingTemplateHelpers(context)).
		Funcs(map[string]any{
			"isDisjunctionOfBuilders": context.IsDisjunctionOfBuilders,
			"formatType":              jenny.typeFormatter.formatType,
			"formatRawType":           jenny.rawTypeFormatter.formatType,
			"formatRawTypeNotNullable": func(def ast.Type) string {
				typeDef := def.DeepCopy()
				typeDef.Nullable = false

				return jenny.rawTypeFormatter.formatType(typeDef)
			},
			"formatValue": func(destinationType ast.Type, value any) string {
				if destinationType.IsRef() {
					referredObj, found := context.LocateObject(destinationType.AsRef().ReferredPkg, destinationType.AsRef().ReferredType)
					if found && referredObj.Type.IsEnum() {
						return jenny.typeFormatter.formatEnumValue(referredObj, value)
					}
				}

				return formatValue(value)
			},
			"defaultForType": func(typeDef ast.Type) string {
				return formatValue(defaultValueForType(context.Schemas, typeDef, jenny.importModule, nil))
			},
		}).
		RenderAsBytes("builders/builder.tmpl", map[string]any{
			"Package":              builder.Package,
			"BuilderSignatureType": buildObjectSignature,
			"BuilderName":          tools.UpperCamelCase(builder.Name),
			"ObjectName":           fullObjectName,
			"Comments":             builder.For.Comments,
			"Constructor":          builder.Constructor,
			"Properties":           builder.Properties,
			"Options":              tools.Map(builder.Options, jenny.generateOption),
		})
}

func (jenny *Builder) generateOption(option ast.Option) ast.Option {
	option.Args = tools.Map(option.Args, func(arg ast.Argument) ast.Argument {
		newArg := arg.DeepCopy()
		newArg.Type.Nullable = false

		return newArg
	})

	return option
}
