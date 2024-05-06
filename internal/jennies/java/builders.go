package java

import (
	"fmt"
	"strings"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/jennies/template"
	"github.com/grafana/cog/internal/tools"
)

type Builders struct {
	context       common.Context
	typeFormatter *typeFormatter
	builders      map[string]map[string]ast.Builder
}

func parseBuilders(context common.Context, formatter *typeFormatter) Builders {
	b := make(map[string]map[string]ast.Builder)
	for _, builder := range context.Builders {
		if _, ok := b[builder.Package]; !ok {
			b[builder.Package] = map[string]ast.Builder{}
		}
		b[builder.Package][builder.Name] = builder
	}

	return Builders{
		context:       context,
		builders:      b,
		typeFormatter: formatter,
	}
}

func (b Builders) genBuilder(pkg string, name string, fields []Field) (Builder, bool) {
	builder, ok := b.getBuilder(pkg, name)
	if !ok {
		return Builder{}, false
	}

	object, _ := b.context.LocateObject(builder.For.SelfRef.ReferredPkg, builder.For.SelfRef.ReferredType)
	return Builder{
		Builder: template.Builder{
			ObjectName:  object.Name,
			BuilderName: builder.Name,
			Constructor: b.genConstructor(builder),
			Options:     b.genOptions(builder.Options),
			Properties:  builder.Properties,
		},
		Fields: fields,
	}, true
}

func (b Builders) getBuilder(pkg string, name string) (ast.Builder, bool) {
	builderMap, ok := b.builders[pkg]
	if !ok {
		return ast.Builder{}, false
	}

	builder, ok := builderMap[name]
	return builder, ok
}

func (b Builders) genConstructor(builder ast.Builder) template.Constructor {
	var argsList []ast.Argument
	var assignmentList []template.Assignment
	for _, opt := range builder.Options {
		if !opt.IsConstructorArg {
			continue
		}

		argsList = append(argsList, opt.Args[0])
		assignmentList = append(assignmentList, b.genAssignments(opt.Assignments)...)
	}

	return template.Constructor{
		Args:        argsList,
		Assignments: assignmentList,
	}
}

func (b Builders) genOptions(opts []ast.Option) []template.Option {
	options := make([]template.Option, len(opts))
	for i, option := range opts {
		options[i] = template.Option{
			Name:        tools.LowerCamelCase(option.Name),
			Args:        b.genArgs(option.Args),
			Assignments: b.genAssignments(option.Assignments),
		}
	}

	return options
}

func (b Builders) genArgs(arguments []ast.Argument) []ast.Argument {
	args := make([]ast.Argument, len(arguments))
	for i, arg := range arguments {
		args[i] = ast.Argument{
			Name: tools.LowerCamelCase(arg.Name),
			Type: arg.Type,
		}
	}

	return args
}

func (b Builders) genAssignments(assignments []ast.Assignment) []template.Assignment {
	assign := make([]template.Assignment, len(assignments))
	for i, assignment := range assignments {
		var constraints []template.Constraint
		if assignment.Value.Argument != nil {
			argName := escapeVarName(tools.LowerCamelCase(assignment.Value.Argument.Name))
			constraints = b.genConstraints(argName, assignment.Constraints)
		}

		assign[i] = template.Assignment{
			Path:           assignment.Path,
			Method:         assignment.Method,
			Constraints:    constraints,
			Value:          assignment.Value,
			InitSafeguards: b.getSafeGuards(assignment),
		}
	}

	return assign
}

func (b Builders) genConstraints(name string, constraints []ast.TypeConstraint) []template.Constraint {
	return tools.Map(constraints, func(t ast.TypeConstraint) template.Constraint {
		return template.Constraint{
			ArgName:   tools.LowerCamelCase(name),
			Op:        t.Op,
			Parameter: t.Args[0],
		}
	})
}

func (b Builders) getSafeGuards(assignment ast.Assignment) []string {
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
		initSafeGuards = append(initSafeGuards, b.initSafeGuard(subPath))
	}

	return initSafeGuards
}

func (b Builders) initSafeGuard(path ast.Path) string {
	parts := formatFieldPath(path)
	valueType := path.Last().Type
	if path.Last().TypeHint != nil {
		valueType = *path.Last().TypeHint
	}

	emptyValue := b.typeFormatter.defaultValueFor(valueType)
	if len(parts) == 1 {
		return fmt.Sprintf(
			`	if (this.%[1]s == null) {
			this.%[1]s = %[2]s;
		}`, tools.LowerCamelCase(parts[0]), emptyValue)
	}

	return fmt.Sprintf(
		`	if (this.%[1]s == null) {
			this.%[1]s = %[2]s;
		}`, strings.Join(parts, "."), emptyValue)
}
