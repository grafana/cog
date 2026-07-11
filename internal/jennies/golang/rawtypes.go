package golang

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/jennies/template"
	"github.com/grafana/cog/internal/languages"
	"github.com/grafana/cog/internal/tools"
)

type RawTypes struct {
	config          Config
	tmpl            *template.Template
	apiRefCollector *common.APIReferenceCollector

	typeFormatter *typeFormatter
	packageMapper func(pkg string) string
}

func (jenny RawTypes) JennyName() string {
	return "GoRawTypes"
}

func (jenny RawTypes) Generate(context languages.Context) (codejen.Files, error) {
	files := make(codejen.Files, 0, len(context.Schemas))

	jenny.tmpl = jenny.tmpl.
		Funcs(common.TypeResolvingTemplateHelpers(context)).
		Funcs(common.TypesTemplateHelpers(context))

	for _, schema := range context.Schemas {
		output, err := jenny.generateSchema(context, schema)
		if err != nil {
			return nil, err
		}

		filename := filepath.Join(
			formatPackageName(schema.Package),
			"types_gen.go",
		)

		files = append(files, *codejen.NewFile(filename, output, jenny))
	}

	return files, nil
}

func (jenny RawTypes) generateSchema(context languages.Context, schema *ast.Schema) ([]byte, error) {
	var buffer strings.Builder
	var err error

	imports := NewImportMap(jenny.config.PackageRoot)
	jenny.packageMapper = func(pkg string) string {
		if imports.IsIdentical(pkg, schema.Package) {
			return ""
		}

		return imports.Add(pkg, jenny.config.importPath(pkg))
	}
	jenny.typeFormatter = defaultTypeFormatter(jenny.config, context, imports, jenny.packageMapper)
	unmarshallerGenerator := newJSONMarshalling(jenny.config, jenny.tmpl, imports, jenny.packageMapper, jenny.typeFormatter, jenny.apiRefCollector)
	strictUnmarshallerGenerator := newStrictJSONUnmarshal(jenny.config, jenny.tmpl, imports, jenny.packageMapper, jenny.typeFormatter, jenny.apiRefCollector)
	equalityMethodsGenerator := newEqualityMethods(jenny.tmpl, jenny.apiRefCollector)
	validationMethodsGenerator := newValidationMethods(jenny.tmpl, jenny.packageMapper, jenny.apiRefCollector)

	schema.Objects.Iterate(func(_ string, object ast.Object) {
		innerErr := jenny.formatObject(&buffer, schema, object)
		if innerErr != nil {
			err = innerErr
			return
		}

		buffer.WriteString("\n")

		jenny.generateConstructor(&buffer, context, object)

		innerErr = unmarshallerGenerator.generateForObject(&buffer, context, object)
		if innerErr != nil {
			err = innerErr
			return
		}

		if !jenny.config.SkipRuntime {
			innerErr = strictUnmarshallerGenerator.generateForObject(&buffer, context, object)
			if innerErr != nil {
				err = innerErr
				return
			}
		}

		if jenny.config.GenerateEqual {
			innerErr = equalityMethodsGenerator.generateForObject(&buffer, context, object, imports)
			if innerErr != nil {
				err = innerErr
				return
			}
		}

		if !jenny.config.SkipRuntime && (jenny.config.generateBuilders || jenny.config.GenerateValidate) {
			innerErr = validationMethodsGenerator.generateForObject(&buffer, context, object, imports)
			if innerErr != nil {
				err = innerErr
				return
			}
		}

		customMethodsBlock := template.CustomObjectMethodsBlock(object)
		if jenny.tmpl.Exists(customMethodsBlock) {
			innerErr = jenny.tmpl.RenderInBuffer(&buffer, customMethodsBlock, map[string]any{
				"Object": object,
			})
			if innerErr != nil {
				err = innerErr
				return
			}
			buffer.WriteString("\n")
		}

		customAllBlock := template.CustomObjectMethodAllBlock()
		if jenny.tmpl.Exists(customAllBlock) {
			innerErr = jenny.tmpl.RenderInBuffer(&buffer, customAllBlock, map[string]any{
				"Object": object,
			})
			if innerErr != nil {
				err = innerErr
				return
			}
			buffer.WriteString("\n")
		}
	})
	if err != nil {
		return nil, err
	}

	customSchemaVariant := template.CustomSchemaVariantBlock(schema)
	if jenny.tmpl.Exists(customSchemaVariant) {
		if err := jenny.tmpl.RenderInBuffer(&buffer, customSchemaVariant, map[string]any{
			"Schema": schema,
			"Config": jenny.config,
		}); err != nil {
			return nil, err
		}
	}

	importStatements := imports.String()
	if importStatements != "" {
		importStatements += "\n\n"
	}

	return []byte(fmt.Sprintf(`package %[1]s

%[2]s%[3]s`, formatPackageName(schema.Package), importStatements, buffer.String())), nil
}

