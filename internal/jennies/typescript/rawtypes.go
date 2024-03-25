package typescript

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/languages"
	"github.com/grafana/cog/internal/orderedmap"
	"github.com/grafana/cog/internal/tools"
)

type raw string

type pkgMapper func(string) string

type RawTypes struct {
	typeFormatter *typeFormatter
	schemas       ast.Schemas
}

func (jenny RawTypes) JennyName() string {
	return "TypescriptRawTypes"
}

func (jenny RawTypes) Generate(context languages.Context) (codejen.Files, error) {
	jenny.schemas = context.Schemas
	files := make(codejen.Files, 0, len(context.Schemas))

	for _, schema := range context.Schemas {
		output, err := jenny.generateSchema(context, schema)
		if err != nil {
			return nil, err
		}

		filename := filepath.Join(
			"src",
			formatPackageName(schema.Package),
			"types.gen.ts",
		)

		files = append(files, *codejen.NewFile(filename, output, jenny))
	}

	return files, nil
}

func (jenny RawTypes) generateSchema(context languages.Context, schema *ast.Schema) ([]byte, error) {
	var buffer strings.Builder
	var err error

	imports := NewImportMap()
	packageMapper := func(pkg string) string {
		if imports.IsIdentical(pkg, schema.Package) {
			return ""
		}

		return imports.Add(pkg, fmt.Sprintf("../%s", pkg))
	}

	jenny.typeFormatter = defaultTypeFormatter(context, packageMapper)

	schema.Objects.Iterate(func(_ string, object ast.Object) {
		typeDefGen, innerErr := jenny.formatObject(object, packageMapper)
		if innerErr != nil {
			err = innerErr
			return
		}

		buffer.Write(typeDefGen)
		buffer.WriteString("\n")
	})
	if err != nil {
		return nil, err
	}

	importStatements := imports.String()
	if importStatements != "" {
		importStatements += "\n\n"
	}

	return []byte(importStatements + buffer.String()), nil
}

func (jenny RawTypes) formatObject(def ast.Object, packageMapper pkgMapper) ([]byte, error) {
	var buffer strings.Builder

	for _, commentLine := range def.Comments {
		buffer.WriteString(fmt.Sprintf("// %s\n", commentLine))
	}

	buffer.WriteString("export ")

	objectName := tools.CleanupNames(def.Name)

	switch def.Type.Kind {
	case ast.KindStruct:
		buffer.WriteString(fmt.Sprintf("interface %s ", objectName))
		buffer.WriteString(jenny.typeFormatter.formatStructFields(def.Type))
		buffer.WriteString("\n")
	case ast.KindEnum:
		buffer.WriteString(fmt.Sprintf("enum %s {\n", objectName))
		for _, val := range def.Type.AsEnum().Values {
			name := tools.CleanupNames(tools.UpperCamelCase(escapeEnumMemberName(val.Name)))

			buffer.WriteString(fmt.Sprintf("\t%s = %s,\n", name, formatValue(val.Value)))
		}
		buffer.WriteString("}\n")
	case ast.KindDisjunction, ast.KindMap, ast.KindArray, ast.KindRef:
		buffer.WriteString(fmt.Sprintf("type %s = %s;\n", objectName, jenny.typeFormatter.formatType(def.Type)))
	case ast.KindScalar:
		scalarType := def.Type.AsScalar()
		typeValue := formatValue(scalarType.Value)

		if !scalarType.IsConcrete() || def.Type.Hints["kind"] == "type" {
			if !scalarType.IsConcrete() {
				typeValue = jenny.typeFormatter.formatScalarKind(scalarType.ScalarKind)
			}

			buffer.WriteString(fmt.Sprintf("type %s = %s;\n", objectName, typeValue))
		} else {
			buffer.WriteString(fmt.Sprintf("const %s = %s;\n", objectName, typeValue))
		}
	case ast.KindIntersection:
		buffer.WriteString(fmt.Sprintf("interface %s ", objectName))
		buffer.WriteString(jenny.typeFormatter.formatType(def.Type))
		buffer.WriteString("\n")
	case ast.KindComposableSlot:
		buffer.WriteString(fmt.Sprintf("interface %s %s\n", objectName, jenny.typeFormatter.variantInterface(string(def.Type.AsComposableSlot().Variant))))
	default:
		return nil, fmt.Errorf("unhandled object of type: %s", def.Type.Kind)
	}
	// generate a "default value factory" for every object, except for constants or composability slots
	if (!def.Type.IsScalar() && !def.Type.IsComposableSlot()) || (def.Type.IsScalar() && !def.Type.AsScalar().IsConcrete()) {
		buffer.WriteString("\n")

		buffer.WriteString(fmt.Sprintf("export const default%[1]s = (): %[2]s => (", tools.UpperCamelCase(objectName), objectName))

		formattedDefaults := formatValue(jenny.defaultValueForObject(def, packageMapper))
		buffer.WriteString(formattedDefaults)

		buffer.WriteString(");\n")
	}

	return []byte(buffer.String()), nil
}

