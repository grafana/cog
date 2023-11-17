package compiler

import (
	"github.com/grafana/cog/internal/ast"
)

const (
	LibraryPanelPackage    = "librarypanel"
	LibraryPanelObject     = "LibraryPanel"
	LibraryPanelModelField = "model"
	DashboardPackage       = "dashboard"
	DashboardPanelObject   = "Panel"
)

type LibraryPanels struct {
}

func (lp *LibraryPanels) Process(schemas []*ast.Schema) ([]*ast.Schema, error) {
	dashboardRef := lp.getDashboardSchema(schemas)
	if dashboardRef == nil {
		return schemas, nil
	}

	newSchemas := make([]*ast.Schema, len(schemas))
	for i, schema := range schemas {
		if schema.Package != LibraryPanelPackage {
			newSchemas[i] = schema
			continue
		}

		newSchemas[i] = lp.parseSchema(schema, dashboardRef)
	}

	return newSchemas, nil
}

func (lp *LibraryPanels) parseSchema(schema *ast.Schema, dashboardRef *ast.RefType) *ast.Schema {
	newSchema := schema.DeepCopy()
	newSchema.Objects = nil

	for _, object := range schema.Objects {
		if object.Name != LibraryPanelObject {
			newSchema.Objects = append(newSchema.Objects, object)
			continue
		}

		newSchema.Objects = append(newSchema.Objects, lp.processObject(object, dashboardRef))
	}

	return &newSchema
}

func (lp *LibraryPanels) processObject(object ast.Object, dashboardRef *ast.RefType) ast.Object {
	if object.Type.Kind != ast.KindStruct {
		return object
	}

	structDef := object.Type.AsStruct()
	fields := make([]ast.StructField, 0, len(structDef.Fields))

	for _, field := range structDef.Fields {
		if field.Name != LibraryPanelModelField {
			fields = append(fields, field)
			continue
		}

		newField := field.DeepCopy()
		newField.Type = ast.NewRef(dashboardRef.ReferredPkg, dashboardRef.ReferredType)
		fields = append(fields, newField)
	}

	object.Type.Struct.Fields = fields
	return object
}

func (lp *LibraryPanels) getDashboardSchema(schemas []*ast.Schema) *ast.RefType {
	for _, schema := range schemas {
		if schema.Package == DashboardPackage {
			return &ast.RefType{ReferredPkg: schema.Package, ReferredType: DashboardPanelObject}
		}
	}

	return nil
}
