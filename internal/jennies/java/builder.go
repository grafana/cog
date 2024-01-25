package java

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/tools"
)

type Builder struct {
	config  Config
	imports *common.DirectImportMap

	typeFormatter *typeFormatter
}

func (jenny Builder) JennyName() string {
	return "JavaBuilder"
}

func (jenny Builder) Generate(context common.Context) (codejen.Files, error) {
	files := make(codejen.Files, len(context.Builders))
	jenny.typeFormatter = createFormatter(context)

	for i, builder := range context.Builders {
		output, err := jenny.genFilesForBuilder(context, builder)
		if err != nil {
			return nil, err
		}

		filename := filepath.Join(
			strings.ToLower(builder.Package),
			fmt.Sprintf("%sBuilder.java", tools.UpperCamelCase(builder.Name)),
		)

		files[i] = *codejen.NewFile(filename, output, jenny)
	}

	return files, nil
}

func (jenny Builder) genFilesForBuilder(context common.Context, builder ast.Builder) ([]byte, error) {
	var buffer strings.Builder

	jenny.imports = NewImportMap()
	jenny.imports.Add("Builder", "cog.variants")
	jenny.typeFormatter = jenny.typeFormatter.withPackageMapper(func(pkg string, class string) string {
		if builder.Package == pkg {
			return ""
		}

		return jenny.imports.Add(class, pkg)
	})

	object, _ := context.LocateObject(builder.For.SelfRef.ReferredPkg, builder.For.SelfRef.ReferredType)
	err := templates.Funcs(map[string]any{
		"formatType":     jenny.typeFormatter.formatFieldType,
		"lowerCamelCase": tools.LowerCamelCase,
		"typeHasBuilder": context.ResolveToBuilder,
		"resolvesToComposableSlot": func(typeDef ast.Type) bool {
			_, found := context.ResolveToComposableSlot(typeDef)
			return found
		},
	}).ExecuteTemplate(&buffer, "builders/builder.tmpl", BuilderTemplate{
		Package:         builder.Package,
		Imports:         jenny.imports,
		Name:            builder.Name,
		ObjectSignature: jenny.getFullObjectName(builder.For.SelfRef),
		Options:         jenny.genOptions(builder.Options),
		Fields:          jenny.genFields(context, object),
		Properties:      jenny.genProperties(builder.Properties),
	})

	if err != nil {
		return nil, err
	}

	return []byte(buffer.String()), nil
}

func (jenny Builder) getFullObjectName(ref ast.RefType) string {
	refType := tools.UpperCamelCase(ref.ReferredType)
	jenny.typeFormatter.packageMapper(ref.ReferredPkg, refType)
	return refType
}

func (jenny Builder) genOptions(opts []ast.Option) []Option {
	options := make([]Option, len(opts))
	for i, option := range opts {
		options[i] = Option{
			Name:        tools.UpperCamelCase(option.Name),
			Args:        jenny.genArgs(option.Args),
			Assignments: jenny.genAssignments(option.Assignments),
		}
	}

	return options
}

func (jenny Builder) genArgs(arguments []ast.Argument) []Arg {
	args := make([]Arg, len(arguments))
	for i, arg := range arguments {
		args[i] = Arg{
			Name: tools.LowerCamelCase(arg.Name),
			Type: jenny.typeFormatter.formatBuilderArgs(arg.Type),
		}
	}

	return args
}

func (jenny Builder) genAssignments(assignments []ast.Assignment) []Assignment {
	assign := make([]Assignment, len(assignments))
	for i, assignment := range assignments {

		var constraints []Constraint
		if assignment.Value.Argument != nil {
			argName := escapeVarName(tools.LowerCamelCase(assignment.Value.Argument.Name))
			constraints = jenny.genConstraints(argName, assignment.Constraints)
		}

		assign[i] = Assignment{
			Path:           assignment.Path,
			Method:         assignment.Method,
			Constraints:    constraints,
			Value:          assignment.Value,
			InitSafeguards: jenny.getSafeGuards(assignment),
		}
	}

	return assign
}

func (jenny Builder) genFields(context common.Context, object ast.Object) []Field {
	fields := make([]Field, 0)
	switch object.Type.Kind {
	case ast.KindStruct:
		for _, field := range object.Type.AsStruct().Fields {
			fields = append(fields, Field{
				Name:     tools.LowerCamelCase(field.Name),
				Type:     jenny.typeFormatter.formatFieldType(field.Type),
				Comments: field.Comments,
			})
		}
	case ast.KindRef:
		ref := object.Type.AsRef()
		obj, ok := context.LocateObject(ref.ReferredPkg, ref.ReferredType)
		if !ok {
			break
		}
		fields = append(fields, jenny.genFields(context, obj)...)
	default:
		// Shouldn't reach here...
		return nil
	}

	return fields
}

func (jenny Builder) genProperties(properties []ast.StructField) []Field {
	fields := make([]Field, len(properties))
	for i, field := range properties {
		fields[i] = Field{
			Name:     tools.LowerCamelCase(field.Name),
			Type:     jenny.typeFormatter.formatFieldType(field.Type),
			Comments: field.Comments,
		}
	}

	return fields
}

func (jenny Builder) genConstraints(name string, constraints []ast.TypeConstraint) []Constraint {
	return tools.Map(constraints, func(t ast.TypeConstraint) Constraint {
		return Constraint{
			ArgName:   tools.LowerCamelCase(name),
			Op:        t.Op,
			Parameter: t.Args[0],
		}
	})
}

func (jenny Builder) getSafeGuards(assignment ast.Assignment) []string {
	var initSafeGuards []string
	for i, chunk := range assignment.Path {
		if i == len(assignment.Path)-1 && assignment.Method != ast.AppendAssignment {
			continue
		}

		canNullPointer := chunk.Type.IsAnyOf(ast.KindMap, ast.KindArray, ast.KindRef, ast.KindStruct) || chunk.Type.IsAny()

		if !canNullPointer {
			continue
		}

		subPath := assignment.Path[:i+1]
		initSafeGuards = append(initSafeGuards, jenny.initSafeGuard(subPath))
	}

	return initSafeGuards
}

func (jenny Builder) initSafeGuard(path ast.Path) string {
	parts := formatFieldPath(path)
	valueType := path.Last().Type
	if path.Last().TypeHint != nil {
		valueType = *path.Last().TypeHint
	}

	emptyValue := jenny.typeFormatter.defaultValueFor(valueType)
	if len(parts) == 1 {
		return fmt.Sprintf(
			`	if (this.%[1]s == null) {
			this.%[1]s = %[2]s;
		}`, tools.LowerCamelCase(parts[0]), emptyValue)
	}

	getters := make([]string, len(parts))
	setters := make([]string, len(parts))
	for i, p := range parts {
		if i == 0 {
			getters[i] = tools.LowerCamelCase(p)
			setters[i] = tools.LowerCamelCase(p)
			continue
		}

		getters[i] = fmt.Sprintf("get%s()", p)
		if i == len(parts)-1 {
			setters[i] = fmt.Sprintf("set%s", p)
		} else {
			setters[i] = fmt.Sprintf("get%s()", p)
		}
	}

	return fmt.Sprintf(
		`	if (this.%[1]s == null) {
			this.%[2]s(%[3]s);
		}`, strings.Join(getters, "."), strings.Join(setters, "."), emptyValue)
}
