package java

import (
	"fmt"
	"strings"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/tools"
)

type equalityMethods struct{}

func newEqualityMethods() equalityMethods {
	return equalityMethods{}
}

// generateForObject generates Java equals() and hashCode() overrides for struct objects.
// Uses Objects.equals() which handles null checks and delegates to each field's equals().
func (jenny equalityMethods) generateForObject(object ast.Object) string {
	if !object.Type.IsStruct() {
		return ""
	}

	className := tools.UpperCamelCase(object.Name)
	fields := object.Type.AsStruct().Fields

	var buffer strings.Builder

	buffer.WriteString("\n\n    @Override\n")
	buffer.WriteString("    public boolean equals(Object other) {\n")
	buffer.WriteString("        if (this == other) return true;\n")
	buffer.WriteString(fmt.Sprintf("        if (!(other instanceof %s)) return false;\n", className))
	buffer.WriteString(fmt.Sprintf("        %s o = (%s) other;\n", className, className))

	for _, field := range fields {
		fieldName := formatFieldName(field.Name)
		buffer.WriteString(fmt.Sprintf("        if (!Objects.equals(this.%s, o.%s)) return false;\n", fieldName, fieldName))
	}

	buffer.WriteString("        return true;\n")
	buffer.WriteString("    }\n")

	buffer.WriteString("\n    @Override\n")
	buffer.WriteString("    public int hashCode() {\n")

	if len(fields) == 0 {
		buffer.WriteString("        return 0;\n")
	} else {
		fieldExprs := make([]string, 0, len(fields))
		for _, field := range fields {
			fieldExprs = append(fieldExprs, "this."+formatFieldName(field.Name))
		}
		buffer.WriteString(fmt.Sprintf("        return Objects.hash(%s);\n", strings.Join(fieldExprs, ", ")))
	}

	buffer.WriteString("    }")

	return buffer.String()
}
