package terraform

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/languages"
)

type Datasources struct {
	Config Config
}

func (jenny Datasources) JennyName() string {
	return "TerraformDatasource"
}

func (jenny Datasources) Generate(context languages.Context) (codejen.Files, error) {
	files := make(codejen.Files, 0, len(context.Schemas))

	for _, schema := range context.Schemas {
		output, err := jenny.generateSchema(context, schema)
		if err != nil {
			return nil, err
		}

		filename := filepath.Join(formatPackageName(schema.Package), "datasource_gen.go")

		files = append(files, *codejen.NewFile(filename, output, jenny))
	}

	return files, nil
}

func (jenny Datasources) generateSchema(_ languages.Context, schema *ast.Schema) ([]byte, error) {
	var buffer strings.Builder

	//
	//structObjects := schema.Objects.Filter(func(_ string, object ast.Object) bool {
	//		return object.Type.IsStruct()
	//	})
	entryPoint, ok := schema.LocateObject(schema.EntryPoint)
	if !ok {
		log.Printf("Skipping schema %s as it lacks an entrypoint", schema.Package)
		return nil, nil
	}
	err := templates.ExecuteTemplate(&buffer, "types/datasources.tmpl", map[string]any{
		"Schema":     schema,
		"Entrypoint": entryPoint,
	})
	if err != nil {
		return nil, err
	}
	return []byte(buffer.String()), nil
}
