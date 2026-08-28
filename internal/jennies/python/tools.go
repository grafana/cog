package python

import (
	"fmt"
	"strings"

	"github.com/grafana/cog/internal/orderedmap"
	"github.com/grafana/cog/internal/tools"
	"github.com/grafana/cog/pkg/ir"
)

type raw string

func formatValue(val any) string {
	if val == nil {
		return "None"
	}

	if rawVal, ok := val.(raw); ok {
		return string(rawVal)
	}

	if asBool, ok := val.(bool); ok {
		if asBool {
			return "True"
		}

		return "False"
	}

	if list, ok := val.([]any); ok {
		items := make([]string, 0, len(list))

		for _, item := range list {
			items = append(items, formatValue(item))
		}

		return fmt.Sprintf("[%s]", strings.Join(items, ", "))
	}

	if mapVal, ok := val.(map[string]any); ok {
		items := make([]string, 0, len(mapVal))

		for key, value := range mapVal {
			items = append(items, fmt.Sprintf("%q: %s", key, formatValue(value)))
		}

		return fmt.Sprintf("{%s}", strings.Join(items, ", "))
	}

	if orderedMap, ok := val.(*orderedmap.Map[string, any]); ok {
		items := make([]string, 0, orderedMap.Len())

		orderedMap.Iterate(func(key string, value any) {
			items = append(items, fmt.Sprintf("%q: %s", key, formatValue(value)))
		})

		return fmt.Sprintf("{%s}", strings.Join(items, ", "))
	}

	return fmt.Sprintf("%#v", val)
}

func formatFieldPath(fieldPath ir.Path) string {
	path := ""

	for i, chunk := range fieldPath {
		last := i == len(fieldPath)-1
		output := formatIdentifier(chunk.Identifier)

		if chunk.Index != nil {
			output += "["
			if chunk.Index.Constant != nil {
				output += formatValue(chunk.Index.Constant)
			} else {
				output += formatIdentifier(chunk.Index.Argument.Name)
			}
			output += "]"
		}

		path += output
		if !last && fieldPath[i+1].Index == nil {
			path += "."
		}
	}

	return path
}

func formatObjectName(name string) string {
	return tools.UpperCamelCase(name)
}

func formatIdentifier(name string) string {
	name = strings.TrimLeft(name, "$_")
	return tools.SnakeCase(escapeIdentifier(name))
}

func formatFunctionName(name string) string {
	name = strings.TrimLeft(name, "$_")
	return tools.SnakeCase(escapeFunctionName(name))
}

func escapeIdentifier(name string) string {
	if isReservedPythonKeyword(name) || isBuiltInFunction(name) {
		return name + "_val"
	}

	return name
}

func escapeFunctionName(name string) string {
	if isReservedPythonKeyword(name) {
		return name + "_val"
	}

	return name
}

func isBuiltInFunction(input string) bool {
	// see: https://docs.python.org/3/library/functions.html
	switch input {
	case "abs", "aiter", "all", "anext", "any", "ascii", "bin", "bool", "breakpoint", "bytearray",
		"bytes", "callable", "chr", "classmethod", "compile", "complex", "delattr", "dict", "dir",
		"divmod", "enumerate", "eval", "exec", "filter", "float", "format", "frozenset", "getattr",
		"globals", "hasattr", "hash", "help", "hex", "id", "input", "int", "isinstance",
		"issubclass", "iter", "len", "list", "locals", "map", "max", "memoryview", "min", "next",
		"object", "oct", "open", "ord", "pow", "print", "property", "range", "repr", "reversed",
		"round", "set", "setattr", "slice", "sorted", "staticmethod", "str", "sum", "super",
		"tuple", "type", "vars", "zip", "__import__":
		return true
	default:
		return false
	}
}

func isReservedPythonKeyword(input string) bool {
	// see: https://docs.python.org/3/reference/lexical_analysis.html#keywords
	switch input {
	case "False", "await", "else", "import", "pass", "None", "break", "except", "in", "raise",
		"True", "class", "finally", "is", "return", "and", "continue", "for", "lambda", "try",
		"as", "def", "from", "nonlocal", "while", "assert", "del", "global", "not", "with",
		"async", "elif", "if", "or", "yield":
		return true

	default:
		return false
	}
}

