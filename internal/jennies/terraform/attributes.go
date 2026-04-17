package terraform

import (
	"fmt"
	"strings"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/tools"
)

type attributes struct {
	cfg           Config
	packageMapper func(pkg string) string
	typeFormatter *typeFormatter
}

func newAttributesGenerator(cfg Config, typeFormatter *typeFormatter, packageMapper func(pkg string) string) *attributes {
	return &attributes{
		cfg:           cfg,
		packageMapper: packageMapper,
		typeFormatter: typeFormatter,
	}
}

func attributesVarName(name string) string {
	return tools.UpperCamelCase(name) + "Attributes"
}

// generateForAllObjects generates a named `var XxxAttributes` for each struct object in the schema,
// skipping the object named skipObjectName (used to avoid duplicating the entry point).
func (a *attributes) generateForAllObjects(schema *ast.Schema, skipObjectName string) (string, error) {
	var buffer strings.Builder

	a.packageMapper("github.com/hashicorp/terraform-plugin-framework/resource/schema")

	schema.Objects.Iterate(func(_ string, obj ast.Object) {
		if obj.Name == skipObjectName {
			return
		}
		if !obj.Type.IsStruct() {
			return
		}

		buffer.WriteString(fmt.Sprintf("var %s = map[string]schema.Attribute{\n", attributesVarName(obj.Name)))
		buffer.WriteString(a.typeFormatter.formatFieldAttributes(obj.Type.AsStruct().Fields))
		buffer.WriteString("}\n\n")
	})

	return buffer.String(), nil
}

func (a *attributes) generateForSchema(schema *ast.Schema) (string, error) {
	var buffer strings.Builder

	a.packageMapper("github.com/hashicorp/terraform-plugin-framework/resource/schema")
	buffer.WriteString("var SpecAttributes = map[string]schema.Attribute{\n")

	schema.Objects.Iterate(func(_ string, obj ast.Object) {
		if obj.Type.IsAnyOf(ast.KindDisjunction, ast.KindRef, ast.KindConstantRef, ast.KindEnum, ast.KindIntersection) || obj.Type.IsDisjunctionOfAnyKind() {
			return
		}
		if obj.Type.IsStruct() {
			required := "Required: true,"
			if obj.Type.Nullable {
				required = "Optional: true,"
			}
			comments := ""
			if c := formatComments(obj.Comments); c != "" {
				comments = fmt.Sprintf("Description: %s,\n", c)
			}
			buffer.WriteString(fmt.Sprintf("\"%s\": schema.SingleNestedAttribute{\n%s\n%sAttributes: %s,\n},\n",
				tools.SnakeCase(obj.Name), required, comments, attributesVarName(obj.Name)))
		} else {
			buffer.WriteString(a.typeFormatter.formatTypeAttributeForObject(obj, formatComments(obj.Comments)))
		}
	})

	buffer.WriteString("}")
	return buffer.String(), nil
}

func (a *attributes) generateForObject(obj ast.Object) (string, error) {
	if !obj.Type.IsStruct() {
		return "", fmt.Errorf("cannot generate attributes for non-struct type %s", obj.Type.Kind)
	}

	var buffer strings.Builder

	a.packageMapper("github.com/hashicorp/terraform-plugin-framework/resource/schema")
	buffer.WriteString(fmt.Sprintf("var %sSpecAttributes = map[string]schema.Attribute{\n", tools.UpperCamelCase(a.cfg.PrefixAttributeSpec)))

	for _, field := range obj.Type.AsStruct().Fields {
		buffer.WriteString(fmt.Sprintf("\"%s\": %s", tools.SnakeCase(field.Name), a.typeFormatter.formatTypeAttribute(field.Type, formatComments(obj.Comments))))
	}

	buffer.WriteString("}")
	return buffer.String(), nil
}
