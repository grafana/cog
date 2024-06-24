package terraform

import (
	"path/filepath"
	"strings"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/languages"
)

type Models struct {
	Config Config
}

func (jenny Models) JennyName() string {
	return "TerraformModels"
}

func (jenny Models) Generate(context languages.Context) (codejen.Files, error) {
	files := make(codejen.Files, 0, len(context.Schemas))

	for _, schema := range context.Schemas {
		output, err := jenny.generateSchema(context, schema)
		if err != nil {
			return nil, err
		}

		filename := filepath.Join(formatPackageName(schema.Package), "models_gen.go")

		files = append(files, *codejen.NewFile(filename, output, jenny))
	}

	return files, nil
}

func (jenny Models) generateSchema(_ languages.Context, schema *ast.Schema) ([]byte, error) {
	var buffer strings.Builder

	// schema.LocateObject(schema.EntryPoint)

	structObjects := schema.Objects.Filter(func(_ string, object ast.Object) bool {
		return object.Type.IsStruct()
	})
	err := templates.ExecuteTemplate(&buffer, "types/models.tmpl", map[string]any{
		"Schema":  schema,
		"Objects": structObjects.Values(),
	})
	if err != nil {
		return nil, err
	}
	return []byte(buffer.String()), nil
}