/******************************************************************************
* 					 Default and "empty" values management 					  *
******************************************************************************/

func defaultValueForType(schemas ir.Schemas, typeDef ir.Type, importModule moduleImporter, defaultsOverrides *orderedmap.Map[string, any]) any {
	if !typeDef.IsRef() && typeDef.Default != nil {
		return typeDef.Default
	}

	switch typeDef.Kind {
	case ir.KindDisjunction:
		if typeDef.AsDisjunction().Branches.HasNullType() {
			return nil
		}

		return defaultValueForType(schemas, typeDef.AsDisjunction().Branches[0], importModule, nil)
	case ir.KindRef:
		ref := typeDef.AsRef()
		referredPkg := ref.ReferredPkg
		referredPkg = importModule(referredPkg, "..models", referredPkg)

		referredObj, found := schemas.GetObject(ref.ReferredPkg, ref.ReferredType)
		if found && referredObj.Type.IsEnum() {
			enumName := tools.UpperSnakeCase(referredObj.Type.AsEnum().Values[0].Name)
			for _, enumValue := range referredObj.Type.AsEnum().Values {
				if enumValue.Value == typeDef.Default {
					enumName = tools.UpperSnakeCase(enumValue.Name)
					break
				}
			}

			objectName := tools.UpperCamelCase(referredObj.Name)

			if referredPkg == "" {
				return raw(objectName + "." + enumName)
			}

			return raw(referredPkg + "." + objectName + "." + enumName)
		} else if found && referredObj.Type.IsDisjunction() {
			return defaultValueForType(schemas, referredObj.Type, importModule, nil)
		}

		var extraDefaults []string

		if defaultsOverrides != nil {
			extraDefaults = make([]string, 0, defaultsOverrides.Len())
			defaultsOverrides.Iterate(func(k string, v any) {
				if !referredObj.Type.IsStruct() {
					return
				}
				field, fieldFound := referredObj.Type.AsStruct().FieldByName(k)
				if !fieldFound {
					return
				}

				value := v
				if field.Type.IsRef() {
					var fieldOverrides *orderedmap.Map[string, any]
					if overrides, ok := value.(map[string]any); ok {
						fieldOverrides = orderedmap.FromMap(overrides)
					}

					value = defaultValueForType(schemas, field.Type, importModule, fieldOverrides)
				}

				extraDefaults = append(extraDefaults, fmt.Sprintf("%s=%s", formatIdentifier(k), formatValue(value)))
			})
		}

		formattedRef := tools.UpperCamelCase(ref.ReferredType)
		if referredPkg != "" {
			formattedRef = referredPkg + "." + formattedRef
		}

		if referredObj.Type.IsConcreteScalar() {
			return raw(formattedRef)
		}

		return raw(fmt.Sprintf("%s(%s)", formattedRef, strings.Join(extraDefaults, ", ")))
	case ir.KindEnum: // anonymous enum
		return typeDef.AsEnum().Values[0].Value
	case ir.KindMap:
		return raw("{}")
	case ir.KindArray:
		return raw("[]")
	case ir.KindScalar:
		return defaultValueForScalar(typeDef.AsScalar())
	default:
		return "unknown"
	}
}

func defaultValueForScalar(scalar ir.ScalarType) any {
	// The scalar represents a constant
	if scalar.Value != nil {
		return scalar.Value
	}

	switch scalar.ScalarKind {
	case ir.KindNull, ir.KindAny:
		return nil

	case ir.KindBytes, ir.KindString:
		return ""

	case ir.KindFloat32, ir.KindFloat64:
		return 0.0

	case ir.KindUint8, ir.KindUint16, ir.KindUint32, ir.KindUint64:
		return 0

	case ir.KindInt8, ir.KindInt16, ir.KindInt32, ir.KindInt64:
		return 0

	case ir.KindBool:
		return false

	default:
		return "unknown"
	}
}
