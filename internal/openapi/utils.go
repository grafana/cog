package openapi

import (
	"errors"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/grafana/cog/internal/ast"
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

func getEnumType(format string) (ast.Type, error) {
	switch format {
	case FormatString:
		return ast.String(), nil
	case FormatInt32:
		return ast.NewScalar(ast.KindInt32), nil
	case FormatInt64:
		return ast.NewScalar(ast.KindInt64), nil
	default:
		return ast.Type{}, errors.New(fmt.Sprintf("Unhandled enum format: %s. Valid formats are string or integers", format))
	}
}