/******************************************************************************
* 					 Default and "empty" values management 					  *
******************************************************************************/

func (jenny RawTypes) defaultValueForObject(object ast.Object, packageMapper pkgMapper) any {
	switch object.Type.Kind {
	case ast.KindEnum:
		enum := object.Type.AsEnum()
		defaultValue := enum.Values[0].Value
		if object.Type.Default != nil {
			defaultValue = object.Type.Default
		}

		return raw(jenny.typeFormatter.formatEnumValue(object, defaultValue))
	default:
		return jenny.defaultValueForType(object.Type, packageMapper)
	}
}

func (jenny RawTypes) defaultValueForType(typeDef ast.Type, packageMapper pkgMapper) any {
	if typeDef.Default != nil {
		return typeDef.Default
	}

	switch typeDef.Kind {
	case ast.KindDisjunction:
		return jenny.defaultValueForType(typeDef.AsDisjunction().Branches[0], packageMapper)
	case ast.KindStruct:
		return jenny.defaultValuesForStructType(typeDef, packageMapper)
	case ast.KindEnum: // anonymous enum
		defaultValue := typeDef.AsEnum().Values[0].Value
		if typeDef.Default != nil {
			defaultValue = typeDef.Default
		}

		return defaultValue
	case ast.KindRef:
		return jenny.defaultValuesForReference(typeDef, packageMapper)
	case ast.KindMap:
		return raw("{}")
	case ast.KindArray:
		return raw("[]")
	case ast.KindScalar:
		return defaultValueForScalar(typeDef.AsScalar())
	case ast.KindIntersection:
		return jenny.defaultValuesForIntersection(typeDef.AsIntersection(), packageMapper)
	default:
		return "unknown"
	}
}

func (jenny RawTypes) defaultValuesForStructType(structType ast.Type, packageMapper pkgMapper) *orderedmap.Map[string, any] {
	defaults := orderedmap.New[string, any]()

	for _, field := range structType.AsStruct().Fields {
		if field.Type.Default != nil {
			switch field.Type.Kind {
			case ast.KindRef:
				defaults.Set(field.Name, jenny.defaultValuesForReference(field.Type, packageMapper))
				continue
			case ast.KindStruct:
				defaultMap := field.Type.Default.(map[string]interface{})
				defaults.Set(field.Name, jenny.defaultValueForStructs(field.Type.AsStruct(), orderedmap.FromMap(defaultMap)))
				continue
			default:
				defaults.Set(field.Name, field.Type.Default)
				continue
			}
		}

		if !field.Required {
			continue
		}

		defaults.Set(field.Name, jenny.defaultValueForType(field.Type, packageMapper))
	}

	if structType.ImplementsVariant() {
		variant := tools.UpperCamelCase(structType.ImplementedVariant())
		defaults.Set("_implements"+variant+"Variant", raw("() => {}"))
	}

	return defaults
}

