package java

import (
	"fmt"
	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/tools"
	"path/filepath"
	"strings"
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
		return jenny.formatStruct(pkg, object.Name, object.Type.AsStruct())
	case ast.KindEnum:
		return jenny.formatEnum(pkg, object.Name, object.Type.AsEnum())
	case ast.KindRef:
		return jenny.formatReference(pkg, object.Name, object.Type.AsRef())
	}

	return nil, nil
}

func (jenny RawTypes) formatEnum(pkg string, name string, enum ast.EnumType) ([]byte, error) {
	var buffer strings.Builder
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

	err := templates.ExecuteTemplate(&buffer, "enum.tmpl", EnumTemplate{
		Package: pkg,
		Name:    name,
		Values:  values,
		Type:    enumType,
	})

	if err != nil {
		return nil, err
	}

	return []byte(buffer.String()), nil
}

func (jenny RawTypes) formatStruct(pkg string, name string, def ast.StructType) ([]byte, error) {
	var buffer strings.Builder

	if err := templates.ExecuteTemplate(&buffer, "class.tmpl", jenny.formatInnerStruct(pkg, name, def)); err != nil {
		return nil, err
	}

	return []byte(buffer.String()), nil
}

func (jenny RawTypes) formatInnerStruct(pkg string, name string, def ast.StructType) ObjectTemplate {
	fields := make([]Field, 0)
	nestedStructs := make([]ObjectTemplate, 0)

	for _, field := range def.Fields {
		if field.Type.Kind == ast.KindStruct {
			nestedStructs = append(nestedStructs, jenny.formatInnerStruct(pkg, field.Name, field.Type.AsStruct()))
		} else {
			fields = append(fields, Field{
				Name: field.Name,
				Type: jenny.typeFormatter.formatFieldType(field.Type),
			})
		}
	}

	return ObjectTemplate{
		Package:              pkg,
		Imports:              jenny.imports,
		Name:                 tools.UpperCamelCase(name),
		Fields:               fields,
		InnerClasses:         nestedStructs,
		GenGettersAndSetters: jenny.config.GenGettersAndSetters,
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

	if err := templates.ExecuteTemplate(&buffer, "constants.tmpl", ConstantTemplate{
		Package:   pkg,
		Name:      "Constants",
		Constants: constants,
	}); err != nil {
		return nil, err
	}

	return []byte(buffer.String()), nil
}

func (jenny RawTypes) formatReference(pkg string, name string, ref ast.RefType) ([]byte, error) {
	var buffer strings.Builder
	reference := jenny.typeFormatter.formatReference(ref)

	if err := templates.ExecuteTemplate(&buffer, "class.tmpl", ObjectTemplate{
		Package: pkg,
		Imports: jenny.imports,
		Name:    name,
		Extends: []string{reference},
	}); err != nil {
		return nil, err
	}

	return []byte(buffer.String()), nil
}
