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
	files := make(codejen.Files, len(schema.Objects))
	packageMapper := func(pkg string, class string) string {
		if pkg == schema.Package {
			return ""
		}

		return jenny.imports.Add(class, pkg)
	}

	jenny.typeFormatter = defaultTypeFormatter(packageMapper)

	for i, object := range schema.Objects {
		jenny.imports = NewImportMap()
		output, err := jenny.generateSchema(schema.Package, object)
		if err != nil {
			return nil, err
		}

		filename := filepath.Join(
			strings.ToLower(schema.Package),
			fmt.Sprintf("%s.java", tools.UpperCamelCase(object.Name)),
		)

		files[i] = *codejen.NewFile(filename, output, jenny)
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
		// TODO
	case ast.KindMap:
		// TODO
	case ast.KindScalar:
		// TODO
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
		Package:      pkg,
		Imports:      jenny.imports,
		Name:         tools.UpperCamelCase(name),
		Fields:       fields,
		InnerClasses: nestedStructs,
	}
}
