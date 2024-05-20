package languages

import (
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/common"
)

func GenerateBuilderNilChecks(language Language, context common.Context) (common.Context, error) {
	var err error
	var languageSpecificNilTypes []ast.Kind
	if nilTypesProvider, ok := language.(NullableKindsProvider); ok {
		languageSpecificNilTypes = nilTypesProvider.NullableKinds()
	}

	nilChecksVisitor := ast.BuilderVisitor{
		OnAssignment: func(_ *ast.BuilderVisitor, _ ast.Schemas, _ ast.Builder, assignment ast.Assignment) (ast.Assignment, error) {
			for i, chunk := range assignment.Path {
				if i == len(assignment.Path)-1 {
					continue
				}

				nullable := chunk.Type.Nullable ||
					chunk.Type.IsAny() || // this assumes that "any" is nullable in every language
					chunk.Type.IsAnyOf(languageSpecificNilTypes...)
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
