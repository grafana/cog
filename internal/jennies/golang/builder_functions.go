package golang

import (
	"fmt"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/template"
	"github.com/grafana/cog/internal/tools"
	"strings"
)

type GoTemplate struct {
	formatter template.TypeFormatter
}

func (g GoTemplate) GetObjectName(builder ast.Builder) string {
	return g.importType(builder.For.SelfRef)
}

func (g GoTemplate) GetBuilderSignature(builder ast.Builder) string {
	buildObjectSignature := g.importType(builder.For.SelfRef)
	if builder.For.Type.ImplementsVariant() {
		buildObjectSignature = g.formatter.VariantInterface(builder.For.Type.ImplementedVariant())
	}

	return buildObjectSignature
}

func (g GoTemplate) GetSafeGuard(path ast.Path) string {
	fieldPath := g.FormatFieldPath(path)
	valueType := path.Last().Type
	if path.Last().TypeHint != nil {
		valueType = *path.Last().TypeHint
	}

	emptyValue := g.emptyValueForType(valueType)
	// This should be alright since there shouldn't be any scalar in the middle of a path
	if emptyValue[0] == '*' {
		emptyValue = "&" + emptyValue[1:]
	}

	if path.Last().Type.IsAny() && emptyValue[0] != '&' {
		emptyValue = "&" + emptyValue
	}

	return fmt.Sprintf(`if builder.internal.%[1]s == nil {
	builder.internal.%[1]s = %[2]s
}`, fieldPath, emptyValue)
}

func (g GoTemplate) GetEscapeVar(varName string) string {
	return escapeVarName(varName)
}

func (g GoTemplate) FormatScalar(val any) string {
	if list, ok := val.([]any); ok {
		items := make([]string, 0, len(list))

		for _, item := range list {
			items = append(items, formatScalar(item))
		}

		// FIXME: this is wrong, we can't just assume a list of strings.
		return fmt.Sprintf("[]string{%s}", strings.Join(items, ", "))
	}

	return fmt.Sprintf("%#v", val)
}

func (g GoTemplate) importType(typeRef ast.RefType) string {
	pkg := g.formatter.PackageMapper(typeRef.ReferredPkg)
	typeName := tools.UpperCamelCase(typeRef.ReferredType)
	if pkg == "" {
		return typeName
	}

	return fmt.Sprintf("%s.%s", pkg, typeName)
}

func (g GoTemplate) FormatFieldPath(fieldPath ast.Path) string {
	parts := make([]string, len(fieldPath))

	for i := range fieldPath {
		output := tools.UpperCamelCase(fieldPath[i].Identifier)

		// don't generate type hints if:
		// * there isn't one defined
		// * the type isn't "any"
		// * as a trailing element in the path
		if !fieldPath[i].Type.IsAny() || fieldPath[i].TypeHint == nil || i == len(fieldPath)-1 {
			parts[i] = output
			continue
		}

		formattedTypeHint := g.formatter.FormatType(*fieldPath[i].TypeHint, true)
		parts[i] = output + fmt.Sprintf(".(*%s)", formattedTypeHint)
	}

	return strings.Join(parts, ".")
}

func (g GoTemplate) ScalarPattern() string {
	return "[]string{%s}"
}

func (g GoTemplate) emptyValueForType(typeDef ast.Type) string {
	switch typeDef.Kind {
	case ast.KindRef:
		return g.formatter.FormatType(typeDef, true) + "{}"
	case ast.KindStruct:
		return g.formatter.FormatType(typeDef, true) + "{}"
	case ast.KindEnum:
		return formatScalar(typeDef.AsEnum().Values[0].Value)
	case ast.KindArray, ast.KindMap:
		return g.formatter.FormatType(typeDef, true) + "{}"
	case ast.KindScalar:
		return "" // no need to do anything here

	default:
		return "unknown"
	}
}
