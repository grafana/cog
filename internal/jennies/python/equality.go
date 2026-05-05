package python

import (
	"fmt"
	"strings"

	"github.com/grafana/cog/internal/ast"
)

type equalityMethods struct{}

func newEqualityMethods() equalityMethods {
	return equalityMethods{}
}

// generateForObject generates a Python __eq__ method for struct objects.
// Python's native != operator handles None, lists, dicts, and custom objects
// with __eq__ defined, so a simple field-by-field != comparison is sufficient.
func (jenny equalityMethods) generateForObject(object ast.Object) string {
	if !object.Type.IsStruct() {
		return ""
	}

	objectName := formatObjectName(object.Name)
	var buffer strings.Builder

	buffer.WriteString(fmt.Sprintf("    def __eq__(self, other: object) -> bool:\n"))
	buffer.WriteString(fmt.Sprintf("        if not isinstance(other, %s):\n", objectName))
	buffer.WriteString("            return False\n")

	for _, field := range object.Type.AsStruct().Fields {
		fieldName := formatIdentifier(field.Name)
		buffer.WriteString(fmt.Sprintf("        if self.%s != other.%s:\n", fieldName, fieldName))
		buffer.WriteString("            return False\n")
	}

	buffer.WriteString("        return True")

	return buffer.String()
}
