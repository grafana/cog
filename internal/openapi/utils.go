package openapi

import (
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/grafana/cog/internal/ast"
	"strconv"
	"strings"
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
