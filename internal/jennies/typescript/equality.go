package typescript

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

// generateForObject generates a TypeScript equality function for struct objects.
// Since TypeScript uses interfaces (not classes), equality is a standalone exported function.
func (jenny equalityMethods) generateForObject(context languages.Context, object ast.Object) string {
	if !object.Type.IsStruct() {
		return ""
	}

	objectName := formatObjectName(object.Name)
	var buffer strings.Builder

	funcName := jenny.equalsFuncNameForObject(objectName)
	buffer.WriteString(fmt.Sprintf("// %s tests the equality of two `%s` objects.\n", funcName, objectName))
	buffer.WriteString(fmt.Sprintf("export const %s = (a: %s, b: %s): boolean => {\n", funcName, objectName, objectName))

	for _, field := range object.Type.AsStruct().Fields {
		fieldName := formatIdentifier(field.Name)
		selfExpr := "a." + fieldName
		otherExpr := "b." + fieldName

		if !field.Required {
			buffer.WriteString(fmt.Sprintf("\tif ((%s === undefined) !== (%s === undefined)) return false;\n", selfExpr, otherExpr))
			buffer.WriteString(fmt.Sprintf("\tif (%s !== undefined) {\n", selfExpr))
			jenny.writeTypeEquality(&buffer, context, field.Type, selfExpr, otherExpr+"!", 2)
			buffer.WriteString("\t}\n")
		} else {
			jenny.writeTypeEquality(&buffer, context, field.Type, selfExpr, otherExpr, 1)
		}
	}

	buffer.WriteString("\treturn true;\n")
	buffer.WriteString("};\n")

	return buffer.String()
}

func (jenny equalityMethods) writeTypeEquality(
	buffer *strings.Builder,
	context languages.Context,
	typeDef ast.Type,
	selfExpr, otherExpr string,
	depth int,
) {
	indent := strings.Repeat("\t", depth)

	if typeDef.IsAny() {
		buffer.WriteString(fmt.Sprintf("%sif (JSON.stringify(%s) !== JSON.stringify(%s)) return false;\n", indent, selfExpr, otherExpr))
		return
	}

	if context.ResolveToStruct(typeDef) {
		funcName := jenny.equalsFuncName(typeDef)
		buffer.WriteString(fmt.Sprintf("%sif (!%s(%s, %s)) return false;\n", indent, funcName, selfExpr, otherExpr))
		return
	}

	resolved := context.ResolveRefs(typeDef)

	if resolved.IsArray() {
		loopVar := fmt.Sprintf("i%d", depth)
		buffer.WriteString(fmt.Sprintf("%sif (%s.length !== %s.length) return false;\n", indent, selfExpr, otherExpr))
		buffer.WriteString(fmt.Sprintf("%sfor (let %s = 0; %s < %s.length; %s++) {\n", indent, loopVar, loopVar, selfExpr, loopVar))
		jenny.writeTypeEquality(buffer, context, resolved.Array.ValueType, fmt.Sprintf("%s[%s]", selfExpr, loopVar), fmt.Sprintf("%s[%s]", otherExpr, loopVar), depth+1)
		buffer.WriteString(fmt.Sprintf("%s}\n", indent))
		return
	}

	if resolved.IsMap() {
		loopVar := fmt.Sprintf("key%d", depth)
		buffer.WriteString(fmt.Sprintf("%sif (Object.keys(%s).length !== Object.keys(%s).length) return false;\n", indent, selfExpr, otherExpr))
		buffer.WriteString(fmt.Sprintf("%sfor (const %s in %s) {\n", indent, loopVar, selfExpr))
		jenny.writeTypeEquality(buffer, context, resolved.Map.ValueType, fmt.Sprintf("%s[%s]", selfExpr, loopVar), fmt.Sprintf("%s[%s]", otherExpr, loopVar), depth+1)
		buffer.WriteString(fmt.Sprintf("%s}\n", indent))
		return
	}

	buffer.WriteString(fmt.Sprintf("%sif (%s !== %s) return false;\n", indent, selfExpr, otherExpr))
}

func (jenny equalityMethods) equalsFuncNameForObject(name string) string {
	return "equals" + formatObjectName(name)
}

func (jenny equalityMethods) equalsFuncName(typeDef ast.Type) string {
	if typeDef.IsRef() {
		return jenny.equalsFuncNameForObject(typeDef.AsRef().ReferredType)
	}
	return "equalsUnknown"
}
