package typescript

import (
	"fmt"
	"strings"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/orderedmap"
	"github.com/grafana/cog/internal/tools"
	"github.com/grafana/cog/pkg/ir"
	"github.com/grafana/cog/pkg/languages"
	"github.com/grafana/cog/pkg/template"
)

type raw string

type RawTypes struct {
	config        Config
	tmpl          *template.Template
	typeFormatter *typeFormatter
	schemas       ir.Schemas
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

		filename := jenny.config.pathWithPrefix(
			formatPackageName(schema.Package),
			"types.gen.ts",
		)

		files = append(files, *codejen.NewFile(filename, output, jenny))
	}

	return files, nil
}

func (jenny RawTypes) generateSchema(context languages.Context, schema *ir.Schema) ([]byte, error) {
	var buffer strings.Builder
	var err error

	imports := NewImportMap(jenny.config.PackagesImportMap)
	pkgMapper := func(pkg string) string {
		if imports.IsIdentical(pkg, schema.Package) {
			return ""
		}

		return imports.Add(pkg, fmt.Sprintf("../%s", pkg))
	}

	jenny.typeFormatter = defaultTypeFormatter(jenny.config, context, pkgMapper)

	schema.Objects.Iterate(func(_ string, object ir.Object) {
		typeDefGen, innerErr := jenny.formatObject(context, object, pkgMapper)
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

func (jenny RawTypes) formatObject(context languages.Context, def ir.Object, packageMapper packageMapper) ([]byte, error) {
	var buffer strings.Builder

	if len(def.Comments) != 0 || def.DeprecationMessage != "" {
		fmt.Fprintf(&buffer, "/**\n")
		for _, commentLine := range def.Comments {
			fmt.Fprintf(&buffer, " * %s\n", commentLine)
		}
		if def.DeprecationMessage != "" {
			fmt.Fprintf(&buffer, " * @deprecated %s\n", def.DeprecationMessage)
		}

		fmt.Fprintf(&buffer, " */\n")
	}

	buffer.WriteString(jenny.typeFormatter.formatTypeDeclaration(def))

	objectName := tools.CleanupNames(def.Name)

	// generate a "default value factory" for every object, except for constants or composability slots
	if (!def.Type.IsScalar() && !def.Type.IsComposableSlot()) || (def.Type.IsScalar() && !def.Type.AsScalar().IsConcrete()) {
		buffer.WriteString("\n")

		fmt.Fprintf(&buffer, "export const default%[1]s = (): %[2]s => (", tools.UpperCamelCase(objectName), objectName)

		formattedDefaults := formatValue(jenny.defaultValueForObject(def, packageMapper))
		buffer.WriteString(formattedDefaults)

		buffer.WriteString(");\n")
	}

	if jenny.config.GenerateEqual && def.Type.IsStruct() {
		eqFunc := newEqualityMethods().generateForObject(context, def)
		if eqFunc != "" {
			buffer.WriteString("\n")
			buffer.WriteString(eqFunc)
		}
	}

	customMethodsBlock := template.CustomObjectMethodsBlock(def)
	if jenny.tmpl.Exists(customMethodsBlock) {
		err := jenny.tmpl.RenderInBuffer(&buffer, customMethodsBlock, map[string]any{
			"Object": def,
		})
		if err != nil {
			return nil, err
		}
		buffer.WriteString("\n")
	}

	customAllBlock := template.CustomObjectMethodAllBlock()
	if jenny.tmpl.Exists(customAllBlock) {
		err := jenny.tmpl.RenderInBuffer(&buffer, customAllBlock, map[string]any{
			"Object": def,
		})
		if err != nil {
			return nil, err
		}
		buffer.WriteString("\n")
	}

	return []byte(buffer.String()), nil
}

func prefixLinesWith(input string, prefix string) string {
	lines := strings.Split(input, "\n")
	prefixed := make([]string, 0, len(lines))

	for _, line := range lines {
		prefixed = append(prefixed, prefix+line)
	}

	return strings.Join(prefixed, "\n")
}

/******************************************************************************
* 					 Default and "empty" values management 					  *
******************************************************************************/

func (jenny RawTypes) defaultValueForObject(object ir.Object, packageMapper packageMapper) any {
	switch object.Type.Kind {
	case ir.KindEnum:
		enum := object.Type.AsEnum()
		defaultValue := enum.Values[0].Value
		if object.Type.Default != nil {
			defaultValue = object.Type.Default
		}

		return raw(jenny.typeFormatter.enums.formatValue(object, defaultValue))
	default:
		return jenny.defaultValueForType(object.Type, packageMapper)
	}
}

func (jenny RawTypes) defaultValueForType(typeDef ir.Type, packageMapper packageMapper) any {
	if typeDef.Default != nil {
		return typeDef.Default
	}

	switch typeDef.Kind {
	case ir.KindDisjunction:
		return jenny.defaultValueForType(typeDef.AsDisjunction().Branches[0], packageMapper)
	case ir.KindStruct:
		return jenny.defaultValuesForStructType(typeDef, packageMapper)
	case ir.KindEnum: // anonymous enum
		defaultValue := typeDef.AsEnum().Values[0].Value
		if typeDef.Default != nil {
			defaultValue = typeDef.Default
		}

		return defaultValue
	case ir.KindRef:
		return jenny.defaultValuesForReference(typeDef, packageMapper)
	case ir.KindMap:
		return raw("{}")
	case ir.KindArray:
		return raw("[]")
	case ir.KindScalar:
		return defaultValueForScalar(typeDef.AsScalar())
	case ir.KindIntersection:
		return jenny.defaultValuesForIntersection(typeDef.AsIntersection(), packageMapper)
	case ir.KindConstantRef:
		return jenny.defaultValueForConstantReferences(typeDef.AsConstantRef())
	default:
		return "unknown"
	}
}

func (jenny RawTypes) defaultValuesForStructType(structType ir.Type, packageMapper packageMapper) *orderedmap.Map[string, any] {
	defaults := orderedmap.New[string, any]()

	for _, field := range structType.AsStruct().Fields {
		if field.Type.Default != nil {
			switch field.Type.Kind {
			case ir.KindRef:
				defaults.Set(field.Name, jenny.defaultValuesForReference(field.Type, packageMapper))
				continue
			case ir.KindStruct:
				defaultMap := field.Type.Default.(map[string]any)
				defaults.Set(field.Name, jenny.defaultValueForStructs(field.Type.AsStruct(), orderedmap.FromMap(defaultMap)))
				continue
			default:
				defaults.Set(field.Name, field.Type.Default)
				continue
			}
		}

		if !field.Required && !field.Type.IsConstantRef() {
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

func defaultValueForScalar(scalar ir.ScalarType) any {
	// The scalar represents a constant
	if scalar.Value != nil {
		return scalar.Value
	}

	switch scalar.ScalarKind {
	case ir.KindNull:
		return raw("null")
	case ir.KindAny:
		return raw("{}")

	case ir.KindBytes, ir.KindString:
		return raw("\"\"")

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

func (jenny RawTypes) defaultValuesForIntersection(intersectDef ir.IntersectionType, packageMapper packageMapper) *orderedmap.Map[string, any] {
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

func (jenny RawTypes) defaultValuesForReference(typeDef ir.Type, packageMapper packageMapper) any {
	ref := typeDef.AsRef()

	pkg := packageMapper(ref.ReferredPkg)
	referredType, _ := jenny.schemas.GetObject(ref.ReferredPkg, ref.ReferredType)
	referredTypeName := formatObjectName(referredType.Name)
	if referredTypeName == "" {
		referredTypeName = formatObjectName(ref.ReferredType)
	}

	// is the reference to a constant?
	if referredType.Type.IsConcreteScalar() {
		if pkg != "" {
			return raw(fmt.Sprintf("%s.%s", pkg, referredTypeName))
		}

		return raw(referredTypeName)
	}

	if referredType.Type.IsEnum() {
		return raw(jenny.typeFormatter.enums.formatValue(referredType, typeDef.Default))
	}

	if hasStructDefaults(referredType.Type, typeDef.Default) {
		defaultMap := typeDef.Default.(map[string]any)
		return jenny.defaultValueForStructs(referredType.Type.AsStruct(), orderedmap.FromMap(defaultMap))
	}

	if pkg != "" {
		return raw(fmt.Sprintf("%s.default%s()", pkg, tools.UpperCamelCase(referredTypeName)))
	}

	return raw(fmt.Sprintf("default%s()", tools.UpperCamelCase(referredTypeName)))
}

func (jenny RawTypes) defaultValueForStructs(def ir.StructType, m *orderedmap.Map[string, any]) any {
	var buffer strings.Builder

	for _, f := range def.Fields {
		if m.Has(f.Name) {
			switch x := m.Get(f.Name).(type) {
			case map[string]any:
				fmt.Fprintf(&buffer, "%s: %v, ", f.Name, jenny.defaultValueForStructs(f.Type.AsStruct(), orderedmap.FromMap(x)))
			case nil:
				fmt.Fprintf(&buffer, "%s: %v, ", f.Name, formatValue([]any{}))
			default:
				if f.Type.IsRef() {
					ref := f.Type.AsRef()
					referredType, refFound := jenny.schemas.GetObject(ref.ReferredPkg, ref.ReferredType)

					if refFound && referredType.Type.IsEnum() {
						fmt.Fprintf(&buffer, "%s: %v, ", f.Name, jenny.typeFormatter.enums.formatValue(referredType, x))
						continue
					}
				}

				fmt.Fprintf(&buffer, "%s: %v, ", f.Name, formatValue(x))
			}
		} else if f.Required {
			switch f.Type.Kind {
			case ir.KindStruct:
				fmt.Fprintf(&buffer, "%s: { %v }, ", f.Name, defaultEmptyValuesForStructs(f.Type.AsStruct()))
			case ir.KindArray:
				fmt.Fprintf(&buffer, "%s: []", f.Name)
			case ir.KindScalar:
				fmt.Fprintf(&buffer, "%s: %v, ", f.Name, defaultValueForScalar(f.Type.AsScalar()))
			}
		}
	}

	return raw(fmt.Sprintf("{ %+v}", buffer.String()))
}

func defaultEmptyValuesForStructs(def ir.StructType) string {
	var buffer strings.Builder

	for _, f := range def.Fields {
		switch f.Type.Kind {
		case ir.KindStruct:
			fmt.Fprintf(&buffer, "%s: { %v }, ", f.Name, defaultEmptyValuesForStructs(f.Type.AsStruct()))
		case ir.KindArray:
			fmt.Fprintf(&buffer, "%s: []", f.Name)
		case ir.KindScalar:
			fmt.Fprintf(&buffer, "%s: %v, ", f.Name, defaultValueForScalar(f.Type.AsScalar()))
		default:
		}
	}

	return buffer.String()
}

func (jenny RawTypes) defaultValueForConstantReferences(def ir.ConstantReferenceType) any {
	referredType, ok := jenny.schemas.GetObject(def.ReferredPkg, def.ReferredType)
	if !ok {
		return "unknown"
	}

	if referredType.Type.IsEnum() {
		return raw(jenny.typeFormatter.enums.formatValue(referredType, def.ReferenceValue))
	}

	if referredType.Type.IsScalar() {
		return raw(def.ReferredType)
	}

	return "unknown"
}

func hasStructDefaults(typeDef ir.Type, defaults any) bool {
	_, ok := defaults.(map[string]any)
	return ok && typeDef.IsStruct()
}
