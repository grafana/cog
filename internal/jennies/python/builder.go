package python

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
	imports          *ModuleImportMap
	typeFormatter    *typeFormatter
	rawTypeFormatter *typeFormatter
	importModule     func(alias string, pkg string, module string) string
}

func (jenny *Builder) JennyName() string {
	return "PythonBuilder"
}

func (jenny *Builder) Generate(context common.Context) (codejen.Files, error) {
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

func (jenny *Builder) generateBuilder(context common.Context, builder ast.Builder) ([]byte, error) {
	var buffer strings.Builder

	// every builder uses the following imports
	jenny.imports.AddPackage("typing", "typing")
	jenny.importModule("cogbuilder", "..cog", "builder")

	fullObjectName := jenny.typeFormatter.formatRef(builder.For.SelfRef)
	buildObjectSignature := fullObjectName

	err := templates.
		Funcs(map[string]any{
			"formatType":     jenny.typeFormatter.formatType,
			"formatRawType":  jenny.rawTypeFormatter.formatType,
			"typeHasBuilder": context.ResolveToBuilder,
			"resolvesToComposableSlot": func(typeDef ast.Type) bool {
				_, found := context.ResolveToComposableSlot(typeDef)
				return found
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
		ExecuteTemplate(&buffer, "builders/builder.tmpl", template.Builder{
			Package:              builder.Package,
			BuilderSignatureType: buildObjectSignature,
			BuilderName:          tools.UpperCamelCase(builder.Name),
			ObjectName:           fullObjectName,
			Comments:             builder.For.Comments,
			Constructor:          jenny.generateConstructor(context, builder),
			Properties:           builder.Properties,
			Options: tools.Map(builder.Options, func(option ast.Option) template.Option {
				return jenny.generateOption(context, option)
			}),
		})
	if err != nil {
		return nil, err
	}

	return []byte(buffer.String()), nil
}

func (jenny *Builder) generateConstructor(context common.Context, builder ast.Builder) template.Constructor {
	var argsList []ast.Argument
	var assignments []template.Assignment
	for _, opt := range builder.Options {
		if !opt.IsConstructorArg {
			continue
		}

		// FIXME: this is assuming that there's only one argument for that option
		argsList = append(argsList, opt.Args[0])
		assignments = append(assignments, jenny.generateAssignment(context, opt.Assignments[0]))
	}

	for _, init := range builder.Initializations {
		assignments = append(assignments, jenny.generateAssignment(context, init))
	}

	return template.Constructor{
		Args:        argsList,
		Assignments: assignments,
	}
}

func (jenny *Builder) generateOption(context common.Context, def ast.Option) template.Option {
	return template.Option{
		Name:     def.Name,
		Comments: def.Comments,
		Args: tools.Map(def.Args, func(arg ast.Argument) ast.Argument {
			newArg := arg.DeepCopy()
			newArg.Type.Nullable = false

			return newArg
		}),
		Assignments: tools.Map(def.Assignments, func(assignment ast.Assignment) template.Assignment {
			return jenny.generateAssignment(context, assignment)
		}),
	}
}

func (jenny *Builder) generatePathInitializationSafeGuard(context common.Context, path ast.Path) string {
	fieldPath := formatFieldPath(path)
	valueType := path.Last().Type
	if path.Last().TypeHint != nil {
		valueType = *path.Last().TypeHint
	}

	emptyValue := formatValue(defaultValueForType(context.Schemas, valueType, jenny.importModule, nil))

	nonOptionalType := valueType.DeepCopy()
	nonOptionalType.Nullable = false

	guard := fmt.Sprintf(`if self._internal.%[1]s is None:
    self._internal.%[1]s = %[2]s
`, fieldPath, emptyValue)

	if !nonOptionalType.IsArray() {
		guard += fmt.Sprintf("\nassert isinstance(self._internal.%s, %s)\n", fieldPath, jenny.rawTypeFormatter.formatType(nonOptionalType))
	}

	return guard
}

func (jenny *Builder) generateAssignment(context common.Context, assignment ast.Assignment) template.Assignment {
	var initSafeGuards []string
	for i, chunk := range assignment.Path {
		if i == len(assignment.Path)-1 && assignment.Method != ast.AppendAssignment {
			continue
		}

		chunkType := chunk.Type
		if chunk.TypeHint != nil {
			chunkType = *chunk.TypeHint
		}

		maybeUndefined := chunkType.Nullable ||
			chunkType.IsAnyOf(ast.KindMap, ast.KindArray, ast.KindRef, ast.KindStruct) ||
			chunk.Type.IsAny()

		if !maybeUndefined {
			continue
		}

		subPath := assignment.Path[:i+1]
		initSafeGuards = append(initSafeGuards, jenny.generatePathInitializationSafeGuard(context, subPath))
	}

	var constraints []template.Constraint
	if assignment.Value.Argument != nil {
		constraints = jenny.constraints(assignment.Value.Argument.Name, assignment.Constraints)
	}

	return template.Assignment{
		Path:           assignment.Path,
		InitSafeguards: initSafeGuards,
		Constraints:    constraints,
		Method:         assignment.Method,
		Value:          assignment.Value,
	}
}

func (jenny *Builder) constraints(argumentName string, constraints []ast.TypeConstraint) []template.Constraint {
	return tools.Map(constraints, func(constraint ast.TypeConstraint) template.Constraint {
		return template.Constraint{
			ArgName:   argumentName,
			Op:        constraint.Op,
			Parameter: constraint.Args[0],
		}
	})
}
