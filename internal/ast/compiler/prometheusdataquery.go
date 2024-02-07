package compiler

import (
	"github.com/grafana/cog/internal/ast"
)

var _ Pass = (*PrometheusDataquery)(nil)

// PrometheusDataquery rewrites the definition of the prometheus.Dataquery type and adds a few missing fields.
// Note: this pass is meant to be removed once the schema is up-to-date.
type PrometheusDataquery struct {
}

func (pass *PrometheusDataquery) Process(schemas []*ast.Schema) ([]*ast.Schema, error) {
	for i, schema := range schemas {
		schemas[i] = pass.processSchema(schema)
	}

	return schemas, nil
}

func (pass *PrometheusDataquery) processSchema(schema *ast.Schema) *ast.Schema {
	for i, object := range schema.Objects {
		if schema.Package == prometheusPackage && object.Name == prometheusDataqueryObject {
			schema.Objects[i] = pass.processDataquery(object)
			continue
		}
	}

	return schema
}

func (pass *PrometheusDataquery) processDataquery(object ast.Object) ast.Object {
	if !object.Type.IsStruct() {
		return object
	}

	if _, exists := object.Type.AsStruct().FieldByName("interval"); !exists {
		object.Type.Struct.Fields = append(object.Type.Struct.Fields,
			ast.NewStructField("interval", ast.String(), ast.Comments([]string{"An additional lower limit for the step parameter of the Prometheus query and for the", "`$__interval` and `$__rate_interval` variables."})),
		)
	}

	return object
}
