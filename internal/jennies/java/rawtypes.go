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
	jsonMarshaller JSONMarshaller
}

func (jenny RawTypes) JennyName() string {
	return "JavaRawTypes"
}

func (jenny RawTypes) Generate(context languages.Context) (codejen.Files, error) {
	files := make(codejen.Files, 0)
	jenny.imports = NewImportMap(jenny.config.PackagePath)
	jenny.typeFormatter = createFormatter(context, jenny.config)
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
		"emptyValueForType":             jenny.typeFormatter.defaultValueFor,
		"shouldCastNilCheck":            jenny.typeFormatter.shouldCastNilCheck,
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
		output, innerErr := jenny.generateSchema(pkg, schema.Metadata.Identifier, object)
		if innerErr != nil {
			err = innerErr
			return
		}

		filename := filepath.Join(jenny.config.ProjectPath, pkg, fmt.Sprintf("%s.java", tools.UpperCamelCase(object.Name)))

		files = append(files, *codejen.NewFile(filename, output, jenny))
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

func (jenny RawTypes) generateSchema(pkg string, identifier string, object ast.Object) ([]byte, error) {
	switch object.Type.Kind {
	case ast.KindStruct:
		return jenny.formatStruct(pkg, identifier, object)
	case ast.KindEnum:
		return jenny.formatEnum(pkg, object)
	case ast.KindRef:
		return jenny.formatReference(pkg, identifier, object)
	case ast.KindIntersection:
		return jenny.formatIntersection(pkg, identifier, object)
	}

	return nil, nil
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
		Package:  jenny.config.formatPackage(pkg),
		Name:     object.Name,
		Values:   values,
		Type:     enumType,
		Comments: object.Comments,
	})
}

func (jenny RawTypes) formatStruct(pkg string, identifier string, object ast.Object) ([]byte, error) {
	return jenny.getTemplate().RenderAsBytes("types/class.tmpl", ClassTemplate{
		Package:                 jenny.config.formatPackage(pkg),
		Imports:                 jenny.imports,
		Name:                    tools.UpperCamelCase(object.Name),
		Fields:                  object.Type.AsStruct().Fields,
		Comments:                object.Comments,
		Variant:                 jenny.getVariant(object.Type),
		Identifier:              identifier,
		Annotation:              jenny.jsonMarshaller.annotation(object.Type),
		ToJSONFunction:          jenny.jsonMarshaller.genToJSONFunction(object.Type),
		ShouldAddSerializer:     jenny.typeFormatter.objectNeedsCustomSerializer(object),
		ShouldAddDeserializer:   jenny.typeFormatter.objectNeedsCustomDeserializer(object),
		ShouldAddFactoryMethods: object.Type.HasHint(ast.HintDisjunctionOfScalars) || object.Type.HasHint(ast.HintDiscriminatedDisjunctionOfRefs),
		DefaultConstructorArgs:  jenny.defaultConstructor(object),
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
		Package:   jenny.config.formatPackage(pkg),
		Name:      "Constants",
		Constants: constants,
	})
}

func (jenny RawTypes) formatReference(pkg string, identifier string, object ast.Object) ([]byte, error) {
	reference := jenny.typeFormatter.formatReference(object.Type.AsRef())

	return jenny.getTemplate().RenderAsBytes("types/class.tmpl", ClassTemplate{
		Package:    jenny.config.formatPackage(pkg),
		Imports:    jenny.imports,
		Name:       tools.UpperCamelCase(object.Name),
		Extends:    []string{reference},
		Comments:   object.Comments,
		Variant:    jenny.getVariant(object.Type),
		Identifier: identifier,
	})
}

func (jenny RawTypes) formatIntersection(pkg string, identifier string, object ast.Object) ([]byte, error) {
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
		Package:    jenny.config.formatPackage(pkg),
		Imports:    jenny.imports,
		Name:       object.Name,
		Extends:    extensions,
		Comments:   object.Comments,
		Fields:     fields,
		Variant:    jenny.getVariant(object.Type),
		Identifier: identifier,
	})
}

func (jenny RawTypes) getVariant(t ast.Type) string {
	variant := ""
	if t.ImplementsVariant() {
		variant = fmt.Sprintf("cog.variants.%s", tools.UpperCamelCase(t.ImplementedVariant()))
		variant = jenny.config.formatPackage(variant)
	}
	return variant
}

func (jenny RawTypes) defaultConstructor(object ast.Object) []ast.Argument {
	if object.Type.IsStructGeneratedFromDisjunction() {
		return nil
	}

	fields := object.Type.AsStruct().Fields

	args := make([]ast.Argument, len(fields))
	for i, field := range object.Type.AsStruct().Fields {
		args[i] = ast.Argument{
			Name: field.Name,
			Type: field.Type,
		}
	}

	return args
}