func (jenny RawTypes) formatObject(buffer *strings.Builder, schema *ast.Schema, object ast.Object) error {
	objectName := formatObjectName(object.Name)

	comments := object.Comments
	if jenny.config.debug {
		passesTrail := tools.Map(object.PassesTrail, func(trail string) string {
			return fmt.Sprintf("Modified by compiler pass '%s'", trail)
		})
		comments = append(comments, passesTrail...)
	}

	for _, commentLine := range comments {
		fmt.Fprintf(buffer, "// %s\n", commentLine)
	}
	if object.DeprecationMessage != "" {
		fmt.Fprintf(buffer, "// Deprecated: %s\n", object.DeprecationMessage)
	}

	buffer.WriteString(jenny.typeFormatter.formatTypeDeclaration(object))
	buffer.WriteString("\n")

	if object.Type.ImplementsVariant() && !object.Type.IsRef() {
		variant := tools.UpperCamelCase(object.Type.ImplementedVariant())

		fmt.Fprintf(buffer, "func (resource %s) Implements%sVariant() {}\n", objectName, variant)
		buffer.WriteString("\n")

		customVariantTmpl := template.CustomObjectVariantBlock(object)
		if jenny.tmpl.Exists(customVariantTmpl) {
			if err := jenny.tmpl.RenderInBuffer(buffer, customVariantTmpl, map[string]any{
				"Object": object,
				"Schema": schema,
			}); err != nil {
				return err
			}
		}
	}

	return nil
}

func (jenny RawTypes) generateConstructor(buffer *strings.Builder, context languages.Context, object ast.Object) {
	objectName := formatObjectName(object.Name)
	constructorName := "New" + formatFunctionName(object.Name)

	declareConstructor := func() {
		jenny.apiRefCollector.RegisterFunction(object.SelfRef.ReferredPkg, common.FunctionReference{
			Name: constructorName,
			Comments: []string{
				fmt.Sprintf("%[1]s creates a new %[2]s object.", constructorName, objectName),
			},
			Return: "*" + objectName,
		})
	}

	if object.Type.IsRef() {
		referredObj, found := context.LocateObjectByRef(*object.Type.Ref)
		if !found || !referredObj.Type.IsStruct() {
			return
		}

		declareConstructor()
		fmt.Fprintf(buffer, "// %[1]s creates a new %[2]s object.\n", constructorName, objectName)
		fmt.Fprintf(buffer, "func %[1]s() *%[2]s {\n", constructorName, objectName)

		delegatedConstructorName := fmt.Sprintf("New%s", formatObjectName(referredObj.Name))
		referredPkg := jenny.packageMapper(referredObj.SelfRef.ReferredPkg)
		if referredPkg != "" {
			delegatedConstructorName = fmt.Sprintf("%s.%s", referredPkg, delegatedConstructorName)
		}

		fmt.Fprintf(buffer, "\treturn %s()", delegatedConstructorName)
		buffer.WriteString("\n}\n")
		return
	}

	if !object.Type.IsStruct() {
		return
	}

	declareConstructor()
	fmt.Fprintf(buffer, "// %[1]s creates a new %[2]s object.\n", constructorName, objectName)
	fmt.Fprintf(buffer, "func %[1]s() *%[2]s {\n", constructorName, objectName)
	fmt.Fprintf(buffer, "\treturn &%s", jenny.defaultsForStruct(context, object.SelfRef, object.Type, nil))
	buffer.WriteString("\n}\n")
}

func (jenny RawTypes) defaultsForStruct(context languages.Context, objectRef ast.RefType, objectType ast.Type, extraTyped *ast.TypedDefault) string {
	var buffer strings.Builder

	objectName := formatObjectName(objectRef.ReferredType)
	if referredPkg := jenny.packageMapper(objectRef.ReferredPkg); referredPkg != "" {
		objectName = referredPkg + "." + objectName
	}

	buffer.WriteString(objectName)
	buffer.WriteString("{\n")

	extraDefaults := structFieldDefaultsAsAny(extraTyped)

	for _, field := range objectType.Struct.Fields {
		resolvedFieldType := context.ResolveRefs(field.Type)

		if !fieldNeedsExplicitDefault(field, resolvedFieldType, extraDefaults) {
			continue
		}

		defaultValue, abort := jenny.defaultValueForField(context, field, resolvedFieldType, extraDefaults)
		if abort {
			break
		}

		fmt.Fprintf(&buffer, "\t\t%s: %s,\n", formatFieldName(field.Name), defaultValue)
	}

	buffer.WriteString("}")

	return buffer.String()
}

