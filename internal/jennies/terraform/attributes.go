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
	validator     *validatorGenerator
}

func newAttributesGenerator(cfg Config, typeFormatter *typeFormatter, packageMapper func(pkg string) string, validator *validatorGenerator) *attributes {
	return &attributes{
		cfg:           cfg,
		packageMapper: packageMapper,
		typeFormatter: typeFormatter,
		validator:     validator,
	}
}

func (a *attributes) generateForSchema(schema *ast.Schema) (string, error) {
	var buffer strings.Builder

	a.packageMapper("github.com/hashicorp/terraform-plugin-framework/resource/schema")
	buffer.WriteString("var SpecAttributes = map[string]schema.Attribute{\n")

	var validatorBuffer strings.Builder
	var innerErr error

	schema.Objects.Iterate(func(_ string, obj ast.Object) {
		if !obj.Type.IsAnyOf(ast.KindDisjunction, ast.KindRef, ast.KindConstantRef, ast.KindEnum, ast.KindIntersection) && !obj.Type.IsDisjunctionOfAnyKind() {
			buffer.WriteString(fmt.Sprintf("\"%s\": %s", tools.SnakeCase(obj.Name), a.typeFormatter.formatTypeAttribute(obj.Type, formatComments(obj.Comments))))

			validator, err := a.validator.generateValidatorIfNeeded(obj.Name, obj.Type)
			if err != nil {
				innerErr = err
			}

			validatorBuffer.WriteString(validator)
		}
	})

	if innerErr != nil {
		return "", innerErr
	}

	buffer.WriteString("}\n")
	buffer.WriteString(validatorBuffer.String())
	return buffer.String(), nil
}

func (a *attributes) generateForObject(obj ast.Object) (string, error) {
	if !obj.Type.IsStruct() {
		return "", fmt.Errorf("cannot generate attributes for non-struct type %s", obj.Type.Kind)
	}

	var buffer strings.Builder
	var validatorBuffer strings.Builder

	a.packageMapper("github.com/hashicorp/terraform-plugin-framework/resource/schema")
	buffer.WriteString(fmt.Sprintf("var %sSpecAttributes = map[string]schema.Attribute{\n", tools.UpperCamelCase(a.cfg.PrefixAttributeSpec)))

	for _, field := range obj.Type.AsStruct().Fields {
		buffer.WriteString(fmt.Sprintf("\"%s\": %s", tools.SnakeCase(field.Name), a.typeFormatter.formatTypeAttribute(field.Type, formatComments(obj.Comments))))
		validator, err := a.validator.generateValidatorIfNeeded(field.Name, field.Type)
		if err != nil {
			return "", err
		}

		validatorBuffer.WriteString(validator)
	}

	buffer.WriteString("}\n")
	buffer.WriteString(validatorBuffer.String())
	return buffer.String(), nil
}

func formatComments(objectComments []string) string {
	comments := ""
	if len(objectComments) > 0 {
		comments += "`\n"

		for _, comment := range objectComments {
			comments += comment + "\n"
		}
		comments += "`"
	}
	return comments
}
