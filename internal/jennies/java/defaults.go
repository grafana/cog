package java

import (
	"fmt"
	"strings"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/tools"
)

func arrayDefaults(t ast.Type, v any) string {
	if t.AsArray().IsArrayOfScalars() {
		elements := make([]string, len(v.([]any)))
		for i, value := range v.([]any) {
			elements[i] = formatType(t.AsArray().ValueType, value)
		}
		return fmt.Sprintf("List.of(%s)", strings.Join(elements, ","))
	}

	// TODO: Rest of types

	return ""
}

func enumDefaults(name string, t ast.Type, v any) string {
	for _, value := range t.AsEnum().Values {
		if v == value.Value {
			return fmt.Sprintf("%s.%s", name, tools.UpperSnakeCase(value.Name))
		}
	}

	return ""
}

func structDefaults(name string, _ ast.Type, _ any) string {
	return fmt.Sprintf("%sResource", tools.LowerCamelCase(name))
}