func fieldNeedsExplicitDefault(field ast.StructField, resolvedFieldType ast.Type, extraDefaults map[string]any) bool {
	return field.Type.TypedDefault != nil ||
		extraDefaults[field.Name] != nil ||
		(field.Required && field.Type.IsRef() && resolvedFieldType.IsStruct()) ||
		(field.Required && field.Type.IsArray()) ||
		(field.Required && field.Type.IsMap()) ||
		field.Type.IsConcreteScalar() ||
		field.Type.IsConstantRef()
}

// defaultValueForField computes the Go literal for a field's default. The
// second return value signals the caller to abort writing further fields
// (used when a constant ref resolves to something we can't render).
func (jenny RawTypes) defaultValueForField(context languages.Context, field ast.StructField, resolvedFieldType ast.Type, extraDefaults map[string]any) (string, bool) {
	if extraDefault, ok := extraDefaults[field.Name]; ok {
		return jenny.defaultFromExtraDefault(field, resolvedFieldType, extraDefault), false
	}

	switch {
	case field.Type.IsConcreteScalar():
		return jenny.maybeValueAsPointer(formatScalar(field.Type.Scalar.Value), field.Type.Nullable, resolvedFieldType), false
	case resolvedFieldType.IsMap() && field.Type.TypedDefault != nil:
		return jenny.maybeValueAsPointer(jenny.defaultForMapWithDefault(field, resolvedFieldType), field.Type.Nullable, resolvedFieldType), false
	case resolvedFieldType.IsArray() && field.Type.TypedDefault != nil:
		return jenny.maybeValueAsPointer(jenny.defaultForArrayWithDefault(field, resolvedFieldType), field.Type.Nullable, resolvedFieldType), false
	case resolvedFieldType.IsScalar() && field.Type.TypedDefault != nil:
		return jenny.maybeValueAsPointer(formatScalar(ast.TypedDefaultToAny(*field.Type.TypedDefault)), field.Type.Nullable, resolvedFieldType), false
	case field.Type.IsRef() && resolvedFieldType.IsStruct() && field.Type.TypedDefault != nil:
		return jenny.defaultForStructRefWithDefault(context, field, resolvedFieldType), false
	case field.Type.IsRef() && resolvedFieldType.IsStruct():
		return jenny.defaultForStructRef(field), false
	case field.Type.IsRef() && resolvedFieldType.IsEnum():
		return jenny.defaultForEnumRef(field, resolvedFieldType), false
	case field.Type.IsConstantRef():
		return jenny.defaultForConstantRef(context, field)
	case field.Type.IsArray():
		return "[]" + jenny.typeFormatter.formatType(field.Type.Array.ValueType) + "{}", false
	case field.Type.IsMap():
		return "map[" + jenny.typeFormatter.formatType(field.Type.Map.IndexType) + "]" + jenny.typeFormatter.formatType(field.Type.Map.ValueType) + "{}", false
	default:
		return "\"unsupported default value case: this is likely a bug in cog\"", false
	}
}

func (jenny RawTypes) defaultFromExtraDefault(field ast.StructField, resolvedFieldType ast.Type, extraDefault any) string {
	defaultValue := formatScalar(extraDefault)

	if !(field.Type.IsRef() && resolvedFieldType.IsStructGeneratedFromDisjunction()) {
		return jenny.maybeValueAsPointer(defaultValue, field.Type.Nullable, resolvedFieldType)
	}

	disjunctionBranchName := formatFieldName(anyToDisjunctionBranchName(extraDefault))
	disjunctionBranch, found := resolvedFieldType.Struct.FieldByName(disjunctionBranchName)
	if !found {
		disjunctionBranchName = "Any"
		disjunctionBranch, _ = resolvedFieldType.Struct.FieldByName(disjunctionBranchName)
	}

	actualDefault := jenny.maybeValueAsPointer(defaultValue, true, disjunctionBranch.Type)

	nonNullableRefType := field.Type.DeepCopy()
	nonNullableRefType.Nullable = false

	defaultValue = fmt.Sprintf("%s{\n\t%s: %s,\n}", jenny.typeFormatter.formatRef(nonNullableRefType, false), formatFieldName(disjunctionBranchName), actualDefault)

	if field.Type.Nullable {
		defaultValue = "&" + defaultValue
	}

	return defaultValue
}

func (jenny RawTypes) defaultForMapWithDefault(field ast.StructField, resolvedFieldType ast.Type) string {
	defaultAny := ast.TypedDefaultToAny(*field.Type.TypedDefault)
	if emptyMap, ok := defaultAny.(map[string]any); ok && len(emptyMap) == 0 {
		return "map[" + jenny.typeFormatter.formatType(resolvedFieldType.Map.IndexType) + "]" + jenny.typeFormatter.formatType(resolvedFieldType.Map.ValueType) + "{}"
	}
	return formatScalar(defaultAny)
}

