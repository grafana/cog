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
	isPanel       map[string]bool
}

func parseBuilders(context common.Context, formatter *typeFormatter) Builders {
	b := make(map[string]map[string]ast.Builder)
	panels := make(map[string]bool)
	for _, builder := range context.Builders {
		if _, ok := b[builder.Package]; !ok {
			b[builder.Package] = map[string]ast.Builder{}
		}
		b[builder.Package][builder.Name] = builder
		panels[builder.Package] = builder.Name == "Panel" && builder.Package != "dashboard" // TODO: Ugh! Maybe a compiler pass??
	}

	return Builders{
		context:       context,
		builders:      b,
		typeFormatter: formatter,
		isPanel:       panels,
	}
}

func (b Builders) genBuilder(pkg string, name string) (template.Builder, bool) {
	builder, ok := b.getBuilder(pkg, name)
	if !ok {
		return template.Builder{}, false
	}

	object, _ := b.context.LocateObject(builder.For.SelfRef.ReferredPkg, builder.For.SelfRef.ReferredType)
	return template.Builder{
		Package:     builder.Package,
		ObjectName:  tools.UpperCamelCase(object.Name),
		BuilderName: builder.Name,
		Constructor: b.genConstructor(builder),
		Options:     b.genOptions(builder.Options),
		Properties:  builder.Properties,
	}, true
}

func (b Builders) genPanelBuilder(pkg string) (template.Builder, bool) {
	if !b.isPanel[pkg] {
		return template.Builder{}, false
	}

	return b.genBuilder(pkg, "Panel")
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
	return template.Constructor{
		Args:        builder.Constructor.Args,
		Assignments: tools.Map(builder.Constructor.Assignments, b.genAssignment),
	}
}

func (b Builders) genOptions(opts []ast.Option) []template.Option {
	options := make([]template.Option, len(opts))
	for i, option := range opts {
		options[i] = template.Option{
			Name:        tools.LowerCamelCase(option.Name),
			Args:        b.genArgs(option.Args),
			Assignments: tools.Map(option.Assignments, b.genAssignment),
		}
	}

	return options
}

func (b Builders) genArgs(arguments []ast.Argument) []ast.Argument {
	args := make([]ast.Argument, len(arguments))
	for i, arg := range arguments {
		// TODO: To force to add imports. In no-panel builders they are already added when set Schema fields
		// TODO: Try to find an elegant way later...
		b.typeFormatter.formatFieldType(arg.Type)
		args[i] = ast.Argument{
			Name: tools.LowerCamelCase(arg.Name),
			Type: arg.Type,
		}
	}

	return args
}

func (b Builders) genAssignment(assignment ast.Assignment) template.Assignment {
	return template.Assignment{
		Path:           assignment.Path,
		Method:         assignment.Method,
		Constraints:    assignment.Constraints,
		Value:          assignment.Value,
		InitSafeguards: b.getSafeGuards(assignment),
	}
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
