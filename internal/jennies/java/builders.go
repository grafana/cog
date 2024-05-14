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

func (b Builders) genBuilder(pkg string, name string) (template.Builder, bool) {
	builder, ok := b.getBuilder(pkg, name)
	if !ok {
		return template.Builder{}, false
	}

	object, _ := b.context.LocateObject(builder.For.SelfRef.ReferredPkg, builder.For.SelfRef.ReferredType)
	return template.Builder{
		ObjectName:  tools.UpperCamelCase(object.Name),
		BuilderName: builder.Name,
		Constructor: b.genConstructor(builder),
		Options:     b.genOptions(builder.Options),
		Properties:  builder.Properties,
		Defaults:    b.genDefaults(builder.Options),
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
			`	if (this.internal.%[1]s == null) {
			this.internal.%[1]s = %[2]s;
		}`, tools.LowerCamelCase(parts[0]), emptyValue)
	}

	return fmt.Sprintf(
		`	if (this.internal.%[1]s == null) {
			this.internal.%[1]s = %[2]s;
		}`, strings.Join(parts, "."), emptyValue)
}

func (b Builders) genDefaults(options []ast.Option) []template.OptionCall {
	calls := make([]template.OptionCall, 0)
	for _, opt := range options {
		if opt.Default == nil || len(opt.Args) == 0 {
			continue
		}

		calls = append(calls, template.OptionCall{
			OptionName: tools.UpperCamelCase(opt.Name),
			Args:       b.formatDefaultValues(opt.Args),
		})
	}

	return calls
}

func (b Builders) formatDefaultValues(args []ast.Argument) []string {
	argumentList := make([]string, 0, len(args))
	for _, arg := range args {
		switch arg.Type.Kind {
		case ast.KindRef:
			argumentList = append(argumentList, b.formatDefaultReference(arg.Type.AsRef(), arg.Type.Default))
		case ast.KindScalar:
			scalar := arg.Type.AsScalar()
			if scalar.ScalarKind == ast.KindFloat32 || scalar.ScalarKind == ast.KindFloat64 {
				argumentList = append(argumentList, fmt.Sprintf("%.1f", float64(arg.Type.Default.(int64))))
			} else {
				argumentList = append(argumentList, fmt.Sprintf("%#v", arg.Type.Default))
			}
		case ast.KindArray:
			array := arg.Type.AsArray()
			if array.IsArrayOfScalars() {
				argumentList = append(argumentList, fmt.Sprintf("List.of(%s)", b.typeFormatter.formatScalar(arg.Type.Default)))
			}
		case ast.KindStruct:
			// TODO: Java is using veneers to avoid anonymous structs but it should be detailed if we need it at any moment.
			argumentList = append(argumentList, "new Object()")
		}
		// TODO: Implement the rest of types if any
	}

	return argumentList
}

func (b Builders) formatDefaultReference(ref ast.RefType, defValue any) string {
	object, _ := b.context.LocateObject(ref.ReferredPkg, ref.ReferredType)
	switch object.Type.Kind {
	case ast.KindEnum:
		for _, v := range object.Type.AsEnum().Values {
			if defValue == v.Value {
				return fmt.Sprintf("%s.%s", object.Name, tools.UpperSnakeCase(v.Name))
			}
		}
	case ast.KindStruct:
		// TODO: Builder could have arguments 🙃
		builder := fmt.Sprintf("new %s.Builder()", tools.UpperCamelCase(object.Name))
		structType := object.Type.AsStruct()
		defValues := defValue.(map[string]interface{})
		for _, field := range structType.Fields {
			if f, ok := defValues[field.Name]; ok {
				builder = fmt.Sprintf("%s.set%s(%#v)", builder, tools.UpperCamelCase(field.Name), f)
			}
		}
		return builder + ".build()"
	}

	return ""
}