func defaultValueForScalar(scalar ast.ScalarType) any {
	// The scalar represents a constant
	if scalar.Value != nil {
		return scalar.Value
	}

	switch scalar.ScalarKind {
	case ast.KindNull:
		return raw("null")
	case ast.KindAny:
		return raw("{}")

	case ast.KindBytes, ast.KindString:
		return raw("\"\"")

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

func (jenny RawTypes) defaultValuesForIntersection(intersectDef ast.IntersectionType, packageMapper pkgMapper) *orderedmap.Map[string, any] {
	defaults := orderedmap.New[string, any]()

	for _, branch := range intersectDef.Branches {
		if branch.Ref != nil {
			continue
		}

		if branch.Struct != nil {
			strctDef := jenny.defaultValuesForStructType(branch, packageMapper)
			strctDef.Iterate(func(key string, value any) {
				defaults.Set(key, value)
			})
		}

		// TODO: Add them for other types?
	}

	return defaults
}

func (jenny RawTypes) defaultValuesForReference(typeDef ast.Type, packageMapper pkgMapper) any {
	ref := typeDef.AsRef()

	pkg := packageMapper(ref.ReferredPkg)
	referredType, _ := jenny.schemas.LocateObject(ref.ReferredPkg, ref.ReferredType)
	referredTypeName := tools.CleanupNames(referredType.Name)

	// is the reference to a constant?
	if referredType.Type.IsConcreteScalar() {
		if pkg != "" {
			return raw(fmt.Sprintf("%s.%s", pkg, referredTypeName))
		}

		return raw(referredTypeName)
	}

	if referredType.Type.IsEnum() {
		return raw(jenny.typeFormatter.formatEnumValue(referredType, typeDef.Default))
	}

	if hasStructDefaults(referredType.Type, typeDef.Default) {
		defaultMap := typeDef.Default.(map[string]any)
		return jenny.defaultValueForStructs(referredType.Type.AsStruct(), orderedmap.FromMap(defaultMap))
	}

	if pkg != "" {
		return raw(fmt.Sprintf("%s.default%s()", pkg, referredTypeName))
	}

	return raw(fmt.Sprintf("default%s()", tools.UpperCamelCase(referredTypeName)))
}

func (jenny RawTypes) defaultValueForStructs(def ast.StructType, m *orderedmap.Map[string, any]) any {
	var buffer strings.Builder

	for _, f := range def.Fields {
		if m.Has(f.Name) {
			switch x := m.Get(f.Name).(type) {
			case map[string]any:
				buffer.WriteString(fmt.Sprintf("%s: %v, ", f.Name, jenny.defaultValueForStructs(f.Type.AsStruct(), orderedmap.FromMap(x))))
			case nil:
				buffer.WriteString(fmt.Sprintf("%s: %v, ", f.Name, formatValue([]any{})))
			default:
				if f.Type.IsRef() {
					ref := f.Type.AsRef()
					referredType, refFound := jenny.schemas.LocateObject(ref.ReferredPkg, ref.ReferredType)

					if refFound && referredType.Type.IsEnum() {
						buffer.WriteString(fmt.Sprintf("%s: %v, ", f.Name, jenny.typeFormatter.formatEnumValue(referredType, x)))
						continue
					}
				}

				buffer.WriteString(fmt.Sprintf("%s: %v, ", f.Name, formatValue(x)))
			}
		} else if f.Required {
			switch f.Type.Kind {
			case ast.KindStruct:
				buffer.WriteString(fmt.Sprintf("%s: { %v }, ", f.Name, defaultEmptyValuesForStructs(f.Type.AsStruct())))
			case ast.KindArray:
				buffer.WriteString(fmt.Sprintf("%s: []", f.Name))
			case ast.KindScalar:
				buffer.WriteString(fmt.Sprintf("%s: %v, ", f.Name, defaultValueForScalar(f.Type.AsScalar())))
			}
		}
	}

	return raw(fmt.Sprintf("{ %+v}", buffer.String()))
}

func defaultEmptyValuesForStructs(def ast.StructType) string {
	var buffer strings.Builder

	for _, f := range def.Fields {
		switch f.Type.Kind {
		case ast.KindStruct:
			buffer.WriteString(fmt.Sprintf("%s: { %v }, ", f.Name, defaultEmptyValuesForStructs(f.Type.AsStruct())))
		case ast.KindArray:
			buffer.WriteString(fmt.Sprintf("%s: []", f.Name))
		case ast.KindScalar:
			buffer.WriteString(fmt.Sprintf("%s: %v, ", f.Name, defaultValueForScalar(f.Type.AsScalar())))
		default:
		}
	}

	return buffer.String()
}

func hasStructDefaults(typeDef ast.Type, defaults any) bool {
	_, ok := defaults.(map[string]interface{})
	return ok && typeDef.IsStruct()
}

func escapeEnumMemberName(identifier string) string {
	if strings.EqualFold("nan", identifier) {
		return "not_a_number"
	}

	return identifier
}
