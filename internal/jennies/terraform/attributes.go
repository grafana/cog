package terraform

import (
	"fmt"
	"strings"

	"github.com/grafana/cog/internal/ast"
)

type attributes struct {
	packageMapper func(pkg string) string
	typeFormatter *typeFormatter
}

func newAttributesGenerator(typeFormatter *typeFormatter, packageMapper func(pkg string) string) *attributes {
	return &attributes{
		packageMapper: packageMapper,
		typeFormatter: typeFormatter,
	}
}

func (a *attributes) generateForSchema(schema *ast.Schema) (string, error) {
	var buffer strings.Builder

	a.packageMapper("github.com/hashicorp/terraform-plugin-framework/resource/schema")
	buffer.WriteString("var SpecAttributes = map[string]schema.Attribute{\n")

	schema.Objects.Iterate(func(_ string, obj ast.Object) {
		if !obj.Type.IsAnyOf(ast.KindDisjunction, ast.KindRef, ast.KindConstantRef, ast.KindEnum, ast.KindIntersection) && !obj.Type.IsDisjunctionOfAnyKind() {
			buffer.WriteString(fmt.Sprintf("\"%s\": %s", strings.ToLower(obj.Name), a.typeFormatter.formatTypeAttribute(obj)))
		}
	})

	buffer.WriteString("}")
	return buffer.String(), nil
}
