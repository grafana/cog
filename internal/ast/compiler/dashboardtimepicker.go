package compiler

import (
	"github.com/grafana/cog/internal/ast"
)

var _ Pass = (*DashboardTimePicker)(nil)

type DashboardTimePicker struct {
}

func (pass *DashboardTimePicker) Process(schemas []*ast.Schema) ([]*ast.Schema, error) {
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

func (pass *DashboardTimePicker) processSchema(schema *ast.Schema) (*ast.Schema, error) {
	var timepickerObject ast.Object
	var dashboardObject ast.Object

	for i, object := range schema.Objects {
		if object.Name != "Dashboard" {
			continue
		}

		dashboardObject, timepickerObject = pass.processDashboard(object)

		schema.Objects[i] = dashboardObject
	}

	// did we actually define a new object?
	if timepickerObject.Name != "" {
		schema.Objects = append(schema.Objects, timepickerObject)
	}

	return schema, nil
}

func (pass *DashboardTimePicker) processDashboard(object ast.Object) (ast.Object, ast.Object) {
	var timepickerObject ast.Object

	pkg := object.SelfRef.ReferredPkg

	for i, field := range object.Type.AsStruct().Fields {
		if field.Name != "timepicker" {
			continue
		}

		// we expect the timepicker field to define an inline struct
		if field.Type.Kind != ast.KindStruct {
			continue
		}

		timepickerObject = ast.NewObject(pkg, "TimePicker", field.Type)
		object.Type.AsStruct().Fields[i].Type = ast.NewRef(pkg, "TimePicker")
	}

	return object, timepickerObject
}
