package golang

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/template"
	"github.com/grafana/cog/internal/languages"
	"github.com/grafana/cog/internal/orderedmap"
	"github.com/grafana/cog/internal/tools"
)

type Builder struct {
	Config Config

	typeImportMapper func(pkg string) string
	pathFormatter    func(path ast.Path) string
	typeFormatter    *typeFormatter
}

func (jenny *Builder) JennyName() string {
	return "GoBuilder"
}

func (jenny *Builder) Generate(context languages.Context) (codejen.Files, error) {
	files := codejen.Files{}

	for _, builder := range context.Builders {
		output, err := jenny.generateBuilder(context, builder)
		if err != nil {
			return nil, err
		}

		filename := filepath.Join(
			formatPackageName(builder.Package),
			fmt.Sprintf("%s_builder_gen.go", strings.ToLower(builder.Name)),
		)

		files = append(files, *codejen.NewFile(filename, output, jenny))
	}

	return files, nil
}

func (jenny *Builder) generateBuilder(context languages.Context, builder ast.Builder) ([]byte, error) {
	var buffer strings.Builder

	imports := NewImportMap()
	jenny.typeImportMapper = func(pkg string) string {
		if imports.IsIdentical(pkg, builder.Package) {
			return ""
		}

		return imports.Add(pkg, jenny.Config.importPath(pkg))
	}
	jenny.typeFormatter = builderTypeFormatter(jenny.Config, context, jenny.typeImportMapper)
	jenny.pathFormatter = makePathFormatter(jenny.typeFormatter)

	// every builder has a dependency on cog's runtime, so let's make sure it's declared.
	jenny.typeImportMapper("cog")

	fullObjectName := jenny.importType(builder.For.SelfRef)
	buildObjectSignature := fullObjectName
	if builder.For.Type.ImplementsVariant() {
		buildObjectSignature = jenny.typeFormatter.variantInterface(builder.For.Type.ImplementedVariant())
	}

	err := templates.
		Funcs(map[string]any{
			"formatPath": jenny.pathFormatter,
			"formatType": jenny.typeFormatter.formatType,
			"formatTypeNoBuilder": func(typeDef ast.Type) string {
				return jenny.typeFormatter.doFormatType(typeDef, false)
			},
			"typeHasBuilder": context.ResolveToBuilder,
			"resolvesToComposableSlot": func(typeDef ast.Type) bool {
				_, found := context.ResolveToComposableSlot(typeDef)
				return found
			},
			"emptyValueForGuard": func(guard ast.AssignmentNilCheck) string {
				emptyValue := jenny.emptyValueForType(guard.EmptyValueType)

				// This should be alright since there shouldn't be any scalar in the middle of a path
				if emptyValue[0] == '*' {
					emptyValue = "&" + emptyValue[1:]
				}

				if guard.Path.Last().Type.IsAny() && emptyValue[0] != '&' {
					emptyValue = "&" + emptyValue
				}

				return emptyValue
			},
		}).
		ExecuteTemplate(&buffer, "builders/builder.tmpl", template.Builder{
			Package:              builder.Package,
			BuilderSignatureType: buildObjectSignature,
			Imports:              imports,
			BuilderName:          tools.UpperCamelCase(builder.Name),
			ObjectName:           fullObjectName,
			Comments:             builder.For.Comments,
			Constructor:          builder.Constructor,
			Properties:           builder.Properties,
			Defaults:             jenny.genDefaultOptionsCalls(context, builder),
			Options:              builder.Options,
		})
	if err != nil {
		return nil, err
	}

	return []byte(buffer.String()), nil
}

func (jenny *Builder) genDefaultOptionsCalls(context languages.Context, builder ast.Builder) []template.OptionCall {
	calls := make([]template.OptionCall, 0)
	for _, opt := range builder.Options {
		if opt.Default == nil {
			continue
		}

		if len(opt.Args) == 0 {
			continue
		}

		if hasTypedDefaults(opt) {
			calls = append(calls, template.OptionCall{
				OptionName: opt.Name,
				Args:       jenny.formatDefaultTypedArgs(context, opt),
			})
			continue
		}

		calls = append(calls, template.OptionCall{
			OptionName: opt.Name,
			Args:       tools.Map(opt.Default.ArgsValues, formatScalar),
		})
	}

	return calls
}

func hasTypedDefaults(opt ast.Option) bool {
	for _, defArg := range opt.Default.ArgsValues {
		if _, ok := defArg.(map[string]any); ok {
			return true
		}
	}

	return false
}

func (jenny *Builder) formatDefaultTypedArgs(context languages.Context, opt ast.Option) []string {
	args := make([]string, 0)
	for i, arg := range opt.Default.ArgsValues {
		val, _ := arg.(map[string]interface{})

		pkg := ""
		refPkg := ""
		if opt.Args[i].Type.IsRef() {
			refPkg = jenny.typeImportMapper(opt.Args[i].Type.AsRef().ReferredPkg)
			pkg = opt.Args[i].Type.AsRef().ReferredType
			_, isBuilder := context.Builders.LocateByObject(opt.Args[i].Type.AsRef().ReferredPkg, pkg)
			obj, ok := context.LocateObject(opt.Args[i].Type.AsRef().ReferredPkg, pkg)
			if !ok {
				return []string{"unknown"}
			}
			args = append(args, formatDefaultReferenceStructForBuilder(refPkg, pkg, isBuilder, obj.Type.AsStruct(), orderedmap.FromMap(val)))
		}

		// Anonymous structs
		if opt.Args[i].Type.IsStruct() {
			def := opt.Args[i].Type.AsStruct()
			args = append(args, formatAnonymousDefaultStruct(def, orderedmap.FromMap(val)))
		}
	}
	return args
}

func (jenny *Builder) emptyValueForType(typeDef ast.Type) string {
	switch typeDef.Kind {
	case ast.KindRef, ast.KindStruct, ast.KindArray, ast.KindMap:
		return jenny.typeFormatter.doFormatType(typeDef, false) + "{}"
	case ast.KindEnum:
		return formatScalar(typeDef.AsEnum().Values[0].Value)
	case ast.KindScalar:
		return "" // no need to do anything here

	default:
		return "unknown"
	}
}

// importType declares an import statement for the type definition of
// the given object and returns a fully qualified type name for it.
func (jenny *Builder) importType(typeRef ast.RefType) string {
	pkg := jenny.typeImportMapper(typeRef.ReferredPkg)
	typeName := tools.UpperCamelCase(typeRef.ReferredType)
	if pkg == "" {
		return typeName
	}

	return fmt.Sprintf("%s.%s", pkg, typeName)
}
