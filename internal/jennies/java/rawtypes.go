package java

import (
	"fmt"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/tools"
)

type RawTypes struct {
	imports *common.DirectImportMap

	typeFormatter *typeFormatter
	builders      Builders
}

func (jenny RawTypes) JennyName() string {
	return "JavaRawTypes"
}

func (jenny RawTypes) Generate(context common.Context) (codejen.Files, error) {
	files := make(codejen.Files, 0)
	jenny.imports = NewImportMap()
	jenny.typeFormatter = createFormatter(context)
	jenny.builders = parseBuilders(context, jenny.typeFormatter)

	for _, schema := range context.Schemas {
		output, err := jenny.genFilesForSchema(schema)
		if err != nil {
			return nil, err
		}

		files = append(files, output...)
	}

	return files, nil
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
		jenny.imports = NewImportMap()
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

		filename := filepath.Join(
			pkg,
			fmt.Sprintf("%s.java", tools.UpperCamelCase(object.Name)),
		)

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

		filename := filepath.Join(strings.ToLower(schema.Package), "Constants.java")
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

	err := templates.ExecuteTemplate(&buffer, "types/enum.tmpl", EnumTemplate{
		Package:  pkg,
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

	if err := templates.Funcs(template.FuncMap{
		"formatBuilderFieldType":   jenny.typeFormatter.formatBuilderFieldType,
		"formatType":               jenny.typeFormatter.formatFieldType,
		"typeHasBuilder":           jenny.typeFormatter.typeHasBuilder,
		"resolvesToComposableSlot": jenny.typeFormatter.resolvesToComposableSlot,
		"emptyValueForType": func(typeDef ast.Type) string {
			return jenny.typeFormatter.defaultValueFor(typeDef)
		},
	}).ExecuteTemplate(&buffer, "types/class.tmpl", ClassTemplate{
		Package:    pkg,
		Imports:    jenny.imports,
		Name:       tools.UpperCamelCase(object.Name),
		Fields:     fields,
		Comments:   object.Comments,
		Variant:    tools.UpperCamelCase(object.Type.ImplementedVariant()),
		Builder:    builder,
		HasBuilder: hasBuilder,
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

	if err := templates.ExecuteTemplate(&buffer, "types/constants.tmpl", ConstantTemplate{
		Package:   pkg,
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

	if err := templates.ExecuteTemplate(&buffer, "types/class.tmpl", ClassTemplate{
		Package:  pkg,
		Imports:  jenny.imports,
		Name:     tools.UpperCamelCase(object.Name),
		Extends:  []string{reference},
		Comments: object.Comments,
		Variant:  tools.UpperCamelCase(object.Type.ImplementedVariant()),
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

	if err := templates.ExecuteTemplate(&buffer, "types/class.tmpl", ClassTemplate{
		Package:  pkg,
		Imports:  jenny.imports,
		Name:     object.Name,
		Extends:  extensions,
		Comments: object.Comments,
		Fields:   fields,
		Variant:  tools.UpperCamelCase(object.Type.ImplementedVariant()),
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

// TODO: Need to say to the serializer the correct name.
func escapeVarName(varName string) string {
	if isReservedJavaKeyword(varName) {
		return varName + "Arg"
	}

	return varName
}

// nolint: gocyclo
func isReservedJavaKeyword(input string) bool {
	// see https://docs.oracle.com/javase/tutorial/java/nutsandbolts/_keywords.html
	switch input {
	case "static", "abstract", "enum", "class", "if", "else", "switch", "final", "public", "private", "protected", "package", "continue", "new", "for", "assert",
		"do", "default", "goto", "synchronized", "boolean", "double", "int", "short", "char", "float", "long", "byte", "break", "throw", "throws", "this",
		"implements", "transient", "return", "catch", "extends", "case", "try", "void", "volatile", "super", "native", "finally", "instanceof", "import", "while":
		return true
	}
	return false
}
