package golang

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/languages"
	"github.com/grafana/cog/internal/tools"
)

type RawTypes struct {
	Config Config

	typeFormatter *typeFormatter
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

	imports := NewImportMap()
	packageMapper := func(pkg string) string {
		if imports.IsIdentical(pkg, schema.Package) {
			return ""
		}

		return imports.Add(pkg, jenny.Config.importPath(pkg))
	}
	jenny.typeFormatter = defaultTypeFormatter(jenny.Config, context, packageMapper)
	unmarshallerGenerator := JSONMarshalling{
		config:        jenny.Config,
		packageMapper: packageMapper,
		typeFormatter: jenny.typeFormatter,
	}

	schema.Objects.Iterate(func(_ string, object ast.Object) {
		objectOutput, innerErr := jenny.formatObject(object)
		if innerErr != nil {
			err = innerErr
			return
		}

		buffer.Write(objectOutput)
		buffer.WriteString("\n")

		innerErr = unmarshallerGenerator.generateForObject(&buffer, context, schema, object)
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

func (jenny RawTypes) formatObject(def ast.Object) ([]byte, error) {
	var buffer strings.Builder

	defName := tools.UpperCamelCase(def.Name)

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

	switch def.Type.Kind {
	case ast.KindEnum:
		buffer.WriteString(jenny.formatEnumDef(def))
	case ast.KindScalar:
		scalarType := def.Type.AsScalar()

		//nolint: gocritic
		if scalarType.Value != nil {
			buffer.WriteString(fmt.Sprintf("const %s = %s", defName, formatScalar(scalarType.Value)))
		} else if scalarType.ScalarKind == ast.KindBytes {
			buffer.WriteString(fmt.Sprintf("type %s %s", defName, "[]byte"))
		} else {
			buffer.WriteString(fmt.Sprintf("type %s %s", defName, jenny.typeFormatter.formatType(def.Type)))
		}
	case ast.KindRef:
		buffer.WriteString(fmt.Sprintf("type %s = %s", defName, jenny.typeFormatter.formatType(def.Type)))
	case ast.KindMap, ast.KindArray, ast.KindStruct, ast.KindIntersection:
		buffer.WriteString(fmt.Sprintf("type %s %s", defName, jenny.typeFormatter.formatType(def.Type)))
	default:
		return nil, fmt.Errorf("unhandled type def kind: %s", def.Type.Kind)
	}

	buffer.WriteString("\n")

	if def.Type.ImplementsVariant() {
		variant := tools.UpperCamelCase(def.Type.ImplementedVariant())

		buffer.WriteString(fmt.Sprintf("func (resource %s) Implements%sVariant() {}\n", defName, variant))
		buffer.WriteString("\n")
	}

	return []byte(buffer.String()), nil
}

func (jenny RawTypes) formatEnumDef(def ast.Object) string {
	var buffer strings.Builder

	enumName := tools.UpperCamelCase(def.Name)
	enumType := def.Type.AsEnum()

	buffer.WriteString(fmt.Sprintf("type %s %s\n", enumName, jenny.typeFormatter.formatType(enumType.Values[0].Type)))

	buffer.WriteString("const (\n")
	for _, val := range enumType.Values {
		name := tools.CleanupNames(tools.UpperCamelCase(val.Name))
		buffer.WriteString(fmt.Sprintf("\t%s %s = %#v\n", name, enumName, val.Value))
	}
	buffer.WriteString(")\n")

	return buffer.String()
}
