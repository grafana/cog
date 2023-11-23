package compiler

import (
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/tools"
)

var _ Pass = (*LibraryPanels)(nil)

// LibraryPanels rewrites the definition of the "LibraryPanel" object in the "librarypanel" package.
//
// In the original schema, the "model" field is left mainly undefined but a comment indicates
// that it should be the same panel schema defined in dashboard with a few fields omitted.
//
// This compiler pass implements the modifications described in that comment to define the
// "model" field as:
//
//	```
//	# In the LibraryPanel object
//	model: Omit<dashboard.Panel, 'gridPos' | 'id' | 'libraryPanel'>
//	```
//
// Note: this pass needs the "dashboard.Panel" schema to be parsed. Barring that, it leaves
// the schemas untouched.
type LibraryPanels struct {
}

func (pass *LibraryPanels) Process(schemas []*ast.Schema) ([]*ast.Schema, error) {
	dashboardPanelObj, found := ast.Schemas(schemas).LocateObject(dashboardPackage, dashboardPanelObject)
	if !found {
		return schemas, nil
	}

	newSchemas := make([]*ast.Schema, len(schemas))
	for i, schema := range schemas {
		if schema.Package != libraryPanelPackage {
			newSchemas[i] = schema
			continue
		}

		newSchemas[i] = pass.parseSchema(schema, dashboardPanelObj)
	}

	return newSchemas, nil
}

func (pass *LibraryPanels) parseSchema(schema *ast.Schema, dashboardPanel ast.Object) *ast.Schema {
	newSchema := schema.DeepCopy()
	newSchema.Objects = nil

	for _, object := range schema.Objects {
		if object.Name != libraryPanelObject {
			newSchema.Objects = append(newSchema.Objects, object)
			continue
		}

		newSchema.Objects = append(newSchema.Objects, pass.processLibraryPanel(object, dashboardPanel))
	}

	return &newSchema
}

func (pass *LibraryPanels) processLibraryPanel(object ast.Object, dashboardPanel ast.Object) ast.Object {
	if object.Type.Kind != ast.KindStruct {
		return object
	}

	structDef := object.Type.AsStruct()

	for i, field := range structDef.Fields {
		if field.Name != libraryPanelModelField {
			continue
		}

		structDef.Fields[i].Type = pass.buildModelType(dashboardPanel)
		break
	}

	return object
}

func (pass *LibraryPanels) buildModelType(dashboardPanel ast.Object) ast.Type {
	panelFields := dashboardPanel.Type.AsStruct().Fields
	fields := make([]ast.StructField, 0, len(panelFields))
	excludedFields := []string{
		dashboardPanelIDField,
		dashboardPanelGridPosField,
		dashboardPanelLibraryPanelField,
	}

	for _, panelField := range panelFields {
		if tools.ItemInList(panelField.Name, excludedFields) {
			continue
		}

		fields = append(fields, panelField)
	}

	return ast.NewStruct(fields...)
}
