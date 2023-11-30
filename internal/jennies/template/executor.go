package template

import (
	"fmt"
	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/tools"
	"path/filepath"
	"strings"
	"text/template"
)

type BuilderFormatter interface {
	GetObjectName(builder ast.Builder) string
	GetBuilderSignature(builder ast.Builder) string
	GetSafeGuard(path ast.Path) string
	FormatFieldPath(fieldPath ast.Path) string
	GetEscapeVar(varName string) string
	ScalarPattern() string
}

type Executor struct {
	cfg Config

	templateBuilder *Template

	typeFormatter   TypeFormatter
	importMapper    *common.DirectImportMap
	builderTemplate BuilderFormatter
}

func NewExecutor(cfg Config) *Executor {
	return &Executor{
		cfg:             cfg,
		importMapper:    cfg.ImportMapper,
		templateBuilder: NewTemplate(cfg.TemplateConfig.Name, cfg.TemplateConfig.FuncMap),
		typeFormatter:   cfg.Formatter,
		builderTemplate: cfg.BuilderTemplate,
	}
}

func (e *Executor) JennyName() string {
	return "TemplateExecutor"
}

func (e *Executor) Generate(context common.Context) (codejen.Files, error) {
	files := codejen.Files{}
	for _, builder := range context.Builders {
		e.typeFormatter = e.typeFormatter.WithContext(context)
		output, err := e.generateBuilder(context, builder)
		if err != nil {
			return nil, err
		}

		filename := filepath.Join(
			strings.ToLower(builder.Package),
			fmt.Sprintf("%s_builder_gen.%s", strings.ToLower(builder.For.Name), e.cfg.FileExtension),
		)

		files = append(files, *codejen.NewFile(filename, output, e))
	}

	return files, nil
}

func (e *Executor) generateBuilder(context common.Context, builder ast.Builder) ([]byte, error) {
	e.typeFormatter = e.typeFormatter.ForBuilder(builder)

	buildObjectSignature := e.builderTemplate.GetBuilderSignature(builder)
	if builder.For.Type.ImplementsVariant() {
		buildObjectSignature = e.typeFormatter.VariantInterface(builder.For.Type.ImplementedVariant())
	}

	comments := builder.For.Comments
	if e.cfg.Debug {
		veneerTrail := tools.Map(builder.VeneerTrail, func(veneer string) string {
			return fmt.Sprintf("Modified by veneer '%s'", veneer)
		})
		comments = append(comments, veneerTrail...)
	}

	res, err := e.templateBuilder.
		AddFuncMap(template.FuncMap{
			"formatType":     e.formatType,
			"typeHasBuilder": context.ResolveToBuilder,
			"formatPath":     e.builderTemplate.FormatFieldPath,
			"formatScalar":   e.formatScalar,
			"formatArgName": func(name string) string {
				return e.builderTemplate.GetEscapeVar(tools.LowerCamelCase(name))
			},
			"resolvesToComposableSlot": func(typeDef ast.Type) bool {
				_, found := context.ResolveToComposableSlot(typeDef)
				return found
			},
			"formatTypeNoBuilder": func(typeDef ast.Type) string {
				return e.typeFormatter.FormatType(typeDef, true)
			},
		}).
		Execute("builders/builder.tmpl", Builder{
			Package:              builder.Package,
			BuilderSignatureType: buildObjectSignature,
			Imports:              e.importMapper,
			BuilderName:          tools.UpperCamelCase(builder.Name),
			ObjectName:           e.builderTemplate.GetObjectName(builder),
			Comments:             comments,
			Constructor:          e.generateConstructor(builder),
			Properties:           builder.Properties,
			Defaults:             e.genDefaultOptionsCalls(builder),
			Options:              tools.Map(builder.Options, e.generateOption),
		})
	if err != nil {
		return nil, err
	}

	return []byte(res), nil
}

func (e *Executor) generateConstructor(builder ast.Builder) Constructor {
	var argsList []ast.Argument
	var assignments []Assignment
	for _, opt := range builder.Options {
		if !opt.IsConstructorArg {
			continue
		}

		// FIXME: this is assuming that there's only one argument for that option
		argsList = append(argsList, opt.Args[0])
		assignments = append(assignments, e.generateAssignment(opt.Assignments[0]))
	}

	for _, init := range builder.Initializations {
		assignments = append(assignments, e.generateAssignment(init))
	}

	return Constructor{
		Args:        argsList,
		Assignments: assignments,
	}
}

func (e *Executor) genDefaultOptionsCalls(builder ast.Builder) []OptionCall {
	calls := make([]OptionCall, 0)
	for _, opt := range builder.Options {
		if opt.Default == nil {
			continue
		}

		calls = append(calls, OptionCall{
			OptionName: opt.Name,
			Args:       tools.Map(opt.Default.ArgsValues, e.formatScalar),
		})
	}

	return calls
}

func (e *Executor) generateOption(def ast.Option) Option {
	comments := def.Comments

	if e.cfg.Debug {
		veneerTrail := tools.Map(def.VeneerTrail, func(veneer string) string {
			return fmt.Sprintf("Modified by veneer '%s'", veneer)
		})
		comments = append(comments, veneerTrail...)
	}

	assignments := tools.Map(def.Assignments, func(assignment ast.Assignment) Assignment {
		return e.generateAssignment(assignment)
	})

	return Option{
		Name:        tools.UpperCamelCase(def.Name), // fix for TS
		Comments:    comments,
		Args:        def.Args,
		Assignments: assignments,
	}
}

func (e *Executor) generateAssignment(assignment ast.Assignment) Assignment {
	var initSafeGuards []string
	for i, chunk := range assignment.Path {
		if i == len(assignment.Path)-1 {
			continue
		}

		nullable := chunk.Type.Nullable ||
			chunk.Type.Kind == ast.KindMap ||
			chunk.Type.Kind == ast.KindArray ||
			chunk.Type.IsAny()
		if nullable {
			subPath := assignment.Path[:i+1]
			initSafeGuards = append(initSafeGuards, e.builderTemplate.GetSafeGuard(subPath))
		}
	}

	var constraints []Constraint
	if assignment.Value.Argument != nil {
		argName := e.builderTemplate.GetEscapeVar(tools.LowerCamelCase(assignment.Value.Argument.Name))
		constraints = e.constraints(argName, assignment.Constraints)
	}

	return Assignment{
		Path:           assignment.Path,
		InitSafeguards: initSafeGuards,
		Constraints:    constraints,
		Method:         assignment.Method,
		Value:          assignment.Value,
	}
}

func (e *Executor) constraints(argumentName string, constraints []ast.TypeConstraint) []Constraint {
	return tools.Map(constraints, func(constraint ast.TypeConstraint) Constraint {
		return Constraint{
			ArgName:   argumentName,
			Op:        constraint.Op,
			Parameter: constraint.Args[0],
		}
	})
}

func (e *Executor) formatType(def ast.Type) string {
	return e.typeFormatter.FormatType(def, true)
}

func (e *Executor) formatScalar(val any) string {
	if list, ok := val.([]any); ok {
		items := make([]string, 0, len(list))

		for _, item := range list {
			items = append(items, e.formatScalar(item))
		}

		// FIXME: this is wrong, we can't just assume a list of strings.
		return fmt.Sprintf(e.builderTemplate.ScalarPattern(), strings.Join(items, ", "))
	}

	return fmt.Sprintf("%#v", val)
}
