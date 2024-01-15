package java

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/tools"
)

type RawTypes struct {
	config  Config
	imports *common.DirectImportMap

	typeFormatter *typeFormatter
}

func (jenny RawTypes) JennyName() string {
	return "JavaBuilder"
}

func (jenny RawTypes) Generate(context common.Context) (codejen.Files, error) {
	files := make(codejen.Files, 0)
	jenny.imports = NewImportMap()
	jenny.typeFormatter = createFormatter(context)

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
	files := make(codejen.Files, 0)
	scalars := make(map[string]ast.ScalarType)

	packageMapper := func(pkg string, class string) string {
		if pkg == schema.Package {
			return ""
		}

		return jenny.imports.Add(class, pkg)
	}

	jenny.typeFormatter = jenny.typeFormatter.withPackageMapper(packageMapper)

	for _, object := range schema.Objects {
		jenny.imports = NewImportMap()
		if object.Type.IsMap() {
			continue
		}
		if object.Type.IsScalar() {
			scalars[object.Name] = object.Type.AsScalar()
			continue
		}

		output, err := jenny.generateSchema(schema.Package, object)
		if err != nil {
			return nil, err
		}

		filename := filepath.Join(
			strings.ToLower(schema.Package),
			fmt.Sprintf("%s.java", tools.UpperCamelCase(object.Name)),
		)

		files = append(files, *codejen.NewFile(filename, output, jenny))
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
		values[i] = EnumValue{
			Name:  strings.ToUpper(value.Name),
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

	if err := templates.ExecuteTemplate(&buffer, "types/class.tmpl", jenny.formatInnerStruct(pkg, object.Name, object.Comments, object.Type.ImplementedVariant(), object.Type.AsStruct())); err != nil {
		return nil, err
	}

	return []byte(buffer.String()), nil
}

func (jenny RawTypes) formatInnerStruct(pkg string, name string, comments []string, variant string, def ast.StructType) ClassTemplate {
	fields := make([]Field, 0)
	nestedStructs := make([]ClassTemplate, 0)

	for _, field := range def.Fields {
		if field.Type.Kind == ast.KindStruct {
			nestedStructs = append(nestedStructs, jenny.formatInnerStruct(pkg, field.Name, field.Comments, field.Type.ImplementedVariant(), field.Type.AsStruct()))
		} else {
			fields = append(fields, Field{
				Name:     field.Name,
				Type:     jenny.typeFormatter.formatFieldType(field.Type),
				Comments: field.Comments,
			})
		}
	}

	return ClassTemplate{
		Package:              pkg,
		Imports:              jenny.imports,
		Name:                 tools.UpperCamelCase(name),
		Fields:               fields,
		InnerClasses:         nestedStructs,
		GenGettersAndSetters: jenny.config.GenGettersAndSetters,
		Comments:             comments,
		Variant:              tools.UpperCamelCase(variant),
	}
}

func (jenny RawTypes) formatScalars(pkg string, scalars map[string]ast.ScalarType) ([]byte, error) {
	var buffer strings.Builder

	constants := make([]Constant, 0)
	for name, scalar := range scalars {
		if scalar.IsConcrete() {
			constants = append(constants, Constant{
				Name:  name,
				Type:  formatScalarType(scalar),
				Value: scalar.Value,
			})
		}
	}

	if len(constants) == 0 {
		return nil, nil
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
		Name:     object.Name,
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
	if isReservedGoKeyword(varName) {
		return varName + "Arg"
	}

	return varName
}

// nolint: gocyclo
func isReservedGoKeyword(input string) bool {
	// see https://docs.oracle.com/javase/tutorial/java/nutsandbolts/_keywords.html
	return input == "static" ||
		input == "abstract" ||
		input == "enum" ||
		input == "class" ||
		input == "if" ||
		input == "else" ||
		input == "switch" ||
		input == "final" ||
		input == "public" ||
		input == "private" ||
		input == "protected" ||
		input == "package" ||
		input == "continue" ||
		input == "new" ||
		input == "for" ||
		input == "assert" ||
		input == "do" ||
		input == "default" ||
		input == "goto" ||
		input == "synchronized" ||
		input == "boolean" ||
		input == "double" ||
		input == "int" ||
		input == "short" ||
		input == "char" ||
		input == "float" ||
		input == "long" ||
		input == "byte" ||
		input == "break" ||
		input == "throw" ||
		input == "throws" ||
		input == "this" ||
		input == "implements" ||
		input == "transient" ||
		input == "return" ||
		input == "catch" ||
		input == "extends" ||
		input == "case" ||
		input == "try" ||
		input == "void" ||
		input == "volatile" ||
		input == "super" ||
		input == "native" ||
		input == "finally" ||
		input == "instanceof" ||
		input == "import" ||
		input == "while"
}