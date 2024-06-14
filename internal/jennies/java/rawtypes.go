package java

import (
	"fmt"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/languages"
	"github.com/grafana/cog/internal/tools"
)

type RawTypes struct {
	config  Config
	imports *common.DirectImportMap

	typeFormatter *typeFormatter
	builders      Builders
}

func (jenny RawTypes) JennyName() string {
	return "JavaRawTypes"
}

func (jenny RawTypes) Generate(context languages.Context) (codejen.Files, error) {
	files := make(codejen.Files, 0)
	jenny.imports = NewImportMap(jenny.config.PackagePath)
	jenny.typeFormatter = createFormatter(context, jenny.config)
	jenny.builders = parseBuilders(jenny.config, context, jenny.typeFormatter)

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
	return templates.Funcs(map[string]any{
		"formatBuilderFieldType":   jenny.typeFormatter.formatBuilderFieldType,
		"formatType":               jenny.typeFormatter.formatFieldType,
		"typeHasBuilder":           jenny.typeFormatter.typeHasBuilder,
		"resolvesToComposableSlot": jenny.typeFormatter.resolvesToComposableSlot,
		"emptyValueForType":        jenny.typeFormatter.defaultValueFor,
		"formatCastValue":          jenny.typeFormatter.formatCastValue,
		"formatAssignmentPath":     jenny.typeFormatter.formatAssignmentPath,
		"formatPath":               jenny.typeFormatter.formatFieldPath,
		"shouldCastNilCheck":       jenny.typeFormatter.shouldCastNilCheck,
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

	alreadyValidatedBuilder := make(map[string]bool)
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

		if !alreadyValidatedBuilder[schema.Package] {
			buildersOutput, innerErr := jenny.generateOutsideBuilder(schema.Package)
			if innerErr != nil {
				err = innerErr
				return
			}

			files = append(files, buildersOutput...)
			alreadyValidatedBuilder[schema.Package] = true
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
func (jenny RawTypes) generateOutsideBuilder(pkg string) (codejen.Files, error) {
	builders, hasBuilder := jenny.builders.genExternalBuilders(pkg)
	if !hasBuilder {
		return nil, nil
	}

	files := make([]codejen.File, 0)
	for name, builder := range builders {
		builder.Imports = jenny.imports
		var buffer strings.Builder
		if err := jenny.getTemplate().ExecuteTemplate(&buffer, "types/external_builder.tmpl", builder); err != nil {
			return nil, err
		}

		filename := filepath.Join(jenny.config.ProjectPath, strings.ToLower(pkg), fmt.Sprintf("%sBuilder.java", name))
		files = append(files, *codejen.NewFile(filename, []byte(buffer.String()), jenny))
	}

	return files, nil
}

func (jenny RawTypes) formatEnum(pkg string, object ast.Object) ([]byte, error) {
	var buffer strings.Builder

	enum := object.Type.AsEnum()
	values := make([]EnumValue, len(enum.Values))
	for i, value := range enum.Values {
		if value.Name == "" {
			value.Name = "None"
		}
		values[i] = EnumValue{
			Name:  tools.UpperSnakeCase(value.Name),
			Value: value.Value,
		}
	}

	enumType := "Integer"
	if enum.Values[0].Type.AsScalar().ScalarKind == ast.KindString {
		enumType = "String"
	}

	err := jenny.getTemplate().ExecuteTemplate(&buffer, "types/enum.tmpl", EnumTemplate{
		Package:  jenny.typeFormatter.formatPackage(pkg),
		Name:     object.Name,
		Values:   values,
		Type:     enumType,
		Comments: object.Comments,
	})

	if err != nil {
		return nil, err
	}

	return []byte(buffer.String()), nil
}

func (jenny RawTypes) formatStruct(pkg string, object ast.Object) ([]byte, error) {
	var buffer strings.Builder

	fields := make([]Field, 0)
	for _, field := range object.Type.AsStruct().Fields {
		fields = append(fields, Field{
			Name:     tools.LowerCamelCase(field.Name),
			Type:     jenny.typeFormatter.formatFieldType(field.Type),
			Comments: field.Comments,
		})
	}

	builder, hasBuilder := jenny.builders.genBuilder(pkg, object.Name)
	jenny.addJSONImportsIfNeeded()

	if err := jenny.getTemplate().ExecuteTemplate(&buffer, "types/class.tmpl", ClassTemplate{
		Package:              jenny.typeFormatter.formatPackage(pkg),
		Imports:              jenny.imports,
		Name:                 tools.UpperCamelCase(object.Name),
		Fields:               fields,
		Comments:             object.Comments,
		Variant:              jenny.getVariant(object.Type),
		Builder:              builder,
		HasBuilder:           hasBuilder,
		ShouldAddMarshalling: jenny.config.GeneratePOM,
	}); err != nil {
		return nil, err
	}

	return []byte(buffer.String()), nil
}

func (jenny RawTypes) formatScalars(pkg string, scalars map[string]ast.ScalarType) ([]byte, error) {
	var buffer strings.Builder

	constants := make([]Constant, 0)
	for name, scalar := range scalars {
		constants = append(constants, Constant{
			Name:  name,
			Type:  formatScalarType(scalar),
			Value: scalar.Value,
		})
	}

	if err := jenny.getTemplate().ExecuteTemplate(&buffer, "types/constants.tmpl", ConstantTemplate{
		Package:   jenny.typeFormatter.formatPackage(pkg),
		Name:      "Constants",
		Constants: constants,
	}); err != nil {
		return nil, err
	}

	return []byte(buffer.String()), nil
}

func (jenny RawTypes) formatReference(pkg string, object ast.Object) ([]byte, error) {
	var buffer strings.Builder
	reference := jenny.typeFormatter.formatReference(object.Type.AsRef())

	if err := jenny.getTemplate().ExecuteTemplate(&buffer, "types/class.tmpl", ClassTemplate{
		Package:  jenny.typeFormatter.formatPackage(pkg),
		Imports:  jenny.imports,
		Name:     tools.UpperCamelCase(object.Name),
		Extends:  []string{reference},
		Comments: object.Comments,
		Variant:  jenny.getVariant(object.Type),
	}); err != nil {
		return nil, err
	}

	return []byte(buffer.String()), nil
}

func (jenny RawTypes) formatIntersection(pkg string, object ast.Object) ([]byte, error) {
	var buffer strings.Builder

	intersection := object.Type.AsIntersection()
	extensions := make([]string, 0)
	fields := make([]Field, 0)

	for _, branch := range intersection.Branches {
		switch branch.Kind {
		case ast.KindRef:
			extensions = append(extensions, jenny.typeFormatter.formatReference(branch.AsRef()))
		case ast.KindStruct:
			fields = append(fields, jenny.formatFields(branch.AsStruct())...)
		}
	}

	if err := jenny.getTemplate().ExecuteTemplate(&buffer, "types/class.tmpl", ClassTemplate{
		Package:  jenny.typeFormatter.formatPackage(pkg),
		Imports:  jenny.imports,
		Name:     object.Name,
		Extends:  extensions,
		Comments: object.Comments,
		Fields:   fields,
		Variant:  jenny.getVariant(object.Type),
	}); err != nil {
		return nil, err
	}

	return []byte(buffer.String()), nil
}

func (jenny RawTypes) formatFields(def ast.StructType) []Field {
	fields := make([]Field, len(def.Fields))
	for i, field := range def.Fields {
		fields[i] = Field{
			Name:     field.Name,
			Type:     jenny.typeFormatter.formatFieldType(field.Type),
			Comments: field.Comments,
		}
	}

	return fields
}

func (jenny RawTypes) getVariant(t ast.Type) string {
	variant := ""
	if t.ImplementsVariant() {
		variant = fmt.Sprintf("cog.variants.%s", tools.UpperCamelCase(t.ImplementedVariant()))
		variant = jenny.typeFormatter.formatPackage(variant)
	}
	return variant
}

func (jenny RawTypes) addJSONImportsIfNeeded() {
	if jenny.config.GeneratePOM {
		jenny.typeFormatter.packageMapper("com.fasterxml.jackson", "annotation.JsonProperty")
		jenny.typeFormatter.packageMapper("com.fasterxml.jackson", "core.JsonProcessingException")
		jenny.typeFormatter.packageMapper("com.fasterxml.jackson", "databind.ObjectMapper")
		jenny.typeFormatter.packageMapper("com.fasterxml.jackson", "databind.ObjectWriter")
	}
}
