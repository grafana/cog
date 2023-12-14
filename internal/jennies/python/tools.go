package python

import (
	"fmt"
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

func formatFieldName(name string) string {
	return tools.SnakeCase(escapeFieldName(name))
}

func escapeFieldName(name string) string {
	if isReservedPythonKeyword(name) || isBuiltInFunction(name) {
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

func defaultValueForType(schemas ast.Schemas, typeDef ast.Type, importModule func(alias string, pkg string, module string) string) any {
	if !typeDef.IsRef() && typeDef.Default != nil {
		return typeDef.Default
	}

	switch typeDef.Kind {
	case ast.KindDisjunction:
		if typeDef.AsDisjunction().Branches.HasNullType() {
			return nil
		}

		return defaultValueForType(schemas, typeDef.AsDisjunction().Branches[0], importModule)
	case ast.KindRef:
		ref := typeDef.AsRef()
		referredPkg := ref.ReferredPkg
		referredPkg = importModule(referredPkg, "..models", referredPkg)

		referredObj, found := schemas.LocateObject(ref.ReferredPkg, ref.ReferredType)
		if found && referredObj.Type.IsEnum() {
			enumName := tools.UpperSnakeCase(referredObj.Type.AsEnum().Values[0].Name)
			for _, enumValue := range referredObj.Type.AsEnum().Values {
				if enumValue.Value == typeDef.Default {
					enumName = tools.UpperSnakeCase(enumValue.Name)
					break
				}
			}

			if referredPkg == "" {
				return raw(fmt.Sprintf("%s.%s", referredObj.Name, enumName))
			}

			return raw(fmt.Sprintf("%s.%s.%s", referredPkg, referredObj.Name, enumName))
		} else if found && referredObj.Type.IsDisjunction() {
			return defaultValueForType(schemas, referredObj.Type, importModule)
		}

		if referredPkg == "" {
			return raw(fmt.Sprintf("%s()", ref.ReferredType))
		}

		return raw(fmt.Sprintf("%s.%s()", referredPkg, ref.ReferredType))
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
	case ast.KindNull:
		return nil
	case ast.KindAny:
		return raw("{}")

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
