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
	Config          Config
	Tmpl            *template.Template
	apiRefCollector *common.APIReferenceCollector

	typeFormatter *typeFormatter
	packageMapper func(pkg string) string
}

func (jenny RawTypes) JennyName() string {
	return "GoRawTypes"
}

func (jenny RawTypes) Generate(context languages.Context) (codejen.Files, error) {
	files := make(codejen.Files, 0, len(context.Schemas))

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

	imports := NewImportMap(jenny.Config.PackageRoot)
	jenny.packageMapper = func(pkg string) string {
		if imports.IsIdentical(pkg, schema.Package) {
			return ""
		}

		return imports.Add(pkg, jenny.Config.importPath(pkg))
	}
	jenny.typeFormatter = defaultTypeFormatter(jenny.Config, context, imports, jenny.packageMapper)
	unmarshallerGenerator := NewJSONMarshalling(jenny.Config, jenny.Tmpl, imports, jenny.packageMapper, jenny.typeFormatter, jenny.apiRefCollector)
	strictUnmarshallerGenerator := newStrictJSONUnmarshal(jenny.Tmpl, imports, jenny.packageMapper, jenny.typeFormatter, jenny.apiRefCollector)
	equalityMethodsGenerator := newEqualityMethods(jenny.Tmpl, jenny.apiRefCollector)
	validationMethodsGenerator := newValidationMethods(jenny.Tmpl, jenny.packageMapper, jenny.apiRefCollector)

	schema.Objects.Iterate(func(_ string, object ast.Object) {
		jenny.formatObject(&buffer, schema, object)
		buffer.WriteString("\n")

		jenny.generateConstructor(&buffer, context, object)

		innerErr := unmarshallerGenerator.generateForObject(&buffer, context, schema, object)
		if innerErr != nil {
			err = innerErr
			return
		}

		innerErr = strictUnmarshallerGenerator.generateForObject(&buffer, context, object)
		if innerErr != nil {
			err = innerErr
			return
		}

		innerErr = equalityMethodsGenerator.generateForObject(&buffer, context, object, imports)
		if innerErr != nil {
			err = innerErr
			return
		}

		innerErr = validationMethodsGenerator.generateForObject(&buffer, context, object, imports)
		if innerErr != nil {
			err = innerErr
			return
		}
	})
	if err != nil {
		return nil, err
	}

	if err := unmarshallerGenerator.generateForSchema(&buffer, schema); err != nil {
		return nil, err
	}

	importStatements := imports.String()
	if importStatements != "" {
		importStatements += "\n\n"
	}

	return []byte(fmt.Sprintf(`package %[1]s

%[2]s%[3]s`, formatPackageName(schema.Package), importStatements, buffer.String())), nil
}

func (jenny RawTypes) formatObject(buffer *strings.Builder, schema *ast.Schema, def ast.Object) {
	defName := formatObjectName(def.Name)

	comments := def.Comments
	if jenny.Config.debug {
		passesTrail := tools.Map(def.PassesTrail, func(trail string) string {
			return fmt.Sprintf("Modified by compiler pass '%s'", trail)
		})
		comments = append(comments, passesTrail...)
	}

	for _, commentLine := range comments {
		buffer.WriteString(fmt.Sprintf("// %s\n", commentLine))
	}

	buffer.WriteString(jenny.typeFormatter.formatTypeDeclaration(def))
	buffer.WriteString("\n")

	if def.Type.ImplementsVariant() && !def.Type.IsRef() {
		variant := tools.UpperCamelCase(def.Type.ImplementedVariant())

		buffer.WriteString(fmt.Sprintf("func (resource %s) Implements%sVariant() {}\n", defName, variant))
		buffer.WriteString("\n")

		if def.Type.ImplementedVariant() == string(ast.SchemaVariantDataQuery) {
			buffer.WriteString(fmt.Sprintf("func (resource %s) DataqueryType() string {\n", defName))
			buffer.WriteString(fmt.Sprintf("\treturn \"%s\"\n", strings.ToLower(schema.Metadata.Identifier)))
			buffer.WriteString("}\n")
		}
	}
}

func (jenny RawTypes) generateConstructor(buffer *strings.Builder, context languages.Context, object ast.Object) {
	defName := formatObjectName(object.Name)
	constructorName := "New" + defName

	declareConstructor := func() {
		jenny.apiRefCollector.RegisterFunction(object.SelfRef.ReferredPkg, common.FunctionReference{
			Name: constructorName,
			Comments: []string{
				fmt.Sprintf("%[1]s creates a new %[2]s object.", constructorName, defName),
			},
			Return: "*" + defName,
		})
	}

	if object.Type.IsRef() {
		referredObj, found := context.LocateObjectByRef(*object.Type.Ref)
		if !found || !referredObj.Type.IsStruct() {
			return
		}

		declareConstructor()
		buffer.WriteString(fmt.Sprintf("// %[1]s creates a new %[2]s object.\n", constructorName, defName))
		buffer.WriteString(fmt.Sprintf("func %[1]s() *%[2]s {\n", constructorName, defName))

		delegatedConstructorName := fmt.Sprintf("New%s", formatObjectName(referredObj.Name))
		referredPkg := jenny.packageMapper(referredObj.SelfRef.ReferredPkg)
		if referredPkg != "" {
			delegatedConstructorName = fmt.Sprintf("%s.%s", referredPkg, delegatedConstructorName)
		}

		buffer.WriteString(fmt.Sprintf("\treturn %s()", delegatedConstructorName))
		buffer.WriteString("\n}\n")
		return
	}

	if !object.Type.IsStruct() {
		return
	}

	declareConstructor()
	buffer.WriteString(fmt.Sprintf("// %[1]s creates a new %[2]s object.\n", constructorName, defName))
	buffer.WriteString(fmt.Sprintf("func %[1]s() *%[2]s {\n", constructorName, defName))
	buffer.WriteString(fmt.Sprintf("\treturn &%s", jenny.defaultsForStruct(context, object.SelfRef, object.Type, nil)))
	buffer.WriteString("\n}\n")
}

