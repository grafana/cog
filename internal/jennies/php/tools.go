package php

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/grafana/cog/internal/orderedmap"
	"github.com/grafana/cog/internal/tools"
	"github.com/grafana/cog/pkg/ir"
)

type raw string

func formatPackageName(pkg string) string {
	rgx := regexp.MustCompile("[^a-zA-Z0-9_]+")

	return tools.UpperCamelCase(rgx.ReplaceAllString(pkg, ""))
}

func formatObjectName(name string) string {
	return tools.UpperCamelCase(name)
}

func formatOptionName(name string) string {
	return tools.LowerCamelCase(name)
}

func formatConstantName(name string) string {
	return tools.UpperSnakeCase(name)
}

func formatFieldName(name string) string {
	return tools.LowerCamelCase(name)
}

func formatArgName(name string) string {
	return tools.LowerCamelCase(name)
}

func formatEnumMemberName(name string) string {
	return tools.LowerCamelCase(name)
}

func formatCommentsBlock(comments []string) string {
	if len(comments) == 0 {
		return ""
	}

	var buffer strings.Builder

	if len(comments) != 0 {
		buffer.WriteString("/**\n")
	}
	for _, commentLine := range comments {
		buffer.WriteString(fmt.Sprintf(" * %s\n", commentLine))
	}
	if len(comments) != 0 {
		buffer.WriteString(" */\n")
	}

	return buffer.String()
}

func formatFieldPath(fieldPath ir.Path) string {
	path := ""

	for i, chunk := range fieldPath {
		last := i == len(fieldPath)-1
		output := formatFieldName(chunk.Identifier)

		if chunk.Index != nil {
			output += "["
			if chunk.Index.Constant != nil {
				output += formatValue(chunk.Index.Constant)
			} else {
				output += "$" + formatArgName(chunk.Index.Argument.Name)
			}
			output += "]"
		}

		path += output
		if !last && fieldPath[i+1].Index == nil {
			path += "->"
		}
	}

	return path
}

func formatValue(val any) string {
	if val == nil {
		return "null"
	}

	if rawVal, ok := val.(raw); ok {
		return string(rawVal)
	}

	if asBool, ok := val.(bool); ok {
		if asBool {
			return "true"
		}

		return "false"
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
			items = append(items, fmt.Sprintf("%q => %s", key, formatValue(value)))
		}

		return fmt.Sprintf("[%s]", strings.Join(items, ", "))
	}

	if orderedMap, ok := val.(*orderedmap.Map[string, any]); ok {
		items := make([]string, 0, orderedMap.Len())

		orderedMap.Iterate(func(key string, value any) {
			items = append(items, fmt.Sprintf("%q => %s", key, formatValue(value)))
		})

		return fmt.Sprintf("[%s]", strings.Join(items, ", "))
	}

	return fmt.Sprintf("%#v", val)
}

func disjunctionCaseForType(typesFormatter *typeFormatter, input string, typeDef ir.Type) string {
	// TODO: shaky at best
	if typeDef.IsAnyOf(ir.KindArray, ir.KindMap) {
		return fmt.Sprintf("is_array(%s)", input)
	}

	if typeDef.IsScalar() {
		testMap := map[ir.ScalarKind]string{
			ir.KindBytes:   "is_string",
			ir.KindString:  "is_string",
			ir.KindFloat32: "is_float",
			ir.KindFloat64: "is_float",
			ir.KindUint8:   "is_int",
			ir.KindUint16:  "is_int",
			ir.KindUint32:  "is_int",
			ir.KindUint64:  "is_int",
			ir.KindInt8:    "is_int",
			ir.KindInt16:   "is_int",
			ir.KindInt32:   "is_int",
			ir.KindInt64:   "is_int",
			ir.KindBool:    "is_bool",
		}

		testFunc := testMap[typeDef.Scalar.ScalarKind]
		if testFunc == "" {
			return "/* unhandled scalar type */"
		}

		return fmt.Sprintf("%s(%s)", testFunc, input)
	}

	if typeDef.IsRef() {
		return fmt.Sprintf("%s instanceof %s", input, typesFormatter.formatRef(typeDef.Ref.AsType(), false))
	}

	return "/* unhandled type */"
}

/******************************************
 *  Default and "empty" values management *
 *****************************************/

func defaultValueForType(config Config, schemas ir.Schemas, typeDef ir.Type, defaultsOverrides *orderedmap.Map[string, any]) any {
	if !typeDef.IsRef() && typeDef.Default != nil {
		return typeDef.Default
	}

	switch typeDef.Kind {
	case ir.KindDisjunction:
		if typeDef.AsDisjunction().Branches.HasNullType() {
			return nil
		}

		return defaultValueForType(config, schemas, typeDef.AsDisjunction().Branches[0], nil)
	case ir.KindRef:
		ref := typeDef.AsRef()
		referredPkg := formatPackageName(ref.ReferredPkg)
		referredObj, found := schemas.GetObject(ref.ReferredPkg, ref.ReferredType)
		if found && referredObj.Type.IsEnum() {
			enumName := formatObjectName(referredObj.Type.AsEnum().Values[0].Name)
			for _, enumValue := range referredObj.Type.AsEnum().Values {
				if enumValue.Value == typeDef.Default {
					enumName = formatEnumMemberName(enumValue.Name)
					break
				}
			}

			return raw(fmt.Sprintf(config.fullNamespaceRef(referredPkg+"\\"+formatObjectName(referredObj.Name))+"::%s()", enumName))
		} else if found && referredObj.Type.IsDisjunction() {
			return defaultValueForType(config, schemas, referredObj.Type, nil)
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

					value = defaultValueForType(config, schemas, field.Type, fieldOverrides)
				}

				extraDefaults = append(extraDefaults, fmt.Sprintf("%s: %s", formatFieldName(k), formatValue(value)))
			})
		}

		formattedRef := formatObjectName(ref.ReferredType)
		if referredPkg != "" {
			formattedRef = config.fullNamespaceRef(referredPkg + "\\" + formattedRef)
		}

		if referredObj.Type.IsConcreteScalar() {
			return raw(formattedRef)
		}

		return raw(fmt.Sprintf("new %s(%s)", formattedRef, strings.Join(extraDefaults, ", ")))
	case ir.KindEnum: // anonymous enum
		return typeDef.AsEnum().Values[0].Value
	case ir.KindMap, ir.KindArray:
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
