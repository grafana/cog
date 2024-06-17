package php

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
	config Config

	typeFormatter *typeFormatter
}

func (jenny RawTypes) JennyName() string {
	return "PHPRawTypes"
}

func (jenny RawTypes) Generate(context languages.Context) (codejen.Files, error) {
	files := make(codejen.Files, 0, len(context.Schemas))

	for _, schema := range context.Schemas {
		schemaFiles, err := jenny.generateSchema(context, schema)
		if err != nil {
			return nil, err
		}

		files = append(files, schemaFiles...)
	}

	return files, nil
}

func (jenny RawTypes) generateSchema(context languages.Context, schema *ast.Schema) (codejen.Files, error) {
	var err error

	files := make(codejen.Files, 0, schema.Objects.Len())
	schema.Objects.Iterate(func(_ string, object ast.Object) {
		file, innerErr := jenny.formatObject(context, schema, object)
		if innerErr != nil {
			err = innerErr
			return
		}

		files = append(files, file)
	})
	if err != nil {
		return nil, err
	}

	return files, nil
}

func (jenny RawTypes) formatObject(context languages.Context, schema *ast.Schema, def ast.Object) (codejen.File, error) {
	var buffer strings.Builder

	jenny.typeFormatter = defaultTypeFormatter(jenny.config, context)

	defName := formatObjectName(def.Name)

	comments := def.Comments
	if jenny.config.debug {
		passesTrail := tools.Map(def.PassesTrail, func(trail string) string {
			return fmt.Sprintf("Modified by compiler pass '%s'", trail)
		})
		comments = append(comments, passesTrail...)
	}

	buffer.WriteString(formatCommentsBlock(comments))

	switch def.Type.Kind {
	case ast.KindEnum:
		enum, err := jenny.formatEnumDef(def)
		if err != nil {
			return codejen.File{}, err
		}

		buffer.WriteString(enum)
	case ast.KindScalar:
		scalarType := def.Type.AsScalar()

		//nolint: gocritic
		if scalarType.IsConcrete() {
			buffer.WriteString(fmt.Sprintf("const %s = %s;", formatConstantName(def.Name), formatScalar(scalarType.Value)))
		} else {
			return codejen.File{}, fmt.Errorf("type aliases on scalar types is not supported")
		}
	case ast.KindRef:
		buffer.WriteString(fmt.Sprintf("class %s extends %s {}", defName, jenny.typeFormatter.formatType(def.Type)))
	case ast.KindStruct:
		buffer.WriteString(jenny.formatStructDef(def))
	default:
		return codejen.File{}, fmt.Errorf("unhandled type def kind: %s", def.Type.Kind)
	}

	buffer.WriteString("\n")

	filename := filepath.Join(
		"src",
		"Types",
		formatPackageName(schema.Package),
		fmt.Sprintf("%s.php", defName),
	)

	output := fmt.Sprintf("<?php\n\nnamespace %s\\%s;\n\n", jenny.config.fullNamespace("Types"), formatPackageName(schema.Package))
	output += buffer.String()

	return *codejen.NewFile(filename, []byte(output), jenny), nil
}

func (jenny RawTypes) formatStructDef(def ast.Object) string {
	var buffer strings.Builder

	variant := ""
	if def.Type.ImplementsVariant() {
		variant = ", " + jenny.config.fullNamespaceRef("Runtime\\Variants\\"+formatObjectName(def.Type.ImplementedVariant()))
	}

	buffer.WriteString(fmt.Sprintf("class %s implements \\JsonSerializable%s {\n", formatObjectName(def.Name), variant))

	for _, fieldDef := range def.Type.Struct.Fields {
		buffer.WriteString(tools.Indent(jenny.typeFormatter.formatField(fieldDef), 4))
		buffer.WriteString("\n\n")
	}

	buffer.WriteString(tools.Indent(jenny.generateJSONSerialize(def), 4))
	buffer.WriteString("\n}")

	return buffer.String()
}

func (jenny RawTypes) generateJSONSerialize(def ast.Object) string {
	var buffer strings.Builder

	buffer.WriteString("/**\n")
	buffer.WriteString(" * @return array<string, mixed>\n")
	buffer.WriteString(" */\n")
	buffer.WriteString("public function jsonSerialize(): array\n")
	buffer.WriteString("{\n")

	buffer.WriteString("    $data = [\n")

	for _, field := range def.Type.AsStruct().Fields {
		if !field.Required {
			continue
		}

		buffer.WriteString(fmt.Sprintf(`        "%s" => $this->%s,`+"\n", field.Name, formatFieldName(field.Name)))
	}

	buffer.WriteString("    ];\n")

	for _, field := range def.Type.AsStruct().Fields {
		if field.Required {
			continue
		}

		fieldName := formatFieldName(field.Name)

		buffer.WriteString(fmt.Sprintf("    if (isset($this->%s)) {\n", fieldName))
		buffer.WriteString(fmt.Sprintf(`        $data["%s"] = $this->%s;`+"\n", field.Name, fieldName))
		buffer.WriteString("    }\n")
	}

	buffer.WriteString("    return $data;\n")

	buffer.WriteString("}")

	return buffer.String()
}

func (jenny RawTypes) formatEnumDef(def ast.Object) (string, error) {
	return renderTemplate("types/enum.tmpl", map[string]any{
		"Object": def,
	})
}
