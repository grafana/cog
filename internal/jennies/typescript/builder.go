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
	Config Config

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
			strings.ToLower(builder.Package),
			fmt.Sprintf("%s_builder_gen.ts", strings.ToLower(builder.For.Name)),
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

	buildObjectSignature := builder.For.SelfRef.ReferredPkg + "." + tools.UpperCamelCase(builder.For.Name)
	if builder.For.Type.ImplementsVariant() {
		buildObjectSignature = jenny.typeFormatter.variantInterface(builder.For.Type.ImplementedVariant())
	}

	comments := builder.For.Comments
	if jenny.Config.Debug {
		veneerTrail := tools.Map(builder.VeneerTrail, func(veneer string) string {
			return fmt.Sprintf("Modified by veneer '%s'", veneer)
		})
		comments = append(comments, veneerTrail...)
	}

	err := templates.
		Funcs(map[string]any{
			"typeHasBuilder": context.ResolveToBuilder,
			"formatType":     jenny.typeFormatter.formatType,
			"resolvesToComposableSlot": func(typeDef ast.Type) bool {
				_, found := context.ResolveToComposableSlot(typeDef)
				return found
			},
			"defaultValueForType": func(typeDef ast.Type) string {
				return formatValue(jenny.rawTypes.defaultValueForType(typeDef, jenny.typeImportMapper))
			},
		}).
		ExecuteTemplate(&buffer, "builder.tmpl", template.Builder{
			BuilderName:          builder.Name,
			ObjectName:           builder.For.Name,
			BuilderSignatureType: buildObjectSignature,
			Imports:              jenny.imports,
			ImportAlias:          jenny.importType(builder.For.SelfRef),
			Comments:             comments,
			Constructor:          jenny.generateConstructor(builder),
			Properties:           builder.Properties,
			Options:              tools.Map(builder.Options, jenny.generateOption),
		})

	return []byte(buffer.String()), err
}

func (jenny *Builder) generateConstructor(builder ast.Builder) template.Constructor {
	var argsList []ast.Argument
	var assignments []template.Assignment
	for _, opt := range builder.Options {
		if !opt.IsConstructorArg {
			continue
		}

		// FIXME: this is assuming that there's only one argument for that option
		argsList = append(argsList, opt.Args[0])
		assignments = append(assignments, jenny.generateAssignment(opt.Assignments[0]))
	}

	for _, init := range builder.Initializations {
		assignments = append(assignments, jenny.generateAssignment(init))
	}

	return template.Constructor{
		Args:        argsList,
		Assignments: assignments,
	}
}

func (jenny *Builder) generateOption(def ast.Option) template.Option {
	comments := def.Comments

	if jenny.Config.Debug {
		veneerTrail := tools.Map(def.VeneerTrail, func(veneer string) string {
			return fmt.Sprintf("Modified by veneer '%s'", veneer)
		})
		comments = append(comments, veneerTrail...)
	}

	assignments := tools.Map(def.Assignments, func(assignment ast.Assignment) template.Assignment {
		return jenny.generateAssignment(assignment)
	})

	return template.Option{
		Name:        def.Name,
		Comments:    comments,
		Args:        def.Args,
		Assignments: assignments,
	}
}

func (jenny *Builder) generatePathInitializationSafeGuard(path ast.Path) string {
	fieldPath := jenny.formatFieldPath(path)
	valueType := path.Last().Type
	if path.Last().TypeHint != nil {
		valueType = *path.Last().TypeHint
	}

	emptyValue := formatValue(jenny.rawTypes.defaultValueForType(valueType, jenny.typeImportMapper))

	return fmt.Sprintf(`		if (!this.internal.%[1]s) {
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
			chunkType.Kind == ast.KindMap ||
			chunkType.Kind == ast.KindArray ||
			chunkType.Kind == ast.KindRef ||
			chunkType.Kind == ast.KindStruct

		if !maybeUndefined {
			continue
		}

		subPath := assign.Path[:i+1]
		initSafeGuards = append(initSafeGuards, jenny.generatePathInitializationSafeGuard(subPath))
	}

	var constraints []template.Constraint
	if assign.Value.Argument != nil {
		argName := tools.LowerCamelCase(assign.Value.Argument.Name)
		constraints = jenny.constraints(argName, assign.Constraints)
	}

	return template.Assignment{
		Path:           assign.Path,
		InitSafeguards: initSafeGuards,
		Constraints:    constraints,
		Method:         assign.Method,
		Value:          assign.Value,
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
