package java

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
	config  Config
	tmpl    *template.Template
	imports *common.DirectImportMap

	typeFormatter  *typeFormatter
	builders       Builders
	jsonMarshaller JSONMarshaller
}

func (jenny RawTypes) JennyName() string {
	return "JavaRawTypes"
}

func (jenny RawTypes) Generate(context languages.Context) (codejen.Files, error) {
	files := make(codejen.Files, 0)
	jenny.imports = NewImportMap(jenny.config.PackagePath)
	jenny.typeFormatter = createFormatter(context, jenny.config)
	jenny.builders = parseBuilders(jenny.config, context, jenny.typeFormatter)
	jenny.jsonMarshaller = JSONMarshaller{
		config:        jenny.config,
		tmpl:          jenny.tmpl,
		typeFormatter: jenny.typeFormatter,
	}

	for _, schema := range context.Schemas {
		output, err := jenny.genFilesForSchema(schema)
		if err != nil {
			return nil, err
		}

		files = append(files, output...)
	}

	return files, nil
}

func (jenny RawTypes) getTemplate() *template.Template {
	return jenny.tmpl.Funcs(map[string]any{
		"formatBuilderFieldType":        jenny.typeFormatter.formatBuilderFieldType,
		"formatType":                    jenny.typeFormatter.formatFieldType,
		"typeHasBuilder":                jenny.typeFormatter.typeHasBuilder,
		"resolvesToComposableSlot":      jenny.typeFormatter.resolvesToComposableSlot,
		"emptyValueForType":             jenny.typeFormatter.defaultValueFor,
		"formatCastValue":               jenny.typeFormatter.formatCastValue,
		"formatAssignmentPath":          jenny.typeFormatter.formatAssignmentPath,
		"formatPath":                    jenny.typeFormatter.formatFieldPath,
		"shouldCastNilCheck":            jenny.typeFormatter.shouldCastNilCheck,
		"formatValue":                   jenny.typeFormatter.formatValue,
		"fillNullableAnnotationPattern": jenny.typeFormatter.fillNullableAnnotationPattern,
	})
}

func (jenny RawTypes) genFilesForSchema(schema *ast.Schema) (codejen.Files, error) {
	var err error
	files := make(codejen.Files, 0)
	scalars := make(map[string]ast.ScalarType)

	packageMapper := func(pkg string, class string) string {
		if jenny.imports.IsIdentical(pkg, schema.Package) {
			return ""
		}

		return jenny.imports.Add(class, pkg)
	}

	jenny.typeFormatter = jenny.typeFormatter.withPackageMapper(packageMapper)

	alreadyValidatedPanel := make(map[string]bool)
	schema.Objects.Iterate(func(_ string, object ast.Object) {
		jenny.imports = NewImportMap(jenny.config.PackagePath)
		if object.Type.IsMap() || object.Type.IsArray() {
			return
		}
		if object.Type.IsScalar() {
			if object.Type.AsScalar().IsConcrete() {
				scalars[object.Name] = object.Type.AsScalar()
			}
			return
		}

		pkg := formatPackageName(schema.Package)
		output, innerErr := jenny.generateSchema(pkg, object)
		if innerErr != nil {
			err = innerErr
			return
		}

		filename := filepath.Join(jenny.config.ProjectPath, pkg, fmt.Sprintf("%s.java", tools.UpperCamelCase(object.Name)))

		files = append(files, *codejen.NewFile(filename, output, jenny))

		// Because we need to check the package only, it could have multiple files, and we want to generate
		// the builder once.
		if !alreadyValidatedPanel[schema.Package] {
			panelOutput, innerErr := jenny.generatePanelBuilder(schema.Package)
			if innerErr != nil {
				err = innerErr
				return
			}

			if panelOutput != nil {
				alreadyValidatedPanel[schema.Package] = true
				filename := filepath.Join(jenny.config.ProjectPath, strings.ToLower(schema.Package), "PanelBuilder.java")
				files = append(files, *codejen.NewFile(filename, panelOutput, jenny))
			}
		}
	})

	if err != nil {
		return nil, err
	}

	if len(scalars) > 0 {
		output, err := jenny.formatScalars(schema.Package, scalars)
		if err != nil {
			return nil, err
		}

		filename := filepath.Join(jenny.config.ProjectPath, strings.ToLower(schema.Package), "Constants.java")
		files = append(files, *codejen.NewFile(filename, output, jenny))
	}

	return files, nil
}

