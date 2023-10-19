package compiler

import (
	"github.com/grafana/cog/internal/ast"
)

var _ Pass = (*DashboardPanelsRewrite)(nil)

type DashboardPanelsRewrite struct {
}

func (pass *DashboardPanelsRewrite) Process(schemas []*ast.Schema) ([]*ast.Schema, error) {
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

func (pass *DashboardPanelsRewrite) processSchema(schema *ast.Schema) (*ast.Schema, error) {
	newSchema := schema.DeepCopy()
	newSchema.Objects = nil

	for _, object := range schema.Objects {
		if object.Name == "Dashboard" {
			disjunction := ast.NewDisjunction([]ast.Type{
				ast.NewRef(schema.Package, "Panel"),
				ast.NewRef(schema.Package, "RowPanel"),
			})
			disjunction.Disjunction.Discriminator = "type"
			disjunction.Disjunction.DiscriminatorMapping = map[string]string{
				"row":                     "RowPanel",
				ast.DiscriminatorCatchAll: "Panel",
			}

			newPanelsFieldType := ast.NewArray(disjunction)

			newSchema.Objects = append(newSchema.Objects, pass.overwritePanelsFieldType(object, newPanelsFieldType))
			continue
		}
		if object.Name == "RowPanel" {
			newPanelsFieldType := ast.NewArray(ast.NewRef(schema.Package, "Panel"))

			newSchema.Objects = append(newSchema.Objects, pass.overwritePanelsFieldType(object, newPanelsFieldType))
			continue
		}

		newSchema.Objects = append(newSchema.Objects, object)
	}

	return &newSchema, nil
}

func (pass *DashboardPanelsRewrite) overwritePanelsFieldType(object ast.Object, newPanelsFieldType ast.Type) ast.Object {
	newFields := make([]ast.StructField, len(object.Type.AsStruct().Fields))
	for i, field := range object.Type.AsStruct().Fields {
		if field.Name != "panels" {
			newFields[i] = field
			continue
		}

		newField := field
		newField.Type = newPanelsFieldType

		newFields[i] = newField
	}

	newObject := object
	newObject.Type.Struct.Fields = newFields

	return newObject
}
