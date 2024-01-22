package compiler

import (
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/tools"
)

var _ Pass = (*Dashboard)(nil)

// Dashboard rewrites the definition of a few dashboard fields, that are incorrectly defined in the schema.
//
// Example: the "annotations" and "templating" fields are required, but the schema says otherwise.
type Dashboard struct {
}

func (pass *Dashboard) Process(schemas []*ast.Schema) ([]*ast.Schema, error) {
	for i, schema := range schemas {
		schemas[i] = pass.processSchema(schema)
	}

	return schemas, nil
}

func (pass *Dashboard) processSchema(schema *ast.Schema) *ast.Schema {
	for i, object := range schema.Objects {
		if schema.Package == dashboardPackage && object.Name == dashboardObject {
			schema.Objects[i] = pass.processDashboard(object)
			continue
		}
		if schema.Package == dashboardPackage && object.Name == dashboardVariableModelObject {
			schema.Objects[i] = pass.processVariableModel(object)
			continue
		}
	}

	return schema
}

func (pass *Dashboard) processDashboard(object ast.Object) ast.Object {
	if !object.Type.IsStruct() {
		return object
	}

	requiredFields := []string{"annotations", "templating"}

	for i, field := range object.Type.AsStruct().Fields {
		if tools.ItemInList(field.Name, requiredFields) {
			field.Type.Nullable = false
			field.Required = true

			object.Type.Struct.Fields[i] = field
		}
	}

	return object
}

// TODO: remove this once https://github.com/grafana/grafana/pull/79236 is merged
func (pass *Dashboard) processVariableModel(object ast.Object) ast.Object {
	if !object.Type.IsStruct() {
		return object
	}

	object.Type.Struct.Fields = append(object.Type.Struct.Fields,
		ast.NewStructField("allValue", ast.String()),
		ast.NewStructField("regex", ast.String()),
		ast.NewStructField("includeAll", ast.Bool(ast.Default(false))),
	)

	return object
}