func (jenny RawTypes) generateSchema(pkg string, object ast.Object) ([]byte, error) {
	switch object.Type.Kind {
	case ast.KindStruct:
		return jenny.formatStruct(pkg, object)
	case ast.KindEnum:
		return jenny.formatEnum(pkg, object)
	case ast.KindRef:
		return jenny.formatReference(pkg, object)
	case ast.KindIntersection:
		return jenny.formatIntersection(pkg, object)
	}

	return nil, nil
}

// generatePanelBuilder generates the builder for the panels. Panel's builders uses generic "Panel" name, and they don't match
// with any Schema name. These builders accept values from the different models of the panels, being easier to make it them independent.
func (jenny RawTypes) generatePanelBuilder(pkg string) ([]byte, error) {
	builder, hasBuilder := jenny.builders.genPanelBuilder(pkg)
	if !hasBuilder {
		return nil, nil
	}

	builder.Imports = jenny.imports

	return jenny.getTemplate().RenderAsBytes("types/panel_builder.tmpl", builder)
}

func (jenny RawTypes) formatEnum(pkg string, object ast.Object) ([]byte, error) {
	enum := object.Type.AsEnum()
	values := make([]EnumValue, 0)
	for _, value := range enum.Values {
		if value.Name == "" {
			value.Name = "None"
		}
		values = append(values, EnumValue{
			Name:  tools.UpperSnakeCase(value.Name),
			Value: value.Value,
		})
	}

	enumType := "Integer"
	if enum.Values[0].Type.AsScalar().ScalarKind == ast.KindString {
		enumType = "String"
	}

	// Adds empty value if it doesn't exist to avoid
	// to break in deserialization.
	if enumType == "String" {
		hasEmptyValue := false
		for _, value := range values {
			if value.Value == "" {
				hasEmptyValue = true
			}
		}

		if !hasEmptyValue {
			values = append(values, EnumValue{
				Name:  "_EMPTY",
				Value: "",
			})
		}
	}

	return jenny.getTemplate().RenderAsBytes("types/enum.tmpl", EnumTemplate{
		Package:  jenny.typeFormatter.formatPackage(pkg),
		Name:     object.Name,
		Values:   values,
		Type:     enumType,
		Comments: object.Comments,
	})
}

func (jenny RawTypes) formatStruct(pkg string, object ast.Object) ([]byte, error) {
	builders, hasBuilder := jenny.builders.genBuilders(pkg, object.Name)

	return jenny.getTemplate().RenderAsBytes("types/class.tmpl", ClassTemplate{
		Package:                 jenny.typeFormatter.formatPackage(pkg),
		Imports:                 jenny.imports,
		Name:                    tools.UpperCamelCase(object.Name),
		Fields:                  object.Type.AsStruct().Fields,
		Comments:                object.Comments,
		Variant:                 jenny.getVariant(object.Type),
		Builders:                builders,
		HasBuilder:              hasBuilder,
		Annotation:              jenny.jsonMarshaller.annotation(object.Type),
		ToJSONFunction:          jenny.jsonMarshaller.genToJSONFunction(object.Type),
		ShouldAddSerializer:     jenny.typeFormatter.objectNeedsCustomSerializer(object),
		ShouldAddDeserializer:   jenny.typeFormatter.objectNeedsCustomDeserializer(object),
		ShouldAddFactoryMethods: object.Type.HasHint(ast.HintDisjunctionOfScalars) || object.Type.HasHint(ast.HintDiscriminatedDisjunctionOfRefs),
		Defaults:                jenny.formatDefaults(object),
	})
}

func (jenny RawTypes) formatScalars(pkg string, scalars map[string]ast.ScalarType) ([]byte, error) {
	constants := make([]Constant, 0)
	for name, scalar := range scalars {
		constants = append(constants, Constant{
			Name:  name,
			Type:  formatScalarType(scalar),
			Value: scalar.Value,
		})
	}

	return jenny.getTemplate().RenderAsBytes("types/constants.tmpl", ConstantTemplate{
		Package:   jenny.typeFormatter.formatPackage(pkg),
		Name:      "Constants",
		Constants: constants,
	})
}

func (jenny RawTypes) formatReference(pkg string, object ast.Object) ([]byte, error) {
	reference := jenny.typeFormatter.formatReference(object.Type.AsRef())

	return jenny.getTemplate().RenderAsBytes("types/class.tmpl", ClassTemplate{
		Package:  jenny.typeFormatter.formatPackage(pkg),
		Imports:  jenny.imports,
		Name:     tools.UpperCamelCase(object.Name),
		Extends:  []string{reference},
		Comments: object.Comments,
		Variant:  jenny.getVariant(object.Type),
	})
}

