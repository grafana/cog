package typescript

import (
	"fmt"
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
	tmpl *template.Template

	imports          *common.DirectImportMap
	typeImportMapper func(string) string
	typeFormatter    *typeFormatter
	rawTypes         RawTypes
}

func (jenny *Builder) JennyName() string {
	return "TypescriptBuilder"
}

func (jenny *Builder) Generate(context languages.Context) (codejen.Files, error) {
	files := codejen.Files{}
	jenny.rawTypes = RawTypes{
		schemas: context.Schemas,
	}

	for _, builder := range context.Builders {
		output, err := jenny.generateBuilder(context, builder)
		if err != nil {
			return nil, err
		}

		filename := filepath.Join(
			"src",
			formatPackageName(builder.Package),
			fmt.Sprintf("%sBuilder.gen.ts", tools.LowerCamelCase(builder.Name)),
		)

		files = append(files, *codejen.NewFile(filename, output, jenny))
	}

	return files, nil
}

func (jenny *Builder) generateBuilder(context languages.Context, builder ast.Builder) ([]byte, error) {
	jenny.imports = NewImportMap()
	jenny.imports.Add("cog", "../cog")
	jenny.typeImportMapper = func(pkg string) string {
		return jenny.imports.Add(pkg, fmt.Sprintf("../%s", pkg))
	}
	jenny.typeFormatter = builderTypeFormatter(context, jenny.typeImportMapper)

	buildObjectSignature := formatPackageName(builder.For.SelfRef.ReferredPkg) + "." + tools.CleanupNames(builder.For.Name)
	if builder.For.Type.ImplementsVariant() {
		buildObjectSignature = jenny.typeFormatter.variantInterface(builder.For.Type.ImplementedVariant())
	}

	return jenny.tmpl.
		Funcs(map[string]any{
			"typeHasBuilder":              context.ResolveToBuilder,
			"typeIsDisjunctionOfBuilders": context.IsDisjunctionOfBuilders,
			"formatType":                  jenny.typeFormatter.formatType,
			"resolvesToComposableSlot": func(typeDef ast.Type) bool {
				_, found := context.ResolveToComposableSlot(typeDef)
				return found
			},
			"defaultValueForType": func(typeDef ast.Type) string {
				return formatValue(jenny.rawTypes.defaultValueForType(typeDef, jenny.typeImportMapper))
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
			"formatPath": jenny.formatFieldPath,
			"emptyValueForGuard": func(guard ast.AssignmentNilCheck) string {
				return formatValue(jenny.rawTypes.defaultValueForType(guard.EmptyValueType, jenny.typeImportMapper))
			},
		}).
		RenderAsBytes("builder.tmpl", template.Builder{
			BuilderName:          builder.Name,
			ObjectName:           tools.CleanupNames(builder.For.Name),
			BuilderSignatureType: buildObjectSignature,
			Imports:              jenny.imports,
			ImportAlias:          jenny.importType(builder.For.SelfRef),
			Comments:             builder.For.Comments,
			Constructor:          builder.Constructor,
			Properties:           builder.Properties,
			Options:              builder.Options,
		})
}

// importType declares an import statement for the type definition of
// the given object and returns an alias to it.
func (jenny *Builder) importType(typeRef ast.RefType) string {
	return jenny.typeImportMapper(typeRef.ReferredPkg)
}

func (jenny *Builder) formatFieldPath(fieldPath ast.Path) string {
	return strings.Join(tools.Map(fieldPath, func(chunk ast.PathItem) string {
		return chunk.Identifier
	}), ".")
}
