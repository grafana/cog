package terraform

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/languages"
)

type RawTypes struct {
	Config Config

	typeFormatter *typeFormatter
}

func (jenny RawTypes) JennyName() string {
	return "TerraformRawTypes"
}

func (jenny RawTypes) Generate(context languages.Context) (codejen.Files, error) {
	files := make(codejen.Files, 0, len(context.Schemas))

	for _, schema := range context.Schemas {
		if schema.Metadata.Identifier != "Dashboard" {
			continue
		}
		log.Printf("SCHEMA: %s / %s / %s", schema.Metadata.Kind, schema.Metadata.Identifier, schema.Metadata.Variant)
		output, err := jenny.generateSchema(context, schema)
		if err != nil {
			return nil, err
		}

		filename := filepath.Join(
			formatPackageName(schema.Package),
			"terraform_provider_gen.go",
		)

		files = append(files, *codejen.NewFile(filename, output, jenny))
	}

	return files, nil
}

func (jenny RawTypes) generateSchema(context languages.Context, schema *ast.Schema) ([]byte, error) {
	var buffer strings.Builder

	err := templates.
		Funcs(map[string]any{
			"formatTerraformType": formatTerraformType,
		}).
		ExecuteTemplate(&buffer, "types/provider.tmpl", map[string]any{
			"Schema":  schema,
			"Objects": schema.Objects.Values(),
		})
	if err != nil {
		return nil, err
	}
	for i, line := range strings.Split(buffer.String(), "\n") {
		log.Printf("  %d: %s", i, line)
	}
	return []byte(buffer.String()), nil
}

func formatTerraformType(t ast.Type) string {
	if t.IsScalar() {
		tt := t.AsScalar()
		scalarType := "unknown"

		switch tt.ScalarKind {
		case ast.KindString, ast.KindBytes:
			scalarType = "types.String"
		case ast.KindInt8, ast.KindUint8, ast.KindInt16, ast.KindUint16:
			scalarType = "types.Int64"
		case ast.KindInt32, ast.KindUint32:
			scalarType = "types.Int64"
		case ast.KindInt64, ast.KindUint64:
			scalarType = "types.Int64"
		case ast.KindFloat32:
			scalarType = "types.Float64"
		case ast.KindFloat64:
			scalarType = "types.Float64"
		case ast.KindBool:
			scalarType = "types.Bool"
		case ast.KindAny:
			scalarType = "types.Object"
		}
		return scalarType
	}
	return string(t.Kind)
}
