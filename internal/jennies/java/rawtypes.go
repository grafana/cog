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

	jenny.imports = NewImportMap()
	packageMapper := func(pkg string, class string) string {
		if pkg == schema.Package {
			return ""
		}

		return jenny.imports.Add(class, pkg)
	}

	jenny.typeFormatter = defaultTypeFormatter(packageMapper)

	for i, object := range schema.Objects {
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
	case ast.KindDisjunction:
	case ast.KindIntersection:
	case ast.KindRef:
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

	err := templates.ExecuteTemplate(&buffer, "enum.tmpl", EnumTemplate{
		Package: pkg,
		Name:    name,
		Values:  values,
		Type:    formatScalarType(enum.Values[0].Type.AsScalar()),
	})

	if err != nil {
		return nil, err
	}

	return []byte(buffer.String()), nil
}

func (jenny RawTypes) formatStruct(pkg string, name string, def ast.StructType) ([]byte, error) {
	var buffer strings.Builder

	fields := make([]Field, len(def.Fields))

	for i, field := range def.Fields {
		fields[i] = Field{
			Name: field.Name,
			Type: jenny.typeFormatter.formatFieldType(field.Type),
		}
	}

	err := templates.ExecuteTemplate(&buffer, "struct.tmpl", ObjectTemplate{
		Package: pkg,
		Imports: jenny.imports,
		Name:    name,
		Fields:  fields,
	})

	if err != nil {
		return nil, err
	}

	return []byte(buffer.String()), nil
}
