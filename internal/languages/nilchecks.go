package languages

import (
	"github.com/grafana/cog/internal/ast"
)

func GenerateBuilderNilChecks(language Language, context Context) (Context, error) {
	var err error
	nullableKinds := NullableConfig{
		Kinds:              nil,
		ProtectArrayAppend: false,
		AnyIsNullable:      true,
	}
	if nilTypesProvider, ok := language.(NullableKindsProvider); ok {
		nullableKinds = nilTypesProvider.NullableKinds()
	}

	nilChecksVisitor := ast.BuilderVisitor{
		OnAssignment: func(_ *ast.BuilderVisitor, _ ast.Schemas, _ ast.Builder, assignment ast.Assignment) (ast.Assignment, error) {
			for i, chunk := range assignment.Path {
				protectArrayAppend := nullableKinds.ProtectArrayAppend && assignment.Method == ast.AppendAssignment
				if i == len(assignment.Path)-1 && !protectArrayAppend {
					continue
				}

				nullable := chunk.Type.Nullable ||
					(chunk.Type.IsAny() && nullableKinds.AnyIsNullable) ||
					chunk.Type.IsAnyOf(nullableKinds.Kinds...)
				if nullable {
					subPath := assignment.Path[:i+1]
					valueType := subPath.Last().Type
					if subPath.Last().TypeHint != nil {
						valueType = *subPath.Last().TypeHint
					}

					assignment.NilChecks = append(assignment.NilChecks, ast.AssignmentNilCheck{
						Path:           subPath,
						EmptyValueType: valueType,
					})
				}
			}

			return assignment, nil
		},
	}
	context.Builders, err = nilChecksVisitor.Visit(context.Schemas, context.Builders)
	if err != nil {
		return context, err
	}

	return context, nil
}
