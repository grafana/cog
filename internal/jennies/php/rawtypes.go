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
	Config Config

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

	jenny.typeFormatter = defaultTypeFormatter(jenny.Config, context)

	defName := formatObjectName(def.Name)

	comments := def.Comments
	if jenny.Config.debug {
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
		variant := ""
		if def.Type.ImplementsVariant() {
			variant = " implements \\Runtime\\Variants\\" + formatObjectName(def.Type.ImplementedVariant())
		}

		buffer.WriteString(fmt.Sprintf("class %s%s\n%s", defName, variant, jenny.typeFormatter.formatType(def.Type)))
	default:
		return codejen.File{}, fmt.Errorf("unhandled type def kind: %s", def.Type.Kind)
	}

	buffer.WriteString("\n")

	filename := filepath.Join(
		"Types",
		formatPackageName(schema.Package),
		fmt.Sprintf("%s.php", defName),
	)

	output := fmt.Sprintf("<?php\n\nnamespace Types\\%s;\n\n", formatPackageName(schema.Package))
	output += buffer.String()

	return *codejen.NewFile(filename, []byte(output), jenny), nil
}

func (jenny RawTypes) formatEnumDef(def ast.Object) (string, error) {
	return renderTemplate("types/enum.tmpl", map[string]any{
		"Object": def,
	})
}
