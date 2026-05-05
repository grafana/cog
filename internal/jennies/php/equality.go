package php

import (
	"fmt"
	"strings"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/languages"
)

type equalityMethods struct{}

func newEqualityMethods() equalityMethods {
	return equalityMethods{}
}

// generateForObject generates a PHP equals() method for struct objects.
func (jenny equalityMethods) generateForObject(context languages.Context, object ast.Object) string {
	if !object.Type.IsStruct() {
		return ""
	}

	var buffer strings.Builder

	buffer.WriteString("public function equals(mixed $other): bool\n")
	buffer.WriteString("{\n")
	buffer.WriteString("    if (!($other instanceof self)) {\n")
	buffer.WriteString("        return false;\n")
	buffer.WriteString("    }\n")

	for _, field := range object.Type.AsStruct().Fields {
		fieldName := formatFieldName(field.Name)
		selfExpr := "$this->" + fieldName
		otherExpr := "$other->" + fieldName

		buffer.WriteString("\n")
		buffer.WriteString(jenny.compareField(context, field.Type, selfExpr, otherExpr))
	}

	buffer.WriteString("\n    return true;\n")
	buffer.WriteString("}")

	return buffer.String()
}

func (jenny equalityMethods) compareField(context languages.Context, typeDef ast.Type, selfExpr, otherExpr string) string {
	if typeDef.Nullable {
		return jenny.compareNullableField(context, typeDef, selfExpr, otherExpr)
	}
	return jenny.compareNonNullField(context, typeDef, selfExpr, otherExpr)
}

func (jenny equalityMethods) compareNullableField(context languages.Context, typeDef ast.Type, selfExpr, otherExpr string) string {
	var buffer strings.Builder

	buffer.WriteString(fmt.Sprintf("    if ((%s === null) !== (%s === null)) {\n", selfExpr, otherExpr))
	buffer.WriteString("        return false;\n")
	buffer.WriteString("    }\n")
	buffer.WriteString(fmt.Sprintf("    if (%s !== null) {\n", selfExpr))

	inner := jenny.compareNonNullField(context, typeDef, selfExpr, otherExpr)
	for _, line := range strings.Split(strings.TrimRight(inner, "\n"), "\n") {
		buffer.WriteString("    " + line + "\n")
	}

	buffer.WriteString("    }\n")
	return buffer.String()
}

func (jenny equalityMethods) compareNonNullField(context languages.Context, typeDef ast.Type, selfExpr, otherExpr string) string {
	if context.ResolveToStruct(typeDef) {
		return fmt.Sprintf("    if (!%s->equals(%s)) {\n        return false;\n    }\n", selfExpr, otherExpr)
	}

	resolved := context.ResolveRefs(typeDef)
	if resolved.IsArray() || resolved.IsMap() {
		// PHP == for arrays/maps does deep recursive comparison including object properties
		return fmt.Sprintf("    if (%s != %s) {\n        return false;\n    }\n", selfExpr, otherExpr)
	}

	return fmt.Sprintf("    if (%s !== %s) {\n        return false;\n    }\n", selfExpr, otherExpr)
}
