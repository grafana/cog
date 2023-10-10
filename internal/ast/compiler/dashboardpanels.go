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
		if object.Name == "dashboard" {
			newSchema.Objects = append(newSchema.Objects, pass.processDashboardObject(schema, object))
			continue
		}
		if object.Name == "RowPanel" {
			newSchema.Objects = append(newSchema.Objects, pass.processRowPanelObject(schema, object))
			continue
		}

		newSchema.Objects = append(newSchema.Objects, object)
	}

	return &newSchema, nil
}

func (pass *DashboardPanelsRewrite) processDashboardObject(schema *ast.Schema, object ast.Object) ast.Object {
	newFields := make([]ast.StructField, len(object.Type.AsStruct().Fields))
	for i, field := range object.Type.AsStruct().Fields {
		if field.Name == "panels" {
			newFields[i] = pass.processPanelsStructField(schema, field)
			continue
		}

		newFields[i] = field
	}

	newObject := object
	newObject.Type.Struct.Fields = newFields

	return newObject
}

func (pass *DashboardPanelsRewrite) processRowPanelObject(schema *ast.Schema, object ast.Object) ast.Object {
	newFields := make([]ast.StructField, len(object.Type.AsStruct().Fields))
	for i, field := range object.Type.AsStruct().Fields {
		if field.Name == "panels" {
			newFields[i] = pass.processPanelsStructField(schema, field)
			continue
		}

		newFields[i] = field
	}

	newObject := object
	newObject.Type.Struct.Fields = newFields

	return newObject
}

func (pass *DashboardPanelsRewrite) processPanelsStructField(schema *ast.Schema, field ast.StructField) ast.StructField {
	disjunction := ast.NewDisjunction([]ast.Type{
		ast.NewRef(schema.Package, "Panel"),
		ast.NewRef(schema.Package, "RowPanel"),
	})
	disjunction.Disjunction.Discriminator = "type"
	disjunction.Disjunction.DiscriminatorMapping = map[string]any{
		"RowPanel": "row",
		"Panel":    ast.DiscriminatorCatchAll,
	}

	newField := field
	newField.Type = ast.NewArray(disjunction)

	return newField
}
