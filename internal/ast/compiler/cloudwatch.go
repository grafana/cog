package compiler

import (
	"github.com/grafana/cog/internal/ast"
)

var _ Pass = (*Cloudwatch)(nil)

// Cloudwatch rewrites a part of the cloudwatch schema.
//
// In that schema, the `QueryEditorExpression` type is defined as a disjunction
// for which the discriminator and mapping can not be inferred.
// This compiler pass is here to define that mapping.
//
// The `QueryEditorArrayExpression` struct type is also modified to simplify the
// definition of its `expression` field from `[...#QueryEditorExpression] | [...#QueryEditorArrayExpression]` to
// `[...#QueryEditorExpression]`.
// This should be semantically equivalent since `#QueryEditorExpression` is a
// union type that includes `#QueryEditorArrayExpression`.
//
// The Cloudwatch pass also alerts the definition of the `#CloudWatchMetricsQuery`, `#CloudWatchLogsQuery` and
// `#CloudWatchAnnotationQuery` types.
// It removes the "dataquery variant" hint they carry, and defines a `CloudWatchQuery` type instead as a disjunction.
// That disjunction serves as "dataquery entrypoint" for cloudwatch.
type Cloudwatch struct {
}

func (pass *Cloudwatch) Process(schemas []*ast.Schema) ([]*ast.Schema, error) {
	newSchemas := make([]*ast.Schema, 0, len(schemas))

	for _, schema := range schemas {
		if schema.Package != "cloudwatch" {
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

func (pass *Cloudwatch) processSchema(schema *ast.Schema) (*ast.Schema, error) {
	schema.Objects = schema.Objects.Map(func(_ string, object ast.Object) ast.Object {
		if object.Name == "QueryEditorExpression" {
			return pass.processQueryEditorExpression(object)
		}
		if object.Name == "QueryEditorArrayExpression" {
			return pass.processQueryEditorArrayExpression(object)
		}

		// types hinted as a dataquery are replaced by a "CloudWatchQuery" disjunction,
		// serving as a "main entrypoint" for cloudwatch queries.
		if object.Type.ImplementsVariant() && object.Type.ImplementedVariant() == string(ast.SchemaVariantDataQuery) {
			object.Type = pass.processDataquery(object.Name, object.Type)
		}

		return object
	})

	schema.AddObject(pass.defineQueryDisjunction(schema))

	return schema, nil
}

func (pass *Cloudwatch) processDataquery(objectName string, typeDef ast.Type) ast.Type {
	typeDef.Hints[ast.HintSkipVariantPluginRegistration] = true

	if typeDef.Kind != ast.KindStruct {
		return typeDef
	}

	for i, field := range typeDef.AsStruct().Fields {
		if field.Name == "queryMode" {
			switch objectName {
			case "CloudWatchMetricsQuery":
				field.Type.Default = "Metrics"
			case "CloudWatchLogsQuery":
				field.Type.Default = "Logs"
			case "CloudWatchAnnotationQuery":
				field.Type.Default = "Annotations"
			}

			field.Type.Nullable = false
			field.Required = true
			typeDef.Struct.Fields[i] = field
		}
	}

	return typeDef
}

func (pass *Cloudwatch) defineQueryDisjunction(schema *ast.Schema) ast.Object {
	cloudwatchQuery := ast.NewDisjunction(ast.Types{
		ast.NewRef(schema.Package, "CloudWatchMetricsQuery"),
		ast.NewRef(schema.Package, "CloudWatchLogsQuery"),
		ast.NewRef(schema.Package, "CloudWatchAnnotationQuery"),
	})
	cloudwatchQuery.Hints[ast.HintImplementsVariant] = string(ast.SchemaVariantDataQuery)

	cloudwatchQuery.Disjunction.Discriminator = "queryMode"
	cloudwatchQuery.Disjunction.DiscriminatorMapping = map[string]string{
		"Metrics":     "CloudWatchMetricsQuery",
		"Logs":        "CloudWatchLogsQuery",
		"Annotations": "CloudWatchAnnotationQuery",
	}

	return ast.NewObject(schema.Package, "CloudWatchQuery", cloudwatchQuery)
}

func (pass *Cloudwatch) processQueryEditorExpression(object ast.Object) ast.Object {
	if object.Type.Kind != ast.KindDisjunction {
		return object
	}

	object.Type.Disjunction.Discriminator = "type"
	object.Type.Disjunction.DiscriminatorMapping = map[string]string{
		"and":               "QueryEditorArrayExpression",
		"or":                "QueryEditorArrayExpression",
		"property":          "QueryEditorPropertyExpression",
		"groupBy":           "QueryEditorGroupByExpression",
		"function":          "QueryEditorFunctionExpression",
		"functionParameter": "QueryEditorFunctionParameterExpression",
		"operator":          "QueryEditorOperatorExpression",
	}

	return object
}

func (pass *Cloudwatch) processQueryEditorArrayExpression(object ast.Object) ast.Object {
	if object.Type.Kind != ast.KindStruct {
		return object
	}

	structDef := object.Type.AsStruct()
	fields := make([]ast.StructField, 0, len(structDef.Fields)-1)

	for _, field := range structDef.Fields {
		if field.Name != "expressions" {
			fields = append(fields, field)
			continue
		}

		newField := field.DeepCopy()
		newField.Type = newField.Type.Disjunction.Branches[0]

		fields = append(fields, newField)
	}

	object.Type.Struct.Fields = fields

	return object
}