func (jenny RawTypes) defaultForArrayWithDefault(field ast.StructField, resolvedFieldType ast.Type) string {
	defaultAny := ast.TypedDefaultToAny(*field.Type.TypedDefault)
	if emptySlice, ok := defaultAny.([]any); ok && len(emptySlice) == 0 {
		return "[]" + jenny.typeFormatter.formatType(resolvedFieldType.Array.ValueType) + "{}"
	}
	return formatScalar(defaultAny)
}

func (jenny RawTypes) defaultForStructRefWithDefault(context languages.Context, field ast.StructField, resolvedFieldType ast.Type) string {
	defaultValue := jenny.defaultsForStruct(context, *field.Type.Ref, resolvedFieldType, field.Type.TypedDefault)
	if field.Type.Nullable {
		defaultValue = "&" + defaultValue
	}
	return defaultValue
}

func (jenny RawTypes) defaultForStructRef(field ast.StructField) string {
	defaultValue := "New" + formatObjectName(field.Type.Ref.ReferredType) + "()"

	if referredPkg := jenny.packageMapper(field.Type.Ref.ReferredPkg); referredPkg != "" {
		defaultValue = referredPkg + "." + defaultValue
	}

	if !field.Type.Nullable {
		defaultValue = "*" + defaultValue
	}

	return defaultValue
}

func (jenny RawTypes) defaultForEnumRef(field ast.StructField, resolvedFieldType ast.Type) string {
	fieldDefault, _ := fieldDefaultAsAny(field.Type)
	memberName := resolvedFieldType.Enum.Values[0].Name
	for _, member := range resolvedFieldType.Enum.Values {
		if member.Value == fieldDefault {
			memberName = member.Name
			break
		}
	}

	defaultValue := memberName
	if referredPkg := jenny.packageMapper(field.Type.Ref.ReferredPkg); referredPkg != "" {
		defaultValue = referredPkg + "." + defaultValue
	}

	return jenny.maybeValueAsPointer(defaultValue, field.Type.Nullable, field.Type)
}

func (jenny RawTypes) defaultForConstantRef(context languages.Context, field ast.StructField) (string, bool) {
	constRef := field.Type.AsConstantRef()
	t := context.ResolveRefs(ast.NewRef(constRef.ReferredPkg, constRef.ReferredType))

	if !t.IsEnum() && !t.IsScalar() {
		return "", true
	}

	var defaultValue string
	if t.IsScalar() && t.AsScalar().ScalarKind == ast.KindString {
		defaultValue = constRef.ReferredType
	}

	if t.IsEnum() {
		for _, member := range t.AsEnum().Values {
			if member.Value == constRef.ReferenceValue {
				defaultValue = member.Name
				break
			}
		}
	}

	if referredPkg := jenny.packageMapper(constRef.ReferredPkg); referredPkg != "" {
		defaultValue = referredPkg + "." + defaultValue
	}

	return jenny.maybeValueAsPointer(defaultValue, field.Type.Nullable, field.Type), false
}

func (jenny RawTypes) maybeValueAsPointer(value string, nullable bool, typeDef ast.Type) string {
	if !nullable {
		return value
	}

	if typeDef.IsAnyOf(ast.KindArray, ast.KindMap) {
		return value
	}

	if typeDef.IsScalar() && typeDef.AsScalar().ScalarKind == ast.KindBytes {
		return value
	}

	nonNullableField := typeDef.DeepCopy()
	nonNullableField.Nullable = false
	typeHint := jenny.typeFormatter.formatType(nonNullableField)

	// we don't use cog.ToPtr() to avoid a dependency on cog's runtime
	return fmt.Sprintf("(func (input %[1]s) *%[1]s { return &input })(%[2]s)", typeHint, value)
}

// fieldDefaultAsAny returns the field's default value as a raw `any`. The
// second return value reports whether a default is present.
func fieldDefaultAsAny(fieldType ast.Type) (any, bool) {
	if fieldType.TypedDefault != nil {
		return ast.TypedDefaultToAny(*fieldType.TypedDefault), true
	}
	return nil, false
}

// structFieldDefaultsAsAny extracts per-field defaults from a struct (or
// ref-to-struct) TypedDefault as raw `any` values keyed by field name.
func structFieldDefaultsAsAny(td *ast.TypedDefault) map[string]any {
	if td == nil {
		return map[string]any{}
	}
	switch td.Kind {
	case ast.KindStruct:
		if td.Struct == nil {
			return nil
		}
		out := make(map[string]any, len(td.Struct.Fields))
		for name, fieldDefault := range td.Struct.Fields {
			out[name] = ast.TypedDefaultToAny(fieldDefault)
		}
		return out
	case ast.KindRef:
		if td.Ref == nil {
			return nil
		}
		return structFieldDefaultsAsAny(&td.Ref.Inner)
	}
	return nil
}
