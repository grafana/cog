package compiler

import (
	"fmt"
	"strings"

	"github.com/grafana/cog/internal/ast"
)

// SetDatasourceToDataquery uses dashboard.DataSourceRef reference for the datasource field in each dataquery.
//
// Depending on the type of schema, this value can be an any or an internal Datasource struct generating an inconsistency
// between them.
type SetDatasourceToDataquery struct {
}

func (d *SetDatasourceToDataquery) Process(schemas []*ast.Schema) ([]*ast.Schema, error) {
	visitor := Visitor{
		OnSchema:      d.processSchema,
		OnStructField: d.processFieldStruct,
	}

	return visitor.VisitSchemas(schemas)
}

func (d *SetDatasourceToDataquery) processSchema(v *Visitor, schema *ast.Schema) (*ast.Schema, error) {
	var err error
	schema.Objects.Iterate(func(key string, value ast.Object) {
		if schema.Package != dashboardPackage && (key == datasourceName || key == dashboardDatasource) {
			schema.Objects.Remove(key)
		} else {
			obj, objErr := v.VisitObject(schema, value)
			if objErr != nil {
				err = objErr
			}

			schema.Objects.Set(key, obj)
		}
	})

	return schema, err
}

func (d *SetDatasourceToDataquery) processFieldStruct(_ *Visitor, schema *ast.Schema, field ast.StructField) (ast.StructField, error) {
	if schema.Package == dashboardPackage || schema.Metadata.Variant != ast.SchemaVariantDataQuery {
		return field, nil
	}

	if field.Name == strings.ToLower(datasourceName) || field.Name == dashboardDatasource {
		field.Type = ast.NewRef(dashboardPackage, dashboardDatasource)
		field.AddToPassesTrail(fmt.Sprintf("SetDatasourceToDataquery[%s.%s]", dashboardPackage, dashboardDatasource))
	}

	return field, nil
}
