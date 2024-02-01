package compiler

import (
	"github.com/grafana/cog/internal/ast"
)

var _ Pass = (*DashboardPanelsRewrite)(nil)

// DashboardPanelsRewrite rewrites the definition of "panels" fields in the "dashboard" package.
//
// In the original schema, panels are defined as follows:
//
//	```
//	# In the Dashboard object
//	panels?: [...#Panel | #RowPanel | #GraphPanel | #HeatmapPanel]
//
//	# In the RowPanel object
//	panels: [...#Panel | #GraphPanel | #HeatmapPanel]
//	```
//
// These definitions become:
//
//	```
//	# In the Dashboard object
//	panels?: [...#Panel | #RowPanel]
//
//	# In the RowPanel object
//	panels: [...#Panel]
//	```
type DashboardPanelsRewrite struct {
}

func (pass *DashboardPanelsRewrite) Process(schemas []*ast.Schema) ([]*ast.Schema, error) {
	newSchemas := make([]*ast.Schema, 0, len(schemas))

	for _, schema := range schemas {
		if schema.Package != dashboardPackage {
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
	schema.Objects = schema.Objects.Map(func(_ string, object ast.Object) ast.Object {
		if object.Name == dashboardObject {
			disjunction := ast.NewDisjunction([]ast.Type{
				ast.NewRef(schema.Package, dashboardPanelObject),
				ast.NewRef(schema.Package, dashboardRowPanelObject),
			})
			disjunction.Disjunction.Discriminator = dashboardPanelTypeField
			disjunction.Disjunction.DiscriminatorMapping = map[string]string{
				"row":                     dashboardRowPanelObject,
				ast.DiscriminatorCatchAll: dashboardPanelObject,
			}

			newPanelsFieldType := ast.NewArray(disjunction)

			return pass.overwritePanelsFieldType(object, newPanelsFieldType)
		}
		if object.Name == dashboardRowPanelObject {
			newPanelsFieldType := ast.NewArray(ast.NewRef(schema.Package, "Panel"))

			return pass.overwritePanelsFieldType(object, newPanelsFieldType)
		}

		return object
	})

	return schema, nil
}

func (pass *DashboardPanelsRewrite) overwritePanelsFieldType(object ast.Object, newPanelsFieldType ast.Type) ast.Object {
	newFields := make([]ast.StructField, len(object.Type.AsStruct().Fields))
	for i, field := range object.Type.AsStruct().Fields {
		if field.Name != dashboardPanelsField {
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