func (jenny RawTypes) formatIntersection(pkg string, object ast.Object) ([]byte, error) {
	intersection := object.Type.AsIntersection()
	extensions := make([]string, 0)
	fields := make([]ast.StructField, 0)

	for _, branch := range intersection.Branches {
		switch branch.Kind {
		case ast.KindRef:
			extensions = append(extensions, jenny.typeFormatter.formatReference(branch.AsRef()))
		case ast.KindStruct:
			fields = append(fields, branch.AsStruct().Fields...)
		}
	}

	return jenny.getTemplate().RenderAsBytes("types/class.tmpl", ClassTemplate{
		Package:  jenny.typeFormatter.formatPackage(pkg),
		Imports:  jenny.imports,
		Name:     object.Name,
		Extends:  extensions,
		Comments: object.Comments,
		Fields:   fields,
		Variant:  jenny.getVariant(object.Type),
	})
}

func (jenny RawTypes) getVariant(t ast.Type) string {
	variant := ""
	if t.ImplementsVariant() {
		variant = fmt.Sprintf("cog.variants.%s", tools.UpperCamelCase(t.ImplementedVariant()))
		variant = jenny.typeFormatter.formatPackage(variant)
	}
	return variant
}

func (jenny RawTypes) formatDefaults(obj ast.Object) []OptionCall {
	calls := make([]OptionCall, 0)
	structDef := obj.Type.AsStruct()

	for _, field := range structDef.Fields {
		if field.Type.Default == nil {
			continue
		}

		calls = append(calls, OptionCall{
			Initializers: jenny.formatInitializers(field),
			OptionName:   tools.UpperCamelCase(field.Name),
			Args:         jenny.formatDefaultValues(field),
		})
	}

	return calls
}

// formatInitializers initialises objects with their defaults before set the value in the corresponding setter.
// TODO: It could have conflicts if we have different fields with the same kind of argument.
// TODO: It means that we need to initialize the objects with different names in that case.
func (jenny RawTypes) formatInitializers(field ast.StructField) []string {
	initializers := make([]string, 0)
	if !field.Type.IsRef() {
		return nil
	}

	ref := field.Type.AsRef()
	object, _ := jenny.typeFormatter.context.LocateObject(ref.ReferredPkg, ref.ReferredType)
	if !object.Type.IsStruct() {
		return nil
	}

	structType := object.Type.AsStruct()
	defValues := field.Type.Default.(map[string]interface{})

	initializers = append(initializers, fmt.Sprintf("%s %sResource = new %s();", ref.ReferredType, tools.LowerCamelCase(ref.ReferredType), ref.ReferredType))
	for _, field := range structType.Fields {
		if defVal, ok := defValues[field.Name]; ok {
			if field.Type.IsScalar() {
				initializers = append(initializers, fmt.Sprintf("%sResource.%s = %s;", tools.LowerCamelCase(ref.ReferredType), field.Name, formatType(field.Type.AsScalar().ScalarKind, defVal)))
			}
			// TODO: Implement lists if needed
		}
	}

	return initializers
}

func (jenny RawTypes) formatDefaultValues(field ast.StructField) []string {
	argumentList := make([]string, 0)
	switch field.Type.Kind {
	case ast.KindRef:
		argumentList = append(argumentList, jenny.formatDefaultReference(field.Type.AsRef(), field.Type.Default))
	case ast.KindScalar:
		scalar := field.Type.AsScalar()
		argumentList = append(argumentList, formatType(scalar.ScalarKind, field.Type.Default))
	case ast.KindArray:
		array := field.Type.AsArray()
		if array.IsArrayOfScalars() {
			argumentList = append(argumentList, fmt.Sprintf("List.of(%s)", jenny.typeFormatter.formatScalar(field.Type.Default)))
		}
	case ast.KindStruct:
		// TODO: Java is using veneers to avoid anonymous structs but it should be detailed if we need it at any moment.
		argumentList = append(argumentList, "new Object()")
	}
	// TODO: Implement the rest of types if any

	return argumentList
}

func (jenny RawTypes) formatDefaultReference(ref ast.RefType, defValue any) string {
	object, _ := jenny.typeFormatter.context.LocateObject(ref.ReferredPkg, ref.ReferredType)
	if object.Type.Kind == ast.KindEnum {
		for _, v := range object.Type.AsEnum().Values {
			if defValue == v.Value {
				return fmt.Sprintf("%s.%s", object.Name, tools.UpperSnakeCase(v.Name))
			}
		}
	}

	return fmt.Sprintf("%sResource", tools.LowerCamelCase(object.Name))
}
