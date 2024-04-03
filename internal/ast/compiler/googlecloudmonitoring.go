package compiler

import (
	"github.com/grafana/cog/internal/ast"
)

var _ Pass = (*GoogleCloudMonitoring)(nil)

// GoogleCloudMonitoring rewrites a part of the googlecloudmonitoring schema.
//
// Older schemas (pre 10.2.x) define `CloudMonitoringQuery.timeSeriesList`
// as a disjunction that cog can't handle: `timeSeriesList?: #TimeSeriesList | #AnnotationQuery`,
// where `AnnotationQuery` is a type that extends `TimeSeriesList` to add two
// fields.
//
// This compiler pass checks for the presence of that disjunction, and rewrites
// it as a reference to `TimeSeriesList`. It also adds the two missing fields
// to this type if they aren't already defined.
type GoogleCloudMonitoring struct {
}

func (pass *GoogleCloudMonitoring) Process(schemas []*ast.Schema) ([]*ast.Schema, error) {
	for i, schema := range schemas {
		if schema.Package != "googlecloudmonitoring" {
			continue
		}

		newSchema, err := pass.processSchema(schema)
		if err != nil {
			return nil, err
		}

		schemas[i] = newSchema
	}

	return schemas, nil
}

func (pass *GoogleCloudMonitoring) processSchema(schema *ast.Schema) (*ast.Schema, error) {
	schema.Objects = schema.Objects.Map(func(_ string, object ast.Object) ast.Object {
		if object.Name == "CloudMonitoringQuery" {
			return pass.processCloudMonitoringQuery(object)
		}
		if object.Name == "TimeSeriesList" {
			return pass.processTimeSeriesList(object)
		}

		return object
	})

	return schema, nil
}

func (pass *GoogleCloudMonitoring) processCloudMonitoringQuery(object ast.Object) ast.Object {
	if !object.Type.IsStruct() {
		return object
	}

	structDef := object.Type.AsStruct()
	fields := make([]ast.StructField, 0, len(structDef.Fields)-1)

	for _, field := range structDef.Fields {
		if field.Name != "timeSeriesList" {
			fields = append(fields, field)
			continue
		}

		if !field.Type.IsDisjunction() {
			fields = append(fields, field)
			continue
		}

		// from `timeSeriesList?: #TimeSeriesList | #AnnotationQuery`
		// to `timeSeriesList?: #TimeSeriesList`
		newField := field.DeepCopy()
		newField.Type = newField.Type.Disjunction.Branches[0]
		newField.AddToPassesTrail("GoogleCloudMonitoring[removed disjunction]")

		fields = append(fields, newField)
	}

	object.Type.Struct.Fields = fields

	return object
}

func (pass *GoogleCloudMonitoring) processTimeSeriesList(object ast.Object) ast.Object {
	if !object.Type.IsStruct() {
		return object
	}

	structDef := object.Type.AsStruct()

	if _, found := structDef.FieldByName("title"); !found {
		field := ast.NewStructField("title", ast.String(ast.Nullable()))
		field.Comments = []string{"Annotation title."}
		field.AddToPassesTrail("GoogleCloudMonitoring[created]")

		structDef.Fields = append(structDef.Fields, field)
	}
	if _, found := structDef.FieldByName("text"); !found {
		field := ast.NewStructField("text", ast.String(ast.Nullable()))
		field.Comments = []string{"Annotation text."}
		field.AddToPassesTrail("GoogleCloudMonitoring[created]")

		structDef.Fields = append(structDef.Fields, field)
	}

	object.Type.Struct = &structDef

	return object
}
