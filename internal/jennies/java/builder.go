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
		output, err := jenny.genFilesForBuilder(builder)
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

func (jenny Builder) genFilesForBuilder(builder ast.Builder) ([]byte, error) {
	var buffer strings.Builder

	jenny.imports = NewImportMap()
	jenny.imports.Add("Builder", "cog.variants")
	jenny.typeFormatter = jenny.typeFormatter.withPackageMapper(func(pkg string, class string) string {
		if builder.Package == pkg {
			return ""
		}

		return jenny.imports.Add(class, pkg)
	})

	object, _ := jenny.typeFormatter.locateObject(builder.For.SelfRef.ReferredPkg, builder.For.SelfRef.ReferredType)
	options := make([]Option, len(builder.Options))
	for i, option := range builder.Options {
		options[i] = Option{
			Name:        tools.UpperCamelCase(option.Name),
			Args:        option.Args,
			Assignments: jenny.genAssignments(option.Assignments),
		}
	}

	err := templates.Funcs(map[string]any{
		"formatType": jenny.typeFormatter.formatBuilderArgs,
	}).ExecuteTemplate(&buffer, "builders/builder.tmpl", BuilderTemplate{
		Package: builder.Package,
		Imports: jenny.imports,
		Name:    builder.Name,
		Options: options,
		Fields:  jenny.genFields(object),
	})

	if err != nil {
		return nil, err
	}

	return []byte(buffer.String()), nil
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
			Path:        assignment.Path,
			Method:      assignment.Method,
			Constraints: constraints,
			Value:       assignment.Value,
		}
	}

	return assign
}

func (jenny Builder) genFields(object ast.Object) []Field {
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
		obj, ok := jenny.typeFormatter.locateObject(ref.ReferredPkg, ref.ReferredType)
		if !ok {
			break
		}
		fields = append(fields, jenny.genFields(obj)...)
	default:
		// Shouldn't reach here...
		return nil
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
