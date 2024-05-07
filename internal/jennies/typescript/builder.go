package typescript

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/jennies/template"
	"github.com/grafana/cog/internal/tools"
)

type Builder struct {
	imports          *common.DirectImportMap
	typeImportMapper func(string) string
	typeFormatter    *typeFormatter
	rawTypes         RawTypes
}

func (jenny *Builder) JennyName() string {
	return "TypescriptBuilder"
}

func (jenny *Builder) Generate(context common.Context) (codejen.Files, error) {
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

func (jenny *Builder) generateBuilder(context common.Context, builder ast.Builder) ([]byte, error) {
	var buffer strings.Builder

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

	err := templates.
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
		}).
		ExecuteTemplate(&buffer, "builder.tmpl", template.Builder{
			BuilderName:          builder.Name,
			ObjectName:           tools.CleanupNames(builder.For.Name),
			BuilderSignatureType: buildObjectSignature,
			Imports:              jenny.imports,
			ImportAlias:          jenny.importType(builder.For.SelfRef),
			Comments:             builder.For.Comments,
			Constructor:          jenny.generateConstructor(builder),
			Properties:           builder.Properties,
			Options:              tools.Map(builder.Options, jenny.generateOption),
		})

	return []byte(buffer.String()), err
}

func (jenny *Builder) generateConstructor(builder ast.Builder) template.Constructor {
	return template.Constructor{
		Args:        builder.Constructor.Args,
		Assignments: tools.Map(builder.Constructor.Assignments, jenny.generateAssignment),
	}
}

func (jenny *Builder) generateOption(def ast.Option) template.Option {
	return template.Option{
		Name:        formatIdentifier(def.Name),
		Comments:    def.Comments,
		Args:        def.Args,
		Assignments: tools.Map(def.Assignments, jenny.generateAssignment),
	}
}

func (jenny *Builder) generatePathInitializationSafeGuard(path ast.Path) string {
	fieldPath := jenny.formatFieldPath(path)
	valueType := path.Last().Type
	if path.Last().TypeHint != nil {
		valueType = *path.Last().TypeHint
	}

	emptyValue := formatValue(jenny.rawTypes.defaultValueForType(valueType, jenny.typeImportMapper))

	return fmt.Sprintf(`        if (!this.internal.%[1]s) {
            this.internal.%[1]s = %[2]s;
        }`, fieldPath, emptyValue)
}

func (jenny *Builder) generateAssignment(assign ast.Assignment) template.Assignment {
	var initSafeGuards []string
	for i, chunk := range assign.Path {
		if i == len(assign.Path)-1 && assign.Method != ast.AppendAssignment {
			continue
		}

		chunkType := chunk.Type
		if chunk.TypeHint != nil {
			chunkType = *chunk.TypeHint
		}

		maybeUndefined := chunkType.Nullable ||
			chunkType.IsAnyOf(ast.KindMap, ast.KindArray, ast.KindRef, ast.KindStruct)

		if !maybeUndefined {
			continue
		}

		subPath := assign.Path[:i+1]
		initSafeGuards = append(initSafeGuards, jenny.generatePathInitializationSafeGuard(subPath))
	}

	return template.Assignment{
		Path:           assign.Path,
		InitSafeguards: initSafeGuards,
		Constraints:    assign.Constraints,
		Method:         assign.Method,
		Value:          assign.Value,
	}
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
