package compiler

import (
	"github.com/grafana/cog/internal/ast"
)

var _ Pass = (*DashboardTargetsRewrite)(nil)

type DashboardTargetsRewrite struct {
}

func (pass *DashboardTargetsRewrite) Process(schemas []*ast.Schema) ([]*ast.Schema, error) {
	newSchemas := make([]*ast.Schema, 0, len(schemas))

	for _, schema := range schemas {
		if schema.Package != "dashboard" {
			newSchemas = append(newSchemas, schema)
			continue
		}

		newSchema, err := pass.processSchema(schema)
		if err != nil {
			return nil, err
		}

		newSchemas = append(newSchemas, newSchema)
	}

	return newSchemas, nil
}

func (pass *DashboardTargetsRewrite) processSchema(schema *ast.Schema) (*ast.Schema, error) {
	newSchema := schema.DeepCopy()
	newSchema.Objects = nil

	for _, object := range schema.Objects {
		if object.Name == "Target" {
			newObject := object.DeepCopy()
			newObject.Type = ast.NewComposabilitySlot(ast.SchemaVariantDataQuery)

			newSchema.Objects = append(newSchema.Objects, newObject)
			continue
		}

		newSchema.Objects = append(newSchema.Objects, object)
	}

	return &newSchema, nil
}
