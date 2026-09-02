package builders

import (
	"github.com/grafana/cog/internal/ir"
	"github.com/grafana/cog/internal/languages"
)

func GenerateNilChecks(nullableKinds languages.NullableConfig, schemas ir.Schemas, builders ir.Builders) (ir.Builders, error) {
	// Allows us to keep track of the checks already performed for the current scope (constructor or option)
	// When a check is generated, the path being checked is stored in this map.
	// Changes in scope must reset this map.
	checks := make(map[string]struct{})

	nilChecksVisitor := ir.BuilderVisitor{
		OnConstructor: func(visitor *ir.BuilderVisitor, schemas ir.Schemas, builder ir.Builder, constructor ir.Constructor) (ir.Constructor, error) {
			checks = make(map[string]struct{})

			return visitor.TraverseConstructor(schemas, builder, constructor)
		},
		OnOption: func(visitor *ir.BuilderVisitor, schemas ir.Schemas, builder ir.Builder, option ir.Option) (ir.Option, error) {
			checks = make(map[string]struct{})

			return visitor.TraverseOption(schemas, builder, option)
		},
		OnAssignment: func(_ *ir.BuilderVisitor, _ ir.Schemas, b ir.Builder, assignment ir.Assignment) (ir.Assignment, error) {
			for i, chunk := range assignment.Path {
				protectArrayAppend := nullableKinds.ProtectArrayAppend && assignment.Method == ir.AppendAssignment
				if i == len(assignment.Path)-1 && !protectArrayAppend {
					continue
				}

				if nullableKinds.TypeIsNullable(chunk.Type) {
					subPath := assignment.Path[:i+1]
					valueType := subPath.Last().Type
					if subPath.Last().TypeHint != nil {
						valueType = *subPath.Last().TypeHint
					}

					// this path already has a nil check: nothing to do.
					if _, found := checks[subPath.String()]; found {
						continue
					}

					assignment.NilChecks = append(assignment.NilChecks, ir.AssignmentNilCheck{
						Path:           subPath,
						EmptyValueType: valueType,
					})
					checks[subPath.String()] = struct{}{}
				}
			}

			return assignment, nil
		},
	}

	return nilChecksVisitor.Visit(schemas, builders)
}
