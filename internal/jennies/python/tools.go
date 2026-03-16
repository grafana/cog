package python

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/tools"
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

	return fmt.Sprintf("%#v", val)
}

func formatFieldPath(fieldPath ast.Path) string {
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

func defaultValueForType(schemas ast.Schemas, typeDef ast.Type, importModule moduleImporter, overrideDefault *ast.TypeDefault) any {
	if !typeDef.IsRef() && overrideDefault != nil {
		if overrideDefault.Scalar != nil {
			return overrideDefault.Scalar.Value
		}
		if overrideDefault.Array != nil {
			items := make([]any, len(overrideDefault.Array))
			for i, elem := range overrideDefault.Array {
				if elem.Scalar != nil {
					items[i] = elem.Scalar.Value
				}
			}
			return items
		}
		// Struct default → fall through
	}

	switch typeDef.Kind {
	case ast.KindDisjunction:
		if typeDef.AsDisjunction().Branches.HasNullType() {
			return nil
		}

		return defaultValueForType(schemas, typeDef.AsDisjunction().Branches[0], importModule, nil)
	case ast.KindRef:
		ref := typeDef.AsRef()
		referredPkg := ref.ReferredPkg
		referredPkg = importModule(referredPkg, "..models", referredPkg)

		referredObj, found := schemas.LocateObject(ref.ReferredPkg, ref.ReferredType)
		if found && referredObj.Type.IsEnum() {
			enumName := tools.UpperSnakeCase(referredObj.Type.AsEnum().Values[0].Name)
			if overrideDefault != nil && overrideDefault.Scalar != nil {
				for _, enumValue := range referredObj.Type.AsEnum().Values {
					if enumValue.Value == overrideDefault.Scalar.Value {
						enumName = tools.UpperSnakeCase(enumValue.Name)
						break
					}
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

		if overrideDefault != nil && overrideDefault.Struct != nil {
			keys := slices.Sorted(maps.Keys(overrideDefault.Struct))
			extraDefaults = make([]string, 0, len(keys))
			for _, k := range keys {
				fieldTD := overrideDefault.Struct[k]
				if !referredObj.Type.IsStruct() {
					continue
				}
				field, fieldFound := referredObj.Type.AsStruct().FieldByName(k)
				if !fieldFound {
					continue
				}

				var value any
				if field.Type.IsRef() {
					value = defaultValueForType(schemas, field.Type, importModule, fieldTD)
				} else if fieldTD.Scalar != nil {
					value = fieldTD.Scalar.Value
				} else if fieldTD.Array != nil {
					items := make([]any, len(fieldTD.Array))
					for i, elem := range fieldTD.Array {
						if elem.Scalar != nil {
							items[i] = elem.Scalar.Value
						}
					}
					value = items
				}

				extraDefaults = append(extraDefaults, fmt.Sprintf("%s=%s", formatIdentifier(k), formatValue(value)))
			}
		}

		formattedRef := tools.UpperCamelCase(ref.ReferredType)
		if referredPkg != "" {
			formattedRef = referredPkg + "." + formattedRef
		}

		if referredObj.Type.IsConcreteScalar() {
			return raw(formattedRef)
		}

		return raw(fmt.Sprintf("%s(%s)", formattedRef, strings.Join(extraDefaults, ", ")))
	case ast.KindEnum: // anonymous enum
		return typeDef.AsEnum().Values[0].Value
	case ast.KindMap:
		return raw("{}")
	case ast.KindArray:
		return raw("[]")
	case ast.KindScalar:
		return defaultValueForScalar(typeDef.AsScalar())
	default:
		return "unknown"
	}
}

func defaultValueForScalar(scalar ast.ScalarType) any {
	// The scalar represents a constant
	if scalar.Value != nil {
		return scalar.Value
	}

	switch scalar.ScalarKind {
	case ast.KindNull, ast.KindAny:
		return nil

	case ast.KindBytes, ast.KindString:
		return ""

	case ast.KindFloat32, ast.KindFloat64:
		return 0.0

	case ast.KindUint8, ast.KindUint16, ast.KindUint32, ast.KindUint64:
		return 0

	case ast.KindInt8, ast.KindInt16, ast.KindInt32, ast.KindInt64:
		return 0

	case ast.KindBool:
		return false

	default:
		return "unknown"
	}
}
