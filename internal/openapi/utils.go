package openapi

import (
	"fmt"
	"strconv"
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

func getEnumType(t string) (ast.Type, error) {
	switch t {
	case openapi3.TypeString:
		return ast.String(), nil
	case openapi3.TypeNumber:
		return ast.NewScalar(ast.KindInt32), nil
	case openapi3.TypeInteger:
		return ast.NewScalar(ast.KindInt64), nil
	case openapi3.TypeBoolean:
		return ast.Bool(), nil
	default:
		// TODO: Handle it correctly
		return ast.String(), nil
	}
}

func parseValue(value interface{}) string {
	if val, ok := value.(string); ok {
		return val
	}
	if val, ok := value.(bool); ok {
		return strconv.FormatBool(val)
	}
	if val, ok := value.(int); ok {
		return strconv.Itoa(val)
	}
	if val, ok := value.(float64); ok {
		return fmt.Sprintf("%f", val)
	}

	return ""
}

func getConstraints(schema *openapi3.Schema) []ast.TypeConstraint {
	constraints := make([]ast.TypeConstraint, 0)

	if schema.MinLength > 0 {
		constraints = append(constraints, ast.TypeConstraint{
			Op:   "minLength",
			Args: []any{schema.MinLength},
		})
	}
	if schema.MaxLength != nil {
		constraints = append(constraints, ast.TypeConstraint{
			Op:   "maxLength",
			Args: []any{schema.MaxLength},
		})
	}

	if schema.MultipleOf != nil {
		constraints = append(constraints, ast.TypeConstraint{
			Op:   "multipleOf",
			Args: []any{schema.MultipleOf},
		})
	}

	if schema.Min != nil {
		op := ">="
		if schema.ExclusiveMin {
			op = ">"
		}
		constraints = append(constraints, ast.TypeConstraint{
			Op:   op,
			Args: []any{schema.Min},
		})
	}

	if schema.Max != nil {
		op := "<="
		if schema.ExclusiveMax {
			op = "<"
		}
		constraints = append(constraints, ast.TypeConstraint{
			Op:   op,
			Args: []any{schema.Max},
		})
	}

	return constraints
}
