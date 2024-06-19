package terraform

import (
	"path/filepath"
	"strings"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/languages"
)

type TerraformModels struct {
	Config Config
}

func (jenny TerraformModels) JennyName() string {
	return "TerraformModels"
}

func (jenny TerraformModels) Generate(context languages.Context) (codejen.Files, error) {
	files := make(codejen.Files, 0, len(context.Schemas))

	for _, schema := range context.Schemas {
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

func (jenny TerraformModels) generateSchema(_ languages.Context, schema *ast.Schema) ([]byte, error) {
	var buffer strings.Builder

	structObjects := schema.Objects.Filter(func(_ string, object ast.Object) bool {
		return object.Type.IsStruct()
	})
	err := templates.
		Funcs(map[string]any{
			"formatTerraformType": formatTerraformType,
		}).
		ExecuteTemplate(&buffer, "types/models.tmpl", map[string]any{
			"Schema":  schema,
			"Objects": structObjects.Values(),
		})
	if err != nil {
		return nil, err
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