func (jenny RawTypes) defaultsForStruct(context languages.Context, objectRef ast.RefType, objectType ast.Type, maybeExtraDefaults any) string {
	var buffer strings.Builder

	defName := formatObjectName(objectRef.ReferredType)
	referredPkg := jenny.packageMapper(objectRef.ReferredPkg)
	if referredPkg != "" {
		defName = fmt.Sprintf("%s.%s", referredPkg, defName)
	}

	buffer.WriteString(fmt.Sprintf("%s{\n", defName))

	extraDefaults := map[string]any{}
	if val, ok := maybeExtraDefaults.(map[string]any); ok {
		extraDefaults = val
	}

	for _, field := range objectType.Struct.Fields {
		resolvedFieldType := context.ResolveRefs(field.Type)

		needsExplicitDefault := field.Type.Default != nil ||
			extraDefaults[field.Name] != nil ||
			(field.Required && field.Type.IsRef() && resolvedFieldType.IsStruct()) ||
			field.Type.IsConcreteScalar()
		if !needsExplicitDefault {
			continue
		}

		fieldName := formatFieldName(field.Name)
		defaultValue := ""

		if extraDefault, ok := extraDefaults[field.Name]; ok {
			defaultValue = formatScalar(extraDefault)

			defaultValue = jenny.maybeScalarValueAsPointer(defaultValue, field.Type.Nullable, resolvedFieldType)
		} else if field.Type.IsConcreteScalar() {
			defaultValue = formatScalar(field.Type.Scalar.Value)

			defaultValue = jenny.maybeScalarValueAsPointer(defaultValue, field.Type.Nullable, resolvedFieldType)
		} else if resolvedFieldType.IsAnyOf(ast.KindScalar, ast.KindMap, ast.KindArray) && field.Type.Default != nil {
			defaultValue = formatScalar(field.Type.Default)

			defaultValue = jenny.maybeScalarValueAsPointer(defaultValue, field.Type.Nullable, resolvedFieldType)
		} else if field.Type.IsRef() && resolvedFieldType.IsStruct() && field.Type.Default != nil {
			referredType := field.Type.Ref.ReferredType
			referredPkg = jenny.packageMapper(field.Type.Ref.ReferredPkg)
			if referredPkg != "" {
				referredType = fmt.Sprintf("%s.%s", referredPkg, referredType)
			}

			defaultValue = jenny.defaultsForStruct(context, *field.Type.Ref, resolvedFieldType, field.Type.Default)
			if field.Type.Nullable {
				defaultValue = "&" + defaultValue
			}
		} else if field.Type.IsRef() && resolvedFieldType.IsStruct() {
			defaultValue = fmt.Sprintf("New%s()", formatObjectName(field.Type.Ref.ReferredType))

			referredPkg = jenny.packageMapper(field.Type.Ref.ReferredPkg)
			if referredPkg != "" {
				defaultValue = fmt.Sprintf("%s.%s", referredPkg, defaultValue)
			}

			if !field.Type.Nullable {
				defaultValue = "*" + defaultValue
			}
		} else if field.Type.IsRef() && resolvedFieldType.IsEnum() {
			memberName := resolvedFieldType.Enum.Values[0].Name
			for _, member := range resolvedFieldType.Enum.Values {
				if member.Value == field.Type.Default {
					memberName = member.Name
					break
				}
			}

			defaultValue = memberName

			referredPkg = jenny.packageMapper(field.Type.Ref.ReferredPkg)
			if referredPkg != "" {
				defaultValue = fmt.Sprintf("%s.%s", referredPkg, defaultValue)
			}

			if field.Type.Nullable {
				jenny.packageMapper("cog")
				defaultValue = fmt.Sprintf("cog.ToPtr(%s)", defaultValue)
			}
		} else {
			defaultValue = "\"unsupported default value case: this is likely a bug in cog\""
		}

		buffer.WriteString(fmt.Sprintf("\t\t%s: %s,\n", fieldName, defaultValue))
	}

	buffer.WriteString("}")

	return buffer.String()
}

func (jenny RawTypes) maybeScalarValueAsPointer(value string, nullable bool, typeDef ast.Type) string {
	if nullable && typeDef.IsScalar() {
		nonNullableField := typeDef.DeepCopy()
		nonNullableField.Nullable = false
		typeHint := jenny.typeFormatter.formatType(nonNullableField)

		jenny.packageMapper("cog")
		return fmt.Sprintf("cog.ToPtr[%s](%s)", typeHint, value)
	}

	return value
}
