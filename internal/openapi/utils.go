package openapi

import (
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/grafana/cog/internal/ast"
)

func schemaComments(schema *openapi3.Schema) []string {
	lines := strings.Split(schema.Description, "\n")
	filtered := make([]string, 0, len(lines))

	for _, line := range lines {
		if line == "" {
			continue
		}

		filtered = append(filtered, line)
	}

	return filtered
}

func getEnumType(t string) ast.Type {
	switch t {
	case openapi3.TypeString:
		return ast.String()
	case openapi3.TypeNumber:
		return ast.NewScalar(ast.KindInt32)
	case openapi3.TypeInteger:
		return ast.NewScalar(ast.KindInt64)
	case openapi3.TypeBoolean:
		return ast.Bool()
	default:
		// TODO: Handle it correctly
		return ast.String()
	}
}

func getConstraints(schema *openapi3.Schema) []ast.TypeConstraint {
	constraints := make([]ast.TypeConstraint, 0)

	if schema.MinLength > 0 {
		constraints = append(constraints, ast.TypeConstraint{
			Op:   ast.MinLengthOp,
			Args: []any{schema.MinLength},
		})
	}
	if schema.MaxLength != nil {
		constraints = append(constraints, ast.TypeConstraint{
			Op:   ast.MaxLengthOp,
			Args: []any{schema.MaxLength},
		})
	}

	if schema.MultipleOf != nil {
		constraints = append(constraints, ast.TypeConstraint{
			Op:   ast.MultipleOfOp,
			Args: []any{schema.MultipleOf},
		})
	}

	if schema.Min != nil {
		op := ast.GreaterThanEqualOp
		if schema.ExclusiveMin {
			op = ast.GreaterThanOp
		}
		constraints = append(constraints, ast.TypeConstraint{
			Op:   op,
			Args: []any{schema.Min},
		})
	}

	if schema.Max != nil {
		op := ast.LessThanEqualOp
		if schema.ExclusiveMax {
			op = ast.LessThanOp
		}
		constraints = append(constraints, ast.TypeConstraint{
			Op:   op,
			Args: []any{schema.Max},
		})
	}

	return constraints
}
